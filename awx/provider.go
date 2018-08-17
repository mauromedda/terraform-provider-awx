package awx

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"TOWER_ENDPOINT",
					"AWX_ENDPOINT",
				}, "http://localhost"),
				Description: descriptions["endpoint"],
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"TOWER_USERNAME",
					"AWX_USERNAME",
				}, "admin"),
				Description: descriptions["username"],
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"TOWER_PASSWORD",
					"AWX_PASSWORD",
				}, "password"),
				Description: descriptions["password"],
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awx_inventory": resourceInventoryObject(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	log.Printf("[INFO] Initializing Tower Client")

	config := &Config{
		Endpoint: d.Get("endpoint").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	return config.Client(), nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"endpoint": "The API Endpoint used to invoke Ansible Tower/AWX",
		"username": "The Ansible Tower API Username",
		"password": "The Ansible Tower API Password",
	}
}
