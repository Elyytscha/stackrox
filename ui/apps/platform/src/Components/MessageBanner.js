import React, { useState } from 'react';
import PropTypes from 'prop-types';
import * as Icon from 'react-feather';

const className = 'w-full flex items-center justify-center leading-normal p-3';
const messageClasses = {
    warn: `${className} bg-warning-300 text-warning-800`,
    error: `${className} bg-alert-300 text-alert-800`,
    info: `${className} bg-tertiary-300 text-tertiary-800`,
};

function MessageBanner({ component, message, type, showCancel, onCancel, dataTestId }) {
    const [isBannerShowing, showBanner] = useState(true);
    function onClickHandler() {
        showBanner(false);
        if (onCancel) onCancel();
    }
    return (
        isBannerShowing && (
            <div data-testid={dataTestId} className={messageClasses[type]}>
                <div className="flex flex-1 justify-center">{component || message}</div>
                {showCancel && (
                    <Icon.X
                        data-testid={dataTestId ? `${dataTestId}-cancel` : null}
                        className="h-6 w-6 cursor-pointer"
                        onClick={onClickHandler}
                    />
                )}
            </div>
        )
    );
}

MessageBanner.defaultProps = {
    component: null,
    message: null,
    type: 'info',
    showCancel: false,
    onCancel: null,
    dataTestId: null,
};

MessageBanner.propTypes = {
    component: PropTypes.element,
    message: PropTypes.string,
    type: PropTypes.oneOf(['warn', 'error', 'info']),
    showCancel: PropTypes.bool,
    onCancel: PropTypes.func,
    dataTestId: PropTypes.string,
};

export default MessageBanner;
