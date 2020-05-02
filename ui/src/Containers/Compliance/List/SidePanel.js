import React from 'react';
import PropTypes from 'prop-types';
import Panel from 'Components/Panel';
import ReactRouterPropTypes from 'react-router-prop-types';
import { resourceTypes, standardEntityTypes } from 'constants/entityTypes';
import { Link, withRouter } from 'react-router-dom';
import URLService from 'modules/URLService';
import getEntityName from 'modules/getEntityName';
import { entityNameQueryMap } from 'modules/queryMap';
import Query from 'Components/CacheFirstQuery';
import * as Icon from 'react-feather';
// TODO: this exception will be unnecessary once Compliance pages are re-structured like Config Management
/* eslint-disable import/no-cycle */
import ControlPage from 'Containers/Compliance/Entity/Control';
import NamespacePage from '../Entity/Namespace';
import ClusterPage from '../Entity/Cluster';
import NodePage from '../Entity/Node';
import DeploymentPage from '../Entity/Deployment';
/* eslint-enable import/no-cycle */

const ComplianceListSidePanel = ({ entityType, entityId, match, location, history }) => {
    function getEntityPage() {
        switch (entityType) {
            case resourceTypes.NODE:
                return <NodePage entityId={entityId} sidePanelMode />;
            case resourceTypes.NAMESPACE:
                return <NamespacePage entityId={entityId} sidePanelMode />;
            case resourceTypes.CLUSTER:
                return <ClusterPage entityId={entityId} sidePanelMode />;
            case resourceTypes.DEPLOYMENT:
                return <DeploymentPage entityId={entityId} sidePanelMode />;
            case standardEntityTypes.CONTROL:
                return <ControlPage entityId={entityId} sidePanelMode />;
            default:
                return null;
        }
    }

    function closeSidePanel() {
        const baseURL = URLService.getURL(match, location).clearSidePanelParams().url();
        history.push(baseURL);
    }
    const headerUrl = URLService.getURL(match, location).base(entityType, entityId).url();

    return (
        <Query query={entityNameQueryMap[entityType]} variables={{ id: entityId }}>
            {({ loading, data }) => {
                let linkText = 'loading...';
                if (!loading && data) {
                    linkText = getEntityName(entityType, data);
                }
                const headerTextComponent = (
                    <div className="w-full flex items-center">
                        <div className="flex items-center">
                            <Link
                                to={headerUrl}
                                className="w-full flex text-primary-700 hover:text-primary-800 focus:text-primary-700"
                            >
                                <div className="flex flex-1 uppercase items-center tracking-wide pl-4 leading-normal font-700">
                                    {linkText}
                                </div>
                            </Link>
                            <Link
                                className="mx-2 text-primary-700 hover:text-primary-800 p-1 bg-primary-300 rounded flex"
                                to={headerUrl}
                                target="_blank"
                            >
                                <Icon.ExternalLink size="14" />
                            </Link>
                        </div>
                    </div>
                );
                const entityPage = getEntityPage();
                return (
                    <Panel
                        className="bg-primary-200 z-40 w-full h-full absolute right-0 top-0 md:w-1/2 min-w-108 md:relative"
                        headerTextComponent={headerTextComponent}
                        onClose={closeSidePanel}
                        id="side-panel"
                    >
                        {entityPage}
                    </Panel>
                );
            }}
        </Query>
    );
};

ComplianceListSidePanel.propTypes = {
    history: ReactRouterPropTypes.history.isRequired,
    match: ReactRouterPropTypes.match.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    entityType: PropTypes.string.isRequired,
    entityId: PropTypes.string.isRequired,
};

export default withRouter(ComplianceListSidePanel);
