export const url = '/main/policies';

export const selectors = {
    configure: 'nav.left-navigation li:contains("Configure") a',
    navLink: '.navigation-panel li:contains("System Policies") a',
    newPolicyButton: 'button:contains("New")',
    editPolicyButton: 'button:contains("Edit")',
    savePolicyButton: 'button:contains("Save")',
    nextButton: '.btn:contains("Next")',
    prevButton: 'button:contains("Previous")',
    cancelButton: 'button[data-test-id="cancel"]',
    reassessAllButton: 'button:contains("Reassess")',
    actionMenuBtn: 'button:contains("Actions")',
    actionMenu: '[data-test-id="menu-list"]',
    checkboxes: 'input:checkbox',
    policies: {
        scanImage: 'div.rt-tr:contains("90-Day")',
        addCapabilities: '.rt-tr:contains("CAP_SYS_ADMIN capability added")',
        disabledPolicyImage: 'div.rt-tr.data-test-disabled:first'
    },
    form: {
        nameInput: 'form input[name=name]',
        enableField: 'form input[name=disabled]',
        required: 'form span[data-test-id="required"]',
        select: 'form select',
        selectValue: 'form .react-select__multi-value__label'
    },
    configurationField: {
        select: '#policyConfigurationSelect',
        selectArrow: '#policyConfigurationSelect .react-select__dropdown-indicator',
        options: '#policyConfigurationSelect .react-select__option',
        numericInput: '[data-test-id="policyConfiguration"] .react-numeric-input input'
    },
    imageRegistry: {
        input: 'input[name="fields.imageName.registry"]',
        deleteButton: 'div:contains("Image Registry")+ div.flex>div.flex>button',
        value: '[data-test-id="imageName"] div.flex'
    },
    scanAgeDays: {
        input:
            'div:contains("Days since image was last scanned") + div.flex>.react-numeric-input>input',
        deleteButton:
            'div:contains("Days since image was last scanned") + div.flex>div.flex>button',
        value: '[data-test-id="scanAgeDays"] div.flex'
    },
    categoriesField: {
        input: 'div:contains("Categories") + div.flex .react-select__input > input',
        valueContainer: 'div:contains("Categories") + div.flex .react-select__value-container'
    },
    policyPreview: {
        message: '.warn-message',
        alertPreview: {
            table: '.alert-preview table'
        }
    },
    policyDetailsPanel: {
        idValueDiv: 'div.text-base-600:contains("ID:") + div',
        enabledValueDiv: 'div.text-base-600:contains("Enabled") + div'
    },
    searchInput: '.react-select__input > input',
    sidePanel: '[data-test-id="side-panel"]',
    sidePanelHeader: '[data-test-id="side-panel-header"]',
    tableFirstRow: 'div.rt-tbody > div.rt-tr-group:first > .rt-tr.-odd',
    tableFirstRowName:
        'div.rt-tbody > div.rt-tr-group:first > .rt-tr.-odd [data-test-id=policy-name]',
    hoverActionButtons: '.rt-tr-actions svg',
    tableContainer: '[data-test-id="policies-table-container"]',
    enableDisableIcon: '[data-test-id="enable-disable-icon"]',
    enabledIconColor: 'bg-success-500',
    enforcement: {
        buildTile: '[data-test-id="policy-enforcement-build-tile"]',
        deployTile: '[data-test-id="policy-enforcement-deploy-tile"]',
        onOffToggle: '[data-test-id="policy-enforcement-on-off"]'
    }
};

export const text = {
    policyLatestTagName: 'Latest tag',
    policyPreview: {
        message:
            'This policy is not currently enabled. If enabled, the policy would generate violations for the following deployments on your system.'
    }
};
