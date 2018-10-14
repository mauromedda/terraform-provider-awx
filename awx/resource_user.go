package awx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceUserObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Delete: resourceUserDelete,
		Update: resourceUserUpdate,

		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Username of this user.",
			},

			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Password of this user.",
			},

			"email": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The email address for the user.",
			},

			"first_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "User first name.",
			},

			"last_name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "User last name.",
			},

			"is_system_auditor": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The user is a system auditor.",
			},
			"is_superuser": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "The user is a system administrator.",
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService
	_, res, err := awxService.ListUsers(map[string]string{
		"username": d.Get("username").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) >= 1 {
		return fmt.Errorf("User with name %s already exists",
			d.Get("username").(string))
	}

	result, err := awxService.CreateUser(map[string]interface{}{
		"username":          d.Get("username").(string),
		"password":          d.Get("password").(string),
		"email":             d.Get("email").(string),
		"first_name":        d.Get("first_name").(string),
		"last_name":         d.Get("last_name").(string),
		"is_superuser":      d.Get("is_superuser").(bool),
		"is_system_auditor": d.Get("is_system_auditor").(bool),
	}, map[string]string{})
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceUserRead(d, m)
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService
	_, res, err := awxService.ListUsers(map[string]string{
		"id": d.Id()},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("User with name %s doesn't exists",
			d.Get("username").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.UpdateUser(id, map[string]interface{}{
		"username":          d.Get("username").(string),
		"password":          d.Get("password").(string),
		"email":             d.Get("email").(string),
		"first_name":        d.Get("first_name").(string),
		"last_name":         d.Get("last_name").(string),
		"is_superuser":      d.Get("is_superuser").(bool),
		"is_system_auditor": d.Get("is_system_auditor").(bool),
	}, map[string]string{})
	if err != nil {
		return err
	}

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService
	_, res, err := awxService.ListUsers(map[string]string{
		"username": d.Get("username").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setUserResourceData(d, res.Results[0])
	return nil
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService
	id, err := strconv.Atoi(d.Id())
	_, res, err := awxService.ListUsers(map[string]string{
		"username": d.Get("username").(string)})
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	if _, err = awxService.DeleteUser(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setUserResourceData(d *schema.ResourceData, r *awxgo.User) *schema.ResourceData {
	d.Set("username", r.Username)
	d.Set("password", d.Get("password").(string))
	d.Set("email", r.Email)
	d.Set("first_name", r.FirstName)
	d.Set("last_name", r.LastName)
	d.Set("is_superuser", r.IsSuperUser)
	d.Set("is_system_auditor", r.IsSystemAuditor)
	return d
}
