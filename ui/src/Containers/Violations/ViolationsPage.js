import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactRouterPropTypes from 'react-router-prop-types';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import { sortDate, sortSeverity } from 'sorters/sorters';
import { actions as alertActions } from 'reducers/alerts';
import { selectors } from 'reducers';
import dateFns from 'date-fns';
import dateTimeFormat from 'constants/dateTimeFormat';
import Tooltip from 'rc-tooltip';
import * as Icon from 'react-feather';

import { severityLabels, lifecycleStageLabels } from 'messages/common';

import NoResultsMessage from 'Components/NoResultsMessage';
import {
    wrapClassName,
    defaultHeaderClassName,
    defaultColumnClassName,
    rtTrActionsClassName
} from 'Components/Table';
import CheckboxTable from 'Components/CheckboxTable';
import { toggleRow, toggleSelectAll } from 'utils/checkboxUtils';
import PageHeader, { PageHeaderComponent } from 'Components/PageHeader';
import SearchInput from 'Components/SearchInput';
import Panel from 'Components/Panel';
import PanelButton from 'Components/PanelButton';
import TablePagination from 'Components/TablePagination';
import Dialog from 'Components/Dialog';
import ViolationsPanel from './ViolationsPanel';

const getSeverityClassName = severityValue => {
    const severityClassMapping = {
        Low: 'px-2 rounded-full bg-base-200 border-2 border-base-300 text-base-600',
        Medium: 'px-2 rounded-full bg-warning-200 border-2 border-warning-300 text-warning-800',
        High: 'px-2 rounded-full bg-caution-200 border-2 border-caution-300 text-caution-800',
        Critical: 'px-2 rounded-full bg-alert-200 border-2 border-alert-300 text-alert-800'
    };
    const res = severityClassMapping[severityValue];
    if (res) return res;
    throw new Error(`Unknown severity: ${severityValue}`);
};

class ViolationsPage extends Component {
    static propTypes = {
        violations: PropTypes.shape({}).isRequired,
        whitelistDeployment: PropTypes.func.isRequired,
        whitelistDeployments: PropTypes.func.isRequired,
        resolveAlerts: PropTypes.func.isRequired,
        history: ReactRouterPropTypes.history.isRequired,
        location: ReactRouterPropTypes.location.isRequired,
        match: ReactRouterPropTypes.match.isRequired,
        searchOptions: PropTypes.arrayOf(PropTypes.object).isRequired,
        searchModifiers: PropTypes.arrayOf(PropTypes.object).isRequired,
        searchSuggestions: PropTypes.arrayOf(PropTypes.object).isRequired,
        setSearchOptions: PropTypes.func.isRequired,
        setSearchModifiers: PropTypes.func.isRequired,
        setSearchSuggestions: PropTypes.func.isRequired,
        isViewFiltered: PropTypes.bool.isRequired
    };

    constructor(props) {
        super(props);
        this.state = {
            page: 0,
            showWhitelistConfirmationDialog: false,
            showResolveConfirmationDialog: false,
            selection: []
        };
    }

    onSearch = searchOptions => {
        if (searchOptions.length && !searchOptions[searchOptions.length - 1].type) {
            this.clearSelection();
            this.props.history.push('/main/violations');
        }
    };

    onPanelClose = () => {
        this.updateSelectedAlert();
    };

    setTablePage = newPage => {
        this.setState({ page: newPage });
    };

    isRunTimeViolation = violation => violation && violation.lifecycleStage === 'RUNTIME';

    getTableHeaderTextComponent = () => {
        const { violations, isViewFiltered } = this.props;
        const { length: selectionCount } = this.state.selection;
        const { length: rowCount } = Object.keys(violations);

        return (
            <PageHeaderComponent
                length={rowCount}
                selectionCount={selectionCount}
                type="Violation"
                isViewFiltered={isViewFiltered}
            />
        );
    };

    clearSelection = () => this.setState({ selection: [] });

    showResolveConfirmationDialog = () => {
        this.setState({ showResolveConfirmationDialog: true });
    };

    showWhitelistConfirmationDialog = () => {
        this.setState({ showWhitelistConfirmationDialog: true });
    };

    hideConfirmationDialog = () => {
        this.setState({
            showResolveConfirmationDialog: false,
            showWhitelistConfirmationDialog: false
        });
    };

    updateSelectedAlert = alert => {
        const urlSuffix = alert && alert.id ? `/${alert.id}` : '';
        this.props.history.push({
            pathname: `/main/violations${urlSuffix}`,
            search: this.props.location.search
        });
    };

