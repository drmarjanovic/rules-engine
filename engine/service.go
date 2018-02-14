package engine

import (
	"errors"
	"github.com/mainflux/mainflux/writer"
)

var (
	// ErrMalformedUrl indicates malformed URL specification.
	ErrMalformedUrl error = errors.New("malformed url specification")

	// ErrMalformedEntity indicates malformed entity specification.
	ErrMalformedEntity error = errors.New("malformed entity specification")

	// ErrNotFound indicates a non-existent entity request.
	ErrNotFound error = errors.New("non-existent entity")
)

// Service specifies an API that must be fulfilled by domain service implementation.
type Service interface {
	// Save specific rule.
	SaveRule(Rule) error

	// ViewRule retrieves specific rule using unique identifiers of user and rule.
	ViewRule(string, string) (*Rule, error)

	// ListRules retrieves data about all rules that belongs to specific user
	// identified by user unique identifier.
	ListRules(string) ([]Rule, error)

	// RemoveRule removes specific rule identified by the user's unique identifier
	// and rule's unique identifier.
	RemoveRule(string, string) error

	// ApplyRules checks which events satisfy which rules and execute related actions
	// for satisfied rules.
	ApplyRules(userId string, events []writer.Message) error
}
