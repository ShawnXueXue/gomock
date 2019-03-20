package api

import (
	"fmt"
)

type Api struct {
	ApiName     string
	Response    interface{}
	RequestType string
	Status      int
}

func getId(apiName, requestType string) string {
	return fmt.Sprintf("%s/%s", requestType, apiName)
}

func (api *Api) getId() string {
	return fmt.Sprintf("%s/%s", api.RequestType, api.ApiName)
}
