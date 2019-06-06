package tower

import (
	"fmt"
	"time"

	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCreate,
		Read:   resourceCredentialRead,
		Update: resourceCredentialUpdate,
		Delete: resourceCredentialDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ssh",
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"security_token": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ssh_key_data": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"ssh_key_unlock": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"become_method": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "One of \"\", sudo, su, pbrun, pfexec",
			},
			"become_username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"become_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"vault_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"subscription": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"tenant": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"secret": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"client": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"authorize": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  "",
			},
			"authorize_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceCredentialCreate(d *schema.ResourceData, meta interface{}) error {
	awx := m.(*awxgo.AWX)
	awxServiceCredential := awx.CredentialService

	request, err := awxServiceCredential.
	if err != nil {
		return err
	}
	i, err := service.Create(request)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(i.ID))
	return resourceCredentialRead(d, meta)
}

func resourceCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Credentials

	if r, err := service.GetByID(d.Id()); err != nil {
		return err
	} else {
		d = setCredentialResourceData(d, r)
	}
	return nil
}

func resourceCredentialUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Credentials
	request, err := buildCredential(d, meta)
	if err != nil {
		return err
	}
	if _, err := service.Update(request); err != nil {
		return err
	}
	return resourceCredentialRead(d, meta)
}

func resourceCredentialDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*towerapi.Client)
	service := client.Credentials
	if err := service.Delete(d.Id()); err != nil {
		return fmt.Errorf("Failed to delete (%s): %s", d.Id(), err)
	}
	return nil
}

func setCredentialResourceData(d *schema.ResourceData, r *credentials.Credential) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("kind", r.Kind)
	d.Set("host", r.Host)
	d.Set("username", r.Username)
	d.Set("password", r.Password)
	d.Set("securityToken", r.SecurityToken)
	d.Set("project_id", r.Project)
	d.Set("domain", r.Domain)
	d.Set("ssh_key_data", r.SSHKeyData)
	d.Set("ssh_key_unlock", r.SSHKeyUnlock)
	d.Set("organization_id", r.Organization)
	d.Set("become_method", r.BecomeMethod)
	d.Set("become_username", r.BecomeUsername)
	d.Set("become_password", r.BecomePassword)
	d.Set("vault_password", r.VaultPassword)
	d.Set("subscription_id", r.Subscription)
	d.Set("tenant", r.Tenant)
	d.Set("secret", r.Secret)
	d.Set("client", r.Client)
	d.Set("authorize", r.Authorize)
	d.Set("authorize_password", r.AuthorizePassword)
	return d
}

func buildCredential(d *schema.ResourceData, meta interface{}) (*credentials.Request, error) {

	request := &credentials.Request{
		Name:              d.Get("name").(string),
		Description:       d.Get("description").(string),
		Kind:              d.Get("kind").(string),
		Host:              d.Get("host").(string),
		Username:          d.Get("username").(string),
		Password:          d.Get("password").(string),
		SecurityToken:     d.Get("security_token").(string),
		Project:           d.Get("project_id").(string),
		Domain:            d.Get("domain").(string),
		SSHKeyData:        d.Get("ssh_key_data").(string),
		SSHKeyUnlock:      d.Get("ssh_key_unlock").(string),
		Organization:      AtoipOr(d.Get("organization_id").(string), nil),
		BecomeMethod:      d.Get("become_method").(string),
		BecomeUsername:    d.Get("become_username").(string),
		BecomePassword:    d.Get("become_password").(string),
		VaultPassword:     d.Get("vault_password").(string),
		Subscription:      d.Get("subscription").(string),
		Tenant:            d.Get("tenant").(string),
		Secret:            d.Get("secret").(string),
		Client:            d.Get("client").(string),
		Authorize:         d.Get("authorize").(bool),
		AuthorizePassword: d.Get("authorize_password").(string),
	}

	return request, nil
}
