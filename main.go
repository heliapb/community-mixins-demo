package main

import (
	"flag"

	"github.com/perses/community-mixins/pkg/dashboards"
	"github.com/perses/community-mixins/pkg/dashboards/blackbox"
	"github.com/perses/community-mixins/pkg/dashboards/perses"
	"github.com/perses/community-mixins/pkg/dashboards/prometheus"
	"github.com/perses/community-mixins/pkg/rules"
	blackboxrules "github.com/perses/community-mixins/pkg/rules/blackbox"
)

var (
	project          string
	datasource       string
	clusterLabelName string
)

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
	dashboardWriter.Add(prometheus.BuildPrometheusRemoteWrite(project, datasource, clusterLabelName))
	dashboardWriter.Add(blackbox.BuildBlackboxExporter(project, datasource, clusterLabelName))

	dashboardWriter.Write()
	ruleWriter.Write()
}
