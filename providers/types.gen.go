// Package providers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package providers

// Defines values for ApiKeyAsBasicOptsFieldUsed.
const (
	PasswordField ApiKeyAsBasicOptsFieldUsed = "password"
	UsernameField ApiKeyAsBasicOptsFieldUsed = "username"
)

// Defines values for ApiKeyOptsAttachmentType.
const (
	Header ApiKeyOptsAttachmentType = "header"
	Query  ApiKeyOptsAttachmentType = "query"
)

// Defines values for AuthType.
const (
	ApiKey AuthType = "apiKey"
	Basic  AuthType = "basic"
	None   AuthType = "none"
	Oauth2 AuthType = "oauth2"
)

// Defines values for Oauth2OptsGrantType.
const (
	AuthorizationCode Oauth2OptsGrantType = "authorizationCode"
	ClientCredentials Oauth2OptsGrantType = "clientCredentials"
	PKCE              Oauth2OptsGrantType = "PKCE"
	Password          Oauth2OptsGrantType = "password"
)

// ApiKeyAsBasicOpts when this object is present, it means that this provider uses Basic Auth to actually collect an API key
type ApiKeyAsBasicOpts struct {
	// FieldUsed whether the API key should be used as the username or password.
	FieldUsed ApiKeyAsBasicOptsFieldUsed `json:"fieldUsed,omitempty"`

	// KeyFormat How to transform the API key in to a basic auth user:pass string. The %s is replaced with the API key value.
	KeyFormat string `json:"keyFormat,omitempty"`
}

// ApiKeyAsBasicOptsFieldUsed whether the API key should be used as the username or password.
type ApiKeyAsBasicOptsFieldUsed string

// ApiKeyOpts Configuration for API key. Must be provided if authType is apiKey.
type ApiKeyOpts struct {
	// AttachmentType How the API key should be attached to requests.
	AttachmentType ApiKeyOptsAttachmentType `json:"attachmentType" validate:"required"`

	// DocsURL URL with more information about how to get or use an API key.
	DocsURL string `json:"docsURL,omitempty"`

	// Header Configuration for API key in header. Must be provided if type is in-header.
	Header *ApiKeyOptsHeader `json:"header,omitempty"`

	// Query Configuration for API key in query parameter. Must be provided if type is in-query.
	Query *ApiKeyOptsQuery `json:"query,omitempty"`
}

// ApiKeyOptsAttachmentType How the API key should be attached to requests.
type ApiKeyOptsAttachmentType string

// ApiKeyOptsHeader Configuration for API key in header. Must be provided if type is in-header.
type ApiKeyOptsHeader struct {
	// Name The name of the header to be used for the API key.
	Name string `json:"name"`

	// ValuePrefix The prefix to be added to the API key value when it is sent in the header.
	ValuePrefix string `json:"valuePrefix,omitempty"`
}

// ApiKeyOptsQuery Configuration for API key in query parameter. Must be provided if type is in-query.
type ApiKeyOptsQuery struct {
	// Name The name of the query parameter to be used for the API key.
	Name string `json:"name"`
}

// AuthType The type of authentication required by the provider.
type AuthType string

// BasicAuthOpts Configuration for Basic Auth. Optional.
type BasicAuthOpts struct {
	// ApiKeyAsBasic If true, the provider uses an API key which then gets encoded as a basic auth user:pass string.
	ApiKeyAsBasic bool `json:"apiKeyAsBasic,omitempty"`

	// ApiKeyAsBasicOpts when this object is present, it means that this provider uses Basic Auth to actually collect an API key
	ApiKeyAsBasicOpts *ApiKeyAsBasicOpts `json:"apiKeyAsBasicOpts,omitempty"`

	// DocsURL URL with more information about how to get or use an API key.
	DocsURL string `json:"docsURL,omitempty"`
}

// BulkWriteSupport defines model for BulkWriteSupport.
type BulkWriteSupport struct {
	Delete bool `json:"delete"`
	Insert bool `json:"insert"`
	Update bool `json:"update"`
	Upsert bool `json:"upsert"`
}

// CatalogType defines model for CatalogType.
type CatalogType map[string]ProviderInfo

// CatalogWrapper defines model for CatalogWrapper.
type CatalogWrapper struct {
	Catalog CatalogType `json:"catalog"`

	// Timestamp An RFC3339 formatted timestamp of when the catalog was generated.
	Timestamp string `json:"timestamp" validate:"required"`
}

