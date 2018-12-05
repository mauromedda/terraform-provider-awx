package awx

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider AWX provider implementation
func Provider() terraform.ResourceProvider {
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
			"ssl_verify": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: descriptions["ssl_verify"],
				Sensitive:   true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awx_inventory":         resourceInventoryObject(),
			"awx_inventory_group":   resourceInventoryGroupObject(),
			"awx_host":              resourceHostObject(),
			"awx_group_association": resourceGroupAssociationObject(),
			"awx_project":           resourceProjectObject(),
			"awx_job_template":      resourceJobTemplateObject(),
			"awx_user":              resourceUserObject(),
			"awx_team":              resourceTeamObject(),
			"awx_user_role":         resourceUserRoleObject(),
			"awx_team_role":         resourceTeamRoleObject(),
			"awx_organization":      resourceOrganizationObject(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	log.Printf("[INFO] Initializing Tower Client")

	config := &Config{
		Endpoint:  d.Get("endpoint").(string),
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
		Sslverify: d.Get("ssl_verify").(bool),
	}

	return config.Client(), nil
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"endpoint":   "The API Endpoint used to invoke Ansible Tower/AWX",
		"username":   "The Ansible Tower API Username",
		"password":   "The Ansible Tower API Password",
		"ssl_verify": "Skip SSL certificate check",
	}
}
