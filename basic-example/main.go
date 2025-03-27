package main

import (
	"encoding/json"
	"log"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/testdata"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

func main() {
	testdataRef := dashboard.DataSourceRef{
		Type: cog.ToPtr("testdata"),
	}

	builder := dashboard.NewDashboardBuilder("My Dashboard").
		WithPanel(
			stat.NewPanelBuilder().
				Title("Version").
				Datasource(testdataRef).
				ReduceOptions(common.NewReduceDataOptionsBuilder().
					Calcs([]string{"lastNotNull"}).
					Fields("/.*/")).
				WithTarget(
					testdata.NewDataqueryBuilder().
						ScenarioId("csv_content").
						CsvContent("version\nv1.2.3"),
				),
		).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("Random Time Series").
				Datasource(testdataRef).
				WithTarget(
					testdata.NewDataqueryBuilder().
						ScenarioId("random_walk"),
				),
		)

	dashboard, err := builder.Build()
	if err != nil {
		log.Fatalf("failed to build dashboard: %v", err)
	}

	dashboardJson, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal dashboard: %v", err)
	}

	log.Printf("Dashboard JSON:\n%s", dashboardJson)
}
