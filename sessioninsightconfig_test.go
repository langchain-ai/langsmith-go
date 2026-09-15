package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionInsightConfigNew(t *testing.T) {
	const sessionID = "475a0601-e042-4e75-96e1-b2776ab07b21"
	const configID = "26f293f6-7c30-4d6a-a072-60994152fbff"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/v1/sessions/"+sessionID+"/insights/configs", r.URL.Path)

		var body map[string]interface{}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "Reliability review", body["name"])
		assert.Equal(t, "Find failure patterns", body["description"])
		config, ok := body["config"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "Summarize {{run.inputs}} and {{run.error}}", config["summary_prompt"])
		assert.Equal(t, "eq(is_root, true)", config["filter"])
		assert.Equal(t, "openai", config["model"])
		assert.Equal(t, float64(24), config["last_n_hours"])

		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{
  "id": "` + configID + `",
  "name": "Reliability review",
  "description": null,
  "config": {
    "name": "Reliability review",
    "hierarchy": null,
    "partitions": null,
    "sample": null,
    "summary_prompt": "Summarize {{run.inputs}} and {{run.error}}",
    "filter": "eq(is_root, true)",
    "attribute_schemas": null,
    "model": "openai",
    "last_n_hours": 24
  },
  "schedule_cron": null
}`))
		require.NoError(t, err)
	}))
	defer server.Close()

	client := langsmith.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test-api-key"),
		option.WithTenantID("test-tenant-id"),
	)
	response, err := client.Sessions.Insights.Configs.New(
		context.Background(),
		sessionID,
		langsmith.SessionInsightConfigNewParams{
			Name:        langsmith.F("Reliability review"),
			Description: langsmith.F("Find failure patterns"),
			Config: langsmith.F(langsmith.CreateRunClusteringJobRequestParam{
				SummaryPrompt: langsmith.F("Summarize {{run.inputs}} and {{run.error}}"),
				Filter:        langsmith.F("eq(is_root, true)"),
				Model:         langsmith.F(langsmith.CreateRunClusteringJobRequestModelOpenAI),
				LastNHours:    langsmith.F(int64(24)),
			}),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, response)
	assert.Equal(t, configID, response.ID)
	assert.Equal(t, "Reliability review", response.Name)
	assert.Equal(t, int64(24), response.Config.LastNHours)
	assert.Equal(t, "Summarize {{run.inputs}} and {{run.error}}", response.Config.SummaryPrompt)
	assert.True(t, response.JSON.Description.IsNull())
	assert.True(t, response.JSON.ScheduleCron.IsNull())
	assert.True(t, response.Config.JSON.Hierarchy.IsNull())
}

func TestSessionInsightConfigNewRejectsMissingSessionID(t *testing.T) {
	client := langsmith.NewClient(option.WithAPIKey("test-api-key"))
	response, err := client.Sessions.Insights.Configs.New(
		context.Background(),
		"",
		langsmith.SessionInsightConfigNewParams{},
	)

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, "missing required session_id parameter", err.Error())
}
