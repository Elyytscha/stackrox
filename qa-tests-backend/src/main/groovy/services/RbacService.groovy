package services

import io.stackrox.proto.api.v1.Common
import io.stackrox.proto.api.v1.RbacServiceGrpc
import io.stackrox.proto.api.v1.SearchServiceOuterClass
import objects.K8sRole
import objects.K8sRoleBinding
import util.Timer

class RbacService extends BaseService {
    static getRbacService() {
        return RbacServiceGrpc.newBlockingStub(getChannel())
    }

    static getRoles(SearchServiceOuterClass.RawQuery query = SearchServiceOuterClass.RawQuery.newBuilder().build()) {
        return getRbacService().listRoles(query).rolesList
    }

    static getRole(String id) {
        try {
            return getRbacService().getRole(
                    Common.ResourceByID.newBuilder().setId(id).build()
            ).role
        } catch (Exception e) {
            println "Error fetching role: ${e.toString()}"
        }
    }

    static boolean waitForRole(K8sRole role) {
        Timer t = new Timer(30, 3)
        while (t.IsValid()) {
            println "Waiting for Role"
            def roles = getRoles()
            def r = roles.find {
                it.name == role.name &&
                        it.namespace == role.namespace
            }

            if (r) {
                return true
            }
        }
        println "Time out for Waiting for Role"
        return false
    }

    static boolean waitForRoleRemoved(K8sRole role) {
        Timer t = new Timer(30, 3)
        while (t.IsValid()) {
            println "Waiting for Role removed"
            def roles = getRoles()
            def r = roles.find {
                it.name == role.name &&
                        it.namespace == role.namespace
            }
            if (!r) {
                return true
            }
        }
        println "Time out for Waiting for Role removal"
        return false
    }

    static getRoleBindings(
            SearchServiceOuterClass.RawQuery query = SearchServiceOuterClass.RawQuery.newBuilder().build()) {
        return getRbacService().listRoleBindings(query).bindingsList
    }

    static getRoleBinding(String id) {
        try {
            return getRbacService().getRoleBinding(
                    Common.ResourceByID.newBuilder().setId(id).build()
            ).binding
        } catch (Exception e) {
            println "Error fetching role binding: ${e.toString()}"
        }
    }

    static boolean waitForRoleBinding(K8sRoleBinding roleBinding) {
        Timer t = new Timer(30, 3)
        while (t.IsValid()) {
            println "Waiting for Role Binding"
            def roleBindings = getRoleBindings()
            def r = roleBindings.find {
                it.name == roleBinding.name &&
                        it.namespace == roleBinding.namespace
            }

            if (r) {
                return true
            }
        }
        println "Time out for Waiting for Role Binding"
        return false
    }

    static boolean waitForRoleBindingRemoved(K8sRoleBinding roleBinding) {
        Timer t = new Timer(30, 3)
        while (t.IsValid()) {
            println "Waiting for Role Binding removed"
            def roleBindings = getRoleBindings()
            def r = roleBindings.find {
                it.name == roleBinding.name &&
                        it.namespace == roleBinding.namespace
            }
            if (!r) {
                return true
            }
        }
        println "Time out for Waiting for Role Binding removal"
        return false
    }

    static getSubjects(SearchServiceOuterClass.RawQuery query = SearchServiceOuterClass.RawQuery.newBuilder().build()) {
        return getRbacService().listSubjects(query).subjectAndRolesList
    }

    static getSubject(String id) {
        try {
            return getRbacService().getSubject(
                    Common.ResourceByID.newBuilder().setId(id).build()
            ).subject
        } catch (Exception e) {
            println "Error fetching subject: ${e.toString()}"
        }
    }
}