    updateSelection = selection => this.setState({ selection });

    toggleRow = id => {
        const selection = toggleRow(id, this.state.selection);
        this.updateSelection(selection);
    };

    toggleSelectAll = () => {
        const { length: rowsLength } = Object.keys(this.props.violations);
        const tableRef = this.checkboxTable.reactTable;
        const selection = toggleSelectAll(rowsLength, this.state.selection, tableRef);
        this.updateSelection(selection);
    };

    resolveAlerts = () => {
        const { selection } = this.state;
        const { violations } = this.props;
        const resolveSelection = selection.filter(id => this.isRunTimeViolation(violations[id]));
        this.props.resolveAlerts(resolveSelection);
        this.hideConfirmationDialog();
        this.clearSelection();
    };

    resolveAlertHandler = (alertId, whitelist) => e => {
        e.stopPropagation();
        this.props.resolveAlerts([alertId], whitelist);
    };

    whitelistDeployments = () => {
        const { selection } = this.state;
        this.props.whitelistDeployments(selection);
        this.hideConfirmationDialog();
        this.clearSelection();
    };

    whitelistDeploymentHandler = alertId => e => {
        e.stopPropagation();
        this.props.whitelistDeployment(alertId);
    };

    renderWhitelistConfirmationDialog = () => {
        const numSelectedRows = this.state.selection.length;
        return (
            <Dialog
                isOpen={this.state.showWhitelistConfirmationDialog}
                text={`Are you sure you want to whitelist ${numSelectedRows} violation${
                    numSelectedRows === 1 ? '' : 's'
                }?`}
                onConfirm={this.whitelistDeployments}
                onCancel={this.hideConfirmationDialog}
            />
        );
    };

    renderResolveConfirmationDialog = () => {
        const { selection } = this.state;
        const { violations } = this.props;
        const numSelectedRows = selection.reduce(
            (acc, id) => (this.isRunTimeViolation(violations[id]) ? acc + 1 : acc),
            0
        );
        return (
            <Dialog
                isOpen={this.state.showResolveConfirmationDialog}
                text={`Are you sure you want to resolve ${numSelectedRows} violation${
                    numSelectedRows === 1 ? '' : 's'
                }?`}
                onConfirm={this.resolveAlerts}
                onCancel={this.hideConfirmationDialog}
            />
        );
    };

    renderRowActionButtons = alert => {
        const isRuntimeAlert = this.isRunTimeViolation(alert);
        return (
            <div
                data-test-id="alerts-hover-actions"
                className="flex border-2 border-r-2 border-base-400 bg-base-100 shadow"
            >
                {isRuntimeAlert && (
                    <div className="flex">
                        <Tooltip
                            placement="top"
                            mouseLeaveDelay={0}
                            overlay={
                                <div>Resolve violation and add all processes to whitelist</div>
                            }
                            overlayClassName="pointer-events-none"
                        >
                            <button
                                type="button"
                                data-test-id="resolve-button"
                                className="p-1 px-4 hover:bg-primary-200 text-primary-600 hover:text-primary-700"
                                onClick={this.resolveAlertHandler(alert.id, true)}
                            >
                                <Icon.ShieldOff className="mt-1 h-4 w-4" />
                            </button>
                        </Tooltip>
                        <Tooltip
                            placement="top"
                            mouseLeaveDelay={0}
                            overlay={<div>Mark as resolved</div>}
                            overlayClassName="pointer-events-none"
                        >
                            <button
                                type="button"
                                data-test-id="resolve-button"
                                className="p-1 px-4 hover:bg-primary-200 text-primary-600 hover:text-primary-700 border-l-2 border-base-400"
                                onClick={this.resolveAlertHandler(alert.id, false)}
                            >
                                <Icon.Check className="mt-1 h-4 w-4" />
                            </button>
                        </Tooltip>
                    </div>
                )}
                <Tooltip
                    placement="top"
                    mouseLeaveDelay={0}
                    overlay={
                        <div>
                            Whitelist <br /> deployment
                        </div>
                    }
                    overlayClassName="pointer-events-none text-center"
                >
                    <button
                        data-test-id="whitelist-deployment-button"
                        type="button"
                        className={`p-1 px-4 hover:bg-primary-200 text-primary-600 hover:text-primary-700 ${
                            isRuntimeAlert ? 'border-l-2 border-base-400' : ''
                        }`}
                        onClick={this.whitelistDeploymentHandler(alert.id)}
                    >
                        <Icon.BellOff className="mt-1 h-4 w-4" />
                    </button>
                </Tooltip>
            </div>
        );
    };

