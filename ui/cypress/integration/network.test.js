import { url as networkUrl, selectors as networkPageSelectors } from '../constants/NetworkPage';
import * as api from '../constants/apiEndpoints';
import withAuth from '../helpers/basicAuth';

const uploadFile = (fileName, selector) => {
    cy.get(selector).then(subject => {
        cy.fixture(fileName).then(content => {
            const el = subject[0];
            const testFile = new File([content], fileName);
            const dataTransfer = new DataTransfer();

            dataTransfer.items.add(testFile);
            el.files = dataTransfer.files;
        });
    });
};

describe('Network page', () => {
    withAuth();

    beforeEach(() => {
        cy.server();
        cy.fixture('network/networkGraph.json').as('networkGraphJson');
        cy.route('GET', api.network.networkGraph, '@networkGraphJson').as('networkGraph');
        cy.visit(networkUrl);
        cy.wait('@networkGraph');
    });

    it('should have selected item in nav bar', () => {
        cy.get(networkPageSelectors.network).click();
        cy.get(networkPageSelectors.network).should('have.class', 'bg-primary-700');
    });

    it('should display a legend', () => {
        cy.get(networkPageSelectors.legend.deployments)
            .eq(0)
            .children()
            .should('have.class', 'icon-node');
        cy.get(networkPageSelectors.legend.deployments)
            .eq(1)
            .children()
            .children()
            .should('have.class', 'icon-potential');
        cy.get(networkPageSelectors.legend.deployments)
            .eq(2)
            .children()
            .should('have.class', 'icon-node');

        cy.get(networkPageSelectors.legend.namespaces)
            .eq(0)
            .children()
            .should('have.attr', 'alt', 'namespace');
        cy.get(networkPageSelectors.legend.namespaces)
            .eq(1)
            .children()
            .should('have.attr', 'alt', 'namespace-allowed-connection');
        cy.get(networkPageSelectors.legend.namespaces)
            .eq(2)
            .children()
            .should('have.attr', 'alt', 'namespace-connection');

        cy.get(networkPageSelectors.legend.connections)
            .eq(0)
            .children()
            .should('have.attr', 'alt', 'active-connection');
        cy.get(networkPageSelectors.legend.connections)
            .eq(1)
            .children()
            .should('have.attr', 'alt', 'allowed-connection');
        cy.get(networkPageSelectors.legend.connections)
            .eq(2)
            .children()
            .should('have.class', 'icon-ingress-egress');
    });

    it('should handle toggle click on simulator network policy button', () => {
        cy.get(networkPageSelectors.buttons.simulatorButtonOff).click();
        cy.get(networkPageSelectors.buttons.viewActiveYamlButton).should('be.visible');
        cy.get(networkPageSelectors.panels.creatorPanel).should('be.visible');
        cy.get(networkPageSelectors.buttons.simulatorButtonOn).click();
        cy.get(networkPageSelectors.panels.creatorPanel).should('not.be.visible');
    });

    it('should display error messages when uploaded wrong yaml', () => {
        cy.get(networkPageSelectors.buttons.simulatorButtonOff).click();
        uploadFile('network/policywithoutnamespace.yaml', 'input[type="file"]');
        cy.get(networkPageSelectors.simulatorSuccessMessage).should('not.be.visible');
    });

    it('should display success messages when uploaded right yaml', () => {
        cy.get(networkPageSelectors.buttons.simulatorButtonOff).click();
        uploadFile('network/policywithnamespace.yaml', 'input[type="file"]');
        cy.get(networkPageSelectors.simulatorSuccessMessage).should('be.visible');
    });
});
