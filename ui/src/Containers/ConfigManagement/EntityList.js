import React from 'react';
import PropTypes from 'prop-types';
import entityTypes from 'constants/entityTypes';

import PageNotFound from 'Components/PageNotFound';
import Namespaces from './List/Namespaces';
import Subjects from './List/Subjects';
import ServiceAccounts from './List/ServiceAccounts';
import Clusters from './List/Clusters';
import Nodes from './List/Nodes';
import Deployments from './List/Deployments';
import Secrets from './List/Secrets';
import Roles from './List/Roles';
import Images from './List/Images';
import Policies from './List/Policies';

const entityComponentMap = {
    [entityTypes.SUBJECT]: Subjects,
    [entityTypes.SERVICE_ACCOUNT]: ServiceAccounts,
    [entityTypes.CLUSTER]: Clusters,
    [entityTypes.NAMESPACE]: Namespaces,
    [entityTypes.NODE]: Nodes,
    [entityTypes.DEPLOYMENT]: Deployments,
    [entityTypes.IMAGE]: Images,
    [entityTypes.SECRET]: Secrets,
    [entityTypes.ROLE]: Roles,
    [entityTypes.POLICY]: Policies
};

const EntityList = ({ className, entityListType, entityId, onRowClick }) => {
    const Component = entityComponentMap[entityListType];
    if (!Component) return <PageNotFound resourceType={entityListType} />;
    return <Component className={className} selectedRowId={entityId} onRowClick={onRowClick} />;
};

EntityList.propTypes = {
    className: PropTypes.string,
    entityListType: PropTypes.string.isRequired,
    entityId: PropTypes.string,
    onRowClick: PropTypes.func.isRequired
};

EntityList.defaultProps = {
    className: '',
    entityId: null
};

export default EntityList;
