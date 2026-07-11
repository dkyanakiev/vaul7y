package models

// SearchMatchType describes what part of a secret matched a recursive search.
const (
	SearchMatchPath  = "path"
	SearchMatchKey   = "key"
	SearchMatchValue = "value"
)

// SearchOptions configures a recursive KV search.
type SearchOptions struct {
	Mount string
	// Path is the sub-path within the mount to start from ("" = mount root).
	Path string
	Term string
	// MatchKeys matches key names inside secrets (requires read access).
	MatchKeys bool
	// MatchValues matches values inside secrets (requires read access).
	MatchValues bool
	// Concurrency bounds parallel Vault requests (default 5).
	Concurrency int
	// MaxResults stops the search early once reached (default 200).
	MaxResults int
}

// SearchMatch is a single recursive search hit.
type SearchMatch struct {
	// Path of the secret within the mount (no mount prefix).
	Path string
	// Key that matched; empty for path matches.
	Key string
	// Type is one of SearchMatchPath, SearchMatchKey, SearchMatchValue.
	Type string
}

// SearchStats summarizes a finished (or aborted) recursive search.
type SearchStats struct {
	SecretsScanned int
	Matches        int
	// ListDenied counts folders that could not be listed (no access or error).
	ListDenied int
	// ReadDenied counts secrets that could not be read (no access or error).
	ReadDenied int
	// Truncated is true when the search stopped at MaxResults.
	Truncated bool
}
