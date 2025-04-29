package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/stat"
	"github.com/grafana/grafana-foundation-sdk/go/testdata"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"
)

type DashboardWrapper struct {
	APIVersion string              `json:"apiVersion"`
	Kind       string              `json:"kind"`
	Metadata   Metadata            `json:"metadata"`
	Spec       dashboard.Dashboard `json:"spec"`
}

type Metadata struct {
	Name string `json:"name"`
}

func main() {
	testdataRef := dashboard.DataSourceRef{
		Type: cog.ToPtr("grafana-testdata-datasource"),
		Uid:  cog.ToPtr("testdata"),
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

	dashboardWrapper := DashboardWrapper{
		APIVersion: "dashboard.grafana.app/v1alpha1",
		Kind:       "Dashboard",
		Metadata: Metadata{
			Name: "sample-dashboard",
		},
		Spec: dashboard,
	}

	dashboardJson, err := json.MarshalIndent(dashboardWrapper, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal dashboard: %v", err)
	}

	err = os.WriteFile("sample-dashboard.json", dashboardJson, 0644)
	if err != nil {
		log.Fatalf("failed to write dashboard to file: %v", err)
	}

	log.Printf("Dashboard JSON:\n%s", dashboardJson)
}
