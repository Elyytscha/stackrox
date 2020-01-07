export const url = '/main/violations';

export const selectors = {
    navLink: 'nav li:contains("Violations") a',
    rows: '.rt-tbody .rt-tr',
    firstTableRow: '.rt-tbody :nth-child(1) > .rt-tr',
    firstPanelTableRow: '.rt-tbody > :nth-child(1) > .rt-tr',
    lastTableRow: '.rt-tr:last',
    panels: '[data-test-id="panel"]',
    sidePanel: {
        header: '[data-test-id="panel-header"]',
        tabs: 'button[data-test-id="tab"]',
        getTabByIndex: index => `button[data-test-id="tab"]:nth(${index})`
    },
    clusterTableHeader: '.rt-thead > .rt-tr > div:contains("Cluster")',
    viewDeploymentsButton: 'button:contains("View Deployments")',
    modal: '.ReactModalPortal > .ReactModal__Overlay',
    clusterFieldInModal: '.ReactModalPortal > .ReactModal__Overlay span:contains("Cluster")',
    collapsible: {
        header: '.Collapsible__trigger',
        body: '.Collapsible__contentInner'
    },
    securityBestPractices: '[data-test-id="deployment-security-practices"]',
    runtimeProcessCards: '[data-testid="runtime-processes"]',
    lifeCycleColumn: '.rt-thead.-header:contains("Lifecycle")',
    whitelistDeploymentButton: '[data-test-id="whitelist-deployment-button"]',
    resolveButton: '[data-test-id="resolve-button"]',
    whitelistDeploymentRow: '.rt-tr:contains("metadata-proxy-v0.1")'
};
