import React from 'react';
import { defaultHeaderClassName, defaultColumnClassName, wrapClassName } from 'Components/Table';
import { resourceTypes } from 'constants/entityTypes';
import { standardLabels } from 'messages/standards';
import { sortVersion } from 'sorters/sorters';
import LabelChip from 'Components/LabelChip';
import { format } from 'date-fns';
import dateTimeFormat from 'constants/dateTimeFormat';

const getNameCell = name => <div data-test-id="table-row-name">{name}</div>;

const controlColumns = [
    {
        accessor: 'id',
        Header: 'id',
        headerClassName: 'hidden',
        className: 'hidden'
    },
    {
        sortMethod: sortVersion,
        Header: 'Standard',
        headerClassName: `w-1/5 ${defaultHeaderClassName}`,
        className: `w-1/5 ${defaultColumnClassName}`,
        accessor: 'control.standardId',
        Cell: ({ original }) => standardLabels[original.control.standardId]
    },
    {
        accessor: 'control',
        sortMethod: sortVersion,
        Header: 'Control',
        headerClassName: `w-1/2 ${defaultHeaderClassName}`,
        className: `w-1/2 ${defaultColumnClassName}`,
        Cell: ({ original }) =>
            getNameCell(`${original.control.name} - ${original.control.description}`)
    },
    {
        accessor: 'value.overallState',
        Header: 'State',
        headerClassName: `w-1/5 ${defaultHeaderClassName}`,
        className: `w-1/5 ${defaultColumnClassName}`,
        // eslint-disable-next-line react/prop-types
        Cell: ({ original }) => {
            const text = original.value.overallState === 'COMPLIANCE_STATE_FAILURE' && 'Fail';
            return <LabelChip text={text} type="alert" />;
        }
    }
];

const nodesAcrossControlsColumns = [
    {
        Header: 'Id',
        headerClassName: 'hidden',
        className: 'hidden',
        accessor: 'id'
    },
    {
        Header: `Node`,
        headerClassName: `w-1/3 ${defaultHeaderClassName}`,
        className: `w-1/3 ${defaultColumnClassName}`,
        accessor: 'name'
    },
    {
        Header: `Cluster`,
        headerClassName: `w-1/3 ${defaultHeaderClassName}`,
        className: `w-1/3 ${defaultColumnClassName}`,
        accessor: 'clusterName'
    },
    {
        Header: `Control Status`,
        headerClassName: `w-1/8 ${defaultHeaderClassName}`,
        className: `w-1/8 ${defaultColumnClassName}`,
        // eslint-disable-next-line
        Cell: ({ original }) => {
            return !original.passing ? <LabelChip text="Fail" type="alert" /> : 'Pass';
        }
    }
];

const imageColumns = [
    {
        expander: true,
        headerClassName: `w-1/8 ${defaultHeaderClassName} pointer-events-none`,
        className: 'w-1/8 pointer-events-none flex items-center justify-end',
        // eslint-disable-next-line react/prop-types
        Expander: ({ isExpanded, ...rest }) => {
            if (rest.original.components.length === 0) return '';
            const className = 'rt-expander w-1 pt-2 pointer-events-auto';
            return <div className={`${className} ${isExpanded ? '-open' : ''}`} />;
        }
    },
    {
        accessor: 'instruction',
        Header: 'Instruction',
        headerClassName: `text-left ${wrapClassName} ${defaultHeaderClassName}`,
        className: `text-left pl-3 ${wrapClassName} ${defaultColumnClassName}`
    },
    {
        accessor: 'value',
        Header: 'Value',
        headerClassName: `w-3/5 text-left ${wrapClassName} ${defaultHeaderClassName}`,
        className: `w-3/5 text-left pl-3 word-break-all ${wrapClassName} ${defaultColumnClassName}`
    },
    {
        accessor: 'created',
        Header: 'Created',
        align: 'right',
        widthClassName: `text-left pr-3 ${wrapClassName} ${defaultHeaderClassName}`,
        className: `text-left pr-3 ${wrapClassName} ${defaultColumnClassName}`,
        Cell: ({ original }) => format(original.created, dateTimeFormat)
    },
    {
        accessor: 'components.length',
        Header: 'Components',
        headerClassName: `text-left ${wrapClassName} ${defaultHeaderClassName}`,
        className: `text-left pl-3 word-break-all ${wrapClassName} ${defaultColumnClassName}`
    },
    {
        accessor: 'cvesCount',
        Header: 'CVEs',
        headerClassName: `text-left ${wrapClassName} ${defaultHeaderClassName}`,
        className: `text-left pl-3 word-break-all ${wrapClassName} ${defaultColumnClassName}`
    }
];

const deploymentViolationsColumns = [
    {
        Header: 'Id',
        headerClassName: 'hidden',
        className: 'hidden',
        accessor: 'deployment.id'
    },
    {
        Header: `Deployment`,
        headerClassName: `w-1/8 ${defaultHeaderClassName}`,
        className: `w-1/8 ${defaultColumnClassName}`,
        accessor: 'deployment.name'
    },
    {
        Header: `Time`,
        headerClassName: `w-1/8 ${defaultHeaderClassName}`,
        className: `w-1/8 ${defaultColumnClassName}`,
        accessor: 'time',
        Cell: ({ original }) => {
            const { time } = original;
            return format(time, dateTimeFormat);
        }
    }
];

export const entityToColumns = {
    [resourceTypes.CONTROL]: controlColumns,
    [resourceTypes.IMAGE]: imageColumns
};

export const entityAcrossControlsColumns = {
    [resourceTypes.NODE]: nodesAcrossControlsColumns
};

export const entityViolationsColumns = {
    [resourceTypes.DEPLOYMENT]: deploymentViolationsColumns
};
