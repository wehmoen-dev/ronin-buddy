package validate

import (
	"github.com/sethvargo/go-githubactions"
	"net/url"
)

func isValidPlatform(platform Platform) bool {
	githubactions.Debugf("Validating platform validity: %s", platform)
	switch platform {
	case PlatformFacebook, PlatformInstagram, PlatformTwitter, PlatformX, PlatformLinkedIn, PlatformThreads, PlatformMastodon, PlatformTelegram, PlatformDiscord:
		return true
	default:
		return false
	}
}

func isValidCategory(category Category) bool {
	githubactions.Debugf("Validating category validity: %s", category)
	switch category {
	case CategoryGame, CategoryNFT, CategoryFinance, CategoryDAO, CategoryTool, CategoryOther, CategoryBot:
		return true
	default:
		return false
	}
}

func isValidPlatformHost(platform Platform, url *url.URL) bool {
	githubactions.Debugf("Validating platform host: %s, %s", platform, url.Host)
	switch platform {
	case PlatformFacebook:
		if url.Host != "facebook.com" && url.Host != "www.facebook.com" && url.Host != "fb.me" {
			return false
		}
	case PlatformInstagram:
		if url.Host != "instagram.com" && url.Host != "www.instagram.com" && url.Host != "instagr.am" {
			return false
		}
	case PlatformTwitter, PlatformX:
		if url.Host != "twitter.com" && url.Host != "www.twitter.com" && url.Host != "x.com" && url.Host != "www.x.com" && url.Host != "t.co" {
			return false
		}
	case PlatformLinkedIn:
		if url.Host != "linkedin.com" && url.Host != "www.linkedin.com" && url.Host != "lnkd.in" {
			return false
		}
	case PlatformThreads:
		if url.Host != "threads.net" && url.Host != "www.threads.net" {
			return false
		}
	case PlatformMastodon:
		if url.Host != "mastodon.social" && url.Host != "www.mastodon.social" {
			return false
		}
	case PlatformTelegram:
		if url.Host != "t.me" && url.Host != "www.t.me" {
			return false
		}
	case PlatformDiscord:
		if url.Host != "discord.com" && url.Host != "www.discord.com" && url.Host != "discord.gg" && url.Host != "www.discord.gg" {
			return false
		}
	default:
		return false
	}

	return true
}
