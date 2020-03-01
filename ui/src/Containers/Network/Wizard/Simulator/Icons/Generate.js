import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import Tooltip from 'Components/Tooltip';
import TooltipOverlay from 'Components/TooltipOverlay';
import { actions as wizardActions } from 'reducers/network/wizard';

import generate from 'images/generate.svg';

class Generate extends Component {
    static propTypes = {
        generatePolicyModification: PropTypes.func.isRequired
    };

    onClick = () => {
        this.props.generatePolicyModification();
    };

    render() {
        return (
            <Tooltip content={<TooltipOverlay>Generate a new YAML</TooltipOverlay>}>
                <button
                    type="button"
                    className="inline-block px-2 py-2 border-r border-base-300 cursor-pointer"
                    onClick={this.onClick}
                >
                    <img
                        className="text-primary-700 h-4 w-4 hover:bg-base-200"
                        alt=""
                        src={generate}
                    />
                </button>
            </Tooltip>
        );
    }
}

const mapDispatchToProps = {
    generatePolicyModification: wizardActions.generateNetworkPolicyModification
};

export default connect(
    null,
    mapDispatchToProps
)(Generate);
