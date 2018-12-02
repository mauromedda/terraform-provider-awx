package awx

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
)

func resourceUserRoleObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserRoleGrant,
		Read:   resourceUserRoleRead,
		Delete: resourceUserRoleRevoke,

		Schema: map[string]*schema.Schema{
			"user_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"organization_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validResourceTypes := map[string]bool{"admin": true, "read": true, "use": true,
						"member": true, "execute": true, "adhoc": true, "update": true, "auditor": true}
					value := v.(string)
					if !validResourceTypes[value] {
						errors = append(errors, fmt.Errorf("%q must match one of the valid vaules", k))
					}
					return
				},
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					validResourceTypes := map[string]bool{"inventory": true, "team": true, "organization": true,
						"job_template": true, "credential": true, "project": true}
					value := v.(string)
					if !validResourceTypes[value] {
						errors = append(errors, fmt.Errorf("%q must match one of inventory, team, organization, job_template, credential or project", k))
					}
					return
				},
			},
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func getRoleID(d *schema.ResourceData, m interface{}) (int, error) {
	awx := m.(*awxgo.AWX)
	switch d.Get("resource_type").(string) {
	case "inventory":
		awxService := awx.InventoriesService
		obj, _, err := awxService.ListInventories(map[string]string{
			"name":         d.Get("resource_name").(string),
			"organization": d.Get("organization_id").(string),
		})
		if err != nil {
			return 0, err
		}
		if d.Get("role").(string) == "admin" {
			return obj[0].SummaryFields.ObjectRoles.AdminRole.ID, nil
		} else if d.Get("role").(string) == "use" {
			return obj[0].SummaryFields.ObjectRoles.UseRole.ID, nil
		} else if d.Get("role").(string) == "read" {
			return obj[0].SummaryFields.ObjectRoles.ReadRole.ID, nil
		} else if d.Get("role").(string) == "update" {
			return obj[0].SummaryFields.ObjectRoles.UpdateRole.ID, nil
		} else {
			return 0, fmt.Errorf("Role not valid for inventory")
		}

	case "team":
		return 0, fmt.Errorf("Team endpoint not implemeneted")
	case "organization":
		return 0, fmt.Errorf("Organization endpoint not implemeneted")
	case "job_template":
		awxService := awx.JobTemplateService
		obj, _, err := awxService.ListJobTemplates(map[string]string{
			"name": d.Get("resource_name").(string),
		})
		if err != nil {
			return 0, err
		}
		if d.Get("role").(string) == "admin" {
			return obj[0].SummaryFields.ObjectRoles.AdminRole.ID, nil
		} else if d.Get("role").(string) == "execute" {
			return obj[0].SummaryFields.ObjectRoles.ExecuteRole.ID, nil
		} else if d.Get("role").(string) == "read" {
			return obj[0].SummaryFields.ObjectRoles.ReadRole.ID, nil
		} else {
			return 0, fmt.Errorf("Role not valid for Job Template")
		}
	case "credential":
		return 0, fmt.Errorf("Credential endpoint not implemeneted")
	case "project":
		awxService := awx.ProjectService
		obj, _, err := awxService.ListProjects(map[string]string{
			"name":         d.Get("resource_name").(string),
			"organization": d.Get("organization_id").(string),
		})
		if err != nil {
			return 0, err
		}
		if d.Get("role").(string) == "admin" {
			return obj[0].SummaryFields.ObjectRoles.AdminRole.ID, nil
		} else if d.Get("role").(string) == "update" {
			return obj[0].SummaryFields.ObjectRoles.UpdateRole.ID, nil
		} else if d.Get("role").(string) == "read" {
			return obj[0].SummaryFields.ObjectRoles.ReadRole.ID, nil
		} else if d.Get("role").(string) == "use" {
			return obj[0].SummaryFields.ObjectRoles.UseRole.ID, nil
		} else {
			return 0, fmt.Errorf("Role not valid for Project")
		}
	}
	return 0, fmt.Errorf("Not implemented API endpoint")
}
func resourceUserRoleGrant(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService
	_, res, err := awxService.ListUsers(map[string]string{
		"id": d.Get("user_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("User with Id %s doesn't exists",
			d.Get("user_id").(string))
	}
	id, _ := strconv.Atoi(d.Get("user_id").(string))
	roleID, err := getRoleID(d, m)
	if err == nil {
		err = awxService.GrantRole(id, roleID)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	d.SetId(d.Get("user_id").(string))
	return resourceUserRoleRead(d, m)

}

func resourceUserRoleRevoke(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService

	_, res, err := awxService.ListUsers(map[string]string{
		"id": d.Get("user_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("User with Id %s doesn't exists",
			d.Get("user_id").(string))
	}
	roleID, err := getRoleID(d, m)
	if err == nil {
		id, _ := strconv.Atoi(d.Get("user_id").(string))
		err = awxService.RevokeRole(id, roleID)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	d.SetId("")
	return resourceUserRoleRead(d, m)
}

func resourceUserRoleRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.UserService
	_, res, err := awxService.ListUsers(map[string]string{
		"id": d.Get("user_id").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setUserRoleResourceData(d, res.Results[0])
	return nil
}

func setUserRoleResourceData(d *schema.ResourceData, r *awxgo.User) *schema.ResourceData {
	d.Set("username", r.Username)
	d.Set("user_id", r.ID)
	d.Set("resource_name", d.Get("resource_name").(string))
	d.Set("resource_type", d.Get("resource_type").(string))
	d.Set("role", d.Get("role").(string))
	return d
}