    renderPanel = () => {
        const { violations } = this.props;
        const { selection, page } = this.state;
        const whitelistCount = selection.length;
        let resolveCount = 0;
        selection.forEach(id => {
            if (this.isRunTimeViolation(violations[id])) resolveCount += 1;
        });
        const panelButtons = (
            <React.Fragment>
                {resolveCount !== 0 && (
                    <PanelButton
                        icon={<Icon.Check className="h-4 w- ml-1" />}
                        text={`Mark as Resolved (${resolveCount})`}
                        className="btn btn-base"
                        onClick={this.showResolveConfirmationDialog}
                    />
                )}
                {whitelistCount !== 0 && (
                    <PanelButton
                        icon={<Icon.BellOff className="h-4 w- ml-1" />}
                        text={`Whitelist (${whitelistCount})`}
                        className="btn btn-base"
                        onClick={this.showWhitelistConfirmationDialog}
                    />
                )}
            </React.Fragment>
        );
        const { length } = Object.keys(violations);
        const headerComponents = (
            <>
                {panelButtons}
                <TablePagination page={page} dataLength={length} setPage={this.setTablePage} />
            </>
        );
        return (
            <Panel
                headerTextComponent={this.getTableHeaderTextComponent()}
                headerComponents={headerComponents}
            >
                <div className="w-full">{this.renderSelectTable()}</div>
            </Panel>
        );
    };

