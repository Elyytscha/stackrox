import React from 'react';
import { useTheme } from 'Containers/ThemeProvider';
import Tooltip from 'Components/Tooltip';
import TooltipOverlay from 'Components/TooltipOverlay';

import { Moon, Sun } from 'react-feather';

const ThemeToggleButton = () => {
    const themeState = useTheme();
    const tooltipText = themeState.isDarkMode ? 'Switch to Light Mode' : 'Switch to Dark Mode';
    return (
        <Tooltip content={<TooltipOverlay>{tooltipText}</TooltipOverlay>}>
            <button
                title="Invert theme"
                onClick={themeState.toggle}
                type="button"
                className="flex h-full items-center border-l border-base-400 border-r-0 pt-3 pb-2 px-4 h-9 hover:bg-base-200 text-base-600 no-underline rounded-l-sm"
            >
                <span className="uppercase text-sm lg:text-base font-700 tracking-wide leading-relaxed flex flex-col">
                    {themeState.isDarkMode ? <Sun size="16" /> : <Moon size="16" />}
                </span>
            </button>
        </Tooltip>
    );
};

export default ThemeToggleButton;
