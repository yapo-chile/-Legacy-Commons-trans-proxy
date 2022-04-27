package services

import (
	"reflect"
	"strconv"

	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/domain"
)

// TransHandler is an interface to use Trans functions
type TransHandler interface {
	SendCommand(string, []domain.TransParams) (map[string]string, error)
}

// TransFactory is an interface that abstracts the Factory Pattern for creating TransHandler objects
type TransFactory interface {
	MakeTransHandler() TransHandler
}

// TransRepo struct definition
type TransRepo struct {
	trans-proxyFactory TransFactory
}

// NewTransRepo instance TransRepo and set handler
func NewTransRepo(trans-proxyFactory TransFactory) *TransRepo {
	return &TransRepo{
		trans-proxyFactory: trans-proxyFactory,
	}
}

// Execute executes the specified trans-proxy command
func (repo *TransRepo) Execute(command domain.TransCommand) (domain.TransResponse, error) {
	response := domain.TransResponse{
		Params: make(map[string]string),
	}
	resp, err := repo.trans-proxyaction(command.Command, command.Params)
	if err != nil {
		response.Params["error"] = err.Error()
		return response, err
	}
	if status, ok := resp["status"]; ok {
		response.Status = status
		delete(resp, "status")
	}
	for key, val := range resp {
		response.Params[key] = val
	}
	return response, nil
}

func (repo *TransRepo) trans-proxyaction(method string, trans-proxyParams []domain.TransParams) (map[string]string, error) {
	trans-proxy := repo.trans-proxyFactory.MakeTransHandler()
	for _, trans-proxyParam := range trans-proxyParams {
		if reflect.TypeOf(trans-proxyParam.Value).Kind() == reflect.Int {
			trans-proxyParam.Value = strconv.Itoa(trans-proxyParam.Value.(int))
		}
	}
	return trans-proxy.SendCommand(method, trans-proxyParams)
}
