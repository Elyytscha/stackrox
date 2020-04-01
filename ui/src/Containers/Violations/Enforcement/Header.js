import React from 'react';
import PropTypes from 'prop-types';

function deployHeader(count) {
    let message = '';
    if (count && count > 0) {
        message = 'on this deployment is enabled for this policy';
    } else {
        message = 'on this deployment is not enabled for this policy';
    }
    return message;
}

function runtimeHeader(count) {
    let message = '';
    if (count && count > 0) {
        if (count === 1) {
            message = '"Kill Pod" has been applied once';
        } else if (count > 1) {
            message = `"Kill Pod" has been applied ${count} times`;
        }
    } else {
        message = `on this deployment was not enabled as of the last known violation of this policy.`;
    }
    return message;
}

function Header({ lifecycleStage, enforcementCount }) {
    let countMessage = '';
    if (lifecycleStage === 'DEPLOY') {
        countMessage = deployHeader(enforcementCount);
    } else if (lifecycleStage === 'RUNTIME') {
        countMessage = runtimeHeader(enforcementCount);
    }

    return (
        <div
            className="p-3 pb-2 border-b border-base-300 text-base-600 font-700 text-lg leading-normal"
            data-test-id="enforcement-detail-message"
        >
            Enforcement {countMessage}
        </div>
    );
}

Header.propTypes = {
    lifecycleStage: PropTypes.string.isRequired,
    enforcementCount: PropTypes.number
};

Header.defaultProps = {
    enforcementCount: 0
};

export default Header;
