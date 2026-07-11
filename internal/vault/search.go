package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/dkyanakiev/vaul7y/internal/models"
)

const (
	defaultSearchConcurrency = 5
	defaultSearchMaxResults  = 200
)

// searchRun holds the shared state of one recursive search. All mutations go
// through the mutex so onMatch is never called concurrently.
type searchRun struct {
	mu      sync.Mutex
	stats   models.SearchStats
	onMatch func(models.SearchMatch)
	cancel  context.CancelFunc
	max     int
}

func (r *searchRun) emit(m models.SearchMatch) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.stats.Truncated {
		return
	}
	r.stats.Matches++
	r.onMatch(m)
	if r.stats.Matches >= r.max {
		r.stats.Truncated = true
		r.cancel()
	}
}

// SearchKV recursively walks a KV mount starting at opts.Path and reports
// every match through onMatch. Secret names always match; key names and
// values match when opts.MatchKeys / opts.MatchValues are set (both require
// read access to each secret).
//
// Folders that cannot be listed and secrets that cannot be read are skipped
// and counted in the returned stats — a token with partial access still gets
// results for everything it can see.
func (v *Vault) SearchKV(ctx context.Context, opts models.SearchOptions, onMatch func(models.SearchMatch)) (models.SearchStats, error) {
	if opts.Term == "" {
		return models.SearchStats{}, fmt.Errorf("empty search term")
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = defaultSearchConcurrency
	}
	if opts.MaxResults <= 0 {
		opts.MaxResults = defaultSearchMaxResults
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	run := &searchRun{
		onMatch: onMatch,
		cancel:  cancel,
		max:     opts.MaxResults,
	}

	kv2 := v.kvVersionForMount(opts.Mount) == "2"
	term := strings.ToLower(opts.Term)
	sem := make(chan struct{}, opts.Concurrency)

	basePath := strings.TrimPrefix(opts.Path, "/")
	if basePath != "" && !strings.HasSuffix(basePath, "/") {
		basePath += "/"
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go v.searchWalk(ctx, &wg, sem, run, opts, kv2, term, basePath)
	wg.Wait()

	run.mu.Lock()
	stats := run.stats
	run.mu.Unlock()

	// Hitting MaxResults cancels the internal context to stop the walkers;
	// that is a successful (truncated) search, not an error.
	err := ctx.Err()
	if stats.Truncated {
		err = nil
	}
	return stats, err
}

// searchWalk lists one folder, matches and reads its leaf secrets, and
// spawns a goroutine per sub-folder.
func (v *Vault) searchWalk(ctx context.Context, wg *sync.WaitGroup, sem chan struct{}, run *searchRun, opts models.SearchOptions, kv2 bool, term, sub string) {
	defer wg.Done()

	listPath := sanitizePath(fmt.Sprintf("%s/%s", opts.Mount, sub))
	if kv2 {
		listPath = sanitizePath(fmt.Sprintf("%s/metadata/%s", opts.Mount, sub))
	}

	if !acquire(ctx, sem) {
		return
	}
	list, err := v.Logical.List(listPath)
	release(sem)

	if err != nil {
		v.Logger.Debug().Err(err).Msgf("search: skipping unlistable path: %s", listPath)
		run.mu.Lock()
		run.stats.ListDenied++
		run.mu.Unlock()
		return
	}

	keys, ok := extractListData(list)
	if !ok {
		return
	}

	for _, k := range keys {
		if ctx.Err() != nil {
			return
		}

		entry, ok := k.(string)
		if !ok {
			continue
		}
		fullPath := sub + entry

		if strings.HasSuffix(entry, "/") {
			wg.Add(1)
			go v.searchWalk(ctx, wg, sem, run, opts, kv2, term, fullPath)
			continue
		}

		if strings.Contains(strings.ToLower(entry), term) {
			run.emit(models.SearchMatch{Path: fullPath, Type: models.SearchMatchPath})
		}

		if opts.MatchKeys || opts.MatchValues {
			v.searchReadSecret(ctx, sem, run, opts, kv2, term, fullPath)
		}
	}
}

// searchReadSecret reads one secret and matches its keys and/or values.
func (v *Vault) searchReadSecret(ctx context.Context, sem chan struct{}, run *searchRun, opts models.SearchOptions, kv2 bool, term, fullPath string) {
	readPath := sanitizePath(fmt.Sprintf("%s/%s", opts.Mount, fullPath))
	if kv2 {
		readPath = sanitizePath(fmt.Sprintf("%s/data/%s", opts.Mount, fullPath))
	}

	if !acquire(ctx, sem) {
		return
	}
	secret, err := v.Logical.Read(readPath)
	release(sem)

	if err != nil {
		v.Logger.Debug().Err(err).Msgf("search: skipping unreadable secret: %s", readPath)
		run.mu.Lock()
		run.stats.ReadDenied++
		run.mu.Unlock()
		return
	}
	if secret == nil || secret.Data == nil {
		// Deleted or empty version — nothing to match.
		return
	}

	data := secret.Data
	if kv2 {
		nested, ok := secret.Data["data"].(map[string]interface{})
		if !ok {
			// Soft-deleted current version: data is nil but metadata remains.
			return
		}
		data = nested
	}

	run.mu.Lock()
	run.stats.SecretsScanned++
	run.mu.Unlock()

	matchSecretData(run, opts, term, fullPath, data)
}

// matchSecretData walks a secret's data map (including nested maps) matching
// key names and stringified values against the term.
func matchSecretData(run *searchRun, opts models.SearchOptions, term, fullPath string, data map[string]interface{}) {
	for key, value := range data {
		if opts.MatchKeys && strings.Contains(strings.ToLower(key), term) {
			run.emit(models.SearchMatch{Path: fullPath, Key: key, Type: models.SearchMatchKey})
			continue
		}

		if nested, ok := value.(map[string]interface{}); ok {
			matchSecretData(run, opts, term, fullPath, nested)
			continue
		}

		if opts.MatchValues && strings.Contains(strings.ToLower(stringifyValue(value)), term) {
			run.emit(models.SearchMatch{Path: fullPath, Key: key, Type: models.SearchMatchValue})
		}
	}
}

func stringifyValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case json.Number:
		return v.String()
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// acquire takes a slot on the semaphore, giving up when the context is
// cancelled. Returns false if the caller should abort.
func acquire(ctx context.Context, sem chan struct{}) bool {
	select {
	case sem <- struct{}{}:
		return true
	case <-ctx.Done():
		return false
	}
}

func release(sem chan struct{}) {
	<-sem
}
