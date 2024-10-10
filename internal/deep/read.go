package deep

import (
	"context"
	"errors"

	"github.com/amp-labs/connectors/common"
	"github.com/amp-labs/connectors/common/handy"
	"github.com/amp-labs/connectors/common/jsonquery"
	"github.com/amp-labs/connectors/common/urlbuilder"
	"github.com/amp-labs/connectors/internal/deep/requirements"
	"github.com/spyzhov/ajson"
)

type Reader struct {
	urlResolver       ObjectURLResolver
	pageStartBuilder  PaginationStartBuilder
	nextPageBuilder   NextPageBuilder
	readObjectLocator ReadObjectLocator
	objectManager     ObjectManager
	requestBuilder    ReadRequestBuilder
	headerSupplements HeaderSupplements

	clients Clients
}

func NewReader(clients *Clients,
	resolver ObjectURLResolver,
	pageStartBuilder PaginationStartBuilder,
	nextPageBuilder *NextPageBuilder,
	objectLocator *ReadObjectLocator,
	objectManager ObjectManager,
	requestBuilder ReadRequestBuilder,
	headerSupplements *HeaderSupplements,
) *Reader {
	return &Reader{
		urlResolver:       resolver,
		pageStartBuilder:  pageStartBuilder,
		nextPageBuilder:   *nextPageBuilder,
		readObjectLocator: *objectLocator,
		objectManager:     objectManager,
		requestBuilder:    requestBuilder,
		headerSupplements: *headerSupplements,
		clients:           *clients,
	}
}

func (r *Reader) Read(ctx context.Context, config common.ReadParams) (*common.ReadResult, error) {
	if err := config.ValidateParams(true); err != nil {
		return nil, err
	}

	if !r.objectManager.IsReadSupported(config.ObjectName) {
		return nil, common.ErrOperationNotSupportedForObject
	}

	url, err := r.buildReadURL(config)
	if err != nil {
		return nil, err
	}

	read, headers := r.requestBuilder.MakeReadRequest(config.ObjectName, r.clients)
	headers = append(headers, r.headerSupplements.ReadHeaders()...)

	rsp, err := read(ctx, url, nil, headers...)
	if err != nil {
		return nil, err
	}

	recordsFunc, err := r.readObjectLocator.getRecordsFunc(config)
	if err != nil {
		return nil, err
	}

	nextPageFunc, err := r.nextPageBuilder.getNextPageFunc(config, url)
	if err != nil {
		return nil, err
	}

	return common.ParseResult(
		rsp,
		recordsFunc,
		nextPageFunc,
		common.GetMarshaledData,
		config.Fields,
	)
}

func (r *Reader) buildReadURL(config common.ReadParams) (*urlbuilder.URL, error) {
	if len(config.NextPage) != 0 {
		// Next page
		return urlbuilder.New(config.NextPage.String())
	}

	// First page
	url, err := r.urlResolver.FindURL(ReadMethod, r.clients.BaseURL(), config.ObjectName)
	if err != nil {
		return nil, err
	}

	return r.pageStartBuilder.FirstPage(config, url)
}

type PaginationStartBuilder interface {
	requirements.ConnectorComponent
	FirstPage(config common.ReadParams, url *urlbuilder.URL) (*urlbuilder.URL, error)
}

var _ PaginationStartBuilder = FirstPageBuilder{}

type FirstPageBuilder struct {
	Build func(config common.ReadParams, url *urlbuilder.URL) (*urlbuilder.URL, error)
}

func (b FirstPageBuilder) FirstPage(config common.ReadParams, url *urlbuilder.URL) (*urlbuilder.URL, error) {
	if b.Build == nil {
		// TODO error
		return nil, errors.New("build method cannot be empty")
	}

	return b.Build(config, url)
}

func (b FirstPageBuilder) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "paginationStartBuilder",
		Constructor: handy.Returner(b),
		Interface:   new(PaginationStartBuilder),
	}
}

var _ PaginationStartBuilder = DefaultPageBuilder{}

type DefaultPageBuilder struct{}

func (b DefaultPageBuilder) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "paginationStartBuilder",
		Constructor: handy.Returner(b),
		Interface:   new(PaginationStartBuilder),
	}
}

func (b DefaultPageBuilder) FirstPage(config common.ReadParams, url *urlbuilder.URL) (*urlbuilder.URL, error) {
	return url, nil
}

type NextPageBuilder struct {
	Build func(config common.ReadParams, previousPage *urlbuilder.URL, node *ajson.Node) (string, error)
}

func (b NextPageBuilder) getNextPageFunc(config common.ReadParams, url *urlbuilder.URL) (common.NextPageFunc, error) {
	if b.Build == nil {
		// TODO error
		return nil, errors.New("build method cannot be empty")
	}

	return func(node *ajson.Node) (string, error) {
		return b.Build(config, url, node)
	}, nil
}

func (b NextPageBuilder) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "nextPageBuilder",
		Constructor: handy.Returner(b),
	}
}

type ReadObjectLocator struct {
	// Locate should return the fieldName where desired list of Objects is located.
	Locate func(config common.ReadParams, node *ajson.Node) string
	// FlattenRecords is optional and will be used after list was located and extra processing is needed.
	// The desired fields could be nested
	FlattenRecords func(arr []*ajson.Node) ([]map[string]any, error)
}

func (l ReadObjectLocator) getRecordsFunc(config common.ReadParams) (common.RecordsFunc, error) {
	if l.Locate == nil {
		// TODO error
		return nil, errors.New("locate method cannot be empty")
	}

	return func(node *ajson.Node) ([]map[string]any, error) {
		fieldName := l.Locate(config, node)

		arr, err := jsonquery.New(node).Array(fieldName, false)
		if err != nil {
			return nil, err
		}

		if l.FlattenRecords != nil {
			return l.FlattenRecords(arr)
		}

		return jsonquery.Convertor.ArrayToMap(arr)
	}, nil
}

func (l ReadObjectLocator) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "readObjectLocator",
		Constructor: handy.Returner(l),
	}
}

type ReadRequestBuilder interface {
	requirements.ConnectorComponent

	MakeReadRequest(objectName string, clients Clients) (common.ReadMethod, []common.Header)
}

var _ ReadRequestBuilder = GetRequestBuilder{}

type GetRequestBuilder struct {
	simpleGetReadRequest
}

func (b GetRequestBuilder) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "readRequestBuilder",
		Constructor: handy.Returner(b),
		Interface:   new(ReadRequestBuilder),
	}
}

var _ ReadRequestBuilder = GetWithHeadersRequestBuilder{}

type GetWithHeadersRequestBuilder struct {
	delegate simpleGetReadRequest
	Headers  []common.Header
}

func (b GetWithHeadersRequestBuilder) MakeReadRequest(
	objectName string, clients Clients,
) (common.ReadMethod, []common.Header) {
	method, _ := b.delegate.MakeReadRequest(objectName, clients)

	return method, b.Headers
}

func (b GetWithHeadersRequestBuilder) Satisfies() requirements.Dependency {
	return requirements.Dependency{
		ID:          "readRequestBuilder",
		Constructor: handy.Returner(b),
		Interface:   new(ReadRequestBuilder),
	}
}

type simpleGetReadRequest struct{}

func (simpleGetReadRequest) MakeReadRequest(
	objectName string, clients Clients,
) (common.ReadMethod, []common.Header) {
	// Wrapper around GET without request body.
	return func(ctx context.Context, url *urlbuilder.URL,
		body any, headers ...common.Header,
	) (*common.JSONHTTPResponse, error) {
		return clients.JSON.Get(ctx, url.String(), headers...)
	}, nil
}
