// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestChartPreviewWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Charts.Preview(context.TODO(), langsmith.ChartPreviewParams{
		BucketInfo: langsmith.F(langsmith.ChartPreviewParamsBucketInfo{
			EndTime:   langsmith.F(time.Now()),
			OmitData:  langsmith.F(true),
			StartTime: langsmith.F(time.Now()),
			Stride: langsmith.F(langsmith.TimedeltaInputParam{
				Days:    langsmith.F(int64(0)),
				Hours:   langsmith.F(int64(0)),
				Minutes: langsmith.F(int64(0)),
			}),
			Timezone: langsmith.F("timezone"),
		}),
		Chart: langsmith.F(langsmith.ChartPreviewParamsChart{
			Series: langsmith.F([]langsmith.ChartPreviewParamsChartSeries{{
				ID:          langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
				Name:        langsmith.F("name"),
				FeedbackKey: langsmith.F("feedback_key"),
				FilterDefinition: langsmith.F[langsmith.ChartPreviewParamsChartSeriesFilterDefinitionUnion](langsmith.ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProject{
					ProjectIDs:  langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
					SourceType:  langsmith.F(langsmith.ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject),
					RunFilter:   langsmith.F("run_filter"),
					TraceFilter: langsmith.F("trace_filter"),
					TreeFilter:  langsmith.F("tree_filter"),
				}),
				Filters: langsmith.F(langsmith.ChartPreviewParamsChartSeriesFilters{
					Filter:      langsmith.F("filter"),
					Session:     langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
					TraceFilter: langsmith.F("trace_filter"),
					TreeFilter:  langsmith.F("tree_filter"),
				}),
				GroupBy: langsmith.F(langsmith.ChartPreviewParamsChartSeriesGroupBy{
					Attribute: langsmith.F(langsmith.ChartPreviewParamsChartSeriesGroupByAttributeName),
					MaxGroups: langsmith.F(int64(0)),
					Path:      langsmith.F("path"),
					SetBy:     langsmith.F(langsmith.ChartPreviewParamsChartSeriesGroupBySetBySection),
				}),
				GroupByDefinitions: langsmith.F([]langsmith.ChartPreviewParamsChartSeriesGroupByDefinitionUnion{langsmith.ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlain{
					Attribute: langsmith.F(langsmith.ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName),
				}}),
				Metadata: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				Metric: langsmith.F(langsmith.ChartPreviewParamsChartSeriesMetricRunCount),
				MetricDefinition: langsmith.F[langsmith.ChartPreviewParamsChartSeriesMetricDefinitionUnion](langsmith.ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetric{
					Entity: langsmith.F(langsmith.ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback),
					Params: langsmith.F(langsmith.ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricParams{
						FeedbackKey: langsmith.F("feedback_key"),
					}),
					Filter: langsmith.F("filter"),
					Type:   langsmith.F(langsmith.ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount),
				}),
				ProjectMetric: langsmith.F(langsmith.ChartPreviewParamsChartSeriesProjectMetricMemoryUsage),
				WorkspaceID:   langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			}}),
			CommonFilters: langsmith.F(langsmith.ChartPreviewParamsChartCommonFilters{
				Filter:      langsmith.F("filter"),
				Session:     langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
				TraceFilter: langsmith.F("trace_filter"),
				TreeFilter:  langsmith.F("tree_filter"),
			}),
		}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
