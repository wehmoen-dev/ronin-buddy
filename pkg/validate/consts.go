package validate

const (
	PlatformFacebook  Platform = "facebook"
	PlatformInstagram Platform = "instagram"
	PlatformTwitter   Platform = "twitter"
	PlatformX         Platform = "x"
	PlatformLinkedIn  Platform = "linkedin"
	PlatformThreads   Platform = "threads"
	PlatformMastodon  Platform = "mastodon"
	PlatformTelegram  Platform = "telegram"
	PlatformDiscord   Platform = "discord"
)

const (
	CategoryGame    Category = "game"
	CategoryNFT     Category = "nft"
	CategoryFinance Category = "finance"
	CategoryDAO     Category = "dao"
	CategoryTool    Category = "tool"
	CategoryBot     Category = "bot"
	CategoryOther   Category = "other"
)

const (
	ImageMinSideLength int   = 256
	ImageMaxSideLength int   = 1024
	ImageMaxFileSize   int64 = 1024 * 1024
)

const (
	ResultFailed string = "❌ Validation failed"
	ResultPassed string = "✅ Validation passed"
)
