package extensions

import (
	"testing"

	centralv1Alpha1 "github.com/stackrox/rox/operator/api/central/v1alpha1"
	securedClusterv1Alpha1 "github.com/stackrox/rox/operator/api/securedcluster/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestUpdateProductVersion_Central(t *testing.T) {
	var status centralv1Alpha1.CentralStatus
	var uSt unstructured.Unstructured

	var err error
	uSt.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&status)
	require.NoError(t, err)

	// Update from empty to 1.2.3
	assert.True(t, updateProductVersion(&uSt, "1.2.3"))
	status = centralv1Alpha1.CentralStatus{}
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(uSt.Object, &status))

	assert.Equal(t, "1.2.3", status.ProductVersion)

	// Update a second time, same value
	uSt = unstructured.Unstructured{}
	uSt.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&status)
	require.NoError(t, err)

	assert.False(t, updateProductVersion(&uSt, "1.2.3"))
	status = centralv1Alpha1.CentralStatus{}
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(uSt.Object, &status))

	assert.Equal(t, "1.2.3", status.ProductVersion)

	// Update a third time, new value
	uSt = unstructured.Unstructured{}
	uSt.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&status)
	require.NoError(t, err)

	assert.True(t, updateProductVersion(&uSt, "4.5.6"))
	status = centralv1Alpha1.CentralStatus{}
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(uSt.Object, &status))

	assert.Equal(t, "4.5.6", status.ProductVersion)
}

func TestUpdateProductVersion_SecuredCluster(t *testing.T) {
	var status securedClusterv1Alpha1.SecuredClusterStatus
	var uSt unstructured.Unstructured

	var err error
	uSt.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&status)
	require.NoError(t, err)

	// Update from empty to 1.2.3
	assert.True(t, updateProductVersion(&uSt, "1.2.3"))
	status = securedClusterv1Alpha1.SecuredClusterStatus{}
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(uSt.Object, &status))

	assert.Equal(t, "1.2.3", status.ProductVersion)

	// Update a second time, same value
	uSt = unstructured.Unstructured{}
	uSt.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&status)
	require.NoError(t, err)

	assert.False(t, updateProductVersion(&uSt, "1.2.3"))
	status = securedClusterv1Alpha1.SecuredClusterStatus{}
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(uSt.Object, &status))

	assert.Equal(t, "1.2.3", status.ProductVersion)

	// Update a third time, new value
	uSt = unstructured.Unstructured{}
	uSt.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(&status)
	require.NoError(t, err)

	assert.True(t, updateProductVersion(&uSt, "4.5.6"))
	status = securedClusterv1Alpha1.SecuredClusterStatus{}
	require.NoError(t, runtime.DefaultUnstructuredConverter.FromUnstructured(uSt.Object, &status))

	assert.Equal(t, "4.5.6", status.ProductVersion)
}
