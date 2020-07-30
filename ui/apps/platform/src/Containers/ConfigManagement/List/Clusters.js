import React from 'react';
import { gql } from '@apollo/client';
import pluralize from 'pluralize';

import {
    defaultHeaderClassName,
    defaultColumnClassName,
    nonSortableHeaderClassName,
} from 'Components/Table';
import LabelChip from 'Components/LabelChip';
import StatusChip from 'Components/StatusChip';
import { entityListPropTypes, entityListDefaultprops } from 'constants/entityPageProps';
import entityTypes from 'constants/entityTypes';
import { CLIENT_SIDE_SEARCH_OPTIONS as SEARCH_OPTIONS } from 'constants/searchOptions';
import { clusterSortFields } from 'constants/sortFields';
import queryService from 'utils/queryService';
import URLService from 'utils/URLService';
import List from './List';
import TableCellLink from './Link';
import filterByPolicyStatus from './utilities/filterByPolicyStatus';

const CLUSTERS_QUERY = gql`
    query clusters($query: String, $pagination: Pagination) {
        results: clusters(query: $query, pagination: $pagination) {
            id
            name
            serviceAccountCount
            k8sRoleCount
            subjectCount
            status {
                orchestratorMetadata {
                    version
                }
            }
            complianceControlCount(query: "Standard:CIS") {
                passingCount
                failingCount
                unknownCount
            }
            policyStatus {
                status
                failingPolicies {
                    id
                    name
                }
            }
        }
        count: clusterCount(query: $query)
    }
`;

export const defaultClusterSort = [
    {
        id: clusterSortFields.CLUSTER,
        desc: false,
    },
];

const buildTableColumns = (match, location) => {
    const tableColumns = [
        {
            Header: 'Id',
            headerClassName: 'hidden',
            className: 'hidden',
            accessor: 'id',
        },
        {
            Header: `Cluster`,
            headerClassName: `w-1/8 ${defaultHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            accessor: 'name',
            id: clusterSortFields.CLUSTER,
            sortField: clusterSortFields.CLUSTER,
        },
        {
            Header: `K8S Version`,
            headerClassName: `w-1/8 ${nonSortableHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            accessor: 'status.orchestratorMetadata.version',
            sortable: false,
        },
        {
            Header: `Policy Status`,
            headerClassName: `w-1/8 ${nonSortableHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => {
                const { policyStatus } = original;
                return <StatusChip status={policyStatus.status} asString={pdf} />;
            },
            id: 'status',
            accessor: (d) => d.policyStatus.status,
            sortable: false,
        },
        {
            Header: `CIS Controls`,
            headerClassName: `w-1/8 ${nonSortableHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            accessor: 'complianceControlCount',
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => {
                const { complianceControlCount } = original;
                const { passingCount, failingCount, unknownCount } = complianceControlCount;
                const totalCount = passingCount + failingCount + unknownCount;
                if (!totalCount) {
                    return <LabelChip text="No Controls" type="alert" />;
                }
                const url = URLService.getURL(match, location)
                    .push(original.id)
                    .push(entityTypes.CONTROL)
                    .url();
                return (
                    <TableCellLink
                        pdf={pdf}
                        url={url}
                        text={`${totalCount} ${pluralize('Controls', totalCount)}`}
                    />
                );
            },
            sortable: false,
        },
        {
            Header: `Users & Groups`,
            headerClassName: `w-1/8 ${nonSortableHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => {
                const { subjectCount } = original;
                if (!subjectCount) {
                    return <LabelChip text="No Users & Groups" type="alert" />;
                }
                const url = URLService.getURL(match, location)
                    .push(original.id)
                    .push(entityTypes.SUBJECT)
                    .url();
                return (
                    <TableCellLink
                        pdf={pdf}
                        url={url}
                        text={`${subjectCount} ${pluralize('Users & Groups', subjectCount)}`}
                    />
                );
            },
            id: 'subjectCount',
            accessor: (d) => d.subjectCount,
            sortable: false,
        },
        {
            Header: `Service Accounts`,
            headerClassName: `w-1/8 ${nonSortableHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => {
                const { serviceAccountCount } = original;
                if (!serviceAccountCount) {
                    return <LabelChip text="No Service Accounts" type="alert" />;
                }
                const url = URLService.getURL(match, location)
                    .push(original.id)
                    .push(entityTypes.SERVICE_ACCOUNT)
                    .url();
                return (
                    <TableCellLink
                        pdf={pdf}
                        url={url}
                        text={`${serviceAccountCount} ${pluralize(
                            'Service Accounts',
                            serviceAccountCount
                        )}`}
                    />
                );
            },
            id: 'serviceAccountCount',
            accessor: (d) => d.serviceAccountCount,
            sortable: false,
        },
        {
            Header: `Roles`,
            headerClassName: `w-1/8 ${nonSortableHeaderClassName}`,
            className: `w-1/8 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => {
                const { k8sRoleCount } = original;
                if (!k8sRoleCount) return <LabelChip text="No Roles" type="alert" />;
                const url = URLService.getURL(match, location)
                    .push(original.id)
                    .push(entityTypes.ROLE)
                    .url();
                return (
                    <TableCellLink
                        pdf={pdf}
                        url={url}
                        text={`${k8sRoleCount} ${pluralize('Roles', k8sRoleCount)}`}
                    />
                );
            },
            id: 'k8sRoleCount',
            accessor: (d) => d.k8sRoleCount,
            sortable: false,
        },
    ];
    return tableColumns;
};

const createTableRows = (data) => data.results;

const Clusters = ({ match, location, className, selectedRowId, onRowClick, query, data }) => {
    const autoFocusSearchInput = !selectedRowId;
    const tableColumns = buildTableColumns(match, location);
    const { [SEARCH_OPTIONS.POLICY_STATUS.CATEGORY]: policyStatus, ...restQuery } = query || {};
    const queryObject = { ...restQuery };
    const queryText = queryService.objectToWhereClause(queryObject);
    const variables = queryText ? { query: queryText } : null;

    function createTableRowsFilteredByPolicyStatus(items) {
        const tableRows = createTableRows(items);
        const filteredTableRows = filterByPolicyStatus(tableRows, policyStatus);
        return filteredTableRows;
    }

    return (
        <List
            className={className}
            query={CLUSTERS_QUERY}
            variables={variables}
            entityType={entityTypes.CLUSTER}
            tableColumns={tableColumns}
            createTableRows={createTableRowsFilteredByPolicyStatus}
            onRowClick={onRowClick}
            selectedRowId={selectedRowId}
            idAttribute="id"
            defaultSorted={defaultClusterSort}
            defaultSearchOptions={[SEARCH_OPTIONS.POLICY_STATUS.CATEGORY]}
            data={filterByPolicyStatus(data, policyStatus)}
            autoFocusSearchInput={autoFocusSearchInput}
        />
    );
};

Clusters.propTypes = entityListPropTypes;
Clusters.defaultProps = entityListDefaultprops;

export default Clusters;
