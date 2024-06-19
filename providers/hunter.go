package providers

const Hunter Provider = "hunter"

func init() {
	// Hunter Connector Configuration
	SetInfo(Hunter, ProviderInfo{
		AuthType: ApiKey,
		BaseURL:  "https://api.hunter.io/",
		ApiKeyOpts: &ApiKeyOpts{
			Type:           InQuery,
			QueryParamName: "api_key",
			DocsURL:        "https://hunter.io/api-documentation#authentication",
		},
		Support: Support{
			BulkWrite: BulkWriteSupport{
				Insert: false,
				Update: false,
				Upsert: false,
				Delete: false,
			},
			Proxy:     false,
			Read:      false,
			Subscribe: false,
			Write:     false,
		},
	})
}
