package providers

// Supported Microsoft Products includes OneDrive Outlook Excel
// Edge Extensions Sharepoint OneNote Notifications Todos Teams Insights
// Planner and Personal Contacts.
const Microsoft Provider = "microsoft"

func init() {
	// Microsoft Office 365 Configuration
	SetInfo(Microsoft, ProviderInfo{
		DisplayName: "Microsoft",
		AuthType:    Oauth2,
		BaseURL:     "https://graph.microsoft.com",
		Oauth2Opts: &Oauth2Opts{
			GrantType:                 AuthorizationCode,
			AuthURL:                   "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
			TokenURL:                  "https://login.microsoftonline.com/common/oauth2/v2.0/token",
			ExplicitScopesRequired:    true,
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
