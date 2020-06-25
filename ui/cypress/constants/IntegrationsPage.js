import table from '../selectors/table';
import toast from '../selectors/toast';
import tooltip from '../selectors/tooltip';

export const url = '/main/integrations';

export const selectors = {
    configure: 'nav.left-navigation li:contains("Platform Configuration") a',
    navLink: '.navigation-panel li:contains("Integrations") a',
    kubernetesTile: 'div[role="button"]:contains("Kubernetes")',
    dockerRegistryTile: 'div[role="button"]:contains("Generic Docker Registry")',
    tiles: 'div[role="button"]',
    clairTile: 'div[role="button"]:contains("CoreOS Clair")',
    clairifyTile: 'div[role="button"]:contains("Clairify")',
    slackTile: 'div[role="button"]:contains("Slack")',
    dockerTrustedRegistryTile: 'div[role="button"]:contains("Docker Trusted Registry")',
    quayTile: 'div[role="button"]:contains("Quay.io")',
    amazonECRTile: 'div[role="button"]:contains("Amazon ECR")',
    tenableTile: 'div[role="button"]:contains("Tenable.io")',
    googleCloudPlatformTile: 'div[role="button"]:contains("Google Cloud Platform")',
    anchoreScannerTile: 'div[role="button"]:contains("Anchore Scanner")',
    ibmCloudTile: 'div[role="button"]:contains("IBM Cloud")',
    apiTokenTile: 'div[role="button"]:contains("API Token")',
    clusters: {
        k8sCluster0: 'div.rt-td:contains("Kubernetes Cluster 0")',
    },
    buttons: {
        new: 'button:contains("New")',
        next: 'button:contains("Next")',
        downloadYAML: 'button:contains("Download YAML")',
        delete: 'button:contains("Delete")',
        test: 'button:contains("Test")',
        create: 'button:contains("Create")',
        save: 'button:contains("Save")',
        confirm: 'button:contains("Confirm")',
        generate: 'button:contains("Generate"):not([disabled])',
        revoke: 'button:contains("Revoke")',
        closePanel: 'button[data-testid="cancel"]',
    },
    apiTokenForm: {
        nameInput: 'form[data-testid="api-token-form"] input[name="name"]',
        roleSelect: 'form[data-testid="api-token-form"] .react-select__control',
    },
    apiTokenBox: 'span:contains("eyJ")', // all API tokens start with eyJ
    apiTokenDetailsDiv: 'div[data-testid="api-token-details"]',
    clusterForm: {
        nameInput: 'form[data-testid="cluster-form"] input[name="name"]',
        imageInput: 'form[data-testid="cluster-form"] input[name="mainImage"]',
        endpointInput: 'form[data-testid="cluster-form"] input[name="centralApiEndpoint"]',
    },
    dockerRegistryForm: {
        nameInput: "form input[name='name']",
        typesSelect: 'form .react-select__control',
        endpointInput: "form input[name='docker.endpoint']",
    },
    slackForm: {
        nameInput: "form input[name='name']",
        defaultWebhook: "form input[name='labelDefault']",
        labelAnnotationKey: "form input[name='labelKey']",
    },
    labeledValue: '[data-testid="labeled-value"]',
    plugins: '.mb-6:first div[role="button"]',
    dialog: '.dialog',
    checkboxes: 'input',
    table,
    toast,
    tooltip,
};
