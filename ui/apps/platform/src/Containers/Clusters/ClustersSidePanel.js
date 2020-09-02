import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { selectors } from 'reducers';
import { format } from 'date-fns';
import { createStructuredSelector } from 'reselect';
import * as Icon from 'react-feather';
import cloneDeep from 'lodash/cloneDeep';
import get from 'lodash/get';
import set from 'lodash/set';

import Message from 'Components/Message';
import Panel from 'Components/Panel';
import PanelButton from 'Components/PanelButton';
import useInterval from 'hooks/useInterval';
import {
    getClusterById,
    saveCluster,
    downloadClusterYaml,
    fetchKernelSupportAvailable,
    rotateClusterCerts,
} from 'services/ClustersService';

import ClusterEditForm from './ClusterEditForm';
import ClusterDeployment from './ClusterDeployment';
import {
    clusterDetailPollingInterval,
    findUpgradeState,
    formatUpgradeMessage,
    getCredentialExpirationProps,
    getUpgradeStatusDetail,
    initiationOfCertRotationIfApplicable,
    isUpToDateStateObject,
    newClusterDefault,
    wizardSteps,
    centralEnvDefault,
} from './cluster.helpers';
import CollapsibleCard from '../../Components/CollapsibleCard';
import Button from '../../Components/Button';
import { generateSecuredClusterCertSecret } from '../../services/CertGenerationService';
import Dialog from '../../Components/Dialog';

function fetchCentralEnv() {
    return fetchKernelSupportAvailable().then((kernelSupportAvailable) => {
        return {
            kernelSupportAvailable,
            successfullyFetched: true,
        };
    });
}

const requiredKeys = ['name', 'type', 'mainImage', 'centralApiEndpoint'];

const validate = (values) => {
    const errors = {};

    requiredKeys.forEach((key) => {
        if (values[key].length === 0) {
            errors[key] = 'This field is required';
        }
    });

    return errors;
};

