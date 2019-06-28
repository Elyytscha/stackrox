package sac

import (
	"context"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/pkg/features"
)

var (
	// noSAC stores whether scoped access control is disabled. This is to prevent
	// a (relatively expensive) call to `Enabled()` on every data access.
	noSAC = !features.ScopedAccessControl.Enabled()
)

type globalAccessScopeContextKey struct{}

// GlobalAccessScopeChecker retrieves the global access scope from the context.
// This function is guaranteed to return a non-nil value.
func GlobalAccessScopeChecker(ctx context.Context) ScopeChecker {
	if noSAC {
		return NewScopeChecker(AllowAllAccessScopeChecker())
	}

	core, _ := ctx.Value(globalAccessScopeContextKey{}).(ScopeCheckerCore)
	if core == nil {
		core = ErrorAccessScopeCheckerCore(errors.New("global access scope was not found in context"))
	}
	return NewScopeChecker(core)
}

// GlobalAccessScopeCheckerOrNil retrieves the global access scope checker if it is stored in the context; otherwise,
// it returns `nil`.
func GlobalAccessScopeCheckerOrNil(ctx context.Context) *ScopeChecker {
	if noSAC {
		return nil
	}

	core, _ := ctx.Value(globalAccessScopeContextKey{}).(ScopeCheckerCore)
	if core == nil {
		return nil
	}
	sc := NewScopeChecker(core)
	return &sc
}

// WithGlobalAccessScopeChecker returns a context that is a child of the given context and contains
// the given global access scope.
func WithGlobalAccessScopeChecker(ctx context.Context, as ScopeCheckerCore) context.Context {
	return context.WithValue(ctx, globalAccessScopeContextKey{}, as)
}

// WithAllAccess returns a context that is a child of the given context and will allow all access
// scope checks.
func WithAllAccess(ctx context.Context) context.Context {
	return WithGlobalAccessScopeChecker(ctx, allowAllScopeCheckerCore)
}

// WithNoAccess returns a context that is a child of the given context and will deny all access
// scope checks.
func WithNoAccess(ctx context.Context) context.Context {
	return WithGlobalAccessScopeChecker(ctx, denyAllScopeCheckerCore)
}
