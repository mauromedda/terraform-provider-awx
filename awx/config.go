package awx

import (
	awxgo "github.com/Colstuwjx/awx-go"
)

// Config of Ansible Tower/AWX
type Config struct {
	Username string
	Password string
	Endpoint string
}

// Client for Tower/AWX API v2

func (c *Config) Client() *awxgo.AWX {
	awx := awxgo.NewAWX(c.Endpoint, c.Username, c.Password, nil)
	return awx
}
