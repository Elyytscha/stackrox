/* eslint-disable react/jsx-no-bind */
import React, { useContext, useState } from 'react';
import PropTypes from 'prop-types';
import gql from 'graphql-tag';
import * as Icon from 'react-feather';
import { connect } from 'react-redux';
import { withRouter } from 'react-router-dom';

import { defaultHeaderClassName, defaultColumnClassName } from 'Components/Table';
import RowActionButton from 'Components/RowActionButton';
import DateTimeField from 'Components/DateTimeField';
import LabelChip from 'Components/LabelChip';
import TableCountLink from 'Components/workflow/TableCountLink';
import TopCvssLabel from 'Components/TopCvssLabel';
import PanelButton from 'Components/PanelButton';
import WorkflowListPage from 'Containers/Workflow/WorkflowListPage';
import workflowStateContext from 'Containers/workflowStateContext';
import entityTypes from 'constants/entityTypes';
import { LIST_PAGE_SIZE } from 'constants/workflowPages.constants';
import queryService from 'modules/queryService';
import { workflowListPropTypes, workflowListDefaultProps } from 'constants/entityPageProps';
import { actions as notificationActions } from 'reducers/notifications';
import { updateCveSuppressedState } from 'services/VulnerabilitiesService';
import removeEntityContextColumns from 'utils/tableUtils';
import { doesSearchContain } from 'utils/searchUtils';
import { truncate } from 'utils/textUtils';
import { cveSortFields } from 'constants/sortFields';

import { VULN_CVE_LIST_FRAGMENT } from 'Containers/VulnMgmt/VulnMgmt.fragments';

import CveBulkActionDialogue from './CveBulkActionDialogue';

export const defaultCveSort = [
    {
        id: cveSortFields.CVSS_SCORE,
        desc: true
    },
    {
        id: cveSortFields.CVE,
        desc: false
    }
];

