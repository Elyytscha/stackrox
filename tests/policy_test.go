package tests

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	v1 "github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/testutils"
	"github.com/stackrox/rox/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
)

const (
	knownPolicyID              = "d3e480c1-c6de-4cd2-9006-9a3eb3ad36b6"
	notAnID                    = "Joseph Rules"
	duplicateName              = "duplicate_name"
	duplicateID                = "duplicate_id"
	removedClustersOrNotifiers = "removed_clusters_or_notifiers"
)

var (
	addedPolicies []string
)

func TestPolicies(t *testing.T) {
	assumeFeatureFlagHasValue(t, features.PolicyImportExport, true)
	defer tearDownTest(t)
	verifyExportNonExistentFails(t)
	verifyExportExistentSucceeds(t)
	verifyMixedExportFails(t)
	verifyImportSucceeds(t)
	verifyDuplicateImportSucceeds(t)
	verifyImportInvalidFails(t)
	verifyImportDuplicateNameFails(t)
	verifyImportDuplicateIDFails(t)
	verifyImportDuplicateNameAndIDFails(t)
	verifyImportNoIDSucceeds(t)
	verifyImportMultipleSucceeds(t)
	verifyImportMixedSuccess(t)
	verifyNotifiersRemoved(t)
	verifyWhitelistsRemoved(t)
	verifyScopesRemoved(t)
	verifyOverwriteNameSucceeds(t)
	verifyOverwriteIDSucceeds(t)
	verifyOverwriteNameAndIDSucceeds(t)
}

func tearDownTest(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	var cleanupErrors []error
	for _, id := range addedPolicies {
		if id == knownPolicyID {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		_, err := service.DeletePolicy(ctx, &v1.ResourceByID{
			Id: id,
		})
		cancel()
		if err != nil {
			cleanupErrors = append(cleanupErrors, err)
		}
	}
	for _, cleanupError := range cleanupErrors {
		if strings.Contains(cleanupError.Error(), "not found") {
			continue
		}
		// If there was any cleanup error other than "not found", log it here.
		assert.Nil(t, cleanupError)
	}
}

func exportPolicy(t *testing.T, service v1.PolicyServiceClient, id string) *storage.Policy {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := service.ExportPolicies(ctx, &v1.ExportPoliciesRequest{
		PolicyIds: []string{id},
	})
	cancel()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetPolicies(), 1)
	require.Equal(t, id, resp.GetPolicies()[0].GetId())

	return resp.Policies[0]
}

func validateExportFails(t *testing.T, service v1.PolicyServiceClient, id string, expectedErrors []*v1.ExportPolicyError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := service.ExportPolicies(ctx, &v1.ExportPoliciesRequest{
		PolicyIds: []string{notAnID},
	})
	cancel()
	require.Nil(t, resp)
	require.Error(t, err)
	compareErrorsToExpected(t, expectedErrors, err)
}

func validateImport(t *testing.T, importResp *v1.ImportPoliciesResponse, policies []*storage.Policy, errIndices []int, expectedErrs [][]string, validateIDs bool) {
	require.NotNil(t, importResp)
	allSucceeded := len(errIndices) == 0
	require.Equal(t, allSucceeded, importResp.GetAllSucceeded())
	require.NotNil(t, importResp.GetResponses())
	require.Len(t, importResp.GetResponses(), len(policies))
	for i, importPolicyResponse := range importResp.GetResponses() {
		if len(errIndices) == 0 || errIndices[0] != i {
			validateSuccess(t, importPolicyResponse, policies[i], validateIDs)
			continue
		}
		errIndices = errIndices[1:]
		expectedErr := expectedErrs[0]
		expectedErrs = expectedErrs[1:]
		validateFailure(t, importPolicyResponse, policies[i], expectedErr, validateIDs)
	}
}

