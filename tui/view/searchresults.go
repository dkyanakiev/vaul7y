package view

import (
	"context"
	"time"

	"github.com/dkyanakiev/vaul7y/internal/models"
	"github.com/dkyanakiev/vaul7y/tui/component"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const searchRefreshInterval = 300 * time.Millisecond

// startRecursiveSearch kicks off a recursive KV search from the current path
// and switches to the results view. Results stream in as the walker finds
// them; folders and secrets the token cannot access are skipped and counted.
func (v *View) startRecursiveSearch(term string) {
	v.state.Lock()
	v.state.SearchTerm = term
	v.state.SearchResults = nil
	v.state.SearchStats = models.SearchStats{}
	v.state.SearchRunning = true
	v.state.SearchBasePath = v.state.SelectedPath
	mount := v.state.SelectedMount
	basePath := v.state.SelectedPath
	includeValues := v.state.SearchIncludeValues
	v.state.Unlock()

	// Switches the view first: its viewSwitch() cancels any previous search,
	// so the context for this run must be created after it.
	v.SearchResults()

	ctx, cancel := context.WithCancel(context.Background())
	v.searchCancel = cancel

	go func() {
		opts := models.SearchOptions{
			Mount:       mount,
			Path:        basePath,
			Term:        term,
			MatchKeys:   true,
			MatchValues: includeValues,
		}

		stats, err := v.Client.SearchKV(ctx, opts, func(m models.SearchMatch) {
			v.state.Lock()
			v.state.SearchResults = append(v.state.SearchResults, m)
			v.state.Unlock()
		})
		if err != nil && ctx.Err() == nil {
			v.logger.Warn().Err(err).Msg("recursive search failed")
		}

		v.state.Lock()
		v.state.SearchStats = stats
		v.state.SearchRunning = false
		v.state.Unlock()

		v.Layout.Container.QueueUpdateDraw(func() {
			// The user may have navigated away while the search was running;
			// don't draw the results table over another view.
			if ctx.Err() != nil {
				return
			}
			v.renderSearchResults()
		})
	}()

	// Re-render periodically while the search streams results in.
	go func() {
		ticker := time.NewTicker(searchRefreshInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				v.state.RLock()
				running := v.state.SearchRunning
				v.state.RUnlock()
				if !running {
					return
				}
				v.Layout.Container.QueueUpdateDraw(func() {
					if ctx.Err() != nil {
						return
					}
					v.renderSearchResults()
				})
			}
		}
	}()
}

func (v *View) SearchResults() {
	v.viewSwitch()
	v.Layout.Body.Clear()
	v.Layout.Body.SetTitle("Search Results")
	v.Layout.Container.SetFocus(v.components.SearchResultsTable.Table.Primitive())
	v.Layout.Container.SetInputCapture(v.InputSearchResults)
	v.components.Commands.Update(component.SearchResultsCommands)

	v.renderSearchResults()

	v.state.Elements.TableMain = v.components.SearchResultsTable.Table.Primitive().(*tview.Table)
}

func (v *View) renderSearchResults() {
	v.state.RLock()
	data := append([]models.SearchMatch{}, v.state.SearchResults...)
	term := v.state.SearchTerm
	running := v.state.SearchRunning
	stats := v.state.SearchStats
	mount := v.state.SelectedMount
	v.state.RUnlock()

	t := v.components.SearchResultsTable
	t.Props.Data = data
	t.Props.Term = term
	t.Props.Running = running
	t.Props.Stats = stats
	t.Props.SelectedMount = mount
	t.Render()
	v.Draw()
}

func (v *View) inputSearchResults(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return event
	}

	switch event.Key() {
	case tcell.KeyEsc:
		v.searchResultsGoBack()
		return nil
	case tcell.KeyEnter:
		v.openSearchResult()
		return nil
	case tcell.KeyRune:
		switch event.Rune() {
		case 'e':
			v.openSearchResult()
			return nil
		case 'b':
			v.searchResultsGoBack()
			return nil
		case 'v':
			v.state.Lock()
			v.state.SearchIncludeValues = !v.state.SearchIncludeValues
			term := v.state.SearchTerm
			v.state.Unlock()
			v.startRecursiveSearch(term)
			return nil
		}
	}

	return event
}

// openSearchResult navigates to the secret of the selected match.
func (v *View) openSearchResult() {
	path := v.components.SearchResultsTable.GetPathForSelection()
	if path == "" {
		return
	}

	v.state.Lock()
	v.state.SelectedPath = path
	mount := v.state.SelectedMount
	v.state.Unlock()

	v.SecretObject(mount, path)
}

// searchResultsGoBack returns to the secrets browser at the path the search
// was started from.
func (v *View) searchResultsGoBack() {
	v.state.Lock()
	base := v.state.SearchBasePath
	v.state.SelectedPath = base
	v.state.SelectedObject = ""
	v.state.Unlock()

	v.components.SecretsTable.Props.SelectedPath = base
	v.Secrets(base, "false")
}
