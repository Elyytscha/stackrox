import React from 'react';
import PropTypes from 'prop-types';
import { format, distanceInWordsStrict } from 'date-fns';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { selectors } from 'reducers';
import { getHasReadWritePermission } from 'reducers/roles';
import {
    createExpirationMessageWithoutLink,
    getExpirationMessageType,
} from 'Containers/License/helpers';

import * as Icon from 'react-feather';
import Widget from 'Components/Widget';
import Message from 'Components/Message';
import UploadLicense from 'Containers/License/UploadLicense';

const LicenseExpiration = ({ expirationDate, userRolePermissions }) => {
    const canUploadLicense = getHasReadWritePermission('Licenses', userRolePermissions);
    const expirationMessage = createExpirationMessageWithoutLink(expirationDate);

    const showTopTimeRemaining = getExpirationMessageType(expirationDate) === 'info';

    return (
        <Widget header="License Expiration">
            <div className="py-4 px-6 w-full">
                <div className="flex items-center text-lg pb-4 border-b border-base-300">
                    <Icon.Clock className="h-5 w-5 text-primary-800 text-4xl mr-4" />
                    <div className="text-primary-800 font-400 text-4xl">
                        {format(expirationDate, 'MM/DD/YY')}
                    </div>
                    {showTopTimeRemaining && (
                        <div className="flex flex-1 justify-end text-base-500">
                            ({distanceInWordsStrict(expirationDate, new Date())} from now)
                        </div>
                    )}
                </div>
                <div className="text-center">
                    {expirationMessage && (
                        <Message
                            type={expirationMessage.type}
                            message={expirationMessage.message}
                        />
                    )}
                    {canUploadLicense && <UploadLicense />}
                </div>
            </div>
        </Widget>
    );
};

LicenseExpiration.propTypes = {
    expirationDate: PropTypes.string,
    userRolePermissions: PropTypes.shape({ globalAccess: PropTypes.string.isRequired }),
};

LicenseExpiration.defaultProps = {
    expirationDate: null,
    userRolePermissions: null,
};

const mapStateToProps = createStructuredSelector({
    expirationDate: selectors.getLicenseExpirationDate,
    userRolePermissions: selectors.getUserRolePermissions,
});

export default connect(mapStateToProps, null)(LicenseExpiration);
