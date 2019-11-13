import React, { useContext } from 'react';
import { Link } from 'react-router-dom';
import workflowStateContext from 'Containers/workflowStateContext';
import pluralize from 'pluralize';
import CollapsibleSection from 'Components/CollapsibleSection';
import StatusChip from 'Components/StatusChip';
import RiskScore from 'Components/RiskScore';
import Metadata from 'Components/Metadata';
import Tabs from 'Components/Tabs';
import TabContent from 'Components/TabContent';
import entityTypes from 'constants/entityTypes';
import TopRiskyEntitiesByVulnerabilities from 'Containers/VulnMgmt/widgets/TopRiskyEntitiesByVulnerabilities';
import RecentlyDetectedVulnerabilities from 'Containers/VulnMgmt/widgets/RecentlyDetectedVulnerabilities';
import TopRiskiestImagesAndComponents from 'Containers/VulnMgmt/widgets/TopRiskiestImagesAndComponents';
import DeploymentsWithMostSeverePolicyViolations from 'Containers/VulnMgmt/widgets/DeploymentsWithMostSeverePolicyViolations';
import { getPolicyTableColumns } from 'Containers/VulnMgmt/List/Policies/VulnMgmtListPolicies';
import { getCveTableColumns } from 'Containers/VulnMgmt/List/Cves/VulnMgmtListCves';
import { entityGridContainerClassName } from 'Containers/Workflow/WorkflowEntityPage';

import RelatedEntitiesSideList from '../RelatedEntitiesSideList';
import TableWidget from '../TableWidget';

const emptyNamespace = {
    deploymentCount: 0,
    componentCount: 0,
    metadata: {
        clusterName: '',
        clusterId: '',
        priority: 0,
        labels: [],
        id: ''
    },
    policyStatus: {
        status: '',
        failingPolicies: []
    },
    vulnCount: 0,
    vulnerabilities: []
};

const VulnMgmtNamespaceOverview = ({ data, entityContext }) => {
    const workflowState = useContext(workflowStateContext);

    // guard against incomplete GraphQL-cached data
    const safeData = { ...emptyNamespace, ...data };

    const { metadata, policyStatus, vulnerabilities } = safeData;

    if (!metadata || !policyStatus) return null;

    const { clusterName, clusterId, priority, labels, id } = metadata;
    const { failingPolicies, status } = policyStatus;
    const fixableCves = vulnerabilities.filter(cve => cve.isFixable);
    const metadataKeyValuePairs = [];

    if (!entityContext[entityTypes.CLUSTER]) {
        const clusterLink = workflowState.pushRelatedEntity(entityTypes.CLUSTER, clusterId).toUrl();
        metadataKeyValuePairs.push({
            key: 'Cluster',
            value: <Link to={clusterLink}>{clusterName}</Link>
        });
    }

    const namespaceStats = [
        <RiskScore key="risk-score" score={priority} />,
        <React.Fragment key="policy-status">
            <span className="pb-2">Policy status:</span>
            <StatusChip status={status} size="large" />
        </React.Fragment>
    ];

    const newEntityContext = { ...entityContext, [entityTypes.NAMESPACE]: id };

    return (
        <div className="flex h-full">
            <div className="flex flex-col flex-grow min-w-0">
                <CollapsibleSection title="Namespace summary">
                    <div className={entityGridContainerClassName}>
                        <div className="s-1">
                            <Metadata
                                className="h-full min-w-48 bg-base-100 bg-counts-widget"
                                keyValuePairs={metadataKeyValuePairs}
                                statTiles={namespaceStats}
                                labels={labels}
                                title="Details & Metadata"
                            />
                        </div>
                        <div className="sx-1 lg:sx-2 sy-1 h-55">
                            <TopRiskyEntitiesByVulnerabilities
                                defaultSelection={entityTypes.DEPLOYMENT}
                                riskEntityTypes={[entityTypes.DEPLOYMENT, entityTypes.IMAGE]}
                                entityContext={newEntityContext}
                                small
                            />
                        </div>
                        <div className="s-1">
                            <RecentlyDetectedVulnerabilities entityContext={newEntityContext} />
                        </div>
                        <div className="s-1">
                            <TopRiskiestImagesAndComponents entityContext={newEntityContext} />
                        </div>
                        <div className="s-1">
                            <DeploymentsWithMostSeverePolicyViolations
                                entityContext={newEntityContext}
                            />
                        </div>
                    </div>
                </CollapsibleSection>
                <CollapsibleSection title="Namespace findings">
                    <div className="flex pdf-page pdf-stretch shadow rounded relative rounded bg-base-100 mb-4 ml-4 mr-4">
                        <Tabs
                            hasTabSpacing
                            headers={[{ text: 'Policies' }, { text: 'Fixable CVEs' }]}
                        >
                            <TabContent>
                                <TableWidget
                                    header={`${failingPolicies.length} failing ${pluralize(
                                        entityTypes.POLICY,
                                        failingPolicies.length
                                    )} across this namespace`}
                                    entityType={entityTypes.POLICY}
                                    rows={failingPolicies}
                                    noDataText="No failing policies"
                                    className="bg-base-100"
                                    columns={getPolicyTableColumns(workflowState)}
                                    idAttribute="id"
                                />
                            </TabContent>
                            <TabContent>
                                <TableWidget
                                    header={`${fixableCves.length} fixable ${pluralize(
                                        entityTypes.CVE,
                                        fixableCves.length
                                    )} found across this namespace`}
                                    rows={fixableCves}
                                    entityType={entityTypes.CVE}
                                    noDataText="No fixable CVEs available in this namespace"
                                    className="bg-base-100"
                                    columns={getCveTableColumns(workflowState)}
                                    idAttribute="cve"
                                />
                            </TabContent>
                        </Tabs>
                    </div>
                </CollapsibleSection>
            </div>

            <RelatedEntitiesSideList
                entityType={entityTypes.NAMESPACE}
                entityContext={newEntityContext}
                data={safeData}
            />
        </div>
    );
};

export default VulnMgmtNamespaceOverview;
