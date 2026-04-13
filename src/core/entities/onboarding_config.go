package entities

// OnboardingConfig represents the "onboarding" block inside the org's JSONB config.
//
// The onboarding form is tenant-configurable: each organization declares its own
// profile fields and tour steps in `organizations.config.onboarding`, and the frontend
// renders the form dynamically from that config. The user's submitted values are
// stored as JSONB in `users.profile_data` and validated against this schema in
// the SaveProfile usecase. This avoids hardcoding the profile shape per tenant.
type OnboardingConfig struct {
	SkipAllowed   bool             `json:"skip_allowed"`
	ProfileFields []ProfileField   `json:"profile_fields"`
	TourSteps     []TourStepConfig `json:"tour_steps"`
}

// ProfileFieldType is the set of supported input types for a profile field.
type ProfileFieldType string

const (
	ProfileFieldText        ProfileFieldType = "text"
	ProfileFieldNumber      ProfileFieldType = "number"
	ProfileFieldSelect      ProfileFieldType = "select"
	ProfileFieldMultiselect ProfileFieldType = "multiselect"
)

// ProfileField defines a single field in the org's onboarding profile form.
// The set of fields is configured per tenant (see OnboardingConfig).
type ProfileField struct {
	Key      string           `json:"key"`
	Label    string           `json:"label"`
	Type     ProfileFieldType `json:"type"`
	Options  []string         `json:"options,omitempty"`
	Required bool             `json:"required"`
}

// TourStepConfig defines a step in the onboarding product tour.
type TourStepConfig struct {
	Key             string   `json:"key"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Order           int      `json:"order"`
	Roles           []string `json:"roles,omitempty"`
	RequiresFeature string   `json:"requires_feature,omitempty"`
}
