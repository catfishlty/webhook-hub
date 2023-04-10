package utils

import (
	"errors"
	"fmt"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"regexp"
	"strings"
)

const VariableTypeHeader = "header"
const VariableTypeQuery = "query"
const VariableTypeForm = "form"
const VariableTypeJson = "json"

func ValidateVariables(variables map[string]types.VariableItem) error {
	errMsg := make([]string, 0)
	for _, v := range variables {
		if v.Validate != "" && v.Validate != v.Value {
			errMsg = append(errMsg, fmt.Sprintf("variable %s:%s validate failed, expect %s, but got %s", v.Type, v.Key, v.Validate, v.Value))
		}
	}
	if len(errMsg) > 0 {
		return errors.New(strings.Join(errMsg, ","))
	}
	return nil
}

func GetVariables(receive types.ReceiveRequest, c *gin.Context) (map[string]types.VariableItem, []byte, error) {
	variables := make(map[string]types.VariableItem)
	data, err := c.GetRawData()
	if err != nil {
		return nil, nil, err
	}
	for _, v := range receive.Variables {
		jsonData, err := jsoniter.Marshal(v)
		if err != nil {
			return nil, nil, err
		}
		var variable types.VariableItem
		err = jsoniter.Unmarshal(jsonData, &variable)
		if err != nil {
			return nil, nil, err
		}
		t, k, err := parseTypeAndKey(variable.Key)
		if err != nil {
			return nil, nil, err
		}
		variable.Type = t
		variable.Key = k
		variable.Value, err = getVariableValue(c.Request, data, variable)
		if err != nil {
			return nil, nil, err
		}
		variables[variable.Assign] = variable
	}
	return variables, data, nil
}

func getVariableValue(request *http.Request, data []byte, variable types.VariableItem) (string, error) {
	switch variable.Type {
	case VariableTypeHeader:
		return getHeaderVariableValue(request, variable)
	case VariableTypeQuery:
		return getQueryVariableValue(request, variable)
	case VariableTypeForm:
		return getFormVariableValue(request, variable)
	case VariableTypeJson:
		return getJsonVariableValue(data, variable)
	default:
		return "", errors.New("unknown variable type " + variable.Type)
	}
}

func getJsonVariableValue(data []byte, variable types.VariableItem) (string, error) {
	keyList := strings.Split(variable.Key, ".")
	var m map[string]any
	err := json.Unmarshal(data, &m)
	if err != nil {
		return "", err
	}
	for i, key := range keyList {
		if v, ok := m[key]; ok {
			if v, ok := v.(string); ok && i == len(keyList)-1 {
				return v, nil
			}
			if v, ok := v.(map[string]any); ok {
				m = v
			}
		} else {
			return "", errors.New("json value not exists:" + variable.Key)
		}
	}
	return "", errors.New("json value not exists:" + variable.Key)
}

func getFormVariableValue(request *http.Request, variable types.VariableItem) (string, error) {
	if exists := request.Form.Has(variable.Key); !exists {
		return "", errors.New("form value not exists:" + variable.Key)
	}
	return request.Form.Get(variable.Key), nil
}

func getQueryVariableValue(request *http.Request, variable types.VariableItem) (string, error) {
	if exists := request.URL.Query().Has(variable.Key); !exists {
		return "", errors.New("query value not exists:" + variable.Key)
	}
	return request.URL.Query().Get(variable.Key), nil
}

func getHeaderVariableValue(request *http.Request, variable types.VariableItem) (string, error) {
	return request.Header.Get(variable.Key), nil
}

func parseTypeAndKey(k string) (string, string, error) {
	args := strings.Split(k, ":")
	if len(args) != 2 {
		return "", "", errors.New("variable key error:" + k)
	}
	switch args[0] {
	case VariableTypeHeader:
		return VariableTypeHeader, args[1], nil
	case VariableTypeQuery:
		return VariableTypeQuery, args[1], nil
	case VariableTypeForm:
		return VariableTypeForm, args[1], nil
	case VariableTypeJson:
		return VariableTypeJson, args[1], nil
	default:
		return "", "", errors.New("unknown variable type " + args[0])
	}
}

func ReplaceVariables(send *types.RestySendRequest, variables map[string]types.VariableItem) {
	for i, v := range send.Header {
		send.Header[i] = replaceVariable(v, variables)
	}
}

func replaceVariable(v string, variables map[string]types.VariableItem) string {
	r := regexp.MustCompile(`\$\{(.+?)\}`)
	if r.MatchString(v) {
		matches := r.FindAllString(v, -1)
		for _, match := range matches {
			return strings.ReplaceAll(v, match, variables[match[2:len(match)-1]].Value)
		}
	}
	return v
}
