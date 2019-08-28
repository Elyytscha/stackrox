// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package protoutils

import (
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
)

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

// *v1.Query represents a generic proto type that we clone.

// CloneV1Query is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneV1Query(val *v1.Query) *v1.Query {
	return protoCloneWrapper(val).(*v1.Query)
}

// *storage.Cluster represents a generic proto type that we clone.

// CloneStorageCluster is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneStorageCluster(val *storage.Cluster) *storage.Cluster {
	return protoCloneWrapper(val).(*storage.Cluster)
}

// *storage.ImageIntegration represents a generic proto type that we clone.

// CloneStorageImageIntegration is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneStorageImageIntegration(val *storage.ImageIntegration) *storage.ImageIntegration {
	return protoCloneWrapper(val).(*storage.ImageIntegration)
}

// *storage.ClusterStatus represents a generic proto type that we clone.

// CloneStorageClusterStatus is a (generic) wrapper around proto.Clone that is strongly typed.
func CloneStorageClusterStatus(val *storage.ClusterStatus) *storage.ClusterStatus {
	return protoCloneWrapper(val).(*storage.ClusterStatus)
}
