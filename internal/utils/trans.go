package utils

import "gorm.io/datatypes"

func JsonToStringMap(data datatypes.JSONMap) map[string]string {
	m := make(map[string]string)
	for k, v := range data {
		m[k] = v.(string)
	}
	return m
}