func validateSuccess(t *testing.T, importPolicyResponse *v1.ImportPolicyResponse, expectedPolicy *storage.Policy, ignoreID bool) {
	require.True(t, importPolicyResponse.GetSucceeded())
	addedPolicies = append(addedPolicies, importPolicyResponse.GetPolicy().GetId())
	if ignoreID {
		expectedPolicy.Id = ""
		importPolicyResponse.GetPolicy().Id = ""
	}
	require.Equal(t, expectedPolicy, importPolicyResponse.GetPolicy())
	require.Empty(t, importPolicyResponse.GetErrors())
}

func validateFailure(t *testing.T, importPolicyResponse *v1.ImportPolicyResponse, policy *storage.Policy, expectedErrTypes []string, validateErr bool) {
	require.False(t, importPolicyResponse.GetSucceeded())
	if !validateErr {
		policy.Id = ""
		importPolicyResponse.GetPolicy().Id = ""
	}
	require.Equal(t, policy, importPolicyResponse.GetPolicy())
	require.Len(t, importPolicyResponse.GetErrors(), len(expectedErrTypes))
	for i, policyErr := range importPolicyResponse.GetErrors() {
		require.Equal(t, policyErr.GetType(), expectedErrTypes[i])
	}
}

func validateImportPoliciesErrors(t *testing.T, importResp *v1.ImportPoliciesResponse, policy *storage.Policy, expectedErrTypes []string) {
	require.NotNil(t, importResp)
	require.False(t, importResp.GetAllSucceeded())
	require.NotNil(t, importResp.GetResponses())
	require.Len(t, importResp.GetResponses(), 1)
	validateFailure(t, importResp.GetResponses()[0], policy, expectedErrTypes, true)
}

func validateImportPoliciesSuccess(t *testing.T, importResp *v1.ImportPoliciesResponse, policies []*storage.Policy, ignoreID bool) {
	require.NotNil(t, importResp)
	require.True(t, importResp.GetAllSucceeded())
	require.NotNil(t, importResp.GetResponses())
	require.Len(t, importResp.GetResponses(), len(policies))
	for i, importPolicyResponse := range importResp.GetResponses() {
		validateSuccess(t, importPolicyResponse, policies[i], ignoreID)
	}
}

func compareErrorsToExpected(t *testing.T, expectedErrors []*v1.ExportPolicyError, apiError error) {
	apiStatus, ok := status.FromError(apiError)
	require.True(t, ok)
	details := apiStatus.Details()
	require.Len(t, details, 1)
	exportErrors, ok := details[0].(*v1.ExportPoliciesErrorList)
	require.True(t, ok)
	// actual errors == expected errors ignoring order
	require.Len(t, exportErrors.GetErrors(), len(expectedErrors))
	for _, expected := range expectedErrors {
		require.Contains(t, exportErrors.GetErrors(), expected)
	}
}

func makeError(errorID, errorString string) *v1.ExportPolicyError {
	return &v1.ExportPolicyError{
		PolicyId: errorID,
		Error: &v1.PolicyError{
			Error: errorString,
		},
	}
}

func createUniquePolicy(t *testing.T, service v1.PolicyServiceClient) *storage.Policy {
	newUniquePolicy := exportPolicy(t, service, knownPolicyID)
	newUniquePolicy.Name = uuid.NewV4().String()
	newUniquePolicy.Id = uuid.NewV4().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{newUniquePolicy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{newUniquePolicy}, false)
	return newUniquePolicy
}

func validateWhitelistOrScopeOrNotiferRemoved(t *testing.T, importResp *v1.ImportPoliciesResponse, expectedPolicy *storage.Policy) {
	require.NotNil(t, importResp)
	require.True(t, importResp.GetAllSucceeded())
	require.NotNil(t, importResp.GetResponses())
	require.Len(t, importResp.GetResponses(), 1)

	importPolicyResponse := importResp.GetResponses()[0]
	require.True(t, importPolicyResponse.GetSucceeded())
	addedPolicies = append(addedPolicies, importPolicyResponse.GetPolicy().GetId())
	require.Equal(t, expectedPolicy, importPolicyResponse.GetPolicy())
	require.Len(t, importPolicyResponse.GetErrors(), 1)

	policyErrors := importResp.GetResponses()[0].Errors
	require.Len(t, policyErrors, 1)
	policyError := policyErrors[0]
	require.Equal(t, removedClustersOrNotifiers, policyError.GetType())
}

