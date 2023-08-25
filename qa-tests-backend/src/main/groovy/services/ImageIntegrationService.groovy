package services

import common.Constants
import groovy.util.logging.Slf4j
import io.grpc.Status
import io.grpc.StatusRuntimeException
import io.stackrox.proto.api.v1.ImageIntegrationServiceGrpc
import io.stackrox.proto.api.v1.ImageIntegrationServiceOuterClass
import io.stackrox.proto.storage.ImageIntegrationOuterClass
import io.stackrox.proto.storage.ImageIntegrationOuterClass.ImageIntegrationCategory
import objects.StackroxScannerIntegration
import util.Timer

@Slf4j
class ImageIntegrationService extends BaseService {
    static final private String AUTO_REGISTERED_GCR_INTEGRATION = "Autogenerated https://us.gcr.io for cluster remote"

    static getImageIntegrationClient() {
        return ImageIntegrationServiceGrpc.newBlockingStub(getChannel())
    }

    static testImageIntegration(ImageIntegrationOuterClass.ImageIntegration integration) {
        try {
            getImageIntegrationClient().testImageIntegration(integration)
            return true
        } catch (Exception e) {
            log.warn("could not get integration", e)
            return false
        }
    }

    static createImageIntegration(ImageIntegrationOuterClass.ImageIntegration integration,
                                  skipTestSideTest = false) {
        if (!skipTestSideTest) {
            Boolean tested = false
            Timer t = new Timer(15, 3)
            while (t.IsValid()) {
                try {
                    getImageIntegrationClient().testImageIntegration(integration)
                    log.debug "Integration tested: ${integration.name}"
                    tested = true
                    break
                } catch (Exception e) {
                    log.warn("Integration test failed: ${integration.name}", e)
                }
            }

            if (!tested) {
                log.error("Integration test failed")
                return ""
            }
        }

        ImageIntegrationOuterClass.ImageIntegration createdIntegration
        Timer t = new Timer(15, 3)
        while (t.IsValid()) {
            try {
                createdIntegration =
                        getImageIntegrationClient().postImageIntegration(integration)
                log.debug "Integration created: ${createdIntegration.name}: ${createdIntegration.id}"
                break
            } catch (Exception e) {
                log.warn( "Unable to create image integration ${integration.name}", e)
            }
        }

        if (!createdIntegration || !createdIntegration.id) {
            log.error( "Unable to create image integration")
            return ""
        }

        ImageIntegrationOuterClass.ImageIntegration foundIntegration
        t = new Timer(15, 3)
        while (t.IsValid()) {
            try {
                foundIntegration =
                        getImageIntegrationClient().getImageIntegration(getResourceByID(createdIntegration.id))
                if (foundIntegration) {
                    log.debug "Integration found after creation: ${foundIntegration.name}: ${foundIntegration.id}"
                    return foundIntegration.id
                }
            } catch (Exception e) {
                log.warn("Unable to find the created image integration ${integration.name}", e)
            }
        }

        log.error "Unable to find the created image integration"
        return ""
    }

    static deleteImageIntegration(String integrationId) {
        try {
            getImageIntegrationClient().deleteImageIntegration(getResourceByID(integrationId))
        } catch (Exception e) {
            log.warn("Failed to delete integration", e)
            return false
        }
        try {
            ImageIntegrationOuterClass.ImageIntegration integration =
                    getImageIntegrationClient().getImageIntegration(getResourceByID(integrationId))
            while (integration) {
                sleep 2000
                integration = getImageIntegrationClient().getImageIntegration(getResourceByID(integrationId))
            }
        } catch (StatusRuntimeException e) {
            if (e.status.code == Status.Code.NOT_FOUND) {
                log.debug "Image integration deleted: ${integrationId}"
                return true
            }
            log.debug("delete image integration", e)
        }
    }

    static getImageIntegrations() {
        return getImageIntegrationClient().getImageIntegrations(
                ImageIntegrationServiceOuterClass.GetImageIntegrationsRequest.newBuilder().build()
        ).integrationsList
    }

    static getImageIntegrationByName(String name) {
        List<ImageIntegrationOuterClass.ImageIntegration> integrations = getImageIntegrations()
        def integrationId = integrations.find { it.name == name }?.id
        return integrationId ?
                getImageIntegrationClient().getImageIntegration(getResourceByID(integrationId)) :
                null
    }

    /*
        Helper functions to simplify creating known integrations
    */

    static boolean deleteStackRoxScannerIntegrationIfExists() {
        try {
            // The Stackrox Scanner integration is auto-added by the product,
            // so we first check whether it already exists.
            def scannerIntegrations = getImageIntegrationClient().getImageIntegrations(
                    ImageIntegrationServiceOuterClass.GetImageIntegrationsRequest.
                            newBuilder().
                            setName(Constants.AUTO_REGISTERED_STACKROX_SCANNER_INTEGRATION).
                            build()
            )
            if (scannerIntegrations.integrationsCount > 1) {
                throw new RuntimeException("UNEXPECTED: Got more than one scanner integration: ${scannerIntegrations}")
            }
            if (scannerIntegrations.integrationsCount == 0) {
                log.debug "No Stackrox scanner integrations were found"
                return false
            }
            def id = scannerIntegrations.getIntegrations(0).id
            // Delete
            getImageIntegrationClient().deleteImageIntegration(getResourceByID(id))
            return true
        } catch (Exception e) {
            log.warn("Unable to delete existing Stackrox scanner integration", e)
            // return false since we are not sure if the integration exists and failed to delete, or did not exist
            return false
        }
    }

    static String addStackroxScannerIntegration() {
        ImageIntegrationOuterClass.ImageIntegration existing =
            getImageIntegrationByName(Constants.AUTO_REGISTERED_STACKROX_SCANNER_INTEGRATION)
        if (existing) {
            log.debug("${Constants.AUTO_REGISTERED_STACKROX_SCANNER_INTEGRATION} already exists")
            return existing.id
        }
        return StackroxScannerIntegration.createDefaultIntegration()
    }

    static boolean deleteAutoRegisteredGCRIntegrationIfExists() {
        ImageIntegrationOuterClass.ImageIntegration integration = getImageIntegrations().find
                { it.name == AUTO_REGISTERED_GCR_INTEGRATION }
        if (integration) {
            deleteImageIntegration(integration.getId())
        }
        else {
            log.debug "There is no auto-registered GCR integration to delete"
        }
    }

    static getIntegrationCategories(boolean includeScanner) {
        return includeScanner ?
                [ImageIntegrationCategory.REGISTRY, ImageIntegrationCategory.SCANNER] :
                [ImageIntegrationCategory.REGISTRY]
    }
}
