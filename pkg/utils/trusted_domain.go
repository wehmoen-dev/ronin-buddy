package utils

import (
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/url"
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

	parsedDomain, err := url.Parse(domain)

	if err != nil {
		return false
	}

	resp, err := c.c.R().Get(parsedDomain.Hostname())

	if err != nil {
		return false
	}
	return resp.StatusCode() == http.StatusOK
}
