package tower

import (
	"encoding/json"
	"fmt"

	"github.com/go-yaml/yaml"
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
