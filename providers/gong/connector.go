package gong

import (
	"github.com/amp-labs/connectors/common/interpreter"
	"github.com/amp-labs/connectors/common/paramsbuilder"
	"github.com/amp-labs/connectors/common/urlbuilder"
	"github.com/amp-labs/connectors/internal/deep"
	"github.com/amp-labs/connectors/providers"
	"github.com/amp-labs/connectors/providers/gong/metadata"
	"github.com/amp-labs/connectors/tools/scrapper"
)

const ApiVersion = "v2"

type Connector struct {
	*deep.Clients
	*deep.EmptyCloser
	*deep.StaticMetadata
}

type parameters struct {
	paramsbuilder.Client
}

func NewConnector(opts ...Option) (*Connector, error) {
	constructor := func(
		clients *deep.Clients,
		closer *deep.EmptyCloser,
		staticMetadata *deep.StaticMetadata) *Connector {
		return &Connector{
			Clients:        clients,
			EmptyCloser:    closer,
			StaticMetadata: staticMetadata,
		}
	}
	errorHandler := interpreter.ErrorHandler{
		JSON: interpreter.NewFaultyResponder(errorFormats, nil),
	}

	return deep.Connector[Connector, parameters](constructor, providers.Gong, &errorHandler, opts,
		deep.Dependency{
			Constructor: func() *scrapper.ObjectMetadataResult {
				return metadata.Schemas
			},
		},
	)
}

func (c *Connector) getURL(arg string) (*urlbuilder.URL, error) {
	return urlbuilder.New(c.BaseURL(), ApiVersion, arg)
}
