export const alerts = {
    countsByCluster: /v1\/alerts\/summary\/counts\?group_by=CLUSTER.*/,
    countsByCategory: /v1\/alerts\/summary\/counts\?group_by=CATEGORY.*/
};

export const clusters = {
    list: 'v1/clusters',
    zip: 'api/extensions/clusters/zip'
};

export const benchmarks = {
    configs: 'v1/benchmarks/configs',
    scans: 'v1/benchmarks/scans/*',
    cisDockerScans: 'v1/benchmarks/scans?benchmarkId=*',
    triggers: 'v1/benchmarks/triggers/*'
};

export const risks = {
    riskyDeployments: 'v1/deployments*'
};

export const search = {
    globalSearchWithResults: '/v1/search?query=Cluster:remote',
    globalSearchWithNoResults: '/v1/search?query=Cluster:'
};

export const images = {
    list: '/v1/images*'
};
