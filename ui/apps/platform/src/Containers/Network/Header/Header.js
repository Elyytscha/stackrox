import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createSelector, createStructuredSelector } from 'reselect';
import { selectors } from 'reducers';
import { actions as pageActions } from 'reducers/network/page';
import { actions as searchActions } from 'reducers/network/search';

import PageHeader from 'Components/PageHeader';
import ReduxSearchInput from 'Containers/Search/ReduxSearchInput';
import FeatureEnabled from 'Containers/FeatureEnabled';
import { knownBackendFlags } from 'utils/featureFlags';
import ClusterSelect from './ClusterSelect';
import SimulatorButton from './SimulatorButton';
import TimeWindowSelector from './TimeWindowSelector';
import CIDRFormButton from './CIDRFormButton';

class Header extends Component {
    static propTypes = {
        searchOptions: PropTypes.arrayOf(PropTypes.object).isRequired,
        searchModifiers: PropTypes.arrayOf(PropTypes.object).isRequired,
        searchSuggestions: PropTypes.arrayOf(PropTypes.object).isRequired,
        setSearchOptions: PropTypes.func.isRequired,
        setSearchModifiers: PropTypes.func.isRequired,
        setSearchSuggestions: PropTypes.func.isRequired,
        isViewFiltered: PropTypes.bool.isRequired,
        closeWizard: PropTypes.func.isRequired,
    };

    onSearch = (searchOptions) => {
        if (searchOptions.length && !searchOptions[searchOptions.length - 1].type) {
            this.props.closeWizard();
        }
    };

    render() {
        const subHeader = this.props.isViewFiltered ? 'Filtered view' : 'Default view';
        return (
            <>
                <PageHeader
                    header="Network Graph"
                    subHeader={subHeader}
                    classes="flex-1 border-none"
                >
                    <ClusterSelect />
                    <ReduxSearchInput
                        id="network"
                        className="w-full pl-2"
                        searchOptions={this.props.searchOptions}
                        searchModifiers={this.props.searchModifiers}
                        searchSuggestions={this.props.searchSuggestions}
                        setSearchOptions={this.props.setSearchOptions}
                        setSearchModifiers={this.props.setSearchModifiers}
                        setSearchSuggestions={this.props.setSearchSuggestions}
                        onSearch={this.onSearch}
                    />
                    <TimeWindowSelector />
                    <SimulatorButton />
                </PageHeader>
                <FeatureEnabled featureFlag={knownBackendFlags.ROX_NETWORK_GRAPH_EXTERNAL_SRCS}>
                    {({ featureEnabled }) => {
                        return featureEnabled && <CIDRFormButton />;
                    }}
                </FeatureEnabled>
            </>
        );
    }
}

const isViewFiltered = createSelector(
    [selectors.getNetworkSearchOptions],
    (searchOptions) => searchOptions.length !== 0
);

const mapStateToProps = createStructuredSelector({
    searchOptions: selectors.getNetworkSearchOptions,
    searchModifiers: selectors.getNetworkSearchModifiers,
    searchSuggestions: selectors.getNetworkSearchSuggestions,
    isViewFiltered,
});

const mapDispatchToProps = {
    setSearchOptions: searchActions.setNetworkSearchOptions,
    setSearchModifiers: searchActions.setNetworkSearchModifiers,
    setSearchSuggestions: searchActions.setNetworkSearchSuggestions,
    closeWizard: pageActions.closeNetworkWizard,
};

export default connect(mapStateToProps, mapDispatchToProps)(Header);
