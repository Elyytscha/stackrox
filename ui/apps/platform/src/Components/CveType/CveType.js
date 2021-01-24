import React from 'react';
import PropTypes from 'prop-types';

import LabelChip from 'Components/LabelChip';

const cveTypes = ['IMAGE_CVE', 'K8S_CVE', 'ISTIO_CVE', 'NODE_CVE'];
const cveTypeMap = {
    IMAGE_CVE: 'Image CVE',
    K8S_CVE: 'Kubernetes CVE',
    ISTIO_CVE: 'Istio CVE',
    NODE_CVE: 'Node CVE',
};

const CveType = ({ type, context }) => {
    const typeText = cveTypeMap[type] || 'Unknown';

    return context === 'callout' ? (
        <LabelChip type="base" text={`Type: ${typeText}`} />
    ) : (
        <span>{typeText}</span>
    );
};

CveType.propTypes = {
    type: PropTypes.oneOf(cveTypes),
    context: PropTypes.oneOf(['callout', 'bare']),
};

CveType.defaultProps = {
    type: '',
    context: 'bare',
};

export default CveType;
