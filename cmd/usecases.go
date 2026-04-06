package main

// UseCases holds all application use cases.
// Wired incrementally as features are implemented.
type UseCases struct {
	// admin usecases (Épica 2-3)
	// coordination usecases (Épica 4)
	// teaching usecases (Épica 5)
	// resources usecases (Épica 8)
}

func NewUseCases(_ *Repositories) *UseCases {
	return &UseCases{}
}
