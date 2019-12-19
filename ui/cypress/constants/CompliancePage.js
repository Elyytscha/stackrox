export const baseURL = '/main/compliance';

export const url = {
    dashboard: baseURL,
    list: {
        clusters: `${baseURL}/clusters`,
        namespaces: `${baseURL}/namespaces`,
        nodes: `${baseURL}/nodes`,
        standards: {
            CIS_Docker_v1_2_0: `${baseURL}/controls?s[standard]=CIS%20Docker%20v1.2.0`,
            CIS_Kubernetes_v1_5: `${baseURL}/controls?s[standard]=CIS%20Kubernetes%20v1.5`,
            HIPAA_164: `${baseURL}/controls?s[standard]=HIPAA%20164`,
            NIST_800_190: `${baseURL}/controls/?s[standard]=NIST%20SP%20800-190`,
            PCI_DSS_3_2: `${baseURL}/controls?s[standard]=PCI%20DSS%203.2`
        }
    },
    entity: {
        cluster: `${baseURL}/cluster`
    }
};

export const selectors = {
    scanButton: "[data-test-id='scan-button']",
    export: {
        exportButton: "button:contains('Export')",
        pdfButton: "[data-test-id='download-pdf-button']",
        csvButton: "[data-test-id='download-csv-button']"
    },
    dashboard: {
        tileLinks: {
            cluster: {
                tile: "[data-test-id='tile-link']:contains('cluster')",
                value:
                    "[data-test-id='tile-link']:contains('cluster') [data-test-id='tile-link-value']"
            },
            namespace: {
                tile: "[data-test-id='tile-link']:contains('namespace')",
                value:
                    "[data-test-id='tile-link']:contains('namespace') [data-test-id='tile-link-value']"
            },
            node: {
                tile: "[data-test-id='tile-link']:contains('node')",
                value:
                    "[data-test-id='tile-link']:contains('node') [data-test-id='tile-link-value']"
            }
        }
    },
    list: {
        panels: '[data-test-id="panel"]',
        sidePanelHeader: '[data-test-id="panel-header"]:last',
        sidePanelCloseBtn: '[data-test-id="panel"] .close-button',
        banner: {
            content: '[data-test-id="collapsible-banner"]',
            collapseButton: '[data-test-id="banner-collapse-button"]'
        },
        table: {
            header: '[data-test-id="panel-header"]',
            firstGroup: '.table-group-active:first',
            firstTableGroup: '.rt-table:first',
            firstRow: 'div.rt-tr-group > .rt-tr.-odd:first',
            firstRowName: 'div.rt-tr-group > .rt-tr.-odd:first [data-test-id="table-row-name"]',
            secondRow: 'div.rt-tr-group > .rt-tr.-even:first',
            secondRowName: 'div.rt-tr-group > .rt-tr.-even:first [data-test-id="table-row-name"]',
            rows: "table tr:has('td')"
        }
    },
    widgets: "[data-test-id='widget']",
    widget: {
        controlsInCompliance: {
            widget: '[data-test-id="compliance-across-entities"]',
            centerLabel:
                '[data-test-id="compliance-across-entities"] svg .rv-xy-plot__series--label text',
            passingControls:
                '[data-test-id="compliance-across-entities"] [data-test-id="passing-controls-value"]',
            failingControls:
                '[data-test-id="compliance-across-entities"] [data-test-id="failing-controls-value"]',
            arcs: '[data-test-id="compliance-across-entities"] svg path'
        },
        passingStandardsAcrossClusters: {
            widget: '[data-test-id="standards-across-cluster"]',
            axisLinks: '[data-test-id="standards-across-cluster"] a',
            barLabels: '[data-test-id="standards-across-cluster"] .rv-xy-plot__series > text'
        },
        passingStandardsByCluster: {
            NISTBarLinks:
                '[data-test-id="passing-standards-by-cluster"] g.vertical-cluster-bar-NIST rect'
        },
        passingStandardsAcrossNamespaces: {
            axisLinks: '[data-test-id="standards-across-namespace"] a'
        },
        passingStandardsAcrossNodes: {
            axisLinks: '[data-test-id="standards-across-node"] a'
        },
        controlsMostFailed: {
            widget: '[data-test-id="link-list-widget"]:contains("failed")',
            listItems: '[data-test-id="link-list-widget"]:contains("failed") a'
        },
        controlDetails: {
            widget: '[data-test-id="control-details"]',
            standardName: '[data-test-id="control-details"] [data-test-id="standard-name"]',
            controlName: '[data-test-id="control-details"] [data-test-id="control-name"]'
        },
        PCICompliance: {
            controls:
                '[data-test-id="PCI-compliance"] .widget-detail-bullet span:contains("Controls")'
        },
        relatedEntities: '[data-test-id="related-resource-list"]'
    }
};
