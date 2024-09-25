package utils

import (
	"github.com/go-resty/resty/v2"
	"net/http"
)

type TrustedDomainClient struct {
	c *resty.Client
}

func NewTrustedDomainClient() *TrustedDomainClient {
	c := resty.New()
	c.SetBaseURL("https://wallet-manager.skymavis.com/wallet-manager/v2/public/trusted-domains/")
	c.SetHeader("user-agent", "ronin-buddy")

	return &TrustedDomainClient{
		c: c,
	}
}

func (c *TrustedDomainClient) IsWhitelisted(domain string) bool {
	resp, err := c.c.R().Get(domain)

	if err != nil {
		return false
	}
	return resp.StatusCode() == http.StatusOK
}
