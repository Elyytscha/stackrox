import queryString from 'qs';
import axios from './instance';

const networkPoliciesBaseUrl = '/v1/networkpolicies';
const networkFlowBaseUrl = '/v1/networkgraph';

// for large clusters network graph requests may take time to process, so
// removing any global default timeout
const NETWORK_GRAPH_REQUESTS_TIMEOUT = 0;

/**
 * Fetches nodes and links for the network graph.
 * Returns response with nodes and links
 *
 * @returns {Promise<Object, Error>}
 */
export function fetchNetworkPolicyGraph(clusterId, query, modification, includePorts) {
    const urlParams = { query };
    if (includePorts) {
        urlParams.includePorts = true;
    }
    const params = queryString.stringify(urlParams, { arrayFormat: 'repeat' });
    let options;
    let getGraph = (data) => data;
    if (modification) {
        options = {
            method: 'POST',
            data: modification,
            url: `${networkPoliciesBaseUrl}/simulate/${clusterId}?${params}`,
        };
        getGraph = ({ simulatedGraph }) => simulatedGraph;
    } else {
        options = {
            method: 'GET',
            url: `${networkPoliciesBaseUrl}/cluster/${clusterId}?${params}`,
        };
    }
    options = {
        ...options,
        timeout: NETWORK_GRAPH_REQUESTS_TIMEOUT,
    };
    return axios(options).then((response) => ({
        response: getGraph(response.data),
    }));
}

/**
 * Fetches nodes and links for the network flow graph.
 * Returns response with nodes and links
 *
 * @returns {Promise<Object, Error>}
 */
export function fetchNetworkFlowGraph(clusterId, query, date, includePorts) {
    const urlParams = { query };
    if (date) {
        urlParams.since = date.toISOString();
    }
    if (includePorts) {
        urlParams.includePorts = true;
    }

    const params = queryString.stringify(urlParams, { arrayFormat: 'repeat' });
    const options = {
        method: 'GET',
        url: `${networkFlowBaseUrl}/cluster/${clusterId}?${params}`,
        timeout: NETWORK_GRAPH_REQUESTS_TIMEOUT,
    };
    return axios(options).then((response) => ({
        response: response.data,
    }));
}

/**
 * Fetches policies details for given array of ids.
 *
 * @param {!array} policyIds
 * @returns {Promise<Object, Error>}
 */
export function fetchNetworkPolicies(policyIds) {
    const networkPoliciesPromises = policyIds.map((policyId) =>
        axios.get(`${networkPoliciesBaseUrl}/${policyId}`)
    );
    return Promise.all(networkPoliciesPromises).then((response) => ({
        response: response.map((networkPolicy) => networkPolicy.data),
    }));
}

/**
 * Fetches Node updates.
 *
 * @returns {Promise<Object, Error>}
 */
export function fetchNodeUpdates(clusterId) {
    return axios
        .get(`${networkPoliciesBaseUrl}/graph/epoch?clusterId=${clusterId}`)
        .then((response) => ({
            response: response.data,
        }));
}

/**
 * Fetches the network policies currently applied to a cluster and set of deployments (defined by query).
 *
 * @param {!String} clusterId
 * @param {!Object} query
 * @returns {Promise<Object, Error>}
 */
export function getActiveNetworkModification(clusterId, deploymentQuery) {
    let params;
    if (deploymentQuery) {
        params = queryString.stringify({ clusterId, deploymentQuery }, { arrayFormat: 'repeat' });
    } else {
        params = queryString.stringify({ clusterId });
    }
    const options = {
        method: 'GET',
        url: `${networkPoliciesBaseUrl}?${params}`,
    };
    return axios(options).then((response) => {
        const policies = response.data.networkPolicies;
        if (policies) {
            return { applyYaml: policies.map((policy) => policy.yaml).join('\n---\n') };
        }
        return null;
    });
}

/**
 * Retrieves the modification that will undo the last action done through the stackrox UI.
 *
 * @param {!String} clusterId
 * @param {!Object} query
 * @returns {Promise<Object, Error>}
 */
export function getUndoNetworkModification(clusterId) {
    const options = {
        method: 'GET',
        url: `${networkPoliciesBaseUrl}/undo/${clusterId}`,
    };
    return axios(options).then((response) => response.data.undoRecord.undoModification);
}

/**
 * Generates a modification to policies based on a graph.
 *
 * @param {!String} clusterId
 * @param {!Object} query
 * @param {!String} date
 * @param {Boolean} excludePortsProtocols
 * @returns {Promise<Object, Error>}
 */
export function generateNetworkModification(clusterId, query, date, excludePortsProtocols = null) {
    const urlParams = { query };
    if (date) {
        urlParams.networkDataSince = date.toISOString();
    }

    if (excludePortsProtocols !== null) {
        urlParams.includePorts = !excludePortsProtocols;
    }

    const params = queryString.stringify(urlParams, { arrayFormat: 'repeat' });
    const options = {
        method: 'GET',
        url: `${networkPoliciesBaseUrl}/generate/${clusterId}?deleteExisting=NONE&${params}`,
    };
    return axios(options).then((response) => response.data.modification);
}

/**
 * Sends a notification of the simulated yaml
 *
 * @param {!String} clusterId
 * @param {!array} notifierIds
 * @param {!Object} modification
 * @returns {Promise<Object, Error>}
 */
export function notifyNetworkPolicyModification(clusterId, notifierIds, modification) {
    const notifiers = queryString.stringify({ notifierIds }, { arrayFormat: 'repeat' });
    const options = {
        method: 'POST',
        data: modification,
        url: `${networkPoliciesBaseUrl}/simulate/${clusterId}/notify?${notifiers}`,
    };
    return axios(options).then((response) => ({
        response: response.data,
    }));
}

/**
 * Sends a yaml to the backed for application to a cluster.
 *
 * @param {!String} clusterId
 * @param {!Object} modification
 * @returns {Promise<Object, Error>}
 */
export function applyNetworkPolicyModification(clusterId, modification) {
    const options = {
        method: 'POST',
        data: modification,
        url: `${networkPoliciesBaseUrl}/apply/${clusterId}`,
        timeout: NETWORK_GRAPH_REQUESTS_TIMEOUT,
    };
    return axios(options).then((response) => ({
        response: response.data,
    }));
}
