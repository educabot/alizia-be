package entities

import (
	"encoding/json"

	"gorm.io/datatypes"
)

// DefaultTopicMaxLevels is the depth cap applied when an organization's config
// does not specify `topic_max_levels` or stores an invalid value (<=0).
const DefaultTopicMaxLevels = 3

// OrgConfig is the typed view of organizations.config (JSONB). All consumers
// that need values from that column MUST go through ParseOrgConfig rather than
// ad-hoc unmarshalling, so defaults and fallbacks stay consistent.
//
// The JSONB column is intentionally open for tenant-specific keys beyond what
// this struct knows about — mapOrganization in entrypoints/admin.go still
// surfaces the raw map on the API. Adding a field here only types one more key
// for internal use; it does not close the shape.
type OrgConfig struct {
	TopicMaxLevels          int              `json:"topic_max_levels,omitempty"`
	SharedClassesEnabled    bool             `json:"shared_classes_enabled,omitempty"`
	DesarrolloMaxActivities int              `json:"desarrollo_max_activities,omitempty"`
	Features                map[string]bool  `json:"features,omitempty"`
	Onboarding              OnboardingConfig `json:"onboarding,omitempty"`
}

// ParseOrgConfig decodes the JSONB column and applies defaults. Malformed JSON
// is not an error — the column is tenant-controlled, and a bad write should
// not brick every usecase that reads it. In that case we return defaults plus
// SkipAllowed=true so the onboarding UI still works while the bad config is
// fixed out-of-band.
func ParseOrgConfig(raw datatypes.JSON) OrgConfig {
	if len(raw) == 0 {
		return OrgConfig{TopicMaxLevels: DefaultTopicMaxLevels}
	}
	var c OrgConfig
	if err := json.Unmarshal(raw, &c); err != nil {
		return OrgConfig{
			TopicMaxLevels: DefaultTopicMaxLevels,
			Onboarding:     OnboardingConfig{SkipAllowed: true},
		}
	}
	if c.TopicMaxLevels <= 0 {
		c.TopicMaxLevels = DefaultTopicMaxLevels
	}
	return c
}

// IsFeatureActive reports whether the named feature flag is set to true in the
// org config. Returns false if Features is nil or the flag is absent.
func (c OrgConfig) IsFeatureActive(name string) bool {
	if c.Features == nil {
		return false
	}
	return c.Features[name]
}