func verifyExportNonExistentFails(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	mockErrors := []*v1.ExportPolicyError{
		makeError(notAnID, "not found"),
	}
	validateExportFails(t, service, notAnID, mockErrors)
}

func verifyExportExistentSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)

	service := v1.NewPolicyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := service.ExportPolicies(ctx, &v1.ExportPoliciesRequest{
		PolicyIds: []string{knownPolicyID},
	})
	cancel()
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Len(t, resp.GetPolicies(), 1)
	require.Equal(t, knownPolicyID, resp.GetPolicies()[0].GetId())
}

func verifyMixedExportFails(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)

	mockErrors := []*v1.ExportPolicyError{
		makeError(notAnID, "not found"),
	}
	service := v1.NewPolicyServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	resp, err := service.ExportPolicies(ctx, &v1.ExportPoliciesRequest{
		PolicyIds: []string{knownPolicyID, notAnID},
	})
	cancel()
	require.Nil(t, resp)
	require.Error(t, err)
	compareErrorsToExpected(t, mockErrors, err)
}

func verifyImportSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	policy := exportPolicy(t, service, knownPolicyID)
	policy.Name = "A new name"
	policy.Id = "integrationtestpolicy"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{policy}, false)
}

func verifyDuplicateImportSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	policy := exportPolicy(t, service, knownPolicyID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{policy}, false)
}

func verifyImportInvalidFails(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	badPolicy := &storage.Policy{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{badPolicy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesErrors(t, importResp, badPolicy, []string{"invalid_policy"})
}

func verifyImportDuplicateNameFails(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	policy := exportPolicy(t, service, knownPolicyID)

	policy.Id = "duplicateNamePolicy"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesErrors(t, importResp, policy, []string{duplicateName})
}

func verifyImportDuplicateIDFails(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	policy := exportPolicy(t, service, knownPolicyID)

	policy.Name = "New name"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesErrors(t, importResp, policy, []string{duplicateID})
}

func verifyImportDuplicateNameAndIDFails(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	policy := exportPolicy(t, service, knownPolicyID)

	policy.Description = "A different description so the policies are not equal"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesErrors(t, importResp, policy, []string{duplicateID, duplicateName})
}

func verifyImportNoIDSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	policy := exportPolicy(t, service, knownPolicyID)
	policy.Name = "Some unique name"
	policy.Id = ""
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{policy}, true)
}

func verifyImportMultipleSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	validPolicy := exportPolicy(t, service, knownPolicyID)

	policy1 := proto.Clone(validPolicy).(*storage.Policy)
	policy1.Id = "new policy ID"
	policy1.Name = "This is a valid policy"
	policy2 := proto.Clone(validPolicy).(*storage.Policy)
	policy2.Id = "another new policy ID"
	policy2.Name = "This is another valid policy"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy1, policy2},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{policy1, policy2}, false)
}

func verifyImportMixedSuccess(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	validPolicy := exportPolicy(t, service, knownPolicyID)

	// Policy 1 should be valid
	policy1 := proto.Clone(validPolicy).(*storage.Policy)
	policy1.Id = "Probably I should make these UUIDs"
	policy1.Name = "This is a valid and totally unique policy"
	// Policy 2 should have a duplicate name error
	policy2 := proto.Clone(validPolicy).(*storage.Policy)
	policy2.Id = "another new entirely different policy ID"

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy1, policy2},
	})
	cancel()
	require.NoError(t, err)
	validateImport(t, importResp, []*storage.Policy{policy1, policy2}, []int{1}, [][]string{{duplicateName}}, true)
}

func verifyNotifiersRemoved(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	validPolicy := exportPolicy(t, service, knownPolicyID)

	// Policy 1 should be valid
	policy := proto.Clone(validPolicy).(*storage.Policy)
	policy.Id = "verifyNotifiersRemoved policy ID"
	policy.Name = "verifyNotifiersRemoved is a valid policy"
	policy.Notifiers = []string{"This is not a notifier"}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)

	// Notifier should have been scraped out
	policy.Notifiers = nil
	validateWhitelistOrScopeOrNotiferRemoved(t, importResp, policy)
}

