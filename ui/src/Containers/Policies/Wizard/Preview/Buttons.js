import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { actions as wizardActions } from 'reducers/policies/wizard';
import * as Icon from 'react-feather';
import wizardStages from 'Containers/Policies/Wizard/wizardStages';

import PanelButton from 'Components/PanelButton';

class Buttons extends Component {
    static propTypes = {
        setWizardStage: PropTypes.func.isRequired
    };

    goBackToEdit = () => this.props.setWizardStage(wizardStages.edit);

    goToEnforcement = () => this.props.setWizardStage(wizardStages.enforcement);

    render() {
        return (
            <React.Fragment>
                <PanelButton
                    icon={<Icon.ArrowLeft className="h-4 w-4" />}
                    className="btn btn-base"
                    onClick={this.goBackToEdit}
                    tooltip="Back to previous step"
                >
                    Previous
                </PanelButton>
                <PanelButton
                    icon={<Icon.ArrowRight className="h-4 w-4" />}
                    className="btn btn-base"
                    onClick={this.goToEnforcement}
                    tooltip="Go to next step"
                >
                    Next
                </PanelButton>
            </React.Fragment>
        );
    }
}

const mapDispatchToProps = {
    setWizardStage: wizardActions.setWizardStage
};

export default connect(
    null,
    mapDispatchToProps
)(Buttons);
