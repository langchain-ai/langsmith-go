// Command langsmith is the LangSmith CLI. See package
// github.com/langchain-ai/langsmith-go/cli.
package main

import (
	"os"

	"github.com/langchain-ai/langsmith-go/cli"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	os.Exit(cli.Main(version, commit, date))
}
