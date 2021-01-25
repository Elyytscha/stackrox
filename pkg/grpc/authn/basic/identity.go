package basic

import (
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/authproviders"
	"github.com/stackrox/rox/pkg/auth/tokens"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/timeutil"
)

var _ authn.Identity = (*identity)(nil)

// IsBasicIdentity returns whether or not the input Identity is a basic identity.
func IsBasicIdentity(id authn.Identity) bool {
	_, isBasic := id.(Identity)
	return isBasic
}

// Basic identity implementation.
type identity struct {
	username     string
	role         *storage.Role
	authProvider authproviders.Provider
}

func (i identity) UID() string {
	return i.username
}

func (i identity) FriendlyName() string {
	return i.username
}

func (i identity) FullName() string {
	return i.username
}

func (i identity) Permissions() *storage.Role {
	return i.role
}

func (i identity) Roles() []*storage.Role {
	return []*storage.Role{i.role}
}

func (i identity) Service() *storage.ServiceIdentity {
	return nil
}

func (i identity) User() *storage.UserInfo {
	return &storage.UserInfo{
		Username:    i.username,
		Role:        i.role,
		Permissions: i.role,
		Roles:       []*storage.Role{i.role},
	}
}

func (i identity) Expiry() time.Time {
	return timeutil.MaxProtoValid
}

func (i identity) ExternalAuthProvider() authproviders.Provider {
	return i.authProvider
}

func (i identity) isBasicAuthIdentity() {}

func (i identity) AsExternalUser() *tokens.ExternalUserClaim {
	return &tokens.ExternalUserClaim{
		UserID:     i.username,
		Attributes: i.Attributes(),
	}
}

func (i identity) Attributes() map[string][]string {
	return map[string][]string{
		"username": {i.username},
		"role":     {i.role.GetName()},
	}
}

// Identity is an extension of the identity interface for user authenticating via Basic authentication.
type Identity interface {
	authn.Identity

	// AsExternalUser returns the claims
	AsExternalUser() *tokens.ExternalUserClaim

	// isBasicAuthIdentity is a marker method.
	isBasicAuthIdentity()
}
