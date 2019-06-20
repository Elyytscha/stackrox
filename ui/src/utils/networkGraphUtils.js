export const nonIsolated = node => node.nonIsolatedIngress && node.nonIsolatedEgress;

/**
 * Iterates through a list of nodes and returns only links in the same namespace
 *
 * @param {!Object[]} nodes list of nodes
 * @returns {!Object[]}
 */
export const getLinks = (nodes, networkEdgeMap, networkNodeMap) => {
    const filteredLinks = [];

    nodes.forEach(node => {
        if (!node.entity || node.entity.type !== 'DEPLOYMENT' || !networkEdgeMap) {
            return;
        }
        const { id: srcDeploymentId, deployment: srcDeployment } = node.entity;
        const sourceNS = srcDeployment && srcDeployment.namespace;

        const isActive = key => !!(networkEdgeMap[key] && networkEdgeMap[key].active);
        const isNonIsolated = id => !!(networkNodeMap[id] && networkNodeMap[id].nonIsolated);
        const isBetweenNonIsolated = (srcId, tgtId) => isNonIsolated(srcId) && isNonIsolated(tgtId);
        const isAllowed = (key, { source, target, targetNS }) =>
            sourceNS === 'stackrox' ||
            targetNS === 'stackrox' ||
            isBetweenNonIsolated(source, target) ||
            !!(networkEdgeMap[key] && networkEdgeMap[key].allowed);
        const isDisallowed = (key, link) => isActive(key) && !isAllowed(key, link);

        // For nodes that are egress non-isolated, add outgoing edges to ingress non-isolated nodes, as long as the pair
        // of nodes is not fully non-isolated. This is a compromise to make the non-isolation highlight only apply in
        // the case when there are neither ingress nor egress policies (the data sent from the backend is optimized to
        // treat both phenomena separately and omit edges from a egress non-isolated to an ingress non-isolated
        // deployment, but that would be to confusing in the UI).
        if (node.nonIsolatedEgress) {
            nodes.forEach(targetNode => {
                if (
                    Object.is(node, targetNode) ||
                    !targetNode.entity ||
                    targetNode.entity.type !== 'DEPLOYMENT' ||
                    !targetNode.nonIsolatedIngress // nodes that are ingress-isolated have explicit incoming edges
                ) {
                    return;
                }

                const { id: tgtDeploymentId, deployment: tgtDeployment } = targetNode.entity;
                const targetNS = tgtDeployment && tgtDeployment.namespace;
                const key = [srcDeploymentId, tgtDeploymentId].sort().join('--');

                const link = {
                    source: srcDeploymentId,
                    target: tgtDeploymentId,
                    sourceName: srcDeployment.name,
                    targetName: tgtDeployment.name,
                    sourceNS,
                    targetNS
                };

                link.isActive = isActive(key);
                link.isBetweenNonIsolated = isBetweenNonIsolated(srcDeploymentId, tgtDeploymentId);
                link.isDisallowed = isDisallowed(key, link);

                // Do not draw implicit links between fully non-isolated nodes unless the connection is active.
                const isImplicit = node.nonIsolatedIngress && targetNode.nonIsolatedEgress;
                if (!isImplicit || link.isActive) {
                    filteredLinks.push(link);
                }
            });
        }

        Object.keys(node.outEdges).forEach(targetIndex => {
            const tgtNode = nodes[targetIndex];
            if (!tgtNode || !tgtNode.entity || tgtNode.entity.type !== 'DEPLOYMENT') {
                return;
            }
            const { id: tgtDeploymentId, deployment: tgtDeployment } = tgtNode.entity;
            const targetNS = tgtDeployment && tgtDeployment.namespace;
            const key = [srcDeploymentId, tgtDeploymentId].sort().join('--');
            const link = {
                source: srcDeploymentId,
                target: tgtDeploymentId,
                sourceName: node.entity.deployment.name,
                targetName: tgtDeployment.name,
                sourceNS,
                targetNS
            };

            link.isActive = isActive(key);
            link.isBetweenNonIsolated = isBetweenNonIsolated(srcDeploymentId, tgtDeploymentId);
            link.isDisallowed = isDisallowed(key, link);

            filteredLinks.push(link);
        });
    });

    return filteredLinks;
};
