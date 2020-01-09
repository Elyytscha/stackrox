import React, { useState, useContext } from 'react';
import PropTypes from 'prop-types';
import gql from 'graphql-tag';
import pluralize from 'pluralize';
import { useQuery } from 'react-apollo';
import sortBy from 'lodash/sortBy';

import queryService from 'modules/queryService';
import workflowStateContext from 'Containers/workflowStateContext';
import Loader from 'Components/Loader';
import NoResultsMessage from 'Components/NoResultsMessage';
import ViewAllButton from 'Components/ViewAllButton';
import Widget from 'Components/Widget';
import Scatterplot from 'Components/visuals/Scatterplot';
import HoverHintListItem from 'Components/visuals/HoverHintListItem';
import TextSelect from 'Components/TextSelect';
import entityTypes from 'constants/entityTypes';
import entityLabels from 'messages/entity';
import { severityLabels } from 'messages/common';
import isGQLLoading from 'utils/gqlLoading';
import {
    severityColorMap,
    severityTextColorMap,
    severityColorLegend
} from 'constants/severityColors';
import { getSeverityByCvss } from 'utils/vulnerabilityUtils';
import { entitySortFieldsMap } from 'constants/sortFields';

const TopRiskyEntitiesByVulnerabilities = ({
    entityContext,
    defaultSelection,
    riskEntityTypes,
    cveFilter,
    small
}) => {
    const workflowState = useContext(workflowStateContext);
    // Entity Type selection
    const [selectedEntityType, setEntityType] = useState(defaultSelection);
    const entityOptions = riskEntityTypes.map(entityType => ({
        label: `top risky ${pluralize(entityLabels[entityType])} by CVE count & CVSS score`,
        value: entityType
    }));
    function onChange(datum) {
        setEntityType(datum);
    }

    // View all button
    const viewAllUrl = workflowState
        .pushList(selectedEntityType)
        .setSort([
            // @TODO to uncomment when Priority sort is available across images/deployments/namespaces/clusters
            // {
            //     id: entitySortFieldsMap[selectedEntityType].PRIORITY,
            //     desc: false
            // },
            {
                id: entitySortFieldsMap[selectedEntityType].NAME,
                desc: false
            }
        ])
        .toUrl();

    const titleComponents = (
        <TextSelect value={selectedEntityType} onChange={onChange} options={entityOptions} />
    );
    const viewAll = <ViewAllButton url={viewAllUrl} />;

    // Data Queries
    const VULN_FRAGMENT = gql`
        fragment vulnFields on EmbeddedVulnerability {
            cve
            cvss
            isFixable
            severity
        }
    `;
    const DEPLOYMENT_QUERY = gql`
        query topRiskyDeployments($query: String, $vulnQuery: String) {
            results: deployments(query: $query) {
                id
                name
                clusterName
                namespaceName: namespace
                priority
                vulnCounter {
                    all {
                        total
                        fixable
                    }
                }
                vulns(query: $vulnQuery) {
                    ...vulnFields
                }
            }
        }
        ${VULN_FRAGMENT}
    `;

    const CLUSTER_QUERY = gql`
        query topRiskyClusters($query: String, $vulnQuery: String) {
            results: clusters(query: $query) {
                id
                name
                priority
                vulnCounter {
                    all {
                        total
                        fixable
                    }
                }
                vulns(query: $vulnQuery) {
                    ...vulnFields
                }
            }
        }
        ${VULN_FRAGMENT}
    `;

    const NAMESPACE_QUERY = gql`
        query topRiskyNamespaces($query: String, $vulnQuery: String) {
            results: namespaces(query: $query) {
                metadata {
                    clusterName
                    name
                    id
                    priority
                }
                vulnCounter {
                    all {
                        total
                        fixable
                    }
                }
                vulns(query: $vulnQuery) {
                    ...vulnFields
                }
            }
        }
        ${VULN_FRAGMENT}
    `;

    const IMAGE_QUERY = gql`
        query topRiskyImages($query: String, $vulnQuery: String) {
            results: images(query: $query) {
                id
                name {
                    fullName
                }
                priority
                vulnCounter {
                    all {
                        total
                        fixable
                    }
                }
                vulns(query: $vulnQuery) {
                    ...vulnFields
                }
            }
        }
        ${VULN_FRAGMENT}
    `;

    const COMPONENT_QUERY = gql`
        query topRiskyComponents($query: String, $vulnQuery: String) {
            results: components(query: $query) {
                id
                name
                priority
                vulnCounter {
                    all {
                        total
                        fixable
                    }
                }
                vulns(query: $vulnQuery) {
                    ...vulnFields
                }
            }
        }
        ${VULN_FRAGMENT}
    `;

    const queryMap = {
        [entityTypes.DEPLOYMENT]: DEPLOYMENT_QUERY,
        [entityTypes.NAMESPACE]: NAMESPACE_QUERY,
        [entityTypes.CLUSTER]: CLUSTER_QUERY,
        [entityTypes.COMPONENT]: COMPONENT_QUERY,
        [entityTypes.IMAGE]: IMAGE_QUERY
    };
    const query = queryMap[selectedEntityType];

    function getAverageSeverity(vulns) {
        if (vulns.length === 0) return 0;

        // 1. sort the vulns in reverse CVSS order
        const sortedVulns = sortBy(vulns, vuln => {
            return vuln.cvss;
        }).reverse();
        const topVulns = sortedVulns.slice(0, 100);

        // 2. grab up to the first 5 vulns (the ones with the highest CVSS)
        const total = topVulns.reduce((acc, curr) => {
            return acc + parseFloat(curr.cvss);
        }, 0);

        // 3. Take the average of those top 5 (or total, if less than 5)
        const avgScore = total / topVulns.length;

        return avgScore.toFixed(1);
    }

    function getHint(datum, filter) {
        let subtitle = '';
        if (selectedEntityType === entityTypes.DEPLOYMENT) {
            subtitle = `${datum.clusterName} / ${datum.namespaceName}`;
        } else if (selectedEntityType === entityTypes.NAMESPACE) {
            subtitle = `${datum.metadata && datum.metadata.clusterName}`;
        }
        const riskPriority =
            selectedEntityType === entityTypes.NAMESPACE
                ? datum.metadata && datum.metadata.priority
                : datum.priority;
        const severityKey = getSeverityByCvss(datum.avgSeverity);
        const severityTextColor = severityTextColorMap[severityKey];
        const severityText = severityLabels[severityKey];

        let cveCountText = filter !== 'Fixable' ? `${datum.vulnCounter.all.total} total / ` : '';
        cveCountText += `${datum.vulnCounter.all.fixable} fixable`;

        return {
            title:
                (datum.name && datum.name.fullName) ||
                datum.name ||
                (datum.metadata && datum.metadata.name),
            body: (
                <ul className="flex-1 list-reset border-base-300 overflow-hidden">
                    <HoverHintListItem
                        key="severity"
                        label="Severity"
                        value={<span style={{ color: severityTextColor }}>{severityText}</span>}
                    />
                    <HoverHintListItem
                        key="riskPriority"
                        label="Risk Priority"
                        value={riskPriority}
                    />
                    <HoverHintListItem
                        key="weightedCvss"
                        label="Weighted CVSS"
                        value={datum.avgSeverity}
                    />
                    <HoverHintListItem key="cves" label="CVEs" value={cveCountText} />
                </ul>
            ),
            subtitle
        };
    }
    function processData(data) {
        if (!data || !data.results) return [];

        const results = data.results
            .filter(datum => datum.vulns && datum.vulns.length > 0)
            .map(result => {
                const entityId = result.id || result.metadata.id;
                const vulnCount = result.vulns.length;
                const url = workflowState.pushRelatedEntity(selectedEntityType, entityId).toUrl();
                const avgSeverity = getAverageSeverity(result.vulns);
                const color = severityColorMap[getSeverityByCvss(avgSeverity)];

                return {
                    x: vulnCount,
                    y: +avgSeverity,
                    color,
                    hint: getHint({ ...result, avgSeverity }, cveFilter),
                    url
                };
            })
            .sort((a, b) => {
                return a.vulnCount - b.vulnCount;
            });

        return results;
    }
    let results = [];

    const vulnQuery = cveFilter === 'Fixable' ? { 'Fixed By': 'r/.*' } : '';
    const variables = {
        query: queryService.entityContextToQueryString(entityContext),
        vulnQuery: queryService.objectToWhereClause(vulnQuery)
    };
    const { data, loading } = useQuery(query, { variables });

    let content = <Loader />;

    if (!isGQLLoading(loading, data)) {
        results = processData(data);
        if (!results || results.length === 0) {
            content = (
                <NoResultsMessage
                    message={`No ${pluralize(
                        selectedEntityType.toLowerCase()
                    )} with vulnerabilities found`}
                    className="p-6"
                    icon="info"
                />
            );
        } else {
            content = (
                <Scatterplot
                    data={results}
                    xMultiple={10}
                    yMultiple={10}
                    yAxisTitle="Average CVSS Score"
                    xAxisTitle="Critical Vulnerabilities & Exposures"
                    legendData={!small ? severityColorLegend : []}
                />
            );
        }
    }

    return (
        <Widget
            className="h-full pdf-page pdf-stretch"
            titleComponents={titleComponents}
            headerComponents={viewAll}
            bodyClassName="pr-2"
        >
            {content}
        </Widget>
    );
};

TopRiskyEntitiesByVulnerabilities.propTypes = {
    entityContext: PropTypes.shape({}),
    defaultSelection: PropTypes.string.isRequired,
    riskEntityTypes: PropTypes.arrayOf(PropTypes.string),
    cveFilter: PropTypes.string,
    small: PropTypes.bool
};

TopRiskyEntitiesByVulnerabilities.defaultProps = {
    entityContext: {},
    riskEntityTypes: [
        entityTypes.DEPLOYMENT,
        entityTypes.NAMESPACE,
        entityTypes.IMAGE,
        entityTypes.CLUSTER
    ],
    cveFilter: 'All',
    small: false
};

export default TopRiskyEntitiesByVulnerabilities;
