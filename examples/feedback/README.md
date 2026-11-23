# Feedback Management Example

This example demonstrates how to use the LangSmith Go SDK to manage feedback for your runs. Feedback is a crucial part of the LLM development lifecycle, allowing you to record scores, comments, and corrections from both human annotators and automated evaluators.

## What this example covers

1. **Creating a Run**: How to ingest a run into LangSmith so you have an entity to provide feedback on.
2. **Submitting Feedback**: Using `client.Feedback.New` to submit scores, categorical values, and comments.
3. **Updating Feedback**: Using `client.Feedback.Update` to modify existing feedback entries.
4. **Listing Feedback**: Querying feedback for specific runs using `client.Feedback.List`.
5. **Deleting Feedback**: Using `client.Feedback.Delete` to remove feedback entries.

## Prerequisites

- Go 1.22+
- A LangSmith API Key (set as `LANGSMITH_API_KEY` environment variable)

## Running the example

```sh
export LANGSMITH_API_KEY="your-api-key"
go run ./examples/feedback
```

## Key Concepts

- **Key**: The name of the feedback metric (e.g., "accuracy", "user_score", "helpfulness").
- **Score**: A numerical value representing the feedback (typically between 0 and 1).
- **Value**: A categorical or string value for the feedback (e.g., "correct", "incorrect").
- **Comment**: An optional text description providing more context.
- **Run ID**: The ID of the run the feedback is associated with.
