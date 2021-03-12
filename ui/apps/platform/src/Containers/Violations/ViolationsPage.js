import React, { useEffect, useState } from 'react';
import ReactRouterPropTypes from 'react-router-prop-types';

import MessageCentered from 'Components/MessageCentered';
import { PageBody } from 'Components/Panel';
import SidePanelAdjacentArea from 'Components/SidePanelAdjacentArea';
import useEntitiesByIdsCache from 'hooks/useEntitiesByIdsCache';
import LIFECYCLE_STAGES from 'constants/lifecycleStages';
import VIOLATION_STATES from 'constants/violationStates';
import { ENFORCEMENT_ACTIONS } from 'constants/enforcementActions';
import dialogues from './dialogues';

import ViolationsPageHeader from './ViolationsPageHeader';
import ViolationsTablePanel from './ViolationsTablePanel';
import ViolationsSidePanel from './SidePanel/ViolationsSidePanel';
import ResolveConfirmation from './Dialogues/ResolveConfirmation';
import ExcludeConfirmation from './Dialogues/ExcludeConfirmation';
import TagConfirmation from './Dialogues/TagConfirmation';

function ViolationsPage({
    history,
    location: { search },
    match: {
        params: { alertId },
    },
}) {
    // Handle changes to applied search options.
    const [isViewFiltered, setIsViewFiltered] = useState(false);

    // Handle changes in the currently selected alert, and checked alerts.
    const [selectedAlertId, setSelectedAlertId] = useState(alertId);
    const [checkedAlertIds, setCheckedAlertIds] = useState([]);

    // Handle changes in the current table page.
    const [currentPage, setCurrentPage] = useState(0);
    const [sortOption, setSortOption] = useState({ field: 'Violation Time', reversed: true });

    // Handle changes in the currently displayed violations.
    const [currentPageAlerts, setCurrentPageAlerts] = useEntitiesByIdsCache();
    const [currentPageAlertsErrorMessage, setCurrentPageAlertsErrorMessage] = useState('');
    const [alertCount, setAlertCount] = useState(0);

    // Handle confirmation dialogue being open.
    const [dialogue, setDialogue] = useState(null);

    // When the selected image changes, update the URL.
    useEffect(() => {
        const urlSuffix = selectedAlertId ? `/${selectedAlertId}` : '';
        history.push({
            pathname: `/main/violations${urlSuffix}`,
            search,
        });
    }, [selectedAlertId, history, search]);

    // We need to be able to identify which alerts are runtime or attempted, and which are not by id.
    const resolvableAlerts = new Set(
        currentPageAlerts
            .filter(
                (alert) =>
                    alert.lifecycleStage === LIFECYCLE_STAGES.RUNTIME ||
                    alert.state === VIOLATION_STATES.ATTEMPTED
            )
            .map((alert) => alert.id)
    );

    const excludableAlerts = currentPageAlerts.filter(
        (alert) =>
            alert.enforcementAction !== ENFORCEMENT_ACTIONS.FAIL_DEPLOYMENT_CREATE_ENFORCEMENT
    );

    const excludableAlertIds = new Set(excludableAlerts.map((alert) => alert.id));

    return (
        <>
            <ViolationsPageHeader
                currentPage={currentPage}
                sortOption={sortOption}
                selectedAlertId={selectedAlertId}
                currentPageAlerts={currentPageAlerts}
                setCurrentPageAlerts={setCurrentPageAlerts}
                setCurrentPageAlertsErrorMessage={setCurrentPageAlertsErrorMessage}
                setSelectedAlertId={setSelectedAlertId}
                setAlertCount={setAlertCount}
                isViewFiltered={isViewFiltered}
                setIsViewFiltered={setIsViewFiltered}
            />
            <PageBody>
                {currentPageAlertsErrorMessage ? (
                    <MessageCentered type="error">{currentPageAlertsErrorMessage}</MessageCentered>
                ) : (
                    <>
                        <div className="flex-shrink-1 overflow-hidden w-full">
                            <ViolationsTablePanel
                                violations={currentPageAlerts}
                                violationsCount={alertCount}
                                isViewFiltered={isViewFiltered}
                                setDialogue={setDialogue}
                                selectedAlertId={selectedAlertId}
                                setSelectedAlertId={setSelectedAlertId}
                                checkedAlertIds={checkedAlertIds}
                                setCheckedAlertIds={setCheckedAlertIds}
                                currentPage={currentPage}
                                setCurrentPage={setCurrentPage}
                                setSortOption={setSortOption}
                                resolvableAlerts={resolvableAlerts}
                                excludableAlertIds={excludableAlertIds}
                            />
                        </div>
                        {selectedAlertId && (
                            <SidePanelAdjacentArea width="2/5">
                                <ViolationsSidePanel
                                    selectedAlertId={selectedAlertId}
                                    setSelectedAlertId={setSelectedAlertId}
                                />
                            </SidePanelAdjacentArea>
                        )}
                    </>
                )}
            </PageBody>
            {dialogue === dialogues.excludeScopes && (
                <ExcludeConfirmation
                    setDialogue={setDialogue}
                    excludableAlerts={excludableAlerts}
                    checkedAlertIds={checkedAlertIds}
                    setCheckedAlertIds={setCheckedAlertIds}
                />
            )}
            {dialogue === dialogues.resolve && (
                <ResolveConfirmation
                    setDialogue={setDialogue}
                    checkedAlertIds={checkedAlertIds}
                    setCheckedAlertIds={setCheckedAlertIds}
                    resolvableAlerts={resolvableAlerts}
                />
            )}
            {dialogue === dialogues.tag && (
                <TagConfirmation
                    setDialogue={setDialogue}
                    checkedAlertIds={checkedAlertIds}
                    setCheckedAlertIds={setCheckedAlertIds}
                />
            )}
        </>
    );
}

ViolationsPage.propTypes = {
    history: ReactRouterPropTypes.history.isRequired,
    location: ReactRouterPropTypes.location.isRequired,
    match: ReactRouterPropTypes.match.isRequired,
};

export default ViolationsPage;
