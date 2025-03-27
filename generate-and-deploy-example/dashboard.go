package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"github.com/go-openapi/strfmt"
	"github.com/grafana/grafana-foundation-sdk/go/cog"
	"github.com/grafana/grafana-foundation-sdk/go/common"
	"github.com/grafana/grafana-foundation-sdk/go/dashboard"
	"github.com/grafana/grafana-foundation-sdk/go/logs"
	"github.com/grafana/grafana-foundation-sdk/go/loki"
	"github.com/grafana/grafana-foundation-sdk/go/prometheus"
	"github.com/grafana/grafana-foundation-sdk/go/timeseries"

	grafanaApi "github.com/grafana/grafana-openapi-client-go/client"
	"github.com/grafana/grafana-openapi-client-go/models"
)

// getEnv retrieves the value of the environment variable named by the key.
// If the variable is not set, it returns the provided default value.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func prometheusDatasourceRef() dashboard.DataSourceRef {
	return dashboard.DataSourceRef{
		Type: cog.ToPtr("prometheus"),
		Uid:  cog.ToPtr("prometheus"),
	}
}

func lokiDatasourceRef() dashboard.DataSourceRef {
	return dashboard.DataSourceRef{
		Type: cog.ToPtr("loki"),
		Uid:  cog.ToPtr("loki"),
	}
}

func GenerateAndPublishDashboard() {
	builder := dashboard.NewDashboardBuilder("Web Service Dashboard").
		Uid("web-service-dashboard").
		Tags([]string{"generated", "foundation-sdk", "go"}).
		Refresh("5m").
		Time("now-1h", "now").
		Timezone(common.TimeZoneBrowser).
		WithRow(dashboard.NewRowBuilder("Overview")).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("HTTP Requests/Second").
				Datasource(prometheusDatasourceRef()).
				WithTarget(
					prometheus.NewDataqueryBuilder().
						Expr(`rate(http_request_duration_seconds_count[$__rate_interval])`).
						Range().
						Format(prometheus.PromQueryFormatTimeSeries).
						LegendFormat("__auto"),
				),
		).
		WithPanel(
			timeseries.NewPanelBuilder().
				Title("Average HTTP Response Times").
				Datasource(prometheusDatasourceRef()).
				WithTarget(
					prometheus.NewDataqueryBuilder().
						Expr(`sum(rate(http_request_duration_seconds_sum[$__rate_interval])) by (endpoint) / sum(rate(http_request_duration_seconds_count[$__rate_interval])) by (endpoint)`).
						Range().
						Format(prometheus.PromQueryFormatTimeSeries).
						LegendFormat("{{endpoint}}"),
				),
		).
		WithPanel(
			logs.NewPanelBuilder().
				Title("Web Service Logs").
				Datasource(lokiDatasourceRef()).
				ShowTime(true).
				SortOrder(common.LogsSortOrderDescending).
				EnableLogDetails(true).
				WithTarget(
					loki.NewDataqueryBuilder().
						Expr(`{service="web-service"}`).
						QueryType(string(loki.LokiQueryTypeRange)).
						LegendFormat("__auto"),
				),
		).
		WithPanel(
			logs.NewPanelBuilder().
				Title("Dashboard Service Logs").
				Datasource(lokiDatasourceRef()).
				ShowTime(true).
				SortOrder(common.LogsSortOrderDescending).
				EnableLogDetails(true).
				WithTarget(
					loki.NewDataqueryBuilder().
						Expr(`{service="dashboard-service"}`).
						QueryType(string(loki.LokiQueryTypeRange)).
						LegendFormat("__auto"),
				),
		)

	dashboardLogger.Info("Generating dashboard...")

	dashboard, err := builder.Build()
	if err != nil {
		panic(err)
	}

	dashboardJson, err := json.Marshal(dashboard)
	if err != nil {
		panic(err)
	}

	dashboardLogger.Info(string(dashboardJson))

	persistDashboard(dashboard)
}

func persistDashboard(dashboard dashboard.Dashboard) {
	dashboardLogger.Info("Persisting dashboard to Grafana")

	host := getEnv("GRAFANA_HOST", "grafana:3000")

	client := grafanaApi.NewHTTPClientWithConfig(strfmt.Default, &grafanaApi.TransportConfig{
		Host:      host,
		BasePath:  "/api",
		Schemes:   []string{"http"},
		BasicAuth: url.UserPassword("admin", "admin"),
	})

	// Find existing provisioned folder
	_, err := client.Folders.GetFolderByUID("provisioned-folder")
	if err != nil {
		dashboardLogger.Info("Folder not found, creating a new one")
		_, err := client.Folders.CreateFolder(&models.CreateFolderCommand{
			Title: "Provisioned",
			UID:   "provisioned-folder",
		})
		if err != nil {
			dashboardLogger.Error(fmt.Sprintf("Error creating folder: %v", err))
			return
		}
	}

	// Create the dashboard
	_, err = client.Dashboards.PostDashboard(&models.SaveDashboardCommand{
		Dashboard: dashboard,
		FolderUID: "provisioned-folder",
		Overwrite: true,
	})

	if err != nil {
		dashboardLogger.Error(fmt.Sprintf("Error creating dashboard: %v", err))
		return
	}

	dashboardLogger.Info("Dashboard created successfully")
}
