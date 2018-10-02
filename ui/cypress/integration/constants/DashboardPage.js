export const url = '/main/dashboard';

const severityColors = {
    CRITICAL_SEVERITY: 'hsl(358, 81%, 80%)',
    HIGH_SEVERITY: 'hsl(16, 81%, 80%)',
    MEDIUM_SEVERITY: 'hsl(39, 80%, 80%)',
    LOW_SEVERITY: 'hsl(230, 43%, 90%)'
};

export const selectors = {
    navLink: 'nav li:contains("Dashboard") a',
    buttons: {
        viewAll: 'button:contains("View All")'
    },
    sectionHeaders: {
        systemViolations: 'h2:contains("System Violations")',
        benchmarks: '.slick-active h2:contains("Benchmarks")',
        violationsByClusters: 'h2:contains("Violations by Cluster")',
        eventsByTime: 'h2:contains("Active Violations by Time")',
        securityBestPractices: 'h2:contains("Security Best Practices")',
        devopsBestPractices: 'h2:contains("DevOps Best Practices")',
        topRiskyDeployments: 'h2:contains("Top Risky Deployments")'
    },
    chart: {
        xAxis: 'g.xAxis',
        grid: 'g.recharts-cartesian-grid',
        medSeverityBar: `g.recharts-bar-rectangle path[fill="${severityColors.MEDIUM_SEVERITY}"]`,
        lowSeverityBar: `g.recharts-bar-rectangle path[fill="${severityColors.LOW_SEVERITY}"]`,
        medSeveritySector: `g.recharts-pie-sector path[fill="${severityColors.MEDIUM_SEVERITY}"]`,
        legendItem: `span.recharts-legend-item-text`
    },
    slick: {
        dashboardBenchmarks: {
            prevButton: '.dashboard-benchmarks .carousel-prev-arrow',
            nextButton: '.dashboard-benchmarks .carousel-next-arrow',
            list: '.dashboard-benchmarks .slick-slide',
            currentSlide: '.dashboard-benchmarks .slick-current',
            track: '.slick-track'
        }
    },
    timeseries: 'svg.recharts-surface',
    searchInput: '.Select-input > input'
};