    renderSelectTable = () => {
        const columns = [
            {
                Header: 'Deployment',
                accessor: 'deployment.name',
                headerClassName: `w-1/6 sticky-column left-checkbox-offset ${defaultHeaderClassName}`,
                className: `w-1/6 sticky-column left-checkbox-offset ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ original }) => (
                    <div className="flex">
                        <span
                            className="pr-2"
                            title={`${
                                original.deployment.inactive ? 'Inactive' : 'Active'
                            } Deployment`}
                        >
                            <Icon.Circle
                                className="h-2 w-2 text-success-600"
                                hidden={original.deployment.inactive}
                            />
                            <Icon.Slash
                                className="h-2 w-2 text-base-500"
                                hidden={!original.deployment.inactive}
                            />
                        </span>
                        <span>{original.deployment.name}</span>
                    </div>
                )
            },
            {
                Header: 'Cluster',
                accessor: 'deployment.clusterName',
                headerClassName: `w-1/7  ${defaultHeaderClassName}`,
                className: `w-1/7 ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ value }) => <span>{value}</span>
            },
            {
                Header: 'Namespace',
                accessor: 'deployment.namespace',
                headerClassName: `w-1/7 ${defaultHeaderClassName}`,
                className: `w-1/7 ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ value }) => <span>{value}</span>
            },
            {
                Header: 'Policy',
                accessor: 'policy.name',
                headerClassName: `w-1/6 ${defaultHeaderClassName}`,
                className: `w-1/6 ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ original }) => (
                    <Tooltip
                        placement="top"
                        mouseLeaveDelay={0}
                        overlay={<div>{original.policy.description}</div>}
                        overlayClassName="pointer-events-none text-white rounded max-w-xs p-2 text-sm text-center"
                    >
                        <div className="inline-block hover:text-primary-700 underline">
                            {original.policy.name}
                        </div>
                    </Tooltip>
                )
            },
            {
                Header: 'Enforced',
                accessor: '',
                headerClassName: `w-1/10  ${defaultHeaderClassName}`,
                className: `w-1/10 ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ original }) => {
                    const count = original.enforcementCount;
                    if (original.lifecycleStage === 'DEPLOY') {
                        const message = count === 0 ? 'No' : 'Yes';
                        return <span>{message}</span>;
                    }
                    const countMessage = count === 1 ? `1 time` : `${count} times`;
                    const message = count === 0 ? 'No' : countMessage;
                    return <span>{message}</span>;
                }
            },
            {
                Header: 'Severity',
                accessor: 'policy.severity',
                headerClassName: `text-center ${defaultHeaderClassName}`,
                className: `text-center ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ value }) => {
                    const severity = severityLabels[value];
                    return <div className={getSeverityClassName(severity)}>{severity}</div>;
                },
                sortMethod: sortSeverity,
                width: 90
            },
            {
                Header: 'Categories',
                accessor: 'policy.categories',
                headerClassName: `w-1/10 ${defaultHeaderClassName}`,
                className: `w-1/10 ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ value }) =>
                    value.length > 1 ? (
                        <Tooltip
                            placement="top"
                            mouseLeaveDelay={0}
                            overlay={<div>{value.join(' | ')}</div>}
                            overlayClassName="pointer-events-none text-white rounded max-w-xs p-2 w-full text-sm text-center"
                        >
                            <div>Multiple</div>
                        </Tooltip>
                    ) : (
                        value[0]
                    )
            },
            {
                Header: 'Lifecycle',
                accessor: 'lifecycleStage',
                headerClassName: `${defaultHeaderClassName}`,
                className: `${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ value }) => lifecycleStageLabels[value]
            },
            {
                Header: 'Time',
                accessor: 'time',
                headerClassName: `w-1/10 ${defaultHeaderClassName}`,
                className: `w-1/10 ${wrapClassName} ${defaultColumnClassName}`,
                Cell: ({ value }) => dateFns.format(value, dateTimeFormat),
                sortMethod: sortDate
            },
            {
                Header: '',
                accessor: '',
                headerClassName: 'hidden',
                className: rtTrActionsClassName,
                Cell: ({ original }) => this.renderRowActionButtons(original)
            }
        ];
        const rows = Object.values(this.props.violations);
        const id = this.props.match.params.alertId;
        if (!rows.length)
            return <NoResultsMessage message="No results found. Please refine your search." />;
        return (
            <CheckboxTable
                ref={r => (this.checkboxTable = r)} // eslint-disable-line
                rows={rows}
                columns={columns}
                onRowClick={this.updateSelectedAlert}
                toggleRow={this.toggleRow}
                toggleSelectAll={this.toggleSelectAll}
                selection={this.state.selection}
                selectedRowId={id}
                noDataText="No results found. Please refine your search."
                page={this.state.page}
            />
        );
    };

    renderSidePanel = () => {
        if (!this.props.match.params.alertId) return null;
        return (
            <ViolationsPanel
                key={this.props.match.params.alertId}
                alertId={this.props.match.params.alertId}
                onClose={this.onPanelClose}
            />
        );
    };

    render() {
        const subHeader = this.props.isViewFiltered ? 'Filtered view' : 'Default view';
        const defaultOption = this.props.searchModifiers.find(x => x.value === 'Deployment:');
        return (
            <section className="flex flex-1 flex-col h-full">
                <div className="flex flex-1 flex-col">
                    <PageHeader header="Violations" subHeader={subHeader}>
                        <SearchInput
                            className="w-full"
                            searchOptions={this.props.searchOptions}
                            searchModifiers={this.props.searchModifiers}
                            searchSuggestions={this.props.searchSuggestions}
                            setSearchOptions={this.props.setSearchOptions}
                            setSearchModifiers={this.props.setSearchModifiers}
                            setSearchSuggestions={this.props.setSearchSuggestions}
                            onSearch={this.onSearch}
                            defaultOption={defaultOption}
                            autoCompleteCategories={['ALERTS']}
                        />
                    </PageHeader>
                    <div className="flex flex-1 relative">
                        <div className="shadow border-primary-300 bg-base-100 w-full overflow-hidden">
                            {this.renderPanel()}
                        </div>
                        {this.renderSidePanel()}
                    </div>
                </div>
                {this.renderWhitelistConfirmationDialog()}
                {this.renderResolveConfirmationDialog()}
            </section>
        );
    }
}

const isViewFiltered = createSelector(
    [selectors.getAlertsSearchOptions],
    searchOptions => searchOptions.length !== 0
);

const mapStateToProps = createStructuredSelector({
    violations: selectors.getFilteredAlertsById,
    searchOptions: selectors.getAlertsSearchOptions,
    searchModifiers: selectors.getAlertsSearchModifiers,
    searchSuggestions: selectors.getAlertsSearchSuggestions,
    isViewFiltered
});

const mapDispatchToProps = (dispatch, props) => ({
    whitelistDeployment: alertId => dispatch(alertActions.whitelistDeployment.request(alertId)),
    whitelistDeployments: alertIds => dispatch(alertActions.whitelistDeployments.request(alertIds)),
    resolveAlerts: (alertIds, whitelist) =>
        dispatch(alertActions.resolveAlerts(alertIds, whitelist)),
    setSearchOptions: searchOptions => {
        if (searchOptions.length && !searchOptions[searchOptions.length - 1].type) {
            props.history.push('/main/violations');
        }
        dispatch(alertActions.setAlertsSearchOptions(searchOptions));
    },
    setSearchModifiers: alertActions.setAlertsSearchModifiers,
    setSearchSuggestions: alertActions.setAlertsSearchSuggestions
});

export default connect(
    mapStateToProps,
    mapDispatchToProps
)(ViolationsPage);
