package main

import (
	"flag"
	"time"

	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/dashboards/blackbox"
	"github.com/perses/community-mixins/pkg/dashboards/perses"
	"github.com/perses/community-mixins/pkg/dashboards/prometheus"
	"github.com/perses/community-mixins/pkg/promql"
	"github.com/perses/community-mixins/pkg/rules"
	blackboxrules "github.com/perses/community-mixins/pkg/rules/blackbox"
	"github.com/perses/perses/go-sdk/dashboard"
	"github.com/perses/perses/go-sdk/panel"
	panelgroup "github.com/perses/perses/go-sdk/panel-group"
	"github.com/perses/plugins/prometheus/sdk/go/query"
	timeSeriesPanel "github.com/perses/plugins/timeserieschart/sdk/go"
	"github.com/perses/promql-builder/label"
	"github.com/perses/promql-builder/vector"
	"github.com/prometheus/prometheus/promql/parser"
)

var (
	project          string
	datasource       string
	clusterLabelName string
)

var DemoAppCommonPanelQueries = map[string]parser.Expr{
	"DemoAppRequestRate": promql.SumByRate(
		"http_requests_total",
		[]string{"code", "method"},
		label.New("job").Equal("demo-app"),
	),
	"DemoAppTotalRequestsByCode": promql.SumBy(
		"http_requests_total",
		[]string{"code"},
		label.New("job").Equal("demo-app"),
	),
	"DemoAppVersion": vector.New(
		vector.WithMetricName("version"),
		vector.WithLabelMatchers(
			label.New("job").Equal("demo-app"),
		),
	),
	"DemoAppUptime": vector.New(
		vector.WithMetricName("up"),
		vector.WithLabelMatchers(
			label.New("job").Equal("demo-app"),
		),
	),
}

func buildDemoAppDashboard(project string, datasource string) dashboards.DashboardResult {
	return dashboards.NewDashboardResult(
		dashboard.New("demo-app",
			dashboard.ProjectName(project),
			dashboard.Name("demo-app"),
			dashboard.Duration(15*time.Minute),
			dashboard.AddPanelGroup("Demo App",
				panelgroup.PanelsPerLine(2),
				panelgroup.AddPanel("HTTP Request Rate",
					panel.Description("HTTP Request Rate"),
					panel.AddQuery(
						query.PromQL(
							DemoAppCommonPanelQueries["DemoAppRequestRate"].Pretty(0),
							query.SeriesNameFormat("{{code}} - {{method}}"),
						),
					),
					timeSeriesPanel.Chart(),
				),
				panelgroup.AddPanel("HTTP Requests by Status Code",
					panel.Description("HTTP Requests by Status Code"),
					panel.AddQuery(
						query.PromQL(
							DemoAppCommonPanelQueries["DemoAppTotalRequestsByCode"].Pretty(0),
							query.SeriesNameFormat("Status {{code}}"),
						),
					),
					timeSeriesPanel.Chart(),
				),
				panelgroup.AddPanel("App Uptime",
					panel.Description("Application availability status"),
					panel.AddQuery(
						query.PromQL(
							DemoAppCommonPanelQueries["DemoAppUptime"].Pretty(0),
							query.SeriesNameFormat("Status"),
						),
					),
					timeSeriesPanel.Chart(),
				),
				panelgroup.AddPanel("App Version",
					panel.Description("App Version"),
					panel.AddQuery(
						query.PromQL(
							DemoAppCommonPanelQueries["DemoAppVersion"].Pretty(0),
							query.SeriesNameFormat("Version {{version}}"),
						),
					),
					timeSeriesPanel.Chart(),
				),
			),
		),
	).Component("demo-app")
}

func main() {
	flag.StringVar(&project, "project", "default", "The project name")
	flag.StringVar(&datasource, "datasource", "", "The datasource name")
	flag.StringVar(&clusterLabelName, "cluster-label-name", "", "The cluster label name")

	flag.String("output-rules", rules.YAMLOutput, "output format of the rule exec")
	flag.String("output-rules-dir", "./built/rules", "output directory of the rule exec")

	flag.String("output", dashboards.YAMLOutput, "output format of the dashboard exec")
	flag.String("output-dir", "./built", "output directory of the dashboard exec")

	flag.Parse()

	ruleWriter := rules.NewRuleWriter()
	dashboardWriter := dashboards.NewDashboardWriter()

	ruleWriter.Add(
		blackboxrules.BuildBlackboxRules(
			project,
			map[string]string{
				"app.kubernetes.io/component": "blackbox-exporter",
				"app.kubernetes.io/name":      "blackbox-exporter-rules",
				"app.kubernetes.io/part-of":   "blackbox-exporter",
				"app.kubernetes.io/version":   "main",
			},
			map[string]string{},
			blackboxrules.WithDashboardURL("https://demo.perses.dev/projects/perses/dashboards/blackboxexporter"),
		),
	)

	dashboardWriter.Add(perses.BuildPersesOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(prometheus.BuildPrometheusOverview(project, datasource, clusterLabelName))
	dashboardWriter.Add(blackbox.BuildBlackboxExporter(project, datasource, clusterLabelName))
	dashboardWriter.Add(buildDemoAppDashboard(project, datasource))

	dashboardWriter.Write()
	ruleWriter.Write()
}
