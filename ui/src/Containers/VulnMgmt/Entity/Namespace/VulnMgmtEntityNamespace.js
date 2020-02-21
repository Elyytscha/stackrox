import React from 'react';
import gql from 'graphql-tag';

import useCases from 'constants/useCaseTypes';
import { workflowEntityPropTypes, workflowEntityDefaultProps } from 'constants/entityPageProps';
import entityTypes from 'constants/entityTypes';
import { defaultCountKeyMap } from 'constants/workflowPages.constants';
import WorkflowEntityPage from 'Containers/Workflow/WorkflowEntityPage';
import { VULN_CVE_LIST_FRAGMENT } from 'Containers/VulnMgmt/VulnMgmt.fragments';
import VulnMgmtNamespaceOverview from './VulnMgmtNamespaceOverview';
import EntityList from '../../List/VulnMgmtList';
import {
    vulMgmtPolicyQuery,
    tryUpdateQueryWithVulMgmtPolicyClause,
    getScopeQuery
} from '../VulnMgmtPolicyQueryUtil';

const VulnMgmtNamespace = ({ entityId, entityListType, search, sort, page, entityContext }) => {
    const overviewQuery = gql`
        query getNamespace($id: ID!, $query: String, $policyQuery: String, $scopeQuery: String) {
            result: namespace(id: $id) {
                metadata {
                    priority
                    name
                    ${entityContext[entityTypes.CLUSTER] ? '' : 'clusterName clusterId'}
                    id
                    labels {
                        key
                        value
                    }
                }
                policyStatus(query: $policyQuery) {
                    status
                    failingPolicies {
                        id
                        name
                        description
                        policyStatus(query: $scopeQuery)
                        latestViolation(query: $scopeQuery)
                        severity
                        deploymentCount(query: $scopeQuery)
                        lifecycleStages
                        enforcementActions
                    }
                }
                policyCount(query: $policyQuery)
                vulnCount
                deploymentCount 
                imageCount 
                componentCount
                vulnerabilities: vulns {
                    ...cveFields
                }
            }
        }
        ${VULN_CVE_LIST_FRAGMENT}
    `;

    function getListQuery(listFieldName, fragmentName, fragment) {
        return gql`
        query getNamespace${entityListType}($id: ID!, $pagination: Pagination, $query: String, $policyQuery: String, $scopeQuery: String) {
            result: namespace(id: $id) {
                metadata {
                    id
                }
                ${defaultCountKeyMap[entityListType]}(query: $query)
                ${listFieldName}(query: $query, pagination: $pagination) { ...${fragmentName} }
                unusedVarSink(query: $policyQuery)
                unusedVarSink(query: $scopeQuery)
            }
        }
        ${fragment}
    `;
    }
    const newEntityContext = { ...entityContext, [entityTypes.NAMESPACE]: entityId };

    const queryOptions = {
        variables: {
            id: entityId,
            query: tryUpdateQueryWithVulMgmtPolicyClause(entityListType, search, newEntityContext),
            ...vulMgmtPolicyQuery,
            scopeQuery: getScopeQuery(entityContext)
        }
    };

    return (
        <WorkflowEntityPage
            entityId={entityId}
            entityType={entityTypes.NAMESPACE}
            entityListType={entityListType}
            useCase={useCases.VULN_MANAGEMENT}
            ListComponent={EntityList}
            OverviewComponent={VulnMgmtNamespaceOverview}
            overviewQuery={overviewQuery}
            getListQuery={getListQuery}
            search={search}
            sort={sort}
            page={page}
            queryOptions={queryOptions}
            entityContext={entityContext}
        />
    );
};

VulnMgmtNamespace.propTypes = workflowEntityPropTypes;
VulnMgmtNamespace.defaultProps = workflowEntityDefaultProps;

export default VulnMgmtNamespace;
