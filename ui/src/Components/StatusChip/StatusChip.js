import React from 'react';
import PropTypes from 'prop-types';
import { capitalize } from 'lodash';

import LabelChip from 'Components/LabelChip';

const successStates = ['active', 'pass'];
const alertStates = ['inactive', 'fail'];

const StatusChip = ({ status, size, asString }) => {
    if (asString) return capitalize(status);
    let type = null;
    if (successStates.includes(status)) {
        type = 'success';
    } else if (alertStates.includes(status)) {
        type = 'alert';
    }

    return type ? <LabelChip text={status} type={type} size={size} /> : '—';
};

StatusChip.propTypes = {
    status: PropTypes.string,
    size: PropTypes.string,
    asString: PropTypes.bool
};

StatusChip.defaultProps = {
    status: '',
    size: 'large',
    asString: false
};

export default StatusChip;
