package providers

import (
	bcerrors "github.com/educabot/team-ai-toolkit/errors"
)

// Re-export shared errors from team-ai-toolkit
var (
	ErrNotFound     = bcerrors.ErrNotFound
	ErrValidation   = bcerrors.ErrValidation
	ErrUnauthorized = bcerrors.ErrUnauthorized
	ErrForbidden    = bcerrors.ErrForbidden
	ErrDuplicate    = bcerrors.ErrDuplicate
	ErrConflict     = bcerrors.ErrConflict
)

// Alizia-specific errors
var (
	ErrDocNotFound      = bcerrors.New("coordination document not found")
	ErrTopicMaxLevel    = bcerrors.New("topic exceeds max level")
	ErrSharedClassLimit = bcerrors.New("shared classes limit exceeded")
)