// Media defines model for Media.
type Media struct {
	// DarkMode Media to be used in dark mode.
	DarkMode *MediaTypeDarkMode `json:"darkMode,omitempty"`

	// Regular Media for light/regular mode.
	Regular *MediaTypeRegular `json:"regular,omitempty"`
}

// MediaTypeDarkMode Media to be used in dark mode.
type MediaTypeDarkMode struct {
	// IconURL URL to the icon for the provider that is to be used in dark mode.
	IconURL string `json:"iconURL,omitempty"`

	// LogoURL URL to the logo for the provider that is to be used in dark mode.
	LogoURL string `json:"logoURL,omitempty"`
}

// MediaTypeRegular Media for light/regular mode.
type MediaTypeRegular struct {
	// IconURL URL to the icon for the provider.
	IconURL string `json:"iconURL,omitempty"`

	// LogoURL URL to the logo for the provider.
	LogoURL string `json:"logoURL,omitempty"`
}

// Oauth2Opts Configuration for OAuth2.0. Must be provided if authType is oauth2.
type Oauth2Opts struct {
	// Audience A list of URLs that represent the audience for the token, which is needed for some client credential grant flows.
	Audience []string `json:"audience,omitempty"`

	// AuthURL The authorization URL.
	AuthURL       string            `json:"authURL,omitempty"`
	AuthURLParams map[string]string `json:"authURLParams,omitempty"`

	// DocsURL URL with more information about where to retrieve Client ID and Client Secret, etc.
	DocsURL string `json:"docsURL,omitempty"`

	// ExplicitScopesRequired Whether scopes are required to be known ahead of the OAuth flow.
	ExplicitScopesRequired bool `json:"explicitScopesRequired"`

	// ExplicitWorkspaceRequired Whether the workspace is required to be known ahead of the OAuth flow.
	ExplicitWorkspaceRequired bool                `json:"explicitWorkspaceRequired"`
	GrantType                 Oauth2OptsGrantType `json:"grantType"`

	// TokenMetadataFields Fields to be used to extract token metadata from the token response.
	TokenMetadataFields TokenMetadataFields `json:"tokenMetadataFields"`

	// TokenURL The token URL.
	TokenURL string `json:"tokenURL" validate:"required"`
}

// Oauth2OptsGrantType defines model for Oauth2Opts.GrantType.
type Oauth2OptsGrantType string

// Provider defines model for Provider.
type Provider = string

// ProviderInfo defines model for ProviderInfo.
type ProviderInfo struct {
	// ApiKeyOpts Configuration for API key. Must be provided if authType is apiKey.
	ApiKeyOpts *ApiKeyOpts `json:"apiKeyOpts,omitempty"`

	// AuthType The type of authentication required by the provider.
	AuthType AuthType `json:"authType" validate:"required"`

	// BaseURL The base URL for making API requests.
	BaseURL string `json:"baseURL" validate:"required"`

	// BasicOpts Configuration for Basic Auth. Optional.
	BasicOpts *BasicAuthOpts `json:"basicOpts,omitempty"`

	// DisplayName The display name of the provider, if omitted, defaults to provider name.
	DisplayName string `json:"displayName,omitempty"`
	Media       *Media `json:"media,omitempty"`
	Name        string `json:"name"`

	// Oauth2Opts Configuration for OAuth2.0. Must be provided if authType is oauth2.
	Oauth2Opts *Oauth2Opts `json:"oauth2Opts,omitempty"`

	// PostAuthInfoNeeded If true, we require additional information after auth to start making requests.
	PostAuthInfoNeeded bool `json:"postAuthInfoNeeded,omitempty"`

	// ProviderOpts Additional provider-specific metadata.
	ProviderOpts ProviderOpts `json:"providerOpts"`

	// Support The supported features for the provider.
	Support Support `json:"support" validate:"required"`
}

// ProviderOpts Additional provider-specific metadata.
type ProviderOpts map[string]string

// Support The supported features for the provider.
type Support struct {
	BulkWrite BulkWriteSupport `json:"bulkWrite" validate:"required"`
	Proxy     bool             `json:"proxy"`
	Read      bool             `json:"read"`
	Subscribe bool             `json:"subscribe"`
	Write     bool             `json:"write"`
}

// TokenMetadataFields Fields to be used to extract token metadata from the token response.
type TokenMetadataFields struct {
	ConsumerRefField  string `json:"consumerRefField,omitempty"`
	ScopesField       string `json:"scopesField,omitempty"`
	WorkspaceRefField string `json:"workspaceRefField,omitempty"`
}
