# LangChain Go API Library

<a href="https://pkg.go.dev/github.com/langchain-ai/langsmith-go"><img src="https://pkg.go.dev/badge/github.com/langchain-ai/langsmith-go.svg" alt="Go Reference"></a>

The LangChain Go library provides convenient access to the LangChain REST API
from applications written in Go.

It is generated with [Stainless](https://www.stainless.com/).

## Installation

<!-- x-release-please-start-version -->

```go
import (
	"github.com/langchain-ai/langsmith-go" // imported as langsmith
)
```

<!-- x-release-please-end -->

Or to pin the version:

<!-- x-release-please-start-version -->

```sh
go get -u 'github.com/langchain-ai/langsmith-go@v0.1.0-alpha.7'
```

<!-- x-release-please-end -->

## Requirements

This library requires Go 1.22+.

## Configuration

The client can be configured using environment variables or by passing options directly to `langsmith.NewClient()`. Environment variables are automatically read when creating a new client.

### Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `LANGSMITH_API_KEY` | Optional* | Your LangSmith API key for authentication |
| `LANGSMITH_BEARER_TOKEN` | Optional* | Bearer token for authentication (alternative to API key) |
| `LANGSMITH_TENANT_ID` | Optional | Your LangSmith tenant ID |
| `LANGSMITH_ORGANIZATION_ID` | Optional | Your LangSmith organization ID |
| `LANGCHAIN_BASE_URL` | Optional | Custom base URL for the LangSmith API (defaults to `https://api.smith.langchain.com`) |

\* Either `LANGSMITH_API_KEY` or `LANGSMITH_BEARER_TOKEN` is required for authentication.

## Examples

This repository includes several examples demonstrating common use cases:

- **[List Runs](./examples/list_runs)** - Query and filter runs from your LangSmith project
- **[Dataset Management](./examples/dataset)** - Create datasets, add examples individually or in bulk, and manage dataset lifecycle
- **[E2E Evaluation](./examples/e2e_eval)** - Run OpenAI experiments with automatic OpenTelemetry tracing linked to dataset examples
- **[Record Experiment](./examples/record_experiment)** - Create datasets, examples, sessions, and batch ingest runs for experiments
- **[Prompt Management](./examples/prompt_management)** - Create and manage prompt repositories, commits, and versions
- **[OpenTelemetry Ingestion](./examples/otel_ingestion)** - Send OpenTelemetry traces to LangSmith with hierarchical span structure
- **[OpenTelemetry + OpenAI](./examples/otel_openai)** - Make OpenAI API calls with manual OpenTelemetry tracing to LangSmith
- **[OpenTelemetry + OpenAI (Go Client)](./examples/otel_go_client_openai)** - Automatic OpenTelemetry tracing for OpenAI API calls using sashabaranov/go-openai client wrapper
- **[OpenTelemetry + Anthropic](./examples/otel_anthropic)** - Automatic OpenTelemetry tracing for Anthropic API calls using HTTP client wrapper

Each example includes detailed documentation in its source code. To run an example:

```sh
go run ./examples/<example-name>
```

Make sure to set the required environment variables `LANGSMITH_API_KEY` before running and `OPENAI_API_KEY` for the OpenAI examples.

## Usage

The full API of this library can be found in [api.md](api.md).

```go
package main

import (
	"context"
	"fmt"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func main() {
	client := langsmith.NewClient(
		option.WithAPIKey("My API Key"), // defaults to os.LookupEnv("LANGSMITH_API_KEY")
	)
	customChartsSection, err := client.Sessions.Dashboard(
		context.TODO(),
		"1ffaeba7-541e-469f-bae7-df3208ea3d45",
		langsmith.SessionDashboardParams{
			CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{},
		},
	)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", customChartsSection.ID)
}

```

### Request fields

All request parameters are wrapped in a generic `Field` type,
which we use to distinguish zero values from null or omitted fields.

This prevents accidentally sending a zero value if you forget a required parameter,
and enables explicitly sending `null`, `false`, `''`, or `0` on optional parameters.
Any field not specified is not sent.

To construct fields with values, use the helpers `String()`, `Int()`, `Float()`, or most commonly, the generic `F[T]()`.
To send a null, use `Null[T]()`, and to send a nonconforming value, use `Raw[T](any)`. For example:

```go
params := FooParams{
	Name: langsmith.F("hello"),

	// Explicitly send `"description": null`
	Description: langsmith.Null[string](),

	Point: langsmith.F(langsmith.Point{
		X: langsmith.Int(0),
		Y: langsmith.Int(1),

		// In cases where the API specifies a given type,
		// but you want to send something else, use `Raw`:
		Z: langsmith.Raw[int64](0.01), // sends a float
	}),
}
```

### Response objects

All fields in response structs are value types (not pointers or wrappers).

If a given field is `null`, not present, or invalid, the corresponding field
will simply be its zero value.

All response structs also include a special `JSON` field, containing more detailed
information about each property, which you can use like so:

```go
if res.Name == "" {
	// true if `"name"` is either not present or explicitly null
	res.JSON.Name.IsNull()

	// true if the `"name"` key was not present in the response JSON at all
	res.JSON.Name.IsMissing()

	// When the API returns data that cannot be coerced to the expected type:
	if res.JSON.Name.IsInvalid() {
		raw := res.JSON.Name.Raw()

		legacyName := struct{
			First string `json:"first"`
			Last  string `json:"last"`
		}{}
		json.Unmarshal([]byte(raw), &legacyName)
		name = legacyName.First + " " + legacyName.Last
	}
}
```

These `.JSON` structs also include an `Extras` map containing
any properties in the json response that were not specified
in the struct. This can be useful for API features not yet
present in the SDK.

```go
body := res.JSON.ExtraFields["my_unexpected_field"].Raw()
```

### RequestOptions

This library uses the functional options pattern. Functions defined in the
`option` package return a `RequestOption`, which is a closure that mutates a
`RequestConfig`. These options can be supplied to the client or at individual
requests. For example:

```go
client := langsmith.NewClient(
	// Adds a header to every request made by the client
	option.WithHeader("X-Some-Header", "custom_header_info"),
)

client.Sessions.Dashboard(context.TODO(), ...,
	// Override the header
	option.WithHeader("X-Some-Header", "some_other_custom_header_info"),
	// Add an undocumented field to the request body, using sjson syntax
	option.WithJSONSet("some.json.path", map[string]string{"my": "object"}),
)
```

See the [full list of request options](https://pkg.go.dev/github.com/langchain-ai/langsmith-go/option).

### Pagination

This library provides some conveniences for working with paginated list endpoints.

You can use `.ListAutoPaging()` methods to iterate through items across all pages:

```go
iter := client.Datasets.ListAutoPaging(context.TODO(), langsmith.DatasetListParams{
	Limit: langsmith.F(int64(100)),
})
// Automatically fetches more pages as needed.
for iter.Next() {
	dataset := iter.Current()
	fmt.Printf("%+v\n", dataset)
}
if err := iter.Err(); err != nil {
	panic(err.Error())
}
```

Or you can use simple `.List()` methods to fetch a single page and receive a standard response object
with additional helper methods like `.GetNextPage()`, e.g.:

```go
page, err := client.Datasets.List(context.TODO(), langsmith.DatasetListParams{
	Limit: langsmith.F(int64(100)),
})
for page != nil {
	for _, dataset := range page.Items {
		fmt.Printf("%+v\n", dataset)
	}
	page, err = page.GetNextPage()
}
if err != nil {
	panic(err.Error())
}
```

### Errors

When the API returns a non-success status code, we return an error with type
`*langsmith.Error`. This contains the `StatusCode`, `*http.Request`, and
`*http.Response` values of the request, as well as the JSON of the error body
(much like other response objects in the SDK).

To handle errors, we recommend that you use the `errors.As` pattern:

```go
_, err := client.Sessions.Dashboard(
	context.TODO(),
	"1ffaeba7-541e-469f-bae7-df3208ea3d45",
	langsmith.SessionDashboardParams{
		CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{},
	},
)
if err != nil {
	var apierr *langsmith.Error
	if errors.As(err, &apierr) {
		println(string(apierr.DumpRequest(true)))  // Prints the serialized HTTP request
		println(string(apierr.DumpResponse(true))) // Prints the serialized HTTP response
	}
	panic(err.Error()) // GET "/api/v1/sessions/{session_id}/dashboard": 400 Bad Request { ... }
}
```

When other errors occur, they are returned unwrapped; for example,
if HTTP transport fails, you might receive `*url.Error` wrapping `*net.OpError`.

### Timeouts

Requests do not time out by default; use context to configure a timeout for a request lifecycle.

Note that if a request is [retried](#retries), the context timeout does not start over.
To set a per-retry timeout, use `option.WithRequestTimeout()`.

```go
// This sets the timeout for the request, including all the retries.
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()
client.Sessions.Dashboard(
	ctx,
	"1ffaeba7-541e-469f-bae7-df3208ea3d45",
	langsmith.SessionDashboardParams{
		CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{},
	},
	// This sets the per-retry timeout
	option.WithRequestTimeout(20*time.Second),
)
```

### File uploads

Request parameters that correspond to file uploads in multipart requests are typed as
`param.Field[io.Reader]`. The contents of the `io.Reader` will by default be sent as a multipart form
part with the file name of "anonymous_file" and content-type of "application/octet-stream".

The file name and content-type can be customized by implementing `Name() string` or `ContentType()
string` on the run-time type of `io.Reader`. Note that `os.File` implements `Name() string`, so a
file returned by `os.Open` will be sent with the file name on disk.

We also provide a helper `langsmith.FileParam(reader io.Reader, filename string, contentType string)`
which can be used to wrap any `io.Reader` with the appropriate file name and content type.

```go
// A file from the file system
file, err := os.Open("/path/to/file")
langsmith.ExampleUploadFromCsvParams{
	File:      langsmith.F[io.Reader](file),
	InputKeys: langsmith.F([]string{"string"}),
}

// A file from a string
langsmith.ExampleUploadFromCsvParams{
	File:      langsmith.F[io.Reader](strings.NewReader("my file contents")),
	InputKeys: langsmith.F([]string{"string"}),
}

// With a custom filename and contentType
langsmith.ExampleUploadFromCsvParams{
	File:      langsmith.FileParam(strings.NewReader(`{"hello": "foo"}`), "file.go", "application/json"),
	InputKeys: langsmith.F([]string{"string"}),
}
```

### Retries

Certain errors will be automatically retried 2 times by default, with a short exponential backoff.
We retry by default all connection errors, 408 Request Timeout, 409 Conflict, 429 Rate Limit,
and >=500 Internal errors.

You can use the `WithMaxRetries` option to configure or disable this:

```go
// Configure the default for all requests:
client := langsmith.NewClient(
	option.WithMaxRetries(0), // default is 2
)

// Override per-request:
client.Sessions.Dashboard(
	context.TODO(),
	"1ffaeba7-541e-469f-bae7-df3208ea3d45",
	langsmith.SessionDashboardParams{
		CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{},
	},
	option.WithMaxRetries(5),
)
```

### Accessing raw response data (e.g. response headers)

You can access the raw HTTP response data by using the `option.WithResponseInto()` request option. This is useful when
you need to examine response headers, status codes, or other details.

```go
// Create a variable to store the HTTP response
var response *http.Response
customChartsSection, err := client.Sessions.Dashboard(
	context.TODO(),
	"1ffaeba7-541e-469f-bae7-df3208ea3d45",
	langsmith.SessionDashboardParams{
		CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{},
	},
	option.WithResponseInto(&response),
)
if err != nil {
	// handle error
}
fmt.Printf("%+v\n", customChartsSection)

fmt.Printf("Status Code: %d\n", response.StatusCode)
fmt.Printf("Headers: %+#v\n", response.Header)
```

### Making custom/undocumented requests

This library is typed for convenient access to the documented API. If you need to access undocumented
endpoints, params, or response properties, the library can still be used.

#### Undocumented endpoints

To make requests to undocumented endpoints, you can use `client.Get`, `client.Post`, and other HTTP verbs.
`RequestOptions` on the client, such as retries, will be respected when making these requests.

```go
var (
    // params can be an io.Reader, a []byte, an encoding/json serializable object,
    // or a "…Params" struct defined in this library.
    params map[string]interface{}

    // result can be an []byte, *http.Response, a encoding/json deserializable object,
    // or a model defined in this library.
    result *http.Response
)
err := client.Post(context.Background(), "/unspecified", params, &result)
if err != nil {
    …
}
```

#### Undocumented request params

To make requests using undocumented parameters, you may use either the `option.WithQuerySet()`
or the `option.WithJSONSet()` methods.

```go
params := FooNewParams{
    ID:   langsmith.F("id_xxxx"),
    Data: langsmith.F(FooNewParamsData{
        FirstName: langsmith.F("John"),
    }),
}
client.Foo.New(context.Background(), params, option.WithJSONSet("data.last_name", "Doe"))
```

#### Undocumented response properties

To access undocumented response properties, you may either access the raw JSON of the response as a string
with `result.JSON.RawJSON()`, or get the raw JSON of a particular field on the result with
`result.JSON.Foo.Raw()`.

Any fields that are not present on the response struct will be saved and can be accessed by `result.JSON.ExtraFields()` which returns the extra fields as a `map[string]Field`.

### Middleware

We provide `option.WithMiddleware` which applies the given
middleware to requests.

```go
func Logger(req *http.Request, next option.MiddlewareNext) (res *http.Response, err error) {
	// Before the request
	start := time.Now()
	LogReq(req)

	// Forward the request to the next handler
	res, err = next(req)

	// Handle stuff after the request
	end := time.Now()
	LogRes(res, err, start - end)

    return res, err
}

client := langsmith.NewClient(
	option.WithMiddleware(Logger),
)
```

When multiple middlewares are provided as variadic arguments, the middlewares
are applied left to right. If `option.WithMiddleware` is given
multiple times, for example first in the client then the method, the
middleware in the client will run first and the middleware given in the method
will run next.

You may also replace the default `http.Client` with
`option.WithHTTPClient(client)`. Only one http client is
accepted (this overwrites any previous client) and receives requests after any
middleware has been applied.

## Semantic versioning

This package generally follows [SemVer](https://semver.org/spec/v2.0.0.html) conventions, though certain backwards-incompatible changes may be released as minor versions:

1. Changes to library internals which are technically public but not intended or documented for external use. _(Please open a GitHub issue to let us know if you are relying on such internals.)_
2. Changes that we do not expect to impact the vast majority of users in practice.

We take backwards-compatibility seriously and work hard to ensure you can rely on a smooth upgrade experience.

We are keen for your feedback; please open an [issue](https://www.github.com/langchain-ai/langsmith-go/issues) with questions, bugs, or suggestions.

## Contributing

See [the contributing documentation](./CONTRIBUTING.md).
