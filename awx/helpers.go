package awx

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "github.com/mauromedda/awx-go"
	"gopkg.in/yaml.v2"
)

func normalizeJsonYaml(s interface{}) string {
	result := string("")
	if j, ok := normalizeJsonOk(s); ok {
		result = j
	} else if y, ok := normalizeYamlOk(s); ok {
		result = y
	} else {
		result = s.(string)
	}
	return result
}
func normalizeJsonOk(s interface{}) (string, bool) {
	if s == nil || s == "" {
		return "", true
	}
	var j interface{}
	err := json.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing JSON: %s", err), false
	}
	b, _ := json.Marshal(j)
	return string(b[:]), true
}

func normalizeJson(s interface{}) string {
	v, _ := normalizeJsonOk(s)
	return v
}

func normalizeYamlOk(s interface{}) (string, bool) {
	if s == nil || s == "" {
		return "", true
	}
	var j interface{}
	err := yaml.Unmarshal([]byte(s.(string)), &j)
	if err != nil {
		return fmt.Sprintf("Error parsing YAML: %s", err), false
	}
	b, _ := yaml.Marshal(j)
	return string(b[:]), true
}

func normalizeYaml(s interface{}) string {
	v, _ := normalizeYamlOk(s)
	return v
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
		awxService := awx.TeamService
		obj, _, err := awxService.ListTeams(map[string]string{
			"name":         d.Get("resource_name").(string),
			"organization": d.Get("organization_id").(string),
		})
		if err != nil {
			return 0, err
		}
		if d.Get("role").(string) == "admin" {
			return obj[0].SummaryFields.ObjectRoles.AdminRole.ID, nil
		} else if d.Get("role").(string) == "member" {
			return obj[0].SummaryFields.ObjectRoles.MemberRole.ID, nil
		} else if d.Get("role").(string) == "read" {
			return obj[0].SummaryFields.ObjectRoles.ReadRole.ID, nil
		} else {
			return 0, fmt.Errorf("Role not valid for team object")
		}
	case "organization":
		awxService := awx.OrganizationService
		obj, _, err := awxService.ListOrganizations(map[string]string{
			"name": d.Get("resource_name").(string),
		})
		if err != nil {
			return 0, err
		}
		if d.Get("role").(string) == "admin" {
			return obj[0].SummaryFields.ObjectRoles.AdminRole.ID, nil
		} else if d.Get("role").(string) == "member" {
			return obj[0].SummaryFields.ObjectRoles.MemberRole.ID, nil
		} else if d.Get("role").(string) == "read" {
			return obj[0].SummaryFields.ObjectRoles.ReadRole.ID, nil
		} else if d.Get("role").(string) == "member" {
			return obj[0].SummaryFields.ObjectRoles.MemberRole.ID, nil
		} else if d.Get("role").(string) == "workflow admin" {
			return obj[0].SummaryFields.ObjectRoles.WorkflowAdminRole.ID, nil
		} else if d.Get("role").(string) == "credential admin" {
			return obj[0].SummaryFields.ObjectRoles.CredentialAdminRole.ID, nil
		} else if d.Get("role").(string) == "job template admin" {
			return obj[0].SummaryFields.ObjectRoles.JobTemplateAdminRole.ID, nil
		} else if d.Get("role").(string) == "project admin" {
			return obj[0].SummaryFields.ObjectRoles.ProjectAdminRole.ID, nil
		} else if d.Get("role").(string) == "auditor" {
			return obj[0].SummaryFields.ObjectRoles.AuditorRole.ID, nil
		} else if d.Get("role").(string) == "inventory admin" {
			return obj[0].SummaryFields.ObjectRoles.InventoryAdminRole.ID, nil
		} else {
			return 0, fmt.Errorf("Role not valid for organization object")
		}
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
