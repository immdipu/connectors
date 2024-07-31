package openapi

import (
	_ "embed"

	"github.com/amp-labs/connectors/tools/fileconv/api3"
)

var (
	// Static file containing openapi spec.
	//
	//go:embed support-api.yaml
	apiFile []byte

	FileManager = api3.NewOpenapiFileManager(apiFile) // nolint:gochecknoglobals
)
