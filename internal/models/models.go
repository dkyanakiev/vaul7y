package models

import "time"

type HandlerFunc func(format string, args ...interface{})

type Handler string

const (
	HandleError Handler = Handler("Error")
	HandleFatal Handler = Handler("Fatal")
	HandleInfo  Handler = Handler("Info")
)

type MountOutput struct {
	UUID                  string            `json:"uuid"`
	Type                  string            `json:"type"`
	Description           string            `json:"description"`
	Accessor              string            `json:"accessor"`
	Config                MountConfigOutput `json:"config"`
	Options               map[string]string `json:"options"`
	Local                 bool              `json:"local"`
	SealWrap              bool              `json:"seal_wrap" mapstructure:"seal_wrap"`
	ExternalEntropyAccess bool              `json:"external_entropy_access" mapstructure:"external_entropy_access"`
	PluginVersion         string            `json:"plugin_version" mapstructure:"plugin_version"`
	RunningVersion        string            `json:"running_plugin_version" mapstructure:"running_plugin_version"`
	RunningSha256         string            `json:"running_sha256" mapstructure:"running_sha256"`
	DeprecationStatus     string            `json:"deprecation_status" mapstructure:"deprecation_status"`
}

type MountConfigOutput struct {
	DefaultLeaseTTL           int                      `json:"default_lease_ttl" mapstructure:"default_lease_ttl"`
	MaxLeaseTTL               int                      `json:"max_lease_ttl" mapstructure:"max_lease_ttl"`
	ForceNoCache              bool                     `json:"force_no_cache" mapstructure:"force_no_cache"`
	AuditNonHMACRequestKeys   []string                 `json:"audit_non_hmac_request_keys,omitempty" mapstructure:"audit_non_hmac_request_keys"`
	AuditNonHMACResponseKeys  []string                 `json:"audit_non_hmac_response_keys,omitempty" mapstructure:"audit_non_hmac_response_keys"`
	ListingVisibility         string                   `json:"listing_visibility,omitempty" mapstructure:"listing_visibility"`
	PassthroughRequestHeaders []string                 `json:"passthrough_request_headers,omitempty" mapstructure:"passthrough_request_headers"`
	AllowedResponseHeaders    []string                 `json:"allowed_response_headers,omitempty" mapstructure:"allowed_response_headers"`
	TokenType                 string                   `json:"token_type,omitempty" mapstructure:"token_type"`
	AllowedManagedKeys        []string                 `json:"allowed_managed_keys,omitempty" mapstructure:"allowed_managed_keys"`
	UserLockoutConfig         *UserLockoutConfigOutput `json:"user_lockout_config,omitempty"`
	// Deprecated: This field will always be blank for newer server responses.
	PluginName string `json:"plugin_name,omitempty" mapstructure:"plugin_name"`
}

type UiMountsResponse struct {
	Data struct {
		Secret map[string]*MountOutput `json:"secret"`
	} `json:"data"`
}

type UserLockoutConfigOutput struct {
	LockoutThreshold    uint  `json:"lockout_threshold,omitempty" structs:"lockout_threshold" mapstructure:"lockout_threshold"`
	LockoutDuration     int   `json:"lockout_duration,omitempty" structs:"lockout_duration" mapstructure:"lockout_duration"`
	LockoutCounterReset int   `json:"lockout_counter_reset,omitempty" structs:"lockout_counter_reset" mapstructure:"lockout_counter_reset"`
	DisableLockout      *bool `json:"disable_lockout,omitempty" structs:"disable_lockout" mapstructure:"disable_lockout"`
}

type KVSecret struct {
	Data            map[string]interface{}
	VersionMetadata *KVVersionMetadata
	CustomMetadata  map[string]interface{}
	Raw             *Secret
}

type KVVersionMetadata struct {
	Version      int       `mapstructure:"version"`
	CreatedTime  time.Time `mapstructure:"created_time"`
	DeletionTime time.Time `mapstructure:"deletion_time"`
	Destroyed    bool      `mapstructure:"destroyed"`
}

