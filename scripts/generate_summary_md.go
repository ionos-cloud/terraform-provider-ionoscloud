package main

import (
	"fmt"
	"github.com/ionos-cloud/terraform-provider-ionoscloud/v6/internal/doc"
	"os"
)

const DocsDir = "docs"

func main() {
	err := doc.GenerateSummary(DocsDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "an error occurred while generating summary.md file, error: %w", err)
		os.Exit(1)
	}
}
