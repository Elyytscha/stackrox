import React, { useState, useEffect } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { selectors } from 'reducers';
import { actions, LICENSE_STATUS } from 'reducers/license';
import { getLicenseStatusMessage } from 'Containers/License/helpers';

import UploadButton from 'Components/UploadButton';
import Dialog from 'Components/Dialog';

const UploadLicense = ({ licenseUploadStatus, activateLicense, isStartUpScreen }) => {
    const defaultDialogState = licenseUploadStatus ? !!licenseUploadStatus.status : false;
    const defaultVerifyingLicenseState = licenseUploadStatus
        ? licenseUploadStatus.status === LICENSE_STATUS.VERIFYING
        : false;

    const status = licenseUploadStatus ? licenseUploadStatus.status : null;
    const message = licenseUploadStatus ? licenseUploadStatus.message : null;

    const [dialogMessage, setDialogMessage] = useState(getLicenseStatusMessage(status, message));
    const [isDialogOpen, openDialog] = useState(defaultDialogState);
    const [isVerifyingLicense, verifyLicense] = useState(defaultVerifyingLicenseState);

    useEffect(() => {
        if (licenseUploadStatus && licenseUploadStatus.status !== LICENSE_STATUS.VERIFYING) {
            setDialogMessage(getLicenseStatusMessage(status, message));
            verifyLicense(false);
        }
    }, [licenseUploadStatus, message, status]);

    function onUploadHandler(data) {
        verifyLicense(true);
        openDialog(true);
        activateLicense(data);
    }

    function onDialogCancel() {
        openDialog(false);
    }

    return (
        <>
            <UploadButton
                className="p-3 px-6 rounded-sm bg-primary-600 hover:bg-primary-700 text-base-100 uppercase text-center tracking-wide mt-4"
                text="Upload New License Key"
                onChange={onUploadHandler}
            />
            {!isStartUpScreen && (
                <Dialog
                    isOpen={isDialogOpen}
                    text={dialogMessage ? dialogMessage.text : ''}
                    cancelText="Ok"
                    onCancel={onDialogCancel}
                    isLoading={isVerifyingLicense}
                    loadingText="Verifying License Key"
                />
            )}
        </>
    );
};

UploadLicense.propTypes = {
    licenseUploadStatus: PropTypes.shape({
        status: PropTypes.string,
        message: PropTypes.string,
    }),
    activateLicense: PropTypes.func.isRequired,
    isStartUpScreen: PropTypes.bool,
};

UploadLicense.defaultProps = {
    licenseUploadStatus: {},
    isStartUpScreen: false,
};
const mapStateToProps = createStructuredSelector({
    licenseUploadStatus: selectors.getLicenseUploadStatus,
});
const mapDispatchToProps = {
    activateLicense: actions.activateLicense,
    isStartUpScreen: false,
};

export default connect(mapStateToProps, mapDispatchToProps)(UploadLicense);