function ClustersSidePanel({ metadata, selectedClusterId, setSelectedClusterId }) {
    const defaultCluster = cloneDeep(newClusterDefault);
    const envAwareClusterDefault = {
        ...defaultCluster,
        mainImage: metadata.releaseBuild ? 'stackrox.io/main' : 'stackrox/main',
        collectorImage: metadata.releaseBuild
            ? 'collector.stackrox.io/collector'
            : 'stackrox/collector',
    };

    const [selectedCluster, setSelectedCluster] = useState(envAwareClusterDefault);
    const [centralEnv, setCentralEnv] = useState(centralEnvDefault);
    const [wizardStep, setWizardStep] = useState(wizardSteps.FORM);
    const [loadingCounter, setLoadingCounter] = useState(0);
    const [messageState, setMessageState] = useState(null);
    const [pollingCount, setPollingCount] = useState(0);
    const [pollingDelay, setPollingDelay] = useState(null);
    const [submissionError, setSubmissionError] = useState('');
    const [freshCredentialsDownloaded, setFreshCredentialsDownloaded] = useState(false);
    const [showCertRotationModal, setShowCertRotationModal] = useState(false);

    const [createUpgraderSA, setCreateUpgraderSA] = useState(true);

    function unselectCluster() {
        setSubmissionError('');
        setSelectedClusterId('');
        setSelectedCluster(envAwareClusterDefault);
        setMessageState(null);
        setWizardStep(wizardSteps.FORM);
        setPollingDelay(null);
        setFreshCredentialsDownloaded(false);
        setShowCertRotationModal(false);
    }

    useEffect(
        () => {
            const clusterIdToRetrieve = selectedClusterId;

            setLoadingCounter((prev) => prev + 1);
            fetchCentralEnv()
                .then((freshCentralEnv) => {
                    setCentralEnv(freshCentralEnv);
                    if (clusterIdToRetrieve === 'new') {
                        const updatedCluster = {
                            ...selectedCluster,
                            slimCollector: freshCentralEnv.kernelSupportAvailable,
                        };
                        setSelectedCluster(updatedCluster);
                    }
                })
                .finally(() => {
                    setLoadingCounter((prev) => prev - 1);
                });

            if (clusterIdToRetrieve && clusterIdToRetrieve !== 'new') {
                setLoadingCounter((prev) => prev + 1);
                setMessageState(null);
                // don't want to cache or memoize, because we always want the latest real-time data
                getClusterById(clusterIdToRetrieve)
                    .then((cluster) => {
                        // TODO: refactor to use useReducer effect
                        setSelectedCluster(cluster);

                        // stop polling after contact is established
                        if (
                            selectedCluster &&
                            selectedCluster.healthStatus &&
                            selectedCluster.healthStatus.lastContact
                        ) {
                            setPollingDelay(null);
                        }
                    })
                    .catch(() => {
                        setMessageState({
                            blocking: true,
                            type: 'error',
                            message: 'There was an error downloading the configuration files.',
                        });
                    })
                    .finally(() => {
                        setLoadingCounter((prev) => prev - 1);
                    });
                // TODO: When rolling out this feature the user should be informed somehow
                // in case this property could not be retrieved.
                // The default slimCollectorMode (false) is sufficient for now.
            }
        },
        // lint rule "exhaustive-deps" wants to add selectedCluster to change-detection
        // but we don't want to fetch while we're editing, so disabled that rule here
        // eslint-disable-next-line react-hooks/exhaustive-deps
        [selectedClusterId, pollingCount]
    );

    // use a custom hook to set up polling, thanks Dan Abramov and Rob Stark
    useInterval(() => {
        setPollingCount(pollingCount + 1);
    }, pollingDelay);

    /**
     * naive implementation of form handler
     *  - replace with more robust system, probably react-final-form
     *
     * @param   {Event}  event  native JS Event object from an onChange event in an input
     *
     * @return  {nothing}       Side effect: change the corresponding property in selectedCluster
     */
    function onChange(event) {
        if (get(selectedCluster, event.target.name) !== undefined) {
            const newClusterSettings = { ...selectedCluster };
            const newValue =
                event.target.type === 'checkbox' ? event.target.checked : event.target.value;
            set(newClusterSettings, event.target.name, newValue);
            setSelectedCluster(newClusterSettings);
        }
    }

    function onNext() {
        if (wizardStep === wizardSteps.FORM) {
            setSubmissionError('');
            saveCluster(selectedCluster)
                .then((response) => {
                    const newId = response.response.result.cluster; // really is nested like this
                    const clusterWithId = { ...selectedCluster, id: newId };
                    setSelectedCluster(clusterWithId);

                    setWizardStep(wizardSteps.DEPLOYMENT);

                    if (
                        !(
                            selectedCluster &&
                            selectedCluster.healthStatus &&
                            selectedCluster.healthStatus.lastContact
                        )
                    ) {
                        setPollingDelay(clusterDetailPollingInterval);
                    }
                })
                .catch((error) => {
                    const serverError = get(
                        error,
                        'response.data.message',
                        'An unknown error has occurred.'
                    );

                    setSubmissionError(serverError);
                });
        } else {
            unselectCluster();
        }
    }

    function toggleSA() {
        setCreateUpgraderSA(!createUpgraderSA);
    }

    function onDownload() {
        setSubmissionError('');
        downloadClusterYaml(selectedCluster.id, createUpgraderSA).catch((error) => {
            const serverError = get(
                error,
                'response.data.message',
                'We could not download the configuration files.'
            );

            setSubmissionError(serverError);
        });
    }

    function openCertRotationModal() {
        setShowCertRotationModal(true);
    }

    function hideCertRotationModal() {
        setShowCertRotationModal(false);
    }

    function triggerClusterCertRotation() {
        rotateClusterCerts(selectedClusterId)
            .then(() => {
                hideCertRotationModal();
            })
            .catch((error) => {
                const serverError = get(
                    error,
                    'response.data.message',
                    'Failed to apply new credentials to the cluster.'
                );

                setSubmissionError(serverError);
            });
    }

    function generateCertSecret() {
        generateSecuredClusterCertSecret(selectedClusterId).catch((error) => {
            const serverError = get(
                error,
                'response.data.message',
                'Failed to regenerate certificates.'
            );

            setSubmissionError(serverError);
            setFreshCredentialsDownloaded(false);
        });
        setFreshCredentialsDownloaded(true);
    }

    /**
     * rendering section
     */
    if (!selectedClusterId) {
        return null;
    }
    const showFormStyles =
        wizardStep === wizardSteps.FORM && !(messageState && messageState.blocking);
    const showDeploymentStyles =
        wizardStep === wizardSteps.DEPLOYMENT && !(messageState && messageState.blocking);
    const selectedClusterName = (selectedCluster && selectedCluster.name) || '';

    // @TODO: improve error handling when adding support for new clusters
    const panelButtons = (
        <PanelButton
            icon={
                showFormStyles ? (
                    <Icon.ArrowRight className="h-4 w-4" />
                ) : (
                    <Icon.Check className="h-4 w-4" />
                )
            }
            className={`mr-2 btn ${showFormStyles ? 'btn-base' : 'btn-success'}`}
            onClick={onNext}
            disabled={showFormStyles && Object.keys(validate(selectedCluster)).length !== 0}
            tooltip={showFormStyles ? 'Next' : 'Finish'}
        >
            {showFormStyles ? 'Next' : 'Finish'}
        </PanelButton>
    );

    const showPanelButtons = !messageState || !messageState.blocking;

    const upgradeStatus = selectedCluster?.status?.upgradeStatus ?? null;
    const certExpiryStatus = selectedCluster?.status?.certExpiryStatus ?? null;

    const upgradeStateObject = findUpgradeState(upgradeStatus);
    const upgradeStatusDetail = upgradeStatus && getUpgradeStatusDetail(upgradeStatus);
    const upgradeMessage =
        upgradeStatus && formatUpgradeMessage(upgradeStateObject, upgradeStatusDetail);
    const credentialExpirationProps = getCredentialExpirationProps(certExpiryStatus);
    const initiationOfCertRotation = initiationOfCertRotationIfApplicable(upgradeStatus);

    return (
        <Panel
            id="clusters-side-panel"
            header={selectedClusterName}
            headerComponents={showPanelButtons ? panelButtons : <div />}
            bodyClassName="pt-4"
            className="w-full h-full absolute right-0 top-0 md:w-1/2 min-w-72 md:relative z-0 bg-base-100"
            onClose={unselectCluster}
        >
            <Dialog
                className="w-1/3"
                isOpen={showCertRotationModal}
                text={`Select "Apply Update" to create new credentials in your cluster immediately. Each StackRox service begins using the new credentials after it restarts.`}
                onConfirm={triggerClusterCertRotation}
                confirmText="Apply update"
                onCancel={hideCertRotationModal}
            />

            {!!messageState && (
                <div className="m-4">
                    <Message type={messageState.type} message={messageState.message} />
                </div>
            )}
            {freshCredentialsDownloaded && (
                <div className="m-4">
                    <Message
                        message={
                            <div className="flex-1">
                                Fresh credentials downloaded. Use{' '}
                                <span className="italic text-accent-800">
                                    {' '}
                                    {selectedCluster.type === 'OPENSHIFT_CLUSTER'
                                        ? 'oc'
                                        : 'kubectl'}{' '}
                                    apply -f
                                </span>{' '}
                                to apply the credentials to your cluster.
                            </div>
                        }
                    />
                </div>
            )}
            {!!credentialExpirationProps &&
                credentialExpirationProps.showExpiringSoon &&
                !freshCredentialsDownloaded && (
                    <div data-testid="credential-expiration-banner" className="m-4">
                        <Message
                            type={credentialExpirationProps.messageType}
                            message={
                                <div className="flex-1">
                                    This cluster’s credentials expire in{' '}
                                    {credentialExpirationProps.diffInWords}. To use renewed
                                    certificates,{' '}
                                    <Button
                                        text="download this YAML file"
                                        className="text-tertiary-700 hover:text-tertiary-800 underline font-700 justify-center"
                                        onClick={generateCertSecret}
                                    />{' '}
                                    and apply it to your cluster
                                    {isUpToDateStateObject(upgradeStateObject) ? (
                                        <>
                                            , or{' '}
                                            <Button
                                                text="apply credentials by using an automatic upgrade"
                                                className="text-tertiary-700 hover:text-tertiary-800 underline font-700 justify-center"
                                                onClick={openCertRotationModal}
                                            />
                                            .
                                            {initiationOfCertRotation && (
                                                <p className="py-2">
                                                    An automatic upgrade applied renewed credentials
                                                    on{' '}
                                                    {format(
                                                        initiationOfCertRotation,
                                                        'MMMM D, YYYY'
                                                    )}
                                                    , at{' '}
                                                    {format(initiationOfCertRotation, 'h:mm a')}.
                                                    Each StackRox service begins using its new
                                                    credentials the next time it restarts.
                                                </p>
                                            )}
                                        </>
                                    ) : (
                                        '.'
                                    )}
                                </div>
                            }
                        />
                    </div>
                )}
            {!!upgradeMessage && (
                <div className="px-4 w-full">
                    <CollapsibleCard
                        title="Upgrade Status"
                        cardClassName="border border-base-400 mb-2"
                        titleClassName="border-b border-base-300 bg-primary-200 leading-normal cursor-pointer flex justify-between items-center hover:bg-primary-300 hover:border-primary-300"
                    >
                        <div className="m-4">
                            <Message type={upgradeMessage.type} message={upgradeMessage.message} />
                            {upgradeMessage.detail !== '' && (
                                <div className="mt-2 flex flex-col items-center">
                                    <div className="bg-base-200">
                                        <div className="whitespace-normal overflow-x-scroll">
                                            {upgradeMessage.detail}
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>
                    </CollapsibleCard>
                </div>
            )}
            {submissionError && submissionError.length > 0 && (
                <div className="w-full">
                    <div className="mb-4 mx-4">
                        <Message type="error" message={submissionError} />
                    </div>
                </div>
            )}
            {showFormStyles && (
                <ClusterEditForm
                    centralEnv={centralEnv}
                    centralVersion={metadata.version}
                    selectedCluster={selectedCluster}
                    handleChange={onChange}
                    isLoading={loadingCounter > 0}
                />
            )}
            {showDeploymentStyles && (
                <ClusterDeployment
                    editing={!!selectedCluster}
                    createUpgraderSA={createUpgraderSA}
                    toggleSA={toggleSA}
                    onFileDownload={onDownload}
                    clusterCheckedIn={
                        !!(
                            selectedCluster &&
                            selectedCluster.healthStatus &&
                            selectedCluster.healthStatus.lastContact
                        )
                    }
                />
            )}
        </Panel>
    );
}

ClustersSidePanel.propTypes = {
    metadata: PropTypes.shape({ version: PropTypes.string, releaseBuild: PropTypes.bool })
        .isRequired,
    setSelectedClusterId: PropTypes.func.isRequired,
    selectedClusterId: PropTypes.string,
};

ClustersSidePanel.defaultProps = {
    selectedClusterId: '',
};

const mapStateToProps = createStructuredSelector({
    metadata: selectors.getMetadata,
});

export default connect(mapStateToProps)(ClustersSidePanel);
