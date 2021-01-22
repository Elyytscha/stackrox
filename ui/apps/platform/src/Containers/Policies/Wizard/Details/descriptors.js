import {
    enforcementActionLabels,
    lifecycleStageLabels,
    portExposureLabels,
    rbacPermissionLabels,
    envVarSrcLabels,
} from 'messages/common';

import { comparatorOp, formatResources, formatScope, formatDeploymentExcludedScope } from './utils';

// JSON value name mapped to formatting for description page.
const fieldsMap = {
    id: {
        label: 'ID',
        formatValue: (d) => d,
    },
    name: {
        label: 'Name',
        formatValue: (d) => d,
    },
    lifecycleStages: {
        label: 'Lifecycle Stage',
        formatValue: (d) => d.map((v) => lifecycleStageLabels[v]).join(', '),
    },
    severity: {
        label: 'Severity',
        formatValue: (d) => {
            switch (d) {
                case 'CRITICAL_SEVERITY':
                    return 'Critical';
                case 'HIGH_SEVERITY':
                    return 'High';
                case 'MEDIUM_SEVERITY':
                    return 'Medium';
                case 'LOW_SEVERITY':
                    return 'Low';
                default:
                    return '';
            }
        },
    },
    description: {
        label: 'Description',
        formatValue: (d) => d,
    },
    rationale: {
        label: 'Rationale',
        formatValue: (r) => r,
    },
    remediation: {
        label: 'Remediation',
        formatValue: (r) => r,
    },
    notifiers: {
        label: 'Notifications',
        formatValue: (d, props) =>
            props.notifiers
                .filter((n) => d.includes(n.id))
                .map((n) => n.name)
                .join(', '),
    },
    scope: {
        label: 'Restricted to Scopes',
        formatValue: (d, props) =>
            d && d.length ? d.map((scope) => formatScope(scope, props)) : null,
    },
    enforcementActions: {
        label: 'Enforcement Action',
        formatValue: (d) => d.map((v) => enforcementActionLabels[v]).join(', '),
    },
    disabled: {
        label: 'Enabled',
        formatValue: (d) => (d !== true ? 'Yes' : 'No'),
    },
    categories: {
        label: 'Categories',
        formatValue: (d) => d.join(', '),
    },
    exclusions: {
        label: 'Exclusions',
        formatValue: (d, props) => {
            const exclusionObj = {};
            const deploymentExcludedScopes = d
                .filter((obj) => obj.deployment && (obj.deployment.name || obj.deployment.scope))
                .map((obj) => obj.deployment);
            if (deploymentExcludedScopes.length > 0) {
                exclusionObj[
                    'Excluded Deployments'
                ] = deploymentExcludedScopes.map((deploymentExcludedScope) =>
                    formatDeploymentExcludedScope(deploymentExcludedScope, props)
                );
            }
            const images = d
                .filter((obj) => obj.image && obj.image.name !== '')
                .map((obj) => obj.image.name);
            if (images.length !== 0) {
                exclusionObj['Excluded Images'] = images;
            }
            return exclusionObj;
        },
    },
    imageName: {
        label: 'Image',
        formatValue: (d) => {
            const remote = d.remote ? `images named ${d.remote}` : 'any image';
            const tag = d.tag ? `tag ${d.tag}` : 'any tag';
            const registry = d.registry ? `registry ${d.registry}` : 'any registry';
            return `Alert on ${remote} using ${tag} from ${registry}`;
        },
    },
    imageAgeDays: {
        label: 'Days since image was created',
        formatValue: (d) => (d !== '0' ? `${Number(d)} Days ago` : ''),
    },
    noScanExists: {
        label: 'Image Scan Status',
        formatValue: () => 'Verify that the image is scanned',
    },
    scanAgeDays: {
        label: 'Days since image was last scanned',
        formatValue: (d) => (d !== '0' ? `${Number(d)} Days ago` : ''),
    },
    imageUser: {
        label: 'Image User',
        formatValue: (d) => d,
    },
    lineRule: {
        label: 'Dockerfile Line',
        formatValue: (d) => `${d.instruction} ${d.value}`,
    },
    cvss: {
        label: 'CVSS',
        formatValue: (d) => `${comparatorOp[d.op]} ${d.value}`,
    },
    cve: {
        label: 'CVE',
        formatValue: (d) => d,
    },
    fixedBy: {
        label: 'Fixed By',
        formatValue: (d) => d,
    },
    component: {
        label: 'Image Component',
        formatValue: (d) => {
            const name = d.name ? `${d.name}` : '';
            const version = d.version ? d.version : '';
            return `"${name}" with version "${version}"`;
        },
    },
    env: {
        label: 'Environment Variable',
        formatValue: (d) => {
            const key = d.key ? `${d.key}` : '';
            const value = d.value ? d.value : '';
            const valueFrom = !d.envVarSource
                ? ''
                : ` Value From: ${envVarSrcLabels[d.envVarSource]}`;
            return `${key}=${value};${valueFrom}`;
        },
    },
    disallowedAnnotation: {
        label: 'Disallowed Annotation',
        formatValue: (d) => {
            const key = d.key ? `key=${d.key}` : '';
            const value = d.value ? `value=${d.value}` : '';
            const comma = d.key && d.value ? ', ' : '';
            return `Alerts on deployments with the disallowed annotation ${key}${comma}${value}`;
        },
    },
    requiredLabel: {
        label: 'Required Label',
        formatValue: (d) => {
            const key = d.key ? `key=${d.key}` : '';
            const value = d.value ? `value=${d.value}` : '';
            const comma = d.key && d.value ? ', ' : '';
            return `Alerts on deployments missing the required label ${key}${comma}${value}`;
        },
    },
    requiredAnnotation: {
        label: 'Required Annotation',
        formatValue: (d) => {
            const key = d.key ? `key=${d.key}` : '';
            const value = d.value ? `value=${d.value}` : '';
            const comma = d.key && d.value ? ', ' : '';
            return `Alerts on deployments missing the required annotation ${key}${comma}${value}`;
        },
    },
    volumePolicy: {
        label: 'Volume Policy',
        formatValue: (d) => {
            const output = [];
            if (d.name) {
                output.push(`Name: ${d.name}`);
            }
            if (d.type) {
                output.push(`Type: ${d.type}`);
            }
            if (d.source) {
                output.push(`Source: ${d.source}`);
            }
            if (d.destination) {
                output.push(`Dest: ${d.destination}`);
            }
            output.push(d.readOnly ? 'Writable: No' : 'Writable: Yes');
            return output.join(', ');
        },
    },
    nodePortPolicy: {
        label: 'Node Port',
        formatValue: (d) => {
            const protocol = d.protocol ? `${d.protocol} ` : '';
            const port = d.port ? d.port : '';
            return `${protocol}${port}`;
        },
    },
    portPolicy: {
        label: 'Port',
        formatValue: (d) => {
            const protocol = d.protocol ? `${d.protocol} ` : '';
            const port = d.port ? d.port : '';
            return `${protocol}${port}`;
        },
    },
    dropCapabilities: {
        label: 'Drop Capabilities',
        formatValue: (d) => d.join(', '),
    },
    addCapabilities: {
        label: 'Add Capabilities',
        formatValue: (d) => d.join(', '),
    },
    privileged: {
        label: 'Privileged',
        formatValue: (d) => (d === true ? 'Yes' : 'No'),
    },
    readOnlyRootFs: {
        label: 'Read Only Root Filesystem',
        formatValue: (d) => (d === true ? 'Yes' : 'Not Enabled'),
    },
    HostPid: {
        label: 'Host PID',
        formatValue: (d) => (d === true ? 'Yes' : 'Not Enabled'),
    },
    containerResourcePolicy: {
        label: 'Container Resources',
        formatValue: formatResources,
    },
    processPolicy: {
        label: 'Process Execution',
        formatValue: (d) => {
            const name = d.name ? `Process matches name "${d.name}"` : 'Process';
            const args = d.args ? `and matches args "${d.args}"` : '';
            const ancestor = d.ancestor ? `and has ancestor matching "${d.ancestor}"` : '';
            const uid = d.uid ? `with uid ${d.uid}` : ``;
            return `${name} ${args} ${ancestor} ${uid}`;
        },
    },
    portExposurePolicy: {
        label: 'Port Exposure',
        formatValue: (d) => {
            const output = d.exposureLevels.map((element) => portExposureLabels[element]);
            return output.join(', ');
        },
    },
    hostMountPolicy: {
        label: 'Host Mount Policy',
        formatValue: (d) => (d.readOnly ? 'Not Enabled' : 'Writable: Yes'),
    },
    whitelistEnabled: {
        label: 'Excluded Scopes Enabled',
        formatValue: (d) => (d ? 'Yes' : 'No'),
    },
    permissionPolicy: {
        label: 'Minimum RBAC Permissions',
        formatValue: (d) => rbacPermissionLabels[d.permissionLevel],
    },
    requiredImageLabel: {
        label: 'Required Image Label',
        formatValue: (d) => {
            const key = d.key ? `key=${d.key}` : '';
            const value = d.value ? `value=${d.value}` : '';
            const comma = d.key && d.value ? ', ' : '';
            return `Alerts on deployments with images missing the required label ${key}${comma}${value}`;
        },
    },
    disallowedImageLabel: {
        label: 'Disallowed Image Label',
        formatValue: (d) => {
            const key = d.key ? `key=${d.key}` : '';
            const value = d.value ? `value=${d.value}` : '';
            const comma = d.key && d.value ? ', ' : '';
            return `Alerts on deployments with disallowed image label ${key}${comma}${value}`;
        },
    },
};

export default fieldsMap;
