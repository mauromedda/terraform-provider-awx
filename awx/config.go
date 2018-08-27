package awx

import (
	"crypto/tls"
	"net/http"

	awxgo "github.com/Colstuwjx/awx-go"
)

// Config of Ansible Tower/AWX
type Config struct {
	Username  string
	Password  string
	Endpoint  string
	Sslverify bool
}

// Client for Tower/AWX API v2
func (c *Config) Client() *awxgo.AWX {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	awx := awxgo.NewAWX(c.Endpoint, c.Username, c.Password, client)

	return awx
}
