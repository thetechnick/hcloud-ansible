package util

// GetID parses the input to get a resource id
func GetID(value interface{}) (id int) {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case map[string]interface{}:
		if id, ok := v["id"]; ok {
			return GetID(id)
		}
	}
	return 0
}

// GetNames parses the input to get resource names
func GetNames(value interface{}) (names []string) {
	if s, ok := value.([]interface{}); ok {
		for _, v := range s {
			names = append(names, GetName(v))
		}
	}
	return
}

// GetName parses the input to get a resource name
func GetName(value interface{}) (name string) {
	switch v := value.(type) {
	case string:
		return v
	case map[string]interface{}:
		if name, ok := v["name"]; ok {
			return GetName(name)
		}
	}
	return
}
