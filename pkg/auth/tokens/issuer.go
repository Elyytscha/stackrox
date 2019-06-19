package tokens

import (
	"context"

	"github.com/pkg/errors"
)

// Issuer is an interface for issuing tokens, tied to a token source.
type Issuer interface {
	// Issue issues a token for the given claims, applying any of the specified options.
	Issue(ctx context.Context, claims RoxClaims, options ...Option) (*TokenInfo, error)
}

type issuerForSource struct {
	source  Source
	factory *issuerFactory
	options []Option
}

func (i *issuerForSource) Issue(ctx context.Context, roxClaims RoxClaims, options ...Option) (*TokenInfo, error) {
	claims := i.factory.createClaims(i.source.ID(), roxClaims)
	// issuer options may override options that were passed as arguments to this function, and issuer factory options
	// may override both of these types of options. Hence, apply options in ascending order of priority.
	for _, opt := range options {
		opt.apply(claims)
	}
	for _, opt := range i.options {
		opt.apply(claims)
	}
	for _, opt := range i.factory.options {
		opt.apply(claims)
	}
	// Sanity check: the source should accept the token.
	if err := i.source.Validate(ctx, claims); err != nil {
		return nil, errors.Wrap(err, "issued token was rejected by source")
	}

	token, err := i.factory.encode(claims)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode token")
	}

	return &TokenInfo{
		Token:   token,
		Claims:  claims,
		Sources: []Source{i.source},
	}, nil
}
