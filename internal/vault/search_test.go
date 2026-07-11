package vault_test

import (
	"context"
	"errors"
	"sort"
	"sync"
	"testing"

	"github.com/dkyanakiev/vaul7y/internal/models"
	"github.com/dkyanakiev/vaul7y/internal/vault"
	"github.com/dkyanakiev/vaul7y/internal/vault/vaultfakes"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

func listSecret(keys ...interface{}) *api.Secret {
	return &api.Secret{Data: map[string]interface{}{"keys": keys}}
}

func kv2Secret(data map[string]interface{}) *api.Secret {
	return &api.Secret{Data: map[string]interface{}{"data": data}}
}

// fakeSearchTree wires a FakeLogical with this layout (KV v2 paths):
//
//	secret/
//	├── top-token                   {token: abc}
//	├── app/
//	│   ├── db-creds                {db_password: hunter2, note: "token here"}
//	│   ├── locked                  read → permission denied
//	│   └── denied-dir/             list → permission denied
func fakeSearchTree() *vaultfakes.FakeLogical {
	fakeLogical := &vaultfakes.FakeLogical{}

	fakeLogical.ListStub = func(path string) (*api.Secret, error) {
		switch path {
		case "secret/metadata":
			return listSecret("app/", "top-token"), nil
		case "secret/metadata/app":
			return listSecret("db-creds", "locked", "denied-dir/"), nil
		case "secret/metadata/app/denied-dir":
			return nil, errors.New("permission denied")
		}
		return nil, nil
	}

	fakeLogical.ReadStub = func(path string) (*api.Secret, error) {
		switch path {
		case "secret/data/top-token":
			return kv2Secret(map[string]interface{}{"token": "abc"}), nil
		case "secret/data/app/db-creds":
			return kv2Secret(map[string]interface{}{"db_password": "hunter2", "note": "token here"}), nil
		case "secret/data/app/locked":
			return nil, errors.New("permission denied")
		}
		return nil, nil
	}

	return fakeLogical
}

func collectMatches() (func(models.SearchMatch), *[]models.SearchMatch) {
	var mu sync.Mutex
	matches := &[]models.SearchMatch{}
	return func(m models.SearchMatch) {
		mu.Lock()
		defer mu.Unlock()
		*matches = append(*matches, m)
	}, matches
}

func sortedMatches(matches []models.SearchMatch) []models.SearchMatch {
	sorted := append([]models.SearchMatch{}, matches...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Path != sorted[j].Path {
			return sorted[i].Path < sorted[j].Path
		}
		return sorted[i].Type < sorted[j].Type
	})
	return sorted
}

func TestSearchKV_KeysAndPaths(t *testing.T) {
	v := &vault.Vault{
		Logical: fakeSearchTree(),
		Logger:  discardLogger(),
	}

	onMatch, matches := collectMatches()
	stats, err := v.SearchKV(context.Background(), models.SearchOptions{
		Mount:     "secret",
		Term:      "TOKEN",
		MatchKeys: true,
	}, onMatch)

	require.NoError(t, err)
	require.Equal(t, []models.SearchMatch{
		{Path: "top-token", Key: "token", Type: models.SearchMatchKey},
		{Path: "top-token", Key: "", Type: models.SearchMatchPath},
	}, sortedMatches(*matches))

	require.Equal(t, 2, stats.Matches)
	require.Equal(t, 2, stats.SecretsScanned)
	require.Equal(t, 1, stats.ListDenied, "denied-dir list failure should be counted, not fatal")
	require.Equal(t, 1, stats.ReadDenied, "locked read failure should be counted, not fatal")
	require.False(t, stats.Truncated)
}

func TestSearchKV_Values(t *testing.T) {
	v := &vault.Vault{
		Logical: fakeSearchTree(),
		Logger:  discardLogger(),
	}

	onMatch, matches := collectMatches()
	stats, err := v.SearchKV(context.Background(), models.SearchOptions{
		Mount:       "secret",
		Term:        "token",
		MatchKeys:   true,
		MatchValues: true,
	}, onMatch)

	require.NoError(t, err)
	require.Equal(t, []models.SearchMatch{
		{Path: "app/db-creds", Key: "note", Type: models.SearchMatchValue},
		{Path: "top-token", Key: "token", Type: models.SearchMatchKey},
		{Path: "top-token", Key: "", Type: models.SearchMatchPath},
	}, sortedMatches(*matches))
	require.Equal(t, 3, stats.Matches)
}

