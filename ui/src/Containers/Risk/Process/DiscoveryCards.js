import React from 'react';
import PropTypes from 'prop-types';
import orderBy from 'lodash/orderBy';

import { knownBackendFlags } from 'utils/featureFlags';

import FeatureEnabled from 'Containers/FeatureEnabled';
import ProcessDiscoveryCard from './DiscoveryCard';
import Binaries from './Binaries';
import RiskProcessComments from './RiskProcessComments';

function DiscoveryCards({ deploymentId, processGroup, processEpoch, setProcessEpoch }) {
    const sortedProcessGroups = orderBy(
        processGroup.groups,
        ['suspicious', 'name'],
        ['desc', 'asc']
    );
    return sortedProcessGroups.map((pg, i, list) => (
        <div className={`px-3 pt-5 ${i === list.length - 1 ? 'pb-5' : ''}`} key={pg.name}>
            <ProcessDiscoveryCard
                process={pg}
                deploymentId={deploymentId}
                processEpoch={processEpoch}
                setProcessEpoch={setProcessEpoch}
            >
                <div className="p-2">
                    <FeatureEnabled featureFlag={knownBackendFlags.ROX_IQT_ANALYST_NOTES_UI}>
                        <RiskProcessComments />
                    </FeatureEnabled>
                </div>
                <Binaries processes={pg.groups} />
            </ProcessDiscoveryCard>
        </div>
    ));
}

DiscoveryCards.propTypes = {
    deploymentId: PropTypes.string.isRequired,
    processGroup: PropTypes.shape({
        groups: PropTypes.arrayOf(PropTypes.object)
    }).isRequired,
    processEpoch: PropTypes.number.isRequired,
    setProcessEpoch: PropTypes.func.isRequired
};

export default DiscoveryCards;
