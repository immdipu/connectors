package providers

const Paypal Provider = "paypal"

func init() {
	SetInfo(Paypal, ProviderInfo{
		DisplayName: "Paypal",
		AuthType:    Oauth2,
		BaseURL:     "https://api-m.paypal.com",

		Oauth2Opts: &Oauth2Opts{
			GrantType:                 ClientCredentials,
			DocsURL:                   "https://developer.paypal.com/api/rest",
			TokenURL:                  "https://api-m.paypal.com/v1/oauth2/token",
			ExplicitScopesRequired:    false,
			ExplicitWorkspaceRequired: false,
		},

		Support: Support{
			BulkWrite: BulkWriteSupport{
				Insert: false,
				Update: false,
				Upsert: false,
				Delete: false,
			},
			Proxy:     true,
			Read:      false,
			Subscribe: false,
			Write:     false,
		},
	})
}