func TestSearchKV_ValuesOff(t *testing.T) {
	v := &vault.Vault{
		Logical: fakeSearchTree(),
		Logger:  discardLogger(),
	}

	onMatch, matches := collectMatches()
	_, err := v.SearchKV(context.Background(), models.SearchOptions{
		Mount:     "secret",
		Term:      "hunter2",
		MatchKeys: true,
	}, onMatch)

	require.NoError(t, err)
	require.Empty(t, *matches, "values must not match when MatchValues is off")
}

func TestSearchKV_NestedData(t *testing.T) {
	fakeLogical := &vaultfakes.FakeLogical{}
	fakeLogical.ListStub = func(path string) (*api.Secret, error) {
		if path == "secret/metadata" {
			return listSecret("nested"), nil
		}
		return nil, nil
	}
	fakeLogical.ReadReturns(kv2Secret(map[string]interface{}{
		"outer": map[string]interface{}{"inner_token": true},
	}), nil)

	v := &vault.Vault{
		Logical: fakeLogical,
		Logger:  discardLogger(),
	}

	onMatch, matches := collectMatches()
	_, err := v.SearchKV(context.Background(), models.SearchOptions{
		Mount:     "secret",
		Term:      "inner_token",
		MatchKeys: true,
	}, onMatch)

	require.NoError(t, err)
	require.Equal(t, []models.SearchMatch{
		{Path: "nested", Key: "inner_token", Type: models.SearchMatchKey},
	}, *matches)
}

func TestSearchKV_Truncation(t *testing.T) {
	v := &vault.Vault{
		Logical: fakeSearchTree(),
		Logger:  discardLogger(),
	}

	onMatch, matches := collectMatches()
	stats, err := v.SearchKV(context.Background(), models.SearchOptions{
		Mount:      "secret",
		Term:       "token",
		MatchKeys:  true,
		MaxResults: 1,
	}, onMatch)

	require.NoError(t, err, "truncation is not an error")
	require.True(t, stats.Truncated)
	require.Len(t, *matches, 1)
}

func TestSearchKV_Cancelled(t *testing.T) {
	v := &vault.Vault{
		Logical: fakeSearchTree(),
		Logger:  discardLogger(),
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	onMatch, matches := collectMatches()
	_, err := v.SearchKV(ctx, models.SearchOptions{
		Mount:     "secret",
		Term:      "token",
		MatchKeys: true,
	}, onMatch)

	require.ErrorIs(t, err, context.Canceled)
	require.Empty(t, *matches)
}

func TestSearchKV_EmptyTerm(t *testing.T) {
	v := &vault.Vault{
		Logical: &vaultfakes.FakeLogical{},
		Logger:  discardLogger(),
	}

	_, err := v.SearchKV(context.Background(), models.SearchOptions{Mount: "secret"}, func(models.SearchMatch) {})
	require.Error(t, err)
}

func TestSearchKV_StartPath(t *testing.T) {
	fakeLogical := fakeSearchTree()
	v := &vault.Vault{
		Logical: fakeLogical,
		Logger:  discardLogger(),
	}

	onMatch, matches := collectMatches()
	stats, err := v.SearchKV(context.Background(), models.SearchOptions{
		Mount:     "secret",
		Path:      "app/",
		Term:      "db",
		MatchKeys: true,
	}, onMatch)

	require.NoError(t, err)
	require.Equal(t, []models.SearchMatch{
		{Path: "app/db-creds", Key: "db_password", Type: models.SearchMatchKey},
		{Path: "app/db-creds", Key: "", Type: models.SearchMatchPath},
	}, sortedMatches(*matches))
	require.Equal(t, 1, stats.ListDenied, "still counts list skips under the sub-tree")
	require.Equal(t, 1, stats.ReadDenied, "still counts read skips under the sub-tree")
	require.Equal(t, "secret/metadata/app", fakeLogical.ListArgsForCall(0), "search starts at the given sub-path")
}
