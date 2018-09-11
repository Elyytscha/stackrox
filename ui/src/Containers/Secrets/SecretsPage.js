import React, { Component } from 'react';
import PropTypes from 'prop-types';
import ReactRouterPropTypes from 'react-router-prop-types';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';

import { selectors } from 'reducers';
import { actions as secretsActions } from 'reducers/secrets';
import NoResultsMessage from 'Components/NoResultsMessage';
import Table from 'Components/Table';
import Panel from 'Components/Panel';
import PageHeader from 'Components/PageHeader';
import SearchInput from 'Components/SearchInput';
import dateFns from 'date-fns';
import dateTimeFormat from 'constants/dateTimeFormat';
import SecretDetails, { secretTypeEnumMapping } from './SecretDetails';

class SecretPage extends Component {
    static propTypes = {
        secrets: PropTypes.arrayOf(PropTypes.object).isRequired,
        selectedSecret: PropTypes.shape({
            id: PropTypes.string.isRequired
        }),
        searchOptions: PropTypes.arrayOf(PropTypes.object).isRequired,
        searchModifiers: PropTypes.arrayOf(PropTypes.object).isRequired,
        searchSuggestions: PropTypes.arrayOf(PropTypes.object).isRequired,
        setSearchOptions: PropTypes.func.isRequired,
        setSearchModifiers: PropTypes.func.isRequired,
        setSearchSuggestions: PropTypes.func.isRequired,
        history: ReactRouterPropTypes.history.isRequired,
        location: ReactRouterPropTypes.location.isRequired,
        isViewFiltered: PropTypes.bool.isRequired
    };

    static defaultProps = {
        selectedSecret: null
    };

    onSearch = searchOptions => {
        if (searchOptions.length && !searchOptions[searchOptions.length - 1].type) {
            this.props.history.push('/main/secrets');
        }
    };

    updateSelectedSecret = secret => {
        const urlSuffix = secret && secret.id ? `/${secret.id}` : '';
        this.props.history.push({
            pathname: `/main/secrets${urlSuffix}`,
            search: this.props.location.search
        });
    };

    renderTable() {
        const columns = [
            { accessor: 'name', Header: 'Name' },
            {
                id: 'createdAt',
                accessor: d => dateFns.format(d.createdAt, dateTimeFormat),
                Header: 'Created'
            },
            {
                id: 'types',
                accessor: d => d.types.map(v => secretTypeEnumMapping[v]).join(', '),
                Header: 'Types'
            },
            { accessor: 'clusterName', Header: 'Cluster' },
            { accessor: 'namespace', Header: 'Namespace' }
        ];
        const { secrets, selectedSecret } = this.props;
        const rows = secrets;
        const id = selectedSecret && selectedSecret.id;
        if (!rows.length)
            return <NoResultsMessage message="No results found. Please refine your search." />;
        return (
            <Table
                rows={rows}
                columns={columns}
                onRowClick={this.updateSelectedSecret}
                selectedRowId={id}
                noDataText="No results found. Please refine your search."
            />
        );
    }

    renderSidePanel = () => {
        const { selectedSecret } = this.props;
        if (!selectedSecret) return null;

        return (
            <div className="w-2/3">
                <Panel header={selectedSecret.name} onClose={this.updateSelectedSecret}>
                    <SecretDetails secret={selectedSecret} />
                </Panel>
            </div>
        );
    };

    render() {
        const subHeader = this.props.isViewFiltered ? 'Filtered view' : 'Default view';
        return (
            <section className="flex flex-1 h-full">
                <div className="flex flex-1 flex-col">
                    <PageHeader header="Secrets" subHeader={subHeader}>
                        <SearchInput
                            className="flex flex-1"
                            id="secrets"
                            searchOptions={this.props.searchOptions}
                            searchModifiers={this.props.searchModifiers}
                            searchSuggestions={this.props.searchSuggestions}
                            setSearchOptions={this.props.setSearchOptions}
                            setSearchModifiers={this.props.setSearchModifiers}
                            setSearchSuggestions={this.props.setSearchSuggestions}
                            onSearch={this.onSearch}
                        />
                    </PageHeader>
                    <div className="flex flex-1">
                        <div className="w-full pl-3 pt-3 pr-3 overflow-y-scroll bg-white rounded-sm shadow bg-base-100">
                            {this.renderTable()}
                        </div>
                        {this.renderSidePanel()}
                    </div>
                </div>
            </section>
        );
    }
}

const isViewFiltered = createSelector(
    [selectors.getSecretsSearchOptions],
    searchOptions => searchOptions.length !== 0
);

const getSelectedSecret = (state, props) => {
    const { secretId } = props.match.params;
    return secretId ? selectors.getSecret(state, secretId) : null;
};

const mapStateToProps = createStructuredSelector({
    secrets: selectors.getFilteredSecrets,
    selectedSecret: getSelectedSecret,
    searchOptions: selectors.getSecretsSearchOptions,
    searchModifiers: selectors.getSecretsSearchModifiers,
    searchSuggestions: selectors.getSecretsSearchSuggestions,
    isViewFiltered
});

const mapDispatchToProps = {
    setSearchOptions: secretsActions.setSecretsSearchOptions,
    setSearchModifiers: secretsActions.setSecretsSearchModifiers,
    setSearchSuggestions: secretsActions.setSecretsSearchSuggestions
};

export default connect(mapStateToProps, mapDispatchToProps)(SecretPage);
