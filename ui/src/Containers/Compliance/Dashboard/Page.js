import React from 'react';
import { resourceTypes, standardTypes } from 'constants/entityTypes';
import { useTheme } from 'Containers/ThemeProvider';

import StandardsByEntity from 'Containers/Compliance/widgets/StandardsByEntity';
import StandardsAcrossEntity from 'Containers/Compliance/widgets/StandardsAcrossEntity';
import ComplianceByStandard from 'Containers/Compliance/widgets/ComplianceByStandard';
import WaveBackground from 'images/wave-bg.svg';
import WaveBackground2 from 'images/wave-bg-2.svg';
import FeatureEnabled from 'Containers/FeatureEnabled';
import { knownBackendFlags } from 'utils/featureFlags';
import DashboardHeader from './Header';

const ComplianceDashboardPage = () => {
    const { isDarkMode } = useTheme();
    const bgStyle = isDarkMode
        ? {}
        : {
              boxShadow: 'hsla(230, 75%, 63%, 0.62) 0px 5px 30px 0px',
              '--start': 'hsl(226, 70%, 60%)',
              '--end': 'hsl(226, 64%, 56%)'
          };

    return (
        <section>
            <DashboardHeader
                classes={`bg-gradient-horizontal z-10 sticky top-0 ${
                    isDarkMode ? 'text-base-600' : 'text-base-100'
                }`}
                bgStyle={bgStyle}
            />
            <img
                className="absolute left-0 pointer-events-none z-10 w-full"
                id="wave-bg2"
                src={WaveBackground2}
                style={{ mixBlendMode: 'lighten', top: '-60px' }}
                alt="Waves"
            />
            <div
                className="flex-1 relative bg-gradient-diagonal p-6 xxxl:p-8"
                style={{ '--start': 'var(--base-200)', '--end': 'var(--primary-100)' }}
                id="capture-dashboard"
            >
                <img
                    className="absolute left-0 pointer-events-none w-full"
                    id="wave-bg"
                    src={WaveBackground}
                    style={{ top: '-130px' }}
                    alt="Wave"
                />
                <div
                    className="grid grid-gap-6 xxxl:grid-gap-8 md:grid-auto-fit xxl:grid-auto-fit-wide md:grid-dense"
                    style={{ '--min-tile-height': '160px' }}
                >
                    <StandardsAcrossEntity
                        entityType={resourceTypes.CLUSTER}
                        bodyClassName="pr-4 py-1"
                        className="pdf-page"
                    />
                    <StandardsByEntity
                        entityType={resourceTypes.CLUSTER}
                        bodyClassName="p-4"
                        className="pdf-page"
                    />
                    <StandardsAcrossEntity
                        entityType={resourceTypes.NAMESPACE}
                        bodyClassName="px-4 pt-1"
                        className="pdf-page"
                    />
                    <StandardsAcrossEntity
                        entityType={resourceTypes.NODE}
                        bodyClassName="pr-4 py-1"
                        className="pdf-page"
                    />
                    <ComplianceByStandard
                        standardType={standardTypes.CIS_Docker_v1_2_0}
                        className="pdf-page"
                    />
                    <ComplianceByStandard
                        standardType={standardTypes.CIS_Kubernetes_v1_5}
                        className="pdf-page"
                    />
                    <ComplianceByStandard
                        standardType={standardTypes.HIPAA_164}
                        className="pdf-page"
                    />
                    <ComplianceByStandard
                        standardType={standardTypes.NIST_800_190}
                        className="pdf-page"
                    />
                    <FeatureEnabled featureFlag={knownBackendFlags.ROX_NIST_800_53}>
                        <ComplianceByStandard
                            standardType={standardTypes.NIST_SP_800_53_Rev_4}
                            className="pdf-page"
                        />
                    </FeatureEnabled>
                    <ComplianceByStandard
                        standardType={standardTypes.PCI_DSS_3_2}
                        className="pdf-page"
                    />
                </div>
            </div>
        </section>
    );
};
export default ComplianceDashboardPage;