export function getCveTableColumns(workflowState) {
    // to determine whether to show the counts as links in the table when not in pure CVE state
    const inFindingsSection = workflowState.getCurrentEntity().entityType !== entityTypes.CVE;
    const tableColumns = [
        {
            expander: true,
            show: false
        },
        {
            Header: 'id',
            headerClassName: 'hidden',
            className: 'hidden',
            accessor: 'id'
        },
        {
            Header: `CVE`,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            accessor: 'cve',
            sortField: cveSortFields.CVE
        },
        {
            Header: `Fixable`,
            headerClassName: `w-20 text-center ${defaultHeaderClassName}`,
            className: `w-20 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => {
                const fixableFlag = original.isFixable ? (
                    <LabelChip text="Fixable" type="success" size="large" />
                ) : (
                    'No'
                );
                return <div className="mx-auto">{fixableFlag}</div>;
            },
            accessor: 'isFixable',
            id: 'isFixable',
            sortField: cveSortFields.FIXABLE
        },
        {
            Header: `CVSS Score`,
            headerClassName: `w-1/10 text-center ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original }) => {
                const { cvss, scoreVersion } = original;
                return <TopCvssLabel cvss={cvss} version={scoreVersion} />;
            },
            accessor: 'cvss',
            id: 'cvss',
            sortField: cveSortFields.CVSS_SCORE
        },
        {
            Header: `Env. Impact`,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            Cell: ({ original }) => {
                const { envImpact } = original;
                // eslint-disable-next-line eqeqeq
                return envImpact == Number(envImpact)
                    ? `${(envImpact * 100).toFixed(0)}% affected`
                    : '-';
            },
            accessor: 'envImpact',
            sortField: cveSortFields.ENV_IMPACT,
            sortable: false
        },
        {
            Header: `Impact Score`,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            Cell: ({ original }) => {
                const { impactScore } = original;
                // eslint-disable-next-line eqeqeq
                return impactScore == Number(impactScore) ? impactScore.toFixed(1) : '-';
            },
            accessor: 'impactScore',
            sortField: cveSortFields.IMPACT_SCORE
        },
        {
            Header: `Deployments`,
            entityType: entityTypes.DEPLOYMENT,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => (
                <TableCountLink
                    entityType={entityTypes.DEPLOYMENT}
                    count={original.deploymentCount}
                    textOnly={inFindingsSection || pdf}
                    selectedRowId={original.cve}
                />
            ),
            accessor: 'deploymentCount',
            id: 'deploymentCount',
            sortField: cveSortFields.DEPLOYMENTS,
            sortable: false
        },
        {
            Header: `Images`,
            entityType: entityTypes.IMAGE,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => (
                <TableCountLink
                    entityType={entityTypes.IMAGE}
                    count={original.imageCount}
                    textOnly={inFindingsSection || pdf}
                    selectedRowId={original.cve}
                />
            ),
            accessor: 'imageCount',
            id: 'imageCount',
            sortField: cveSortFields.IMAGES,
            sortable: false
        },
        {
            Header: `Components`,
            entityType: entityTypes.COMPONENT,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            // eslint-disable-next-line
            Cell: ({ original, pdf }) => (
                <TableCountLink
                    entityType={entityTypes.COMPONENT}
                    count={original.componentCount}
                    textOnly={inFindingsSection || pdf}
                    selectedRowId={original.cve}
                />
            ),
            accessor: 'componentCount',
            id: 'componentCount',
            sortField: cveSortFields.COMPONENTS,
            sortable: false
        },
        {
            Header: `Scanned`,
            headerClassName: `w-1/10 text-left ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => (
                <DateTimeField date={original.lastScanned} asString={pdf} />
            ),
            accessor: 'lastScanned',
            id: 'lastScanned',
            sortField: cveSortFields.SCANNED
        },
        {
            Header: `Published`,
            headerClassName: `w-1/10 ${defaultHeaderClassName}`,
            className: `w-1/10 ${defaultColumnClassName}`,
            Cell: ({ original, pdf }) => (
                <DateTimeField date={original.publishedOn} asString={pdf} />
            ),
            accessor: 'publishedOn',
            id: 'published',
            sortField: cveSortFields.PUBLISHED
        }
    ];

    return removeEntityContextColumns(tableColumns, workflowState);
}

const maxLengthForSummary = 360; // based on showing up to approximately 2 lines before table starts scrolling horizontally

export function renderCveDescription(row) {
    const { original } = row;
    const truncatedSummary = truncate(original.summary, maxLengthForSummary);
    return (
        <div className="hover:bg-base-100 px-2 pb-4 pt-1 text-base-500">
            {truncatedSummary || 'No description available.'}
        </div>
    );
}

const VulnMgmtCves = ({
    history,
    selectedRowId,
    search,
    sort,
    page,
    data,
    totalResults,
    addToast,
    removeToast,
    refreshTrigger,
    setRefreshTrigger
}) => {
    const [selectedCveIds, setSelectedCveIds] = useState([]);
    const [bulkActionCveIds, setBulkActionCveIds] = useState([]);

    const workflowState = useContext(workflowStateContext);

    const CVES_QUERY = gql`
        query getCves($query: String, $pagination: Pagination) {
            results: vulnerabilities(query: $query, pagination: $pagination) {
                ...cveFields
            }
            count: vulnerabilityCount(query: $query)
        }
        ${VULN_CVE_LIST_FRAGMENT}
    `;

    const viewingSuppressed = doesSearchContain(search, cveSortFields.SUPPRESSED);

    const tableSort = sort || defaultCveSort;
    const queryOptions = {
        variables: {
            query: queryService.objectToWhereClause(search),
            cachebuster: refreshTrigger,
            pagination: queryService.getPagination(tableSort, page, LIST_PAGE_SIZE)
        }
    };

    const addToPolicy = cveId => e => {
        e.stopPropagation();

        const cveIdsToAdd = cveId ? [cveId] : selectedCveIds;

        if (cveIdsToAdd.length) {
            setBulkActionCveIds(cveIdsToAdd);
        } else {
            throw new Error(
                'Logic error: tried to open Add to Policy dialog without any policy selected.'
            );
        }
    };

    const toggleCveSuppression = cveId => e => {
        e.stopPropagation();

        const cveIdsToToggle = cveId ? [cveId] : selectedCveIds;

        const suppressionState = !viewingSuppressed;
        updateCveSuppressedState(cveIdsToToggle, suppressionState)
            .then(() => {
                setSelectedCveIds([]);

                // changing this param value on the query vars, to force the query to refetch
                setRefreshTrigger(Math.random());

                // can't use pluralize() because of this bug: https://github.com/blakeembrey/pluralize/issues/127
                const pluralizedCVEs = cveIdsToToggle.length === 1 ? 'CVE' : 'CVEs';

                addToast(`Successfully suppressed ${cveIdsToToggle.length} ${pluralizedCVEs}`);
                setTimeout(removeToast, 2000);
            })
            .catch(evt => {
                addToast(`Could not suppress all of the selected CVEs: ${evt.message}`);
                setTimeout(removeToast, 2000);
            });
    };
    const toggleSuppressedView = () => {
        const currentSearchState = workflowState.getCurrentSearchState();

        const targetSearchState = { ...currentSearchState };
        if (viewingSuppressed) {
            delete targetSearchState[cveSortFields.SUPPRESSED];
        } else {
            targetSearchState[cveSortFields.SUPPRESSED] = true;
        }

        const newWorkflowState = workflowState.setSearch(targetSearchState);
        const newUrl = newWorkflowState.toUrl();
        history.push(newUrl);
    };

    function closeDialog(idsToStaySelected = []) {
        setBulkActionCveIds([]);
        setSelectedCveIds(idsToStaySelected);
    }

    const renderRowActionButtons = ({ id }) => (
        <div className="flex border-2 border-r-2 border-base-400 bg-base-100">
            <RowActionButton
                text="Add to Policy"
                onClick={addToPolicy(id)}
                icon={<Icon.Plus className="mt-1 h-4 w-4" />}
            />
            <RowActionButton
                text={`${viewingSuppressed ? 'Unsuppress CVE' : 'Suppress CVE'}`}
                border="border-l-2 border-base-400"
                onClick={toggleCveSuppression(id)}
                icon={
                    viewingSuppressed ? (
                        <Icon.Bell className="mt-1 h-4 w-4" />
                    ) : (
                        <Icon.BellOff className="mt-1 h-4 w-4" />
                    )
                }
            />
        </div>
    );

    const toggleButtonText = viewingSuppressed ? 'Unsuppress' : 'Suppress';
    const viewButtonText = viewingSuppressed ? 'View Unsuppressed' : 'View Suppressed';

    const tableHeaderComponents = (
        <React.Fragment>
            <PanelButton
                icon={<Icon.Plus className="h-4 w-4" />}
                className="btn-icon btn-tertiary"
                onClick={addToPolicy()}
                disabled={selectedCveIds.length === 0}
                tooltip="Add Selected CVEs to Policy"
            >
                Add to Policy
            </PanelButton>
            <PanelButton
                icon={
                    viewingSuppressed ? (
                        <Icon.Bell className="h-4 w-4" />
                    ) : (
                        <Icon.BellOff className="h-4 w-4" />
                    )
                }
                className="btn-icon btn-tertiary ml-2"
                onClick={toggleCveSuppression()}
                disabled={selectedCveIds.length === 0}
                tooltip={`${toggleButtonText} Selected CVEs`}
            >
                {toggleButtonText}
            </PanelButton>
            <span className="w-px bg-base-400 ml-2" />
            <PanelButton
                icon={
                    viewingSuppressed ? (
                        <Icon.Zap className="h-4 w-4" />
                    ) : (
                        <Icon.Archive className="h-4 w-4" />
                    )
                }
                className="btn-icon btn-tertiary ml-2"
                onClick={toggleSuppressedView}
                tooltip={`${viewButtonText} CVEs`}
            >
                {viewButtonText}
            </PanelButton>
        </React.Fragment>
    );

    return (
        <>
            <WorkflowListPage
                data={data}
                totalResults={totalResults}
                query={CVES_QUERY}
                queryOptions={queryOptions}
                idAttribute="cve"
                entityListType={entityTypes.CVE}
                getTableColumns={getCveTableColumns}
                selectedRowId={selectedRowId}
                search={search}
                page={page}
                showSubrows
                SubComponent={renderCveDescription}
                checkbox
                tableHeaderComponents={tableHeaderComponents}
                selection={selectedCveIds}
                setSelection={setSelectedCveIds}
                renderRowActionButtons={renderRowActionButtons}
            />
            {bulkActionCveIds.length > 0 && (
                <CveBulkActionDialogue
                    closeAction={closeDialog}
                    bulkActionCveIds={bulkActionCveIds}
                />
            )}
        </>
    );
};

VulnMgmtCves.propTypes = {
    ...workflowListPropTypes,
    refreshTrigger: PropTypes.number,
    setRefreshTrigger: PropTypes.func
};
VulnMgmtCves.defaultProps = {
    ...workflowListDefaultProps,
    sort: null,
    refreshTrigger: 0,
    setRefreshTrigger: null
};

const mapDispatchToProps = {
    addToast: notificationActions.addNotification,
    removeToast: notificationActions.removeOldestNotification
};

export default withRouter(
    connect(
        null,
        mapDispatchToProps
    )(VulnMgmtCves)
);
