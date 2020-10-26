import {
    TEXT_MAX_WIDTH,
    NODE_WIDTH,
    NODE_PADDING,
    SIDE_NODE_PADDING,
    nodeTypes,
} from 'constants/networkGraph';
import entityTypes from 'constants/entityTypes';

const nodeWidth = TEXT_MAX_WIDTH + NODE_WIDTH;
const nodeHeight = NODE_WIDTH + NODE_PADDING;

const avgNSDimensions = { width: [], height: [] };

// Gets dimension metadata for a parent node given # of nodes
function getParentDimensions(nodeCount) {
    const cols = Math.ceil(Math.sqrt(nodeCount));
    const rows = Math.ceil(nodeCount / cols);
    const width = cols * nodeWidth;
    const height = rows * nodeHeight;
    avgNSDimensions.width.push(width);
    if (!Number.isNaN(height)) {
        avgNSDimensions.height.push(height);
    }
    return {
        width,
        height,
        rows,
        cols,
    };
}

// Gets positions and dimensions for all parent nodes
export function getParentPositions(nodes, padding) {
    const namespaceNames = nodes
        .filter((node) => node.data().type === entityTypes.NAMESPACE)
        .map((parent) => parent.data().id);
    const externalEntitiesNames = nodes
        .filter((node) => node.data().type === nodeTypes.EXTERNAL_ENTITIES)
        .map((parent) => parent.data().id);

    // Get namespace dimensions sorted by width
    const namespaces = namespaceNames
        .map((id) => {
            const nodeCount = nodes.filter((node) => {
                const data = node.data();
                return data.parent && !data.side && data.parent === id;
            }).length;

            return { ...getParentDimensions(nodeCount), id, nodeCount };
        })
        .sort((a, b) => b.cols - a.cols);
    const externalEntities = externalEntitiesNames
        .map((id) => {
            const nodeCount = 1;
            return { ...getParentDimensions(nodeCount), id, nodeCount };
        })
        .sort((a, b) => b.cols - a.cols);

    const parents = [...namespaces, ...externalEntities];

    // lay out namespaces
    let y = 0;
    let rowNum = 0;
    let colNum = 0;
    const maxNSWidth = Math.max(...avgNSDimensions.width);
    const maxRowWidth = Math.floor(Math.sqrt(parents.length) + 1) * maxNSWidth;
    return parents.map((NS) => {
        const { id, width, height } = NS;
        const newX = (maxNSWidth + padding.x) * colNum;
        const result = {
            id,
            x: newX,
            y,
            width,
            height,
        };

        if (maxRowWidth < newX) {
            // if newX is past maxRowWidth, reset to new row
            y += height + padding.y;
            rowNum += 1;
            colNum = rowNum % 2 ? 1 : 0;
        } else {
            colNum += 2;
        }

        return result;
    });
}
// Can't use this.options inside prototypal function constructor in strict mode, so using a closure instead
let edgeGridOptions = {};

export function edgeGridLayout(options) {
    const defaults = {
        parentPadding: { bottom: 0, top: 0, left: 0, right: 0 },
        position: { x: 0, y: 0 },
    };
    edgeGridOptions = { ...defaults, ...options };
}

// eslint-disable-next-line func-names
edgeGridLayout.prototype.run = function () {
    const options = edgeGridOptions;
    const { parentPadding, position, eles } = options;

    const nodes = eles.nodes().not('.namespace').not('.cluster');

    const renderNodes = nodes.not('[side]');
    const sideNodes = eles.nodes('[side]');

    const isExternalEntities = !renderNodes.length;
    const numNodes = isExternalEntities ? 1 : renderNodes.length;

    const { width, height, cols } = getParentDimensions(numNodes);

    // Calculate cell dimensions
    const cellWidth = nodeWidth;
    const cellHeight = nodeHeight;

    // Midpoints for sidewall nodes
    const midHeight = position.y + height / 2;
    const midWidth = position.x + width / 2;

    let currentRow = 0;
    let currentCol = 0;
    function incrementCell() {
        currentCol += 1;
        if (currentCol >= cols) {
            currentCol = 0;
            currentRow += 1;
        }
    }

    function getRenderNodePos(element) {
        if (element.locked() || element.isParent()) {
            return false;
        }
        const x = currentCol * cellWidth + cellWidth / 2 + position.x;
        const y = currentRow * cellHeight + cellHeight / 2 + position.y;
        incrementCell();
        return { x, y };
    }

    function getSideNodePos(element) {
        const { side } = element.data();
        switch (side) {
            case 'top':
                return {
                    x: midWidth,
                    y: position.y - parentPadding.top - SIDE_NODE_PADDING,
                };
            case 'bottom':
                return {
                    x: midWidth,
                    y: position.y + height + parentPadding.bottom + SIDE_NODE_PADDING,
                };
            case 'left':
                return {
                    x: position.x - parentPadding.left - SIDE_NODE_PADDING,
                    y: midHeight,
                };
            case 'right':
                return {
                    x: position.x + width + parentPadding.right + SIDE_NODE_PADDING,
                    y: midHeight,
                };
            default:
                return { x: position.x, y: position.y };
        }
    }

    renderNodes.layoutPositions(this, options, getRenderNodePos);
    sideNodes.layoutPositions(this, options, getSideNodePos);
    return this; // chaining
};
