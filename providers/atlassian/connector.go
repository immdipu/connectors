package atlassian

import (
	"errors"
	"fmt"
	"time"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/interpreter"
	"github.com/amp-labs/connectors/common/jsonquery"
	"github.com/amp-labs/connectors/common/paramsbuilder"
	"github.com/amp-labs/connectors/common/urlbuilder"
	"github.com/amp-labs/connectors/internal/deep"
	"github.com/amp-labs/connectors/internal/deep/dpobjects"
	"github.com/amp-labs/connectors/internal/deep/dpread"
	"github.com/amp-labs/connectors/internal/deep/dpvars"
	"github.com/amp-labs/connectors/internal/deep/dpwrite"
	"github.com/amp-labs/connectors/providers"
	"github.com/spyzhov/ajson"
)

// ErrMissingCloudId happens when cloud id was not provided via WithMetadata.
var ErrMissingCloudId = errors.New("connector missing cloud id")

type Connector struct {
	Data dpvars.ConnectorData[parameters, *AuthMetadataVars]
	deep.Clients
	deep.EmptyCloser
	deep.Reader
	// Write will either create or update a Jira issue.
	// Create issue docs:
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-post
	// Update issue docs:
	// https://developer.atlassian.com/cloud/jira/platform/rest/v3/api-group-issues/#api-rest-api-3-issue-issueidorkey-put
	deep.Writer
	deep.Remover

	urlBuilder *customURLBuilder
}

type parameters struct {
	paramsbuilder.Client
	paramsbuilder.Workspace
	paramsbuilder.Module
	paramsbuilder.Metadata
}

func NewConnector(opts ...Option) (*Connector, error) {
	constructor := func(
		clients *deep.Clients,
		closer *deep.EmptyCloser,
		data *dpvars.ConnectorData[parameters, *AuthMetadataVars],
		urlResolver dpobjects.URLResolver,
		reader *deep.Reader,
		writer *deep.Writer,
		remover *deep.Remover,
	) *Connector {
		return &Connector{
			Data:        *data,
			Clients:     *clients,
			EmptyCloser: *closer,
			Reader:      *reader,
			Writer:      *writer,
			Remover:     *remover,
			urlBuilder:  urlResolver.(*customURLBuilder),
		}
	}
	errorHandler := interpreter.ErrorHandler{
		JSON: interpreter.NewFaultyResponder(errorFormats, nil),
	}
	firstPage := dpread.FirstPageBuilder{
		Build: func(config common.ReadParams, url *urlbuilder.URL) (*urlbuilder.URL, error) {
			if !config.Since.IsZero() {
				// Read URL supports time scoping. common.ReadParams.Since is used to get relative time frame.
				// Here is an API example on how to request issues that were updated in the last 30 minutes.
				// search?jql=updated > "-30m"
				// The reason we use minutes is that it is the most granular API permits.
				diff := time.Since(config.Since)

				minutes := int64(diff.Minutes())
				if minutes > 0 {
					url.WithQueryParam("jql", fmt.Sprintf(`updated > "-%vm"`, minutes))
				}
			}

			return url, nil
		},
	}
	nextPage := dpread.NextPageBuilder{
		Build: func(config common.ReadParams, url *urlbuilder.URL, node *ajson.Node) (string, error) {
			startAt, err := getNextRecords(node)
			if err != nil {
				return "", err
			}

			if len(startAt) != 0 {
				url.WithQueryParam("startAt", startAt)

				return url.String(), nil
			}

			return "", nil
		},
	}
	readObjectLocator := dpread.ResponseLocator{
		Locate: func(config common.ReadParams, node *ajson.Node) string {
			return "issues"
		},
		Process: flattenRecords,
	}
	writeResultBuilder := dpwrite.ResponseBuilder{
		Build: func(config common.WriteParams, body *ajson.Node) (*common.WriteResult, error) {
			recordID, err := jsonquery.New(body).Str("id", false)
			if err != nil {
				return nil, err
			}

			return &common.WriteResult{
				Success:  true,
				RecordId: *recordID,
				Errors:   nil,
				Data:     nil,
			}, nil
		},
	}

	return deep.ExtendedConnector[Connector, parameters, *AuthMetadataVars](
		constructor, providers.Atlassian, &AuthMetadataVars{}, opts,
		errorHandler,
		customURLBuilder{},
		firstPage,
		nextPage,
		readObjectLocator,
		writeResultBuilder,
	)
}
