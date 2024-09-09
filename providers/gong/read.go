package gong

import (
	"context"

	"github.com/amp-labs/connectors/common"
)

func (c *Connector) Read(ctx context.Context, config common.ReadParams) (*common.ReadResult, error) {
	if err := config.ValidateParams(true); err != nil {
		return nil, err
	}

	url, err := c.getURL(config.ObjectName)
	if err != nil {
		return nil, err
	}

	if len(config.NextPage) != 0 { // not the first page, add a cursor
		url.WithQueryParam("cursor", config.NextPage.String())
	}

	res, err := c.Client.Get(ctx, url.String())
	if err != nil {
		return nil, err
	}

	return common.ParseResult(res,
		common.GetRecordsUnderJSONPath(config.ObjectName),
		getNextRecordsURL,
		common.GetMarshaledData,
		config.Fields,
	)
}
