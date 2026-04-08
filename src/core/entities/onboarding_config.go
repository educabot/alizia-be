package entities

// OnboardingConfig represents the "onboarding" block inside the org's JSONB config.
type OnboardingConfig struct {
	SkipAllowed   bool             `json:"skip_allowed"`
	ProfileFields []ProfileField   `json:"profile_fields"`
	TourSteps     []TourStepConfig `json:"tour_steps"`
}

// ProfileField defines a dynamic field for the onboarding profile form.
type ProfileField struct {
	Key      string   `json:"key"`
	Label    string   `json:"label"`
	Type     string   `json:"type"` // text, number, select, multiselect
	Options  []string `json:"options,omitempty"`
	Required bool     `json:"required"`
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