func verifyWhitelistsRemoved(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	validPolicy := exportPolicy(t, service, knownPolicyID)

	// Policy 1 should be valid
	policy := proto.Clone(validPolicy).(*storage.Policy)
	policy.Id = "verifyWhitelistsRemoved policy ID"
	policy.Name = "verifyWhitelistsRemoved is a valid policy"
	policy.Whitelists = []*storage.Whitelist{
		{
			Deployment: &storage.Whitelist_Deployment{
				Scope: &storage.Scope{
					Cluster: "This is not a cluster",
				},
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)

	// Whitelist should have been scraped out
	policy.Whitelists = nil
	validateWhitelistOrScopeOrNotiferRemoved(t, importResp, policy)
}

func verifyScopesRemoved(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	validPolicy := exportPolicy(t, service, knownPolicyID)

	// Policy 1 should be valid
	policy := proto.Clone(validPolicy).(*storage.Policy)
	policy.Id = "verifyScopesRemoved policy ID"
	policy.Name = "verifyScopesRemoved is a valid policy"
	policy.Scope = []*storage.Scope{
		{
			Cluster: "This is not a cluster",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{policy},
	})
	cancel()
	require.NoError(t, err)

	// Scope should have been scraped out
	policy.Scope = nil
	validateWhitelistOrScopeOrNotiferRemoved(t, importResp, policy)
}

func verifyOverwriteNameSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	// Create an existing policy so we don't change default policies
	existingPolicy := createUniquePolicy(t, service)

	newPolicy := proto.Clone(existingPolicy).(*storage.Policy)
	newPolicy.Id = uuid.NewV4().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{newPolicy},
		Metadata: &v1.ImportPoliciesMetadata{
			Overwrite: true,
		},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{newPolicy}, false)

	mockErrors := []*v1.ExportPolicyError{
		makeError(notAnID, "not found"),
	}
	validateExportFails(t, service, existingPolicy.GetId(), mockErrors)

	dbPolicy := exportPolicy(t, service, newPolicy.GetId())
	require.Equal(t, newPolicy, dbPolicy)
}

func verifyOverwriteIDSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	// Create an existing policy so we don't change default policies
	existingPolicy := createUniquePolicy(t, service)

	newPolicy := proto.Clone(existingPolicy).(*storage.Policy)
	newPolicy.Name = uuid.NewV4().String()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{newPolicy},
		Metadata: &v1.ImportPoliciesMetadata{
			Overwrite: true,
		},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{newPolicy}, false)

	dbPolicy := exportPolicy(t, service, existingPolicy.GetId())
	require.Equal(t, newPolicy, dbPolicy)
}

func verifyOverwriteNameAndIDSucceeds(t *testing.T) {
	conn := testutils.GRPCConnectionToCentral(t)
	service := v1.NewPolicyServiceClient(conn)

	// Create an existing policy so we don't change default policies
	existingPolicyDuplicateName := createUniquePolicy(t, service)
	existingPolicyDuplicateID := createUniquePolicy(t, service)

	newPolicy := proto.Clone(existingPolicyDuplicateID).(*storage.Policy)
	newPolicy.Name = existingPolicyDuplicateName.GetName()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	importResp, err := service.ImportPolicies(ctx, &v1.ImportPoliciesRequest{
		Policies: []*storage.Policy{newPolicy},
		Metadata: &v1.ImportPoliciesMetadata{
			Overwrite: true,
		},
	})
	cancel()
	require.NoError(t, err)
	validateImportPoliciesSuccess(t, importResp, []*storage.Policy{newPolicy}, false)

	mockErrors := []*v1.ExportPolicyError{
		makeError(notAnID, "not found"),
	}
	validateExportFails(t, service, existingPolicyDuplicateName.GetId(), mockErrors)

	dbPolicy := exportPolicy(t, service, existingPolicyDuplicateID.GetId())
	require.Equal(t, newPolicy, dbPolicy)
}
