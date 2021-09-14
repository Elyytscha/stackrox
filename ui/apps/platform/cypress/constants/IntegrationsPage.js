import table from '../selectors/table';
import toast from '../selectors/toast';
import tooltip from '../selectors/tooltip';
import navigationSelectors from '../selectors/navigation';

export const url = '/main/integrations';

export const selectors = {
    configure: `${navigationSelectors.navExpandable}:contains("Platform Configuration")`,
    navLink: `${navigationSelectors.navLinks}:contains("Integrations")`,
    kubernetesTile: 'a[data-testid="integration-tile"]:contains("Kubernetes")',
    dockerRegistryTile: 'a[data-testid="integration-tile"]:contains("Generic Docker Registry")',
    tiles: 'a[data-testid="integration-tile"]',
    clairTile: 'a[data-testid="integration-tile"]:contains("CoreOS Clair")',
    clairifyTile: 'a[data-testid="integration-tile"]:contains("StackRox Scanner")',
    slackTile: 'a[data-testid="integration-tile"]:contains("Slack")',
    dockerTrustedRegistryTile:
        'a[data-testid="integration-tile"]:contains("Docker Trusted Registry")',
    quayTile: 'a[data-testid="integration-tile"]:contains("Quay.io")',
    amazonECRTile: 'a[data-testid="integration-tile"]:contains("Amazon ECR")',
    tenableTile: 'a[data-testid="integration-tile"]:contains("Tenable.io")',
    googleContainerRegistryTile:
        'a[data-testid="integration-tile"]:contains("Google Container Registry")',
    anchoreScannerTile: 'a[data-testid="integration-tile"]:contains("Anchore Scanner")',
    ibmCloudTile: 'a[data-testid="integration-tile"]:contains("IBM Cloud")',
    microsoftACRTile: 'a[data-testid="integration-tile"]:contains("Microsoft ACR")',
    jFrogArtifactoryTile: 'a[data-testid="integration-tile"]:contains("JFrog Artifactory")',
    sonatypeNexusTile: 'a[data-testid="integration-tile"]:contains("Sonatype Nexus")',
    redHatTile: 'a[data-testid="integration-tile"]:contains("Red Hat")',
    googleCloudSCCTile: 'a[data-testid="integration-tile"]:contains("Google Cloud SCC")',
    jiraTile: 'a[data-testid="integration-tile"]:contains("Jira")',
    emailTile: 'a[data-testid="integration-tile"]:contains("Email")',
    splunkTile: 'a[data-testid="integration-tile"]:contains("Splunk")',
    pagerDutyTile: 'a[data-testid="integration-tile"]:contains("PagerDuty")',
    awsSecurityHubTile: 'a[data-testid="integration-tile"]:contains("AWS Security Hub")',
    sumologicTile: 'a[data-testid="integration-tile"]:contains("Sumo Logic")',
    syslogTile: 'a[data-testid="integration-tile"]:contains("Syslog")',
    teamsTile: 'a[data-testid="integration-tile"]:contains("Teams")',
    genericWebhookTile: 'a[data-testid="integration-tile"]:contains("Generic Webhook")',
    amazonS3Tile: 'a[data-testid="integration-tile"]:contains("Amazon S3")',
    googleCloudStorageTile: 'a[data-testid="integration-tile"]:contains("Google Cloud Storage")',
    scopedAccessPluginTile: 'a[data-testid="integration-tile"]:contains("Scoped Access Plugin")',
    apiTokenTile: 'a[data-testid="integration-tile"]:contains("API Token")',
    clusterInitBundleTile: 'a[data-testid="integration-tile"]:contains("Cluster Init Bundle")',
    clusters: {
        k8sCluster0: 'div.rt-td:contains("Kubernetes Cluster 0")',
    },
    buttons: {
        new: 'button:contains("New")',
        newApiToken: 'button:contains("Generate token")',
        newClusterInitBundle: 'button:contains("Generate bundle")',
        next: 'button:contains("Next")',
        downloadYAML: 'button:contains("Download YAML")',
        delete: 'button:contains("Delete")',
        test: 'button:contains("Test")',
        create: 'button:contains("Create")',
        save: 'button:contains("Save")',
        confirm: 'button:contains("Confirm")',
        generate: 'button:contains("Generate")',
        revoke: 'button:contains("Revoke")',
        closePanel: 'button[data-testid="cancel"]',
        newIntegration: 'button:contains("New integration")',
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
    awsSecurityHubForm: {
        nameInput: "form input[name='name']",
        awsAccountNumber: "form input[name='awsSecurityHub.accountId']",
        awsRegion: 'form .react-select__control',
        awsRegionListItems: '.react-select__menu-list > div',
        awsAccessKeyId: "form input[name='awsSecurityHub.credentials.accessKeyId']",
        awsSecretAccessKey: "form input[name='awsSecurityHub.credentials.secretAccessKey']",
    },
    syslogForm: {
        nameInput: "form input[name='name']",
        localFacility: 'form .react-select__control',
        localFacilityListItems: '.react-select__menu-list > div',
        receiverHost: "form input[name='syslog.tcpConfig.hostname']",
        receiverPort: 'form .react-numeric-input input',
        useTls: "form input[name='syslog.tcpConfig.useTls']",
        disableTlsValidation: "form input[name='syslog.tcpConfig.skipTlsVerify']",
    },
    modalHeader: '.ReactModal__Content header',
    formSaveButton: 'button[data-testid="save-integration"]',
    resultsSection: '[data-testid="results-message"]',
    labeledValue: '[data-testid="labeled-value"]',
    plugins: '#image-integrations a[data-testid="integration-tile"]',
    dialog: '.dialog',
    checkboxes: 'input',
    table,
    toast,
    tooltip,
};
