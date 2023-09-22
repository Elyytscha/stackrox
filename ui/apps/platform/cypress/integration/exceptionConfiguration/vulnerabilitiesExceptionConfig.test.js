import withAuth from '../../helpers/basicAuth';
import { hasFeatureFlag } from '../../helpers/features';
import {
    resetExceptionConfig,
    visitExceptionConfig,
    visitExceptionConfigWithPermissions,
    vulnerabilitiesConfigAlias,
} from './ExceptionConfig.helpers';
import { vulnerabilitiesConfigSelectors as selectors } from './ExceptionConfig.selectors';

describe('Vulnerabilities Exception Configuration', () => {
    withAuth();

    before(function () {
        if (!hasFeatureFlag('ROX_VULN_MGMT_UNIFIED_CVE_DEFERRAL')) {
            this.skip();
        }
    });

    afterEach(() => {
        resetExceptionConfig();
    });

    it('should correctly handle RBAC for the vulnerability exception config', () => {
        cy.fixture('auth/mypermissionsMinimalAccess.json').then(({ resourceToAccess }) => {
            // User with no access
            visitExceptionConfigWithPermissions('vulnerabilities', {
                ...resourceToAccess,
                Administration: 'NO_ACCESS',
            });

            cy.get('h1:contains("Cannot find the page")');

            // User with read-only access
            visitExceptionConfigWithPermissions('vulnerabilities', {
                ...resourceToAccess,
                Administration: 'READ_ACCESS',
            });

            cy.get(`h1:contains("Exception configuration")`);
            cy.get(selectors.saveButton).should('not.exist');

            for (let i = 0; i < 4; i += 1) {
                cy.get(selectors.dayOptionInput(i)).should('be.disabled');
                cy.get(selectors.dayOptionEnabledSwitch(i)).should('be.disabled');
            }

            // TODO cy.get(selectors.indefiniteOptionEnabledSwitch).should('be.disabled');
            cy.get(selectors.whenAllCveFixableSwitch).should('be.disabled');
            cy.get(selectors.whenAnyCveFixableSwitch).should('be.disabled');
            cy.get(selectors.customDateSwitch).should('be.disabled');
        });
    });

    it('should load the default config and allow modification', () => {
        visitExceptionConfig('vulnerabilities');

        cy.get(selectors.dayOptionEnabledSwitch(0)).check({ force: true });
        cy.get(selectors.dayOptionInput(0)).type('{selectall}20');
        cy.get(selectors.dayOptionEnabledSwitch(1)).check({ force: true });
        cy.get(selectors.dayOptionInput(1)).type('{selectall}40');
        cy.get(selectors.dayOptionEnabledSwitch(2)).check({ force: true });
        cy.get(selectors.dayOptionInput(2)).type('{selectall}60');

        cy.get(selectors.dayOptionEnabledSwitch(3)).uncheck({ force: true });
        cy.get(selectors.dayOptionInput(3)).should('be.disabled');

        cy.get(selectors.whenAllCveFixableSwitch).check({ force: true });
        cy.get(selectors.whenAnyCveFixableSwitch).uncheck({ force: true });
        cy.get(selectors.customDateSwitch).check({ force: true });

        cy.get(selectors.saveButton).click();
        cy.get('.pf-c-alert:contains("The configuration was updated successfully")');

        // Refresh the page to make sure options are persisted
        visitExceptionConfig('vulnerabilities');

        cy.get(selectors.dayOptionInput(0)).should('have.value', '20');
        cy.get(selectors.dayOptionEnabledSwitch(0)).should('be.checked');
        cy.get(selectors.dayOptionInput(1)).should('have.value', '40');
        cy.get(selectors.dayOptionEnabledSwitch(1)).should('be.checked');
        cy.get(selectors.dayOptionInput(2)).should('have.value', '60');
        cy.get(selectors.dayOptionEnabledSwitch(2)).should('be.checked');

        cy.get(selectors.dayOptionInput(3)).should('not.be.checked');

        cy.get(selectors.whenAllCveFixableSwitch).should('be.checked');
        cy.get(selectors.whenAnyCveFixableSwitch).should('not.be.checked');
        cy.get(selectors.customDateSwitch).should('be.checked');
    });

    it('should handle null exception configuration from the server and display a valid UI', () => {
        visitExceptionConfig('vulnerabilities', {
            [vulnerabilitiesConfigAlias]: {
                body: { config: null },
            },
        });

        // Check that all form elements are present and options are disabled (unchecked)
        cy.get(selectors.saveButton).should('be.disabled');

        for (let i = 0; i < 4; i += 1) {
            cy.get(selectors.dayOptionInput(i)).should('be.disabled');
            cy.get(selectors.dayOptionEnabledSwitch(i)).should('not.be.checked');
        }

        // TODO cy.get(selectors.indefiniteOptionEnabledSwitch).should('not.be.checked');
        cy.get(selectors.whenAllCveFixableSwitch).should('not.be.checked');
        cy.get(selectors.whenAnyCveFixableSwitch).should('not.be.checked');
        cy.get(selectors.customDateSwitch).should('not.be.checked');
    });

    it.skip('should reflect an updated exception config in the Workload CVE exception flow', () => {
        // TODO Implement once the exception flow is implemented
    });
});
