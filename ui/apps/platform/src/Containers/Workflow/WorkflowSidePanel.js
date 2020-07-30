import React from 'react';
import { withRouter, Link } from 'react-router-dom';
import parseURL from 'utils/URLParser';
import onClickOutside from 'react-onclickoutside';
import { useTheme } from 'Containers/ThemeProvider';
import workflowStateContext from 'Containers/workflowStateContext';
import { ExternalLink as ExternalLinkIcon } from 'react-feather';
import Panel from 'Components/Panel';
import SidePanelAnimatedDiv from 'Components/animations/SidePanelAnimatedDiv';
import EntityBreadCrumbs from 'Containers/BreadCrumbs/EntityBreadCrumbs';

const WorkflowSidePanel = ({ history, location, children, isOpen }) => {
    const { isDarkMode } = useTheme();
    const workflowState = parseURL(location);
    const pageStack = workflowState.getPageStack();
    const breadCrumbEntities = workflowState.stateStack.slice(pageStack.length);

    function onClose() {
        const url = workflowState.removeSidePanelParams().toUrl();
        history.push(url);
    }

    WorkflowSidePanel.handleClickOutside = () => {
        const btn = document.getElementById('panel-close-button');
        if (btn) btn.click();
    };

    const url = workflowState.getSkimmedStack().toUrl();
    const externalLink = (
        <div className="flex items-center h-full hover:bg-base-300">
            <Link
                to={url}
                data-testid="external-link"
                className={`${
                    !isDarkMode ? 'border-base-100' : 'border-base-400'
                } border-l h-full p-4`}
            >
                <ExternalLinkIcon className="h-6 w-6 text-base-600" />
            </Link>
        </div>
    );

    return (
        <workflowStateContext.Provider value={workflowState}>
            <SidePanelAnimatedDiv isOpen={isOpen}>
                <div
                    className={`w-full h-full bg-base-100 rounded-tl-lg shadow-sidepanel ${
                        !isDarkMode ? '' : 'border-l border-base-400'
                    }`}
                >
                    <Panel
                        id="side-panel"
                        headerClassName={`flex w-full h-14 rounded-tl-lg overflow-y-hidden border-b ${
                            !isDarkMode
                                ? 'bg-side-panel-wave border-base-100'
                                : 'bg-primary-200 border-primary-400'
                        }`}
                        bodyClassName={isDarkMode ? 'bg-base-0' : 'bg-base-100'}
                        headerTextComponent={
                            <EntityBreadCrumbs workflowEntities={breadCrumbEntities} />
                        }
                        headerComponents={externalLink}
                        onClose={onClose}
                        closeButtonClassName={
                            isDarkMode ? 'border-l border-base-400' : 'border-l border-base-100'
                        }
                    >
                        {children}
                    </Panel>
                </div>
            </SidePanelAnimatedDiv>
        </workflowStateContext.Provider>
    );
};

const clickOutsideConfig = {
    handleClickOutside: () => WorkflowSidePanel.handleClickOutside,
};

/*
 * If more than one SidePanel is rendered, this Pure Functional Component will need to be converted to
 * a Class Component in order to work correctly. See https://github.com/stackrox/rox/pull/3090#pullrequestreview-274948849
 */
export default onClickOutside(withRouter(WorkflowSidePanel), clickOutsideConfig);
