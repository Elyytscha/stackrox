// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package protoutils

import "github.com/stackrox/rox/generated/storage"

// *storage.Policy represents a generic proto type that we clone.

// CloneStoragePolicy is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneStoragePolicy(val *storage.Policy) *storage.Policy {
	return protoCloneWrapper(val).(*storage.Policy)
}

// *storage.Deployment represents a generic proto type that we clone.

// CloneStorageDeployment is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneStorageDeployment(val *storage.Deployment) *storage.Deployment {
	return protoCloneWrapper(val).(*storage.Deployment)
}

// *storage.Alert represents a generic proto type that we clone.

// CloneStorageAlert is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneStorageAlert(val *storage.Alert) *storage.Alert {
	return protoCloneWrapper(val).(*storage.Alert)
}
