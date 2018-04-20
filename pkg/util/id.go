package util

import (
	"strconv"
)

// GetIDs parses the input to get resource ids
func GetIDs(value interface{}) (ids []int) {
	if i, ok := value.([]interface{}); ok {
		for _, v := range i {
			id := GetID(v)
			if id != 0 {
				ids = append(ids, id)
			}
		}
	} else {
		id := GetID(value)
		if id != 0 {
			ids = append(ids, id)
		}
	}
	return
}

// GetID parses the input to get a resource id
func GetID(value interface{}) (id int) {
	switch v := value.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case string:
		id, _ = strconv.Atoi(v)
		return
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
			n := GetName(v)
			if n != "" {
				names = append(names, n)
			}
		}
	} else {
		n := GetName(value)
		if n != "" {
			names = append(names, n)
		}
	}
	return
}

// GetName parses the input to get a resource name
func GetName(value interface{}) (name string) {
	switch v := value.(type) {
	case string:
		return v
	case float64:
		return strconv.Itoa(int(v))
	case map[string]interface{}:
		if name, ok := v["name"]; ok {
			return GetName(name)
		}
	}
	return
}
