import { all, fork } from 'redux-saga/effects';

import alerts from './alertSagas';
import apiTokens from './apiTokenSagas';
import authProviders from './authSagas';
import clusterInitBundles from './clusterInitBundleSagas';
import clusters from './clusterSagas';
import deployments from './deploymentSagas';
import featureFlags from './featureFlagSagas';
import images from './imageSagas';
import policies from './policiesSagas';
import integrations from './integrationSagas';
import globalSearch from './globalSearchSagas';
import roles from './roleSagas';
import searches from './searchSagas';
import searchAutoComplete from './searchAutocompleteSagas';
import secrets from './secretSagas';
import network from './networkSagas';
import metadata from './metadataSagas';
import processes from './processSagas';
import groups from './groupSagas';
import attributes from './attributesSagas';
import cli from './cliSagas';
import license from './licenseSagas';
import systemConfig from './systemConfig';
import telemetryConfig from './telemetryConfig';

export default function* root() {
    yield all([
        fork(license),
        fork(alerts),
        fork(apiTokens),
        fork(authProviders),
        fork(clusterInitBundles),
        fork(cli),
        fork(clusters),
        fork(deployments),
        fork(images),
        fork(policies),
        fork(featureFlags),
        fork(integrations),
        fork(globalSearch),
        fork(roles),
        fork(searches),
        fork(searchAutoComplete),
        fork(secrets),
        fork(network),
        fork(metadata),
        fork(processes),
        fork(groups),
        fork(attributes),
        fork(systemConfig),
        fork(telemetryConfig),
    ]);
}
