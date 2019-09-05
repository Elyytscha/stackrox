import dateFns from 'date-fns';
import get from 'lodash/get';

import dateTimeFormat from 'constants/dateTimeFormat';

export const runtimeOptions = [
    {
        label: 'No Runtime Support',
        tableDisplay: 'None',
        value: 'NO_COLLECTION'
    },
    {
        label: 'Kernel Module Support',
        tableDisplay: 'Kernel Module',
        value: 'KERNEL_MODULE'
    },
    {
        label: 'eBPF Support',
        tableDisplay: 'eBPF',
        value: 'EBPF'
    }
];

export const clusterTypeOptions = [
    {
        label: 'Kubernetes',
        tableDisplay: 'Kubernetes',
        value: 'KUBERNETES_CLUSTER'
    },
    {
        label: 'OpenShift',
        tableDisplay: 'OpenShift',
        value: 'OPENSHIFT_CLUSTER'
    }
];

export const clusterTablePollingInterval = 5000; // milliseconds

const defaultNewClusterType = 'KUBERNETES_CLUSTER';
const defaultCollectionMethod = 'NO_COLLECTION';

export const newClusterDefault = {
    id: null,
    name: '',
    type: defaultNewClusterType,
    mainImage: '',
    collectorImage: '',
    centralApiEndpoint: '',
    runtimeSupport: false,
    monitoringEndpoint: '',
    collectionMethod: defaultCollectionMethod,
    DEPRECATEDProviderMetadata: null,
    admissionController: false,
    DEPRECATEDOrchestratorMetadata: null,
    status: null,
    dynamicConfig: {
        admissionControllerConfig: {
            enabled: false,
            timeoutSeconds: 3,
            scanInline: false,
            disableBypass: false
        }
    }
};

// @TODO: add optional button text and func
const upgradeStates = {
    UP_TO_DATE: {
        displayValue: 'On the latest version',
        type: 'current'
    },
    MANUAL_UPGRADE_REQUIRED: {
        displayValue: 'Manual upgrade required',
        type: 'intervention'
    },
    UNSET: {
        type: 'download',
        action: {
            actionText: 'Upgrade available'
        }
    },
    UPGRADE_TRIGGER_SENT: {
        displayValue: 'Upgrade trigger sent',
        type: 'progress'
    },
    UPGRADER_LAUNCHING: {
        displayValue: 'Upgrader launching',
        type: 'progress'
    },
    UPGRADER_LAUNCHED: {
        displayValue: 'Upgrader launched',
        type: 'progress'
    },
    PRE_FLIGHT_CHECKS_COMPLETE: {
        displayValue: 'Pre-flight checks complete',
        type: 'progress'
    },
    PRE_FLIGHT_CHECKS_FAILED: {
        displayValue: 'Pre-flight checks failed.',
        type: 'failure'
    },
    UPGRADE_OPERATIONS_DONE: {
        displayValue: 'Upgrade Operations Done',
        type: 'progress'
    },
    UPGRADE_OPERATIONS_COMPLETE: {
        displayValue: 'Upgrade Operations Complete',
        type: 'current'
    },
    UPGRADE_ERROR_ROLLED_BACK: {
        displayValue: 'Upgrade failed. Rolled back.',
        type: 'failure',
        action: {
            actionText: 'Retry upgrade'
        }
    },
    UPGRADE_ERROR_ROLLBACK_FAILED: {
        displayValue: 'Upgrade failed. Rollback failed.',
        type: 'failure'
    },
    unknown: {
        displayValue: 'Undeterminate upgrade state!',
        type: 'intervention'
    }
};

function findOptionInList(options, value) {
    return options.find(opt => opt.value === value);
}

export function formatClusterType(value) {
    const match = findOptionInList(clusterTypeOptions, value);

    return match.tableDisplay;
}

export function formatCollectionMethod(value) {
    const match = findOptionInList(runtimeOptions, value);

    return match.tableDisplay;
}

export function formatEnabledDisabledField(value) {
    return value ? 'Enabled' : 'Disabled';
}
export function formatLastCheckIn(status) {
    if (status && status.lastContact) {
        return dateFns.format(status.lastContact, dateTimeFormat);
    }

    return 'N/A';
}

export function formatSensorVersion(status) {
    return (status && status.sensorVersion) || 'Not Running';
}

export function parseUpgradeStatus(cluster) {
    const upgradability = get(cluster, 'status.upgradeStatus.upgradability', undefined);
    switch (upgradability) {
        case 'UP_TO_DATE':
        case 'MANUAL_UPGRADE_REQUIRED': {
            return upgradeStates[upgradability];
        }
        case 'AUTO_UPGRADE_POSSIBLE': {
            const upgradeState = get(
                cluster,
                'status.upgradeStatus.upgradeProgress.upgradeState',
                'unknown'
            );

            return upgradeStates[upgradeState] || upgradeStates.unknown;
        }
        default: {
            return upgradeStates.unknown;
        }
    }
}

export const wizardSteps = Object.freeze({
    FORM: 'FORM',
    DEPLOYMENT: 'DEPLOYMENT'
});

export default {
    runtimeOptions,
    clusterTypeOptions,
    clusterTablePollingInterval,
    newClusterDefault,
    formatClusterType,
    formatCollectionMethod,
    formatEnabledDisabledField,
    formatLastCheckIn,
    parseUpgradeStatus,
    wizardSteps
};
