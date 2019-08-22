import React from 'react';
import { ClipLoader } from 'react-spinners';
import PropTypes from 'prop-types';

const Loader = ({ message, transparent }) => (
    <div
        className={`flex flex-col items-center justify-center min-h-full w-full ${
            transparent ? '' : 'bg-base-100'
        }`}
    >
        <ClipLoader loading size={14} color="currentColor" />
        {message && <div className="text-lg font-sans font-600 tracking-wide mt-4">{message}</div>}
    </div>
);

Loader.propTypes = {
    message: PropTypes.string,
    transparent: PropTypes.bool
};

Loader.defaultProps = {
    message: 'Loading...',
    transparent: false
};

export default Loader;
