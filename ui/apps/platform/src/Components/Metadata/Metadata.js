import React from 'react';
import PropTypes from 'prop-types';

import MetadataStatsList from 'Components/MetadataStatsList';
import Widget from 'Components/Widget';
import ResourceCountPopper from 'Components/ResourceCountPopper';
import { useTheme } from 'Containers/ThemeProvider';

const renderName = (data) => {
    return data.map(({ name }) => (
        <div className="mt-2" key={name}>
            {name}
        </div>
    ));
};

const Metadata = ({
    keyValuePairs,
    title,
    statTiles,
    labels,
    annotations,
    whitelists,
    secrets,
    bgClass,
    className,
    ...rest
}) => {
    const { isDarkMode } = useTheme();
    const keyValueList = keyValuePairs.map(({ key, value }) => (
        <li className="border-b border-base-300 py-3" key={key}>
            <span className="text-base-600 font-700 mr-2">{key}:</span>
            <span className="font-600" data-testid={`${key}-value`}>
                {value}
            </span>
        </li>
    ));

    const keyValueClasses = `flex-1 last:border-0 border-base-300 overflow-hidden px-3 ${
        labels || annotations || whitelists || secrets ? ' border-r' : ''
    }`;

    const widgetClassName = `${className} ${!isDarkMode && bgClass ? 'bg-counts-widget' : ''}`;

    return (
        <Widget header={title} className={widgetClassName} {...rest}>
            <div className="flex flex-col w-full">
                {statTiles && statTiles.length > 0 && <MetadataStatsList statTiles={statTiles} />}
                <div className="flex w-full h-full">
                    <ul className={keyValueClasses}>{keyValueList}</ul>
                    <ul>
                        {labels && (
                            <li className="m-4">
                                <ResourceCountPopper
                                    data={labels}
                                    reactOutsideClassName="ignore-react-onclickoutside"
                                    label="Label"
                                />
                            </li>
                        )}
                        {annotations && (
                            <li className="m-4">
                                <ResourceCountPopper
                                    data={annotations}
                                    reactOutsideClassName="ignore-react-onclickoutside"
                                    label="Annotation"
                                />
                            </li>
                        )}
                        {whitelists && (
                            <li className="m-4">
                                <ResourceCountPopper
                                    data={whitelists}
                                    reactOutsideClassName="ignore-react-onclickoutside"
                                    label="Excluded Scopes"
                                    renderContent={renderName}
                                />
                            </li>
                        )}
                        {secrets && (
                            <li className="m-4">
                                <ResourceCountPopper
                                    data={secrets}
                                    reactOutsideClassName="ignore-react-onclickoutside"
                                    label="Image Pull Secret"
                                    renderContent={renderName}
                                />
                            </li>
                        )}
                    </ul>
                </div>
            </div>
        </Widget>
    );
};

Metadata.propTypes = {
    keyValuePairs: PropTypes.arrayOf(
        PropTypes.shape({
            key: PropTypes.string.isRequired,
            value: PropTypes.oneOfType([PropTypes.string, PropTypes.element, PropTypes.number]),
        })
    ).isRequired,
    title: PropTypes.string,
    statTiles: PropTypes.arrayOf(PropTypes.node),
    labels: PropTypes.arrayOf(PropTypes.shape({})),
    annotations: PropTypes.arrayOf(PropTypes.shape({})),
    whitelists: PropTypes.arrayOf(PropTypes.shape({})),
    secrets: PropTypes.arrayOf(PropTypes.shape({})),
    bgClass: PropTypes.bool,
    className: PropTypes.string,
};

Metadata.defaultProps = {
    title: 'Metadata',
    statTiles: null,
    labels: null,
    annotations: null,
    whitelists: null,
    secrets: null,
    bgClass: false,
    className: '',
};

export default Metadata;
