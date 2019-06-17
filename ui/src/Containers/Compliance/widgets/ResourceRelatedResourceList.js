import React from 'react';
import PropTypes from 'prop-types';
import LinkListWidget from 'Components/LinkListWidget';
import URLService from 'modules/URLService';
import pluralize from 'pluralize';
import entityTypes from 'constants/entityTypes';
import contextTypes from 'constants/contextTypes';
import { resourceLabels } from 'messages/common';
import { RELATED_SECRETS, RELATED_DEPLOYMENTS, ALL_NAMESPACES } from 'queries/namespace';
import queryService from 'modules/queryService';
import { NODES_BY_CLUSTER } from 'queries/node';
import ReactRouterPropTypes from 'react-router-prop-types';
import { withRouter, Link } from 'react-router-dom';

const queryMap = {
    [entityTypes.NAMESPACE]: ALL_NAMESPACES,
    [entityTypes.NODE]: NODES_BY_CLUSTER,
    [entityTypes.SECRET]: RELATED_SECRETS,
    [entityTypes.DEPLOYMENT]: RELATED_DEPLOYMENTS
};

function getPageContext(entityType) {
    switch (entityType) {
        case entityTypes.DEPLOYMENT:
            return contextTypes.RISK;
        case entityTypes.SECRET:
            return contextTypes.SECRET;
        default:
            return contextTypes.COMPLIANCE;
    }
}

const ResourceRelatedEntitiesList = ({
    match,
    location,
    listEntityType,
    pageEntityType,
    pageEntity,
    clusterName,
    className,
    limit
}) => {
    const linkContext = getPageContext(listEntityType);
    const resourceLabel = resourceLabels[listEntityType];

    function processData(data) {
        if (!data || !data.results) return [];

        let items = data.results;
        if (listEntityType === entityTypes.NAMESPACE) {
            items = items
                .map(item => ({
                    ...item.metadata,
                    name: `${item.metadata.clusterName}/${item.metadata.name}`
                }))
                .filter(item => item.clusterName === pageEntity.name);
        }
        if (listEntityType === entityTypes.NODE) {
            items = data.results.nodes;
        }

        return items.map(item => ({
            label: item.name,
            link: URLService.getURL(match, location)
                .base(listEntityType, item.id)
                .url()
        }));
    }

    const viewAllLink =
        pageEntity && pageEntity.id ? (
            <Link
                to={URLService.getURL(match, location)
                    .base(listEntityType, null, linkContext)
                    .query({
                        [pageEntityType]: pageEntity.name,
                        [entityTypes.CLUSTER]: clusterName
                    })
                    .url()}
                className="no-underline"
            >
                <button className="btn-sm btn-base btn-sm" type="button">
                    View All
                </button>
            </Link>
        ) : null;

    function getVariables() {
        if (listEntityType === entityTypes.NAMESPACE) {
            return null;
        }

        const variables = {
            query: queryService.objectToWhereClause({
                [pageEntityType]: pageEntity.name,
                [entityTypes.CLUSTER]: clusterName
            })
        };

        if (listEntityType === entityTypes.NODE) {
            variables.id = pageEntity.id;
        }

        return variables;
    }

    function getHeadline(list) {
        if (!list) return `Related ${pluralize(resourceLabel)}`;
        return `${list.length} Related ${pluralize(resourceLabel)}`;
    }

    return (
        <LinkListWidget
            query={queryMap[listEntityType]}
            variables={getVariables()}
            processData={processData}
            getHeadline={getHeadline}
            headerComponents={viewAllLink}
            className={className}
            id="related-resource-list"
            limit={limit}
        />
    );
};

ResourceRelatedEntitiesList.propTypes = {
    match: ReactRouterPropTypes.match.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    listEntityType: PropTypes.string.isRequired,
    pageEntityType: PropTypes.string.isRequired,
    className: PropTypes.string,
    pageEntity: PropTypes.shape({
        id: PropTypes.string,
        name: PropTypes.string
    }),
    clusterName: PropTypes.string,
    limit: PropTypes.number
};

ResourceRelatedEntitiesList.defaultProps = {
    pageEntity: null,
    className: null,
    limit: 20,
    clusterName: null
};

export default withRouter(ResourceRelatedEntitiesList);