type Secret struct {
	// The request ID that generated this response
	RequestID string `json:"request_id"`

	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`

	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data map[string]interface{} `json:"data"`

	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`

	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	Auth *SecretAuth `json:"auth,omitempty"`

	// WrapInfo, if non-nil, means that the initial response was wrapped in the
	// cubbyhole of the given token (which has a TTL of the given number of
	// seconds)
	WrapInfo *SecretWrapInfo `json:"wrap_info,omitempty"`
}

type SecretPath struct {
	PathName string
	IsSecret bool
}

type SecretAuth struct {
	ClientToken      string            `json:"client_token"`
	Accessor         string            `json:"accessor"`
	Policies         []string          `json:"policies"`
	TokenPolicies    []string          `json:"token_policies"`
	IdentityPolicies []string          `json:"identity_policies"`
	Metadata         map[string]string `json:"metadata"`
	Orphan           bool              `json:"orphan"`
	EntityID         string            `json:"entity_id"`

	LeaseDuration int  `json:"lease_duration"`
	Renewable     bool `json:"renewable"`

	//MFARequirement *MFARequirement `json:"mfa_requirement"`
}

type SecretWrapInfo struct {
	Token           string    `json:"token"`
	Accessor        string    `json:"accessor`
	TTL             int       `json:"ttl"`
	CreationTime    time.Time `json:"creation_time"`
	CreationPath    string    `json:"creation_path"`
	WrappedAccessor string    `json:"wrapped_accessor"`
}

type VaultListRequest struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Keys []string `json:"keys"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type ApiSecret struct {
	Data struct {
		Keys []string `json:"keys"`
	} `json:"data"`
}
type sentinelPolicy struct{}

type Policy struct {
	sentinelPolicy
	Name   string `json:"name"`
	Policy string `json:"policy"`
}

type PolicyInfo struct {
	Name          string `json:"name"`
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Data          struct {
		Name   string `json:"name"`
		Policy string `json:"policy"`
	} `json:"data"`
	WrapInfo interface{} `json:"wrap_info"`
	Warnings interface{} `json:"warnings"`
	Auth     interface{} `json:"auth"`
}

type Comp string

func (s Comp) Error() string {
	return string(s)
}

const (
	MountTypeKV        = "kv"
	MountTypeSystem    = "system"
	MountTypePki       = "pki"
	MountTypeIdentity  = "identity"
	MountTypeCubbyhole = "cubbyhole"
)

type MountMap map[string]Mount
type MountConfig struct {
	DefaultLeaseTTL int  `json:"default_lease_ttl"`
	ForceNoCache    bool `json:"force_no_cache"`
	MaxLeaseTTL     int  `json:"max_lease_ttl"`
}

type Mount struct {
	Accessor              string      `json:"accessor"`
	Config                MountConfig `json:"config"`
	Description           string      `json:"description"`
	ExternalEntropyAccess bool        `json:"external_entropy_access"`
	Local                 bool        `json:"local"`
	Options               interface{} `json:"options"`
	PluginVersion         string      `json:"plugin_version"`
	RunningPluginVersion  string      `json:"running_plugin_version"`
	RunningSha256         string      `json:"running_sha256"`
	SealWrap              bool        `json:"seal_wrap"`
	Type                  string      `json:"type"`
	UUID                  string      `json:"uuid"`
}

type Namespace struct {
	Name        string
	Description string
}

type MetaResponse struct {
	Data Metadata `json:"data"`
}

type Metadata struct {
	CasRequired        bool                   `json:"cas_required"`
	CreatedTime        string                 `json:"created_time"`
	CurrentVersion     int                    `json:"current_version"`
	DeleteVersionAfter string                 `json:"delete_version_after"`
	MaxVersions        int                    `json:"max_versions"`
	OldestVersion      int                    `json:"oldest_version"`
	UpdatedTime        string                 `json:"updated_time"`
	CustomMetadata     map[string]interface{} `json:"custom_metadata"`
	Versions           map[string]Version     `json:"versions"`
}

type Version struct {
	CreatedTime  string `json:"created_time"`
	DeletionTime string `json:"deletion_time"`
	Destroyed    bool   `json:"destroyed"`
}
