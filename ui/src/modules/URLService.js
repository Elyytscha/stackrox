import qs from 'qs';
import pageTypes from 'constants/pageTypes';
import { resourceTypes } from 'constants/entityTypes';
import contextTypes from 'constants/contextTypes';
import { generatePath } from 'react-router-dom';
import { nestedCompliancePaths, resourceTypesToUrl, riskPath } from '../routePaths';

function isResource(type) {
    return Object.values(resourceTypes).includes(type);
}

function getEntityTypeFromMatch(match) {
    if (!match || !match.params || !match.params.entityType) return null;

    const { entityType } = match.params;

    // Handle url to resourceType mapping for resources
    const entityEntry = Object.entries(resourceTypesToUrl).find(
        entry => entry[1] === match.params.entityType
    );

    return entityEntry ? entityEntry[0] : entityType;
}

function getPath(context, pageType, urlParams) {
    const isResourceType = urlParams.entityType ? isResource(urlParams.entityType) : false;
    const pathMap = {
        [contextTypes.COMPLIANCE]: {
            [pageTypes.DASHBOARD]: nestedCompliancePaths.DASHBOARD,
            [pageTypes.ENTITY]: isResourceType
                ? nestedCompliancePaths.RESOURCE
                : nestedCompliancePaths.CONTROL,
            [pageTypes.LIST]: nestedCompliancePaths.LIST
        },
        [contextTypes.RISK]: {
            [pageTypes.ENTITY]: riskPath,
            [pageTypes.LIST]: '/main/risk'
        }
    };

    const contextData = pathMap[context];
    if (!contextData) return null;

    const path = contextData[pageType];
    if (!path) return null;

    const params = { ...urlParams };
    if (isResourceType) {
        params.entityType = resourceTypesToUrl[urlParams.entityType];
    }
    return generatePath(path, params);
}

function getContext(match) {
    if (match.url.includes('/compliance')) return contextTypes.COMPLIANCE;
    if (match.url.includes('/risk')) return contextTypes.RISK;
    return null;
}

function getPageType(match) {
    if (match.params.entityId) return pageTypes.ENTITY;
    if (match.params.entityType) return pageTypes.LIST;
    return pageTypes.DASHBOARD;
}

function getParams(match, location) {
    const newParams = { ...match.params };
    newParams.entityType = getEntityTypeFromMatch(match);

    return {
        ...newParams,
        context: getContext(match),
        pageType: getPageType(match),
        query: qs.parse(location.search, { ignoreQueryPrefix: true })
    };
}

// this just lowercases stuff, use queryservice
function keysToLowerCase(query) {
    if (!query) return null;

    return Object.entries(query).reduce((acc, entry) => {
        const key = entry[0];
        // eslint-disable-next-line
        acc[key] = entry[1];
        return acc;
    }, {});
}

function getLinkTo(context, pageType, params) {
    const { query, ...urlParams } = params;
    const pathname = getPath(context, pageType, urlParams);
    const search = query ? qs.stringify(keysToLowerCase(query), { addQueryPrefix: true }) : null;

    return {
        pathname,
        search,
        url: pathname + search
    };
}

export default {
    getParams,
    getLinkTo
};
