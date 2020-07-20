package services

import io.stackrox.proto.api.v1.Common
import io.stackrox.proto.api.v1.PolicyServiceGrpc
import io.stackrox.proto.storage.PolicyOuterClass.Policy
import io.stackrox.proto.api.v1.PolicyServiceOuterClass.PatchPolicyRequest
import io.stackrox.proto.api.v1.PolicyServiceOuterClass.PostPolicyRequest

class CreatePolicyService extends BaseService {

    static getPolicyClient() {
        return PolicyServiceGrpc.newBlockingStub(getChannel())
    }

    static String createNewPolicy(Policy policy) {
        String policyID = ""

        try {
            policyID = getPolicyClient().postPolicy(
                PostPolicyRequest.newBuilder().setPolicy(policy).setEnableStrictValidation(true).build()).getId()
        } catch (Exception e) {
            println e.toString()
        }

        return policyID
    }

    static deletePolicy(String policyID) {
        try {
            getPolicyClient().deletePolicy(
                    Common.ResourceByID.newBuilder()
                            .setId(policyID)
                            .build()
            )
        } catch (Exception e) {
            println e.toString()
        }
    }

    static patchPolicy(PatchPolicyRequest pr) {
        try {
            getPolicyClient().patchPolicy(pr).newBuilder().build()
        }
        catch (Exception e) {
            println e.toString()
        }
    }
}
