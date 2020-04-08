import React, { useState, useEffect } from 'react';
import { connect } from 'react-redux';
import { withRouter, NavLink as Link } from 'react-router-dom';
import ReactRouterPropTypes from 'react-router-prop-types';
import { createStructuredSelector } from 'reselect';
import find from 'lodash/find';
import PropTypes from 'prop-types';
import { selectors } from 'reducers';

import { useTheme } from 'Containers/ThemeProvider';
import NavigationPanel from './NavigationPanel';
import ApiDocsNavigation from './ApiDocsNavigation';
import LeftSideNavLinks, { navLinks } from './LeftSideNavLinks';
import { getDarkModeLinkClassName } from './navHelpers';

const versionString = metadata => {
    let result = `v${metadata.version}`;
    if (metadata.releaseBuild === false) {
        result += ' [DEV BUILD]';
    }
    return result;
};

const LeftNavigation = ({ location, metadata }) => {
    const { isDarkMode } = useTheme();

    const [panelType, setPanelType] = useState(null);
    const [clickOnPanelItem, setClickOnPanelItem] = useState(false);
    const [selectedPanel, setSelectedPanel] = useState('');

    const linkClassName = `flex font-condensed leading-normal font-700 text-primary-400 px-3 no-underline h-18 items-center border-b ${getDarkModeLinkClassName(
        isDarkMode
    )}`;

    function getActiveClassName(navLink) {
        const { pathname } = location;
        const navText = navLink.text.toLowerCase();
        const baseActiveClass = isDarkMode
            ? 'text-primary-500 bg-primary-200 hover:bg-primary-200'
            : 'bg-primary-700 hover:bg-primary-700 text-base-100';

        if (
            (pathname.includes('main/policies') ||
                pathname.includes('main/integrations') ||
                pathname.includes('main/access')) &&
            navText === 'platform configuration'
        ) {
            return baseActiveClass;
        }

        if (navLink.to !== '') {
            return baseActiveClass;
        }
        if (navLink.to === '') {
            const baseFocusClass = isDarkMode
                ? 'text-primary-500 bg-primary-200 hover:bg-primary-300'
                : 'text-base-100 bg-base-800 hover:bg-base-800';
            if (panelType && panelType === navLink.panelType) {
                return baseFocusClass;
            }
            if (!panelType && clickOnPanelItem && selectedPanel === navText) {
                return baseFocusClass;
            }
            return isDarkMode ? 'bg-base-100' : 'bg-primary-800';
        }
        return '';
    }

    function closePanel(newClickOnPanelItem, newSelectedPanel) {
        return () => {
            if (newClickOnPanelItem) {
                setClickOnPanelItem(newClickOnPanelItem);
                setSelectedPanel(newSelectedPanel);
            }
            setPanelType(null);
        };
    }

    function showNavigationPanel(navLink) {
        return e => {
            if (navLink.panelType && panelType !== navLink.panelType) {
                e.preventDefault();
                setPanelType(navLink.panelType);
            } else {
                if (panelType === navLink.panelType) {
                    e.preventDefault();
                }
                setPanelType(null);
                setClickOnPanelItem(false);
            }
        };
    }

    function renderLink(navLink) {
        return (
            <Link
                to={navLink.to}
                activeClassName={getActiveClassName(navLink)}
                onClick={showNavigationPanel(navLink)}
                className={linkClassName}
                data-test-id={navLink.data || navLink.text}
            >
                <div className="text-center pr-2">{navLink.renderIcon()}</div>
                <div className={`${isDarkMode ? 'text-base-600' : 'text-base-100'}`}>
                    {navLink.text}
                </div>
            </Link>
        );
    }

    function renderNavigationPanel() {
        if (!panelType) return '';
        return <NavigationPanel panelType={panelType} onClose={closePanel} />;
    }

    function componentDidMount() {
        window.onpopstate = e => {
            const url = e.srcElement.location.pathname;
            const link = find(navLinks, navLink => url === navLink.to);
            if (panelType || link) {
                setPanelType(null);
            }
        };
    }

    useEffect(componentDidMount, []);

    const darkModeClasses = isDarkMode
        ? 'bg-base-100 border-t -mt-px border-r border-base-300'
        : 'bg-primary-800';
    return (
        <>
            <div
                className={`flex flex-col justify-between flex-none overflow-auto z-60 ignore-react-onclickoutside ${darkModeClasses}`}
            >
                <nav className="left-navigation">
                    <LeftSideNavLinks renderLink={renderLink} />
                </nav>
                <div
                    className="flex flex-col h-full justify-end text-center text-xs font-700"
                    data-test-id="nav-footer"
                >
                    <ApiDocsNavigation onClick={closePanel()} />
                    <span
                        className={`left-navigation p-3 leading-normal ${
                            isDarkMode ? 'text-base-500' : 'text-primary-400'
                        } word-break-all`}
                    >
                        {versionString(metadata)}
                    </span>
                </div>
            </div>
            {renderNavigationPanel()}
        </>
    );
};

LeftNavigation.propTypes = {
    location: ReactRouterPropTypes.location.isRequired,
    metadata: PropTypes.shape({
        version: PropTypes.string,
        releaseBuild: PropTypes.bool,
        licenseStatus: PropTypes.string
    })
};

LeftNavigation.defaultProps = {
    metadata: {
        version: 'latest',
        releaseBuild: false
    }
};

const mapStateToProps = createStructuredSelector({
    metadata: selectors.getMetadata
});

export default withRouter(connect(mapStateToProps)(LeftNavigation));
