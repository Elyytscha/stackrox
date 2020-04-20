import React from 'react';
import PropTypes from 'prop-types';
import { useQuery } from 'react-apollo';
import { ArrowLeft } from 'react-feather';

import getPaginatedList from 'utils/getPaginatedList';
import captureGraphQLErrors from 'modules/captureGraphQLErrors';
import queryService from 'modules/queryService';
import Button from 'Components/Button';
import Panel from 'Components/Panel';
import HeaderWithSubText from 'Components/HeaderWithSubText';
import TimelineGraph from 'Components/TimelineGraph';
import Loader from 'Components/Loader';
import { getPod, getContainerEvents } from './getContainerEvents';
import getLargestDifferenceInMilliseconds from '../eventTimelineUtils/getLargestDifferenceInMilliseconds';
import { GET_POD_EVENT_TIMELINE } from '../timelineQueries';

const PodEventTimeline = ({
    id,
    goToNextView,
    goToPreviousView,
    selectedEventType,
    headerComponents,
    currentPage,
    pageSize,
    onPageChange,
}) => {
    const { loading, error, data } = useQuery(GET_POD_EVENT_TIMELINE, {
        variables: {
            podId: id,
            // TODO: We should standardize on using Id vs. ID. Change this once backend makes the change
            containersQuery: queryService.objectToWhereClause({ 'Pod ID': id }),
        },
    });

    captureGraphQLErrors([error]);

    if (loading)
        return (
            <div className="flex flex-1 items-center justify-center py-4">
                <Loader message="Loading Event Timeline..." />
            </div>
        );

    const { name, subText } = getPod(data.pod);

    const headerTextComponent = (
        <>
            <Button
                dataTestId="timeline-back-button"
                className="border-base-300 border-r px-3 hover:bg-base-200"
                icon={<ArrowLeft className="h-4 w-4 text-base-600" />}
                onClick={goToPreviousView}
            />
            <HeaderWithSubText header={name} subText={subText} />
        </>
    );

    // Adding pagination for Grouped Container Instances required a substantial amount of work, so we're going with pagination on the frontend for now
    const paginatedContainers = getPaginatedList(data.containers, currentPage, pageSize);
    const timelineData = getContainerEvents(paginatedContainers, selectedEventType);
    const absoluteMaxTimeRange = getLargestDifferenceInMilliseconds(timelineData);

    const numTotalContainers = data?.pod?.containerCount || 0;

    return (
        <Panel headerTextComponent={headerTextComponent} headerComponents={headerComponents}>
            <TimelineGraph
                data={timelineData}
                goToNextView={goToNextView}
                currentPage={currentPage}
                totalSize={numTotalContainers}
                pageSize={pageSize}
                onPageChange={onPageChange}
                absoluteMaxTimeRange={absoluteMaxTimeRange}
            />
        </Panel>
    );
};

PodEventTimeline.propTypes = {
    id: PropTypes.string.isRequired,
    goToNextView: PropTypes.func.isRequired,
    goToPreviousView: PropTypes.func.isRequired,
    selectedEventType: PropTypes.string.isRequired,
    headerComponents: PropTypes.node.isRequired,
    currentPage: PropTypes.number.isRequired,
    pageSize: PropTypes.number.isRequired,
    onPageChange: PropTypes.func.isRequired,
};

export default PodEventTimeline;
