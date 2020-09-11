import {
    NS_FONT_SIZE,
    TEXT_MAX_WIDTH,
    NODE_WIDTH,
    NODE_SOLID_BORDER_WIDTH,
    COLORS,
} from 'constants/networkGraph';

const deploymentStyle = {
    width: NODE_WIDTH,
    height: NODE_WIDTH,
    label: 'data(name)',
    'font-size': '6px',
    'text-max-width': TEXT_MAX_WIDTH,
    'text-wrap': 'ellipsis',
    'text-margin-y': '5px',
    'text-valign': 'bottom',
    'font-weight': 'bold',
    'font-family': 'Open Sans',
    'min-zoomed-font-size': '20px',
    color: COLORS.label,
    'z-compound-depth': 'top',
};

// Note: there is no specificity in cytoscape style
// the order of the styles in this array matters
const styles = [
    {
        selector: ':parent',
        style: {
            'background-color': '#fff',
            'border-width': '1.5px',
            'border-color': COLORS.inactiveNS,
            shape: 'roundrectangle',
            'compound-sizing-wrt-labels': 'exclude',
            'font-family': 'stackrox, Open Sans',
            'text-margin-y': '8px',
            'text-valign': 'bottom',
            'font-size': NS_FONT_SIZE,
            color: COLORS.label,
            'font-weight': 700,
            label: 'data(name)',
            padding: '0px',
            'text-transform': 'uppercase',
            'z-compound-depth': 'auto',
        },
    },
    {
        selector: 'node.nsHovered',
        style: {
            opacity: 1,
            'border-style': 'solid',
            'border-color': COLORS.hovered,
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 70%, 1)',
            'overlay-opacity': 0.05,
            'z-compound-depth': 'auto',
        },
    },
    {
        selector: 'node.nsSelected',
        style: {
            opacity: 1,
            'border-style': 'solid',
            'border-color': COLORS.selected,
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 60%, 1)',
            'overlay-opacity': 0.05,
            'z-compound-depth': 'auto',
        },
    },
    {
        selector: 'node.nsActive',
        style: {
            'border-style': 'dashed',
            'border-color': COLORS.active,
        },
    },
    {
        selector: 'node.nsActive.nsHovered',
        style: {
            opacity: 1,
            'border-style': 'dashed',
            'border-color': COLORS.hoveredActive,
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 60%, 1)',
            'overlay-opacity': 0.1,
            'z-compound-depth': 'auto',
        },
    },
    {
        selector: 'node.nsActive.nsSelected',
        style: {
            opacity: 1,
            'border-style': 'dashed',
            'border-color': COLORS.selectedActive,
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 50%, 1)',
            'overlay-opacity': 0.1,
            'z-compound-depth': 'auto',
        },
    },
    {
        selector: ':parent > node.deployment',
        style: {
            'background-color': COLORS.inactive,
            ...deploymentStyle,
        },
    },
    {
        selector: 'node.externallyConnected',
        style: {
            width: NODE_WIDTH + NODE_SOLID_BORDER_WIDTH,
            height: NODE_WIDTH + NODE_SOLID_BORDER_WIDTH,
            'background-color': COLORS.externallyConnectedNode,
            'border-style': 'solid',
            'border-width': NODE_SOLID_BORDER_WIDTH,
            'border-color': COLORS.externallyConnectedBorder,
            'text-margin-y': '4px',
        },
    },
    {
        selector: 'node.active',
        style: {
            ...deploymentStyle,
            'background-color': COLORS.active,
            'border-style': 'double',
            'border-width': '1px',
            'border-color': COLORS.active,
        },
    },
    {
        selector: 'node.nonIsolated',
        style: {
            ...deploymentStyle,
            'background-color': COLORS.nonIsolated,
            'border-style': 'double',
            'border-width': '1px',
            'border-color': COLORS.nonIsolated,
        },
    },
    {
        selector: 'node.disallowed',
        style: {
            ...deploymentStyle,
            'background-color': COLORS.disallowed,
        },
    },
    {
        selector: 'node.hovered',
        style: {
            opacity: 1,
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 60%, 1)',
            'overlay-opacity': 0.1,
        },
    },
    {
        selector: 'node.selected',
        style: {
            opacity: 1,
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 50%, 1)',
            'overlay-opacity': 0.1,
        },
    },
    {
        selector: ':parent > node.background',
        style: {
            ...deploymentStyle,
        },
    },
    {
        selector: ':parent.background',
        style: {
            opacity: 0.5,
            'z-compound-depth': 'auto',
        },
    },
    {
        selector: ':parent > node.nsEdge',
        style: {
            width: 0.5,
            height: 0.5,
            padding: '0px',
            'background-color': 'white',
        },
    },
    {
        selector: 'edge',
        style: {
            width: 1,
            'line-style': 'dashed',
            'line-color': 'hsla(231, 74%, 82%, 1.00)',
        },
    },
    {
        selector: 'edge.namespace',
        style: {
            'curve-style': 'unbundled-bezier',
            'line-color': COLORS.NSEdge,
            'edge-distances': 'node-position',
            label: 'data(count)',
            'font-size': '8px',
            color: COLORS.NSEdge,
            'font-weight': 500,
            'text-background-opacity': 1,
            'text-background-color': 'white',
            'text-background-shape': 'roundrectangle',
            'text-background-padding': '3px',
            'text-border-opacity': 1,
            'text-border-color': 'hsla(230, 51%, 75%, 1.00)',
            'text-border-width': 1,
            width: 3,
        },
    },
    {
        selector: 'edge.taxi-vertical',
        style: {
            'taxi-direction': 'vertical',
        },
    },
    {
        selector: 'edge.taxi-horizontal',
        style: {
            'taxi-direction': 'horizontal',
        },
    },
    {
        selector: 'edge.inner',
        style: {
            'curve-style': 'haystack',
            'line-style': 'dashed',
            'target-endpoint': 'inside-to-node',
            'z-index': 1000,
            'z-index-compare': 'manual',
        },
    },

    {
        selector: 'edge.nonIsolated',
        style: {
            display: 'none',
        },
    },
    {
        selector: 'edge.active',
        style: {
            'line-style': 'solid',
            'line-color': 'hsla(229, 76%, 87%, 1)',
            'z-compound-depth': 'top',
        },
    },
    {
        selector: 'edge.disallowed',
        style: {
            'line-style': 'solid',
            'line-color': COLORS.disallowed,
            'z-compound-depth': 'top',
            color: COLORS.disallowed,
            'text-border-color': COLORS.disallowed,
        },
    },
    {
        selector: 'edge.unidirectional',
        style: {
            'mid-target-arrow-shape': 'triangle',
            'mid-target-arrow-fill': 'filled',
            'mid-target-arrow-color': 'hsla(231, 74%, 82%, 1.00)',
            'arrow-scale': 0.5,
        },
    },
    {
        selector: 'edge.bidirectional',
        style: {
            'mid-source-arrow-shape': 'triangle',
            'mid-source-arrow-fill': 'filled',
            'mid-source-arrow-color': 'hsla(231, 74%, 82%, 1.00)',
            'mid-target-arrow-shape': 'triangle',
            'mid-target-arrow-fill': 'filled',
            'mid-target-arrow-color': 'hsla(231, 74%, 82%, 1.00)',
            'arrow-scale': 0.5,
        },
    },
    {
        selector: 'edge.inner.withinNS',
        style: {
            'mid-target-arrow-shape': 'none',
            'mid-source-arrow-shape': 'none',
        },
    },
    {
        selector: 'edge.hovered',
        style: {
            opacity: 1,
            color: 'hsl(228, 56%, 63%)',
            'line-color': 'hsl(228, 56%, 63%)',
            'text-border-color': 'hsl(228, 56%, 63%)',
            'overlay-padding': '3px',
            'mid-source-arrow-color': 'hsl(228, 56%, 63%)',
            'mid-target-arrow-color': 'hsl(228, 56%, 63%)',
        },
    },
    {
        selector: ':active',
        style: {
            'overlay-padding': '3px',
            'overlay-color': 'hsla(227, 85%, 50%, 1)',
            'overlay-opacity': 0.1,
        },
    },
];

export default styles;
