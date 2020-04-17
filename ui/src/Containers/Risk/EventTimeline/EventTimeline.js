import React, { useState } from 'react';
import PropTypes from 'prop-types';

import { eventTypes, rootTypes } from 'constants/timelineTypes';
import NotFoundMessage from 'Components/NotFoundMessage';
import TimelineLegend from 'Components/TimelineLegend';
import EventTypeSelect from './EventTypeSelect';
import DeploymentEventTimeline from './DeploymentEventTimeline';
import PodEventTimeline from './PodEventTimeline';

const PAGE_SIZE = 10;

const EventTimelineComponentMap = {
    [rootTypes.DEPLOYMENT]: DeploymentEventTimeline,
    [rootTypes.POD]: PodEventTimeline
};

const EventTimeline = ({ deploymentId }) => {
    const rootView = [
        {
            type: rootTypes.DEPLOYMENT,
            id: deploymentId
        }
    ];
    const [currentPage, setPage] = useState(1);
    const [selectedEventType, selectEventType] = useState(eventTypes.ALL);
    const [view, setView] = useState(rootView);

    function getCurrentView() {
        return view[view.length - 1];
    }

    function resetSelectedEventType() {
        selectEventType(eventTypes.ALL);
    }

    function goToRootView() {
        setView(rootView);
        resetSelectedEventType();
    }

    function goToNextView(type, id) {
        const newView = [...view, { type, id }];
        setView(newView);
        resetSelectedEventType();
    }

    function goToPreviousView() {
        if (view.length <= 1) return;
        setView(view.slice(0, -1));
        resetSelectedEventType();
    }

    const currentView = getCurrentView();

    const headerComponents = (
        <>
            <EventTypeSelect
                selectedEventType={selectedEventType}
                selectEventType={selectEventType}
            />
            <div className="ml-3">
                <TimelineLegend />
            </div>
        </>
    );

    const Component = EventTimelineComponentMap[currentView.type];
    if (!Component)
        return (
            <NotFoundMessage
                message="The Event Timeline for this view was not found."
                actionText="Go back"
                onClick={goToRootView}
            />
        );
    return (
        <Component
            id={currentView.id}
            goToNextView={goToNextView}
            goToPreviousView={goToPreviousView}
            selectedEventType={selectedEventType}
            headerComponents={headerComponents}
            currentPage={currentPage}
            pageSize={PAGE_SIZE}
            onPageChange={setPage}
        />
    );
};

EventTimeline.propTypes = {
    deploymentId: PropTypes.string.isRequired
};

export default EventTimeline;
