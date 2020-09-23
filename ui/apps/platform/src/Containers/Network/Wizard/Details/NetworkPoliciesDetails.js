import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { selectors } from 'reducers';

import Button from 'Components/Button';
import CollapsibleCard from 'Components/CollapsibleCard';
import NoResultsMessage from 'Components/NoResultsMessage';
import download from 'utils/download';
import * as Icon from 'react-feather';

class NetworkPoliciesDetails extends Component {
    static propTypes = {
        networkPolicies: PropTypes.arrayOf(PropTypes.shape({ id: PropTypes.string.isRequired }))
            .isRequired,
    };

    downloadYamlFile = (name, content, type) => () => {
        download(name, content, type);
    };

    renderOverview() {
        const { networkPolicies } = this.props;
        if (!networkPolicies.length) {
            return <NoResultsMessage message="No network policies have been applied" />;
        }
        return (
            <div>
                {networkPolicies.map((networkPolicy) => {
                    const { id, name, yaml } = networkPolicy;
                    return (
                        <div className="px-3 py-5" key={id}>
                            <CollapsibleCard title={name}>
                                <pre className="font-600 font-sans h-full leading-normal p-3">
                                    {yaml}
                                </pre>
                                <div className="flex justify-center p-3 border-t border-base-400">
                                    <Button
                                        className="download uppercase text-primary-600 p-2 text-center text-sm border border-solid bg-primary-200 border-primary-300 hover:bg-primary-100"
                                        onClick={this.downloadYamlFile(
                                            `${name}.yaml`,
                                            yaml,
                                            'yaml'
                                        )}
                                        tabIndex="-1"
                                        icon={<Icon.Download className="h-3 w-3 mr-4" />}
                                        text="Download YAML file"
                                    />
                                </div>
                            </CollapsibleCard>
                        </div>
                    );
                })}
            </div>
        );
    }

    render() {
        return <div className="w-full h-full">{this.renderOverview()}</div>;
    }
}

const mapStateToProps = createStructuredSelector({
    networkPolicies: selectors.getNetworkPolicies,
});

export default connect(mapStateToProps)(NetworkPoliciesDetails);
