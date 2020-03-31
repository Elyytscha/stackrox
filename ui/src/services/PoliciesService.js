import { normalize } from 'normalizr';
import queryString from 'qs';
import axios from './instance';

import { policy as policySchema } from './schemas';

const baseUrl = '/v1/policies';
const policyCategoriesUrl = '/v1/policyCategories';

/**
 * Fetches policy summary for a given policy ID.
 * Returns normalized response with policy entity extracted.
 *
 * @param {!string} policyId
 * @returns {Promise<Object, Error>} fulfilled with normalized response
 */
export function fetchPolicy(policyId) {
    return axios.get(`${baseUrl}/${policyId}`).then(response => ({
        response: normalize(response.data, policySchema)
    }));
}

/**
 * Fetches a list of policies.
 *
 * @param {!string} filters
 * @returns {Promise<Object, Error>} fulfilled with array of policies (as defined in .proto)
 */
export function fetchPolicies(filters) {
    const params = queryString.stringify({ ...filters }, { arrayFormat: 'repeat' });
    return axios.get(`${baseUrl}?${params}`).then(response => ({
        response: normalize(response.data, { policies: [policySchema] })
    }));
}

/**
 * Fetches a list of policy categories.
 *
 * @returns {Promise<Object, Error>}
 */
export function fetchPolicyCategories() {
    return axios.get(policyCategoriesUrl).then(response => ({
        response: response.data
    }));
}

/**
 * Reassesses policies.
 *
 * @returns {Promise<AxiosResponse, Error>}
 */
export function reassessPolicies() {
    return axios.post(`${baseUrl}/reassess`);
}

/**
 * Deletes a policy with a given id.
 *
 * @param {!string} policyId
 * @returns {Promise<AxiosResponse, Error>}
 */
export function deletePolicy(policyId) {
    return axios.delete(`${baseUrl}/${policyId}`);
}

/**
 * Deletes a list of policies by policyId.
 *
 * @param {string[]} policyIds
 * @returns {Promise<AxiosResponse, Error>}
 */
export function deletePolicies(policyIds = []) {
    return Promise.all(policyIds.map(policyId => deletePolicy(policyId)));
}

/**
 * Enable / Disable notification to notifiers given by notifierIds for policy given by policyId.
 *
 * @param {!string} policyId
 * @param {!object} data
 * @returns {Promise<AxiosResponse, Error>}
 */
export function enableDisablePolicyNotifications(policyId, data) {
    return axios.patch(`${baseUrl}/${policyId}/notifiers`, data);
}

/**
 * Enable / Disable notification to notifiers given by notifierIds for list of policies given by policyIds.
 *
 * @param {!string[]} policyIds
 * @param {!string[]} notifierIds
 * @param {!boolean} disable
 * @returns {Promise<AxiosResponse, Error>}
 */
export function enableDisableNotificationsForPolicies(policyIds, notifierIds, disable) {
    const data = { notifierIds, disable };
    return Promise.all(policyIds.map(policyId => enableDisablePolicyNotifications(policyId, data)));
}

/**
 * Saves a given policy.
 *
 * @param {!object} policy
 * @returns {Promise<AxiosResponse, Error>}
 */
export function savePolicy(policy) {
    if (!policy.id) throw new Error('Policy entity must have an id to be saved');
    return axios.put(`${baseUrl}/${policy.id}`, policy);
}

/**
 * Creates a new policy.
 *
 * @param {!object} policy
 * @returns {Promise<AxiosResponse, Error>}
 */
export function createPolicy(policy) {
    return axios.post(`${baseUrl}`, policy);
}

/**
 * Gets a dry run for a given policy.
 *
 * @param {!object} policy
 * @returns {Promise<AxiosResponse, Error>}
 */
export function getDryRun(policy) {
    return axios.post(`${baseUrl}/dryrun`, policy);
}

/**
 * Starts a dry run for a given policy.
 *
 * @param {!object} policy
 * @returns {Promise<AxiosResponse, Error>}
 */
export function startDryRun(policy) {
    return axios.post(`${baseUrl}/submitdryrunjob`, policy);
}

/**
 * Gets a dry run for a given policy.
 *
 * @param {!object} policy
 * @returns {Promise<AxiosResponse, Error>}
 */
export function checkDryRun(jobId) {
    return axios.get(`${baseUrl}/dryrunjob/${jobId}`);
}

/**
 * Updates policy with a given ID to add deployment into the whitelisted entries.
 *
 * @param {!string} policyId
 * @param {!string[]} deploymentNames
 * @returns {Promise<AxiosResponse, Error>} fulfilled in case of success or rejected with an error
 */
export async function whitelistDeployments(policyId, deploymentNames) {
    const { response } = await fetchPolicy(policyId);
    const policy = response.entities.policy[policyId];

    const deploymentEntries = deploymentNames.map(name => ({
        deployment: { name }
    }));
    policy.whitelists = [...policy.whitelists, ...deploymentEntries];
    return axios.put(`${baseUrl}/${policy.id}`, policy);
}

/**
 * Send request to enable / disable policy with a given ID.
 *
 * @param {!string} policyId
 * @param {!boolean} disabled if policy should be disabled
 * @returns {Promise<AxiosResponse, Error>} fulfilled in case of success or rejected with an error
 */
export function updatePolicyDisabledState(policyId, disabled) {
    return axios.patch(`${baseUrl}/${policyId}`, { disabled });
}
