package dynamicscrm

import (
	"fmt"
	"strings"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/interpreter"
	"github.com/amp-labs/connectors/common/jsonquery"
	"github.com/amp-labs/connectors/common/naming"
	"github.com/amp-labs/connectors/common/paramsbuilder"
	"github.com/amp-labs/connectors/common/urlbuilder"
	"github.com/amp-labs/connectors/internal/deep"
	"github.com/amp-labs/connectors/providers"
	"github.com/spyzhov/ajson"
)

const apiVersion = "v9.2"

type Connector struct {
	deep.Clients
	deep.EmptyCloser
	// nolint:lll
	// Microsoft API supports other capabilities like filtering, grouping, and sorting which we can potentially tap into later.
	// See https://learn.microsoft.com/en-us/power-apps/developer/data-platform/webapi/query-data-web-api#odata-query-options
	deep.Reader
	deep.Writer
	deep.Remover
}

type parameters struct {
	paramsbuilder.Client
	paramsbuilder.Workspace
}

func NewConnector(opts ...Option) (*Connector, error) {
	constructor := func(
		clients *deep.Clients,
		closer *deep.EmptyCloser,
		reader *deep.Reader,
		writer *deep.Writer,
		remover *deep.Remover,
	) *Connector {
		return &Connector{
			Clients:     *clients,
			EmptyCloser: *closer,
			Reader:      *reader,
			Writer:      *writer,
			Remover:     *remover,
		}
	}
	errorHandler := interpreter.ErrorHandler{
		JSON: interpreter.NewFaultyResponder(errorFormats, nil),
	}
	firstPage := deep.FirstPageBuilder{
		Build: func(config common.ReadParams, url *urlbuilder.URL) (*urlbuilder.URL, error) {
			fields := config.Fields.List()
			if len(fields) != 0 {
				url.WithQueryParam("$select", strings.Join(fields, ","))
			}

			return url, nil
		},
	}
	nextPage := deep.NextPageBuilder{
		Build: func(config common.ReadParams, previousPage *urlbuilder.URL, node *ajson.Node) (string, error) {
			nextLink, err := jsonquery.New(node).StrWithDefault("@odata.nextLink", "")
			if err != nil {
				return "", err
			}

			if len(nextLink) != 0 {
				return nextLink, nil
			}

			return "", nil
		},
	}
	readRequestBuilder := deep.GetWithHeadersRequestBuilder{
		Headers: []common.Header{
			{
				Key:   "Prefer",
				Value: fmt.Sprintf("odata.maxpagesize=%v", DefaultPageSize),
			},
			{
				Key:   "Prefer",
				Value: `odata.include-annotations="*"`,
			},
		},
	}
	readObjectLocator := deep.ReadObjectLocator{
		Locate: func(config common.ReadParams) string {
			return "value"
		},
	}
	urlResolver := deep.SingleURLFormat{
		Produce: func(method deep.Method, baseURL, objectName string) (*urlbuilder.URL, error) {
			// Despite the "Method" type the relationship between objectName and
			// URL path is that it must be in singular word case.
			// Ex: objectName=Orders, then url will be http://base/v9.2/Order
			path := naming.NewSingularString(objectName)

			return constructURL(baseURL, apiVersion, path.String())
		},
	}
	writeResultBuilder := deep.WriteResultBuilder{
		Build: func(config common.WriteParams, body *ajson.Node) (*common.WriteResult, error) {
			// Neither Post nor Patch return any response data on successful completion
			// Both complete with 204 NoContent
			return &common.WriteResult{
				Success: true,
			}, nil
		},
	}

	return deep.Connector[Connector, parameters](constructor, providers.DynamicsCRM, opts,
		errorHandler,
		firstPage,
		nextPage,
		readRequestBuilder,
		readObjectLocator,
		urlResolver,
		customWriterRequestBuilder{},
		writeResultBuilder,
		customRemoveRequestBuilder{},
	)
}
