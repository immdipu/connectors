package providers

const Iterable Provider = "iterable"

func init() {
	// Iterable API Key authentication
	SetInfo(Iterable, ProviderInfo{
		DisplayName: "Iterable",
		AuthType:    ApiKey,
		BaseURL:     "https://api.iterable.com",
		ApiKeyOpts: &ApiKeyOpts{
			AttachmentType: Header,
			Header: &ApiKeyOptsHeader{
				Name: "Api-Key",
			},
			DocsURL: "https://app.iterable.com/settings/apiKeys",
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
		Media: &Media{
			DarkMode: &MediaTypeDarkMode{
				IconURL: "https://res.cloudinary.com/dycvts6vp/image/upload/v1722065197/media/iterable_1722065196.svg",
				LogoURL: "https://res.cloudinary.com/dycvts6vp/image/upload/v1722065173/media/iterable_1722065172.svg",
			},
			Regular: &MediaTypeRegular{
				IconURL: "https://res.cloudinary.com/dycvts6vp/image/upload/v1722065197/media/iterable_1722065196.svg",
				LogoURL: "https://res.cloudinary.com/dycvts6vp/image/upload/v1722065173/media/iterable_1722065172.svg",
			},
		},
	})
}
