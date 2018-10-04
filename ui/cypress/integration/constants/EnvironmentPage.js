export const url = '/main/network';

export const selectors = {
    network: 'nav.left-navigation li:contains("Network") a',
    simulatorSuccessMessage: 'div:contains("YAML file uploaded successfully")',
    panels: {
        simulatorPanel: '[data-test-id="network-simulator-panel"]',
        uploadPanel: '[data-test-id="upload-yaml-panel"]'
    },
    legend: {
        deployment: '.env-graph-legend > :nth-child(1)',
        namespace: '.env-graph-legend > :nth-child(2)',
        ingressEgress: '.env-graph-legend > :nth-child(3)',
        internetEgress: '.env-graph-legend > :nth-child(4)'
    },
    namespaces: {
        all: 'g.container > rect',
        getNamespace: namespace => `g.container > rect.namespace-${namespace}`
    },
    services: {
        all: 'g.namespace .node',
        getServicesForNamespace: namespace => `g.namespace-${namespace} .node`
    },
    links: {
        all: '.link',
        bidirectional: '.link[marker-start="url(#start)"]',
        namespaces: '.link.namespace',
        services: '.link.service'
    },
    buttons: {
        simulatorButtonOn: '[data-test-id="simulator-button-on"]',
        simulatorButtonOff: '[data-test-id="simulator-button-off"]'
    }
};
