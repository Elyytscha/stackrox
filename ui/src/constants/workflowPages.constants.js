import entityTypes from 'constants/entityTypes';

export const WIDGET_PAGINATION_START_OFFSET = 0;

export const OVERVIEW_LIMIT = 5;

export const DASHBOARD_LIMIT = 8;

export const LIST_PAGE_SIZE = 50;

export const defaultCountKeyMap = {
    [entityTypes.COMPONENT]: 'componentCount',
    [entityTypes.CVE]: 'vulnCount',
    [entityTypes.DEPLOYMENT]: 'deploymentCount',
    [entityTypes.NAMESPACE]: 'namespaceCount',
    [entityTypes.IMAGE]: 'imageCount',
    [entityTypes.POLICY]: 'policyCount'
};
