// Package nebraska provides types for the Nebraska update server's data model.
// These types map to the Nebraska Publisher API at /api/*.
package nebraska

import "time"

// Package represents a versioned release artifact registered in Nebraska.
type Package struct {
	ID                string `json:"id,omitempty"`
	Type              int    `json:"type"`
	Version           string `json:"version"`
	URL               string `json:"url"`
	Filename          string `json:"filename"`
	Description       string `json:"description"`
	Size              string `json:"size"`
	Hash              string `json:"hash"`
	ApplicationID     string `json:"application_id"`
	FlatcarAction     *FlatcarAction `json:"flatcar_action,omitempty"`
	Arch              int    `json:"arch"`
	ChannelsBlacklist []string `json:"channels_blacklist"`
	CreatedTs         string `json:"created_ts,omitempty"`
}

// FlatcarAction describes the update action for Omaha-compatible clients.
type FlatcarAction struct {
	ID                  string `json:"id,omitempty"`
	Event               string `json:"event"`                // "postinstall"
	Sha256              string `json:"sha256"`
	NeedsAdmin          bool   `json:"needs_admin"`
	IsDelta             bool   `json:"is_delta"`
	DisablePayloadBackoff bool `json:"disable_payload_backoff"`
}

// Application represents a product registered in Nebraska.
type Application struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedTs   string    `json:"created_ts,omitempty"`
	Channels    []Channel `json:"channels,omitempty"`
	Groups      []Group   `json:"groups,omitempty"`
	Packages    []Package `json:"packages,omitempty"`
}

// Channel represents a release track within an application.
// Each channel points to exactly one package (the current version for that track).
type Channel struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Color         string `json:"color"`
	Arch          int    `json:"arch"`
	ApplicationID string `json:"application_id"`
	PackageID     string `json:"package_id"`
	CreatedTs     string `json:"created_ts,omitempty"`
}

// Group represents a set of instances that query a specific channel for updates.
// Nodes reference a group ID in their Omaha update requests.
type Group struct {
	ID                          string `json:"id,omitempty"`
	Name                        string `json:"name"`
	Description                 string `json:"description"`
	ApplicationID               string `json:"application_id"`
	ChannelID                   string `json:"channel_id"`
	PolicyUpdatesEnabled        bool   `json:"policy_updates_enabled"`
	PolicySafeMode              bool   `json:"policy_safe_mode"`
	PolicyOfficeHours           bool   `json:"policy_office_hours"`
	PolicyTimezone              string `json:"policy_timezone"`
	PolicyPeriodInterval        string `json:"policy_period_interval"`
	PolicyMaxUpdatesPerPeriod   int    `json:"policy_max_updates_per_period"`
	PolicyUpdateTimeout         string `json:"policy_update_timeout"`
	CreatedTs                   string `json:"created_ts,omitempty"`
}

// ArchAMD64 is the Nebraska architecture code for x86_64.
const ArchAMD64 = 1

// ArchARM64 is the Nebraska architecture code for ARM64.
const ArchARM64 = 2

// PackageTypeFlatcar is the package type for Omaha/flatcar-style updates.
const PackageTypeFlatcar = 1

// Timestamp helper for Nebraska API format.
func FormatTimestamp(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
