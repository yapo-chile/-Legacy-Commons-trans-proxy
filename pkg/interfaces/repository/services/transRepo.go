package services

import (
	"reflect"
	"strconv"

	"gitlab.com/yapo_team/legacy/commons/trans/pkg/domain"
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
	transFactory TransFactory
}

// NewTransRepo instance TransRepo and set handler
func NewTransRepo(transFactory TransFactory) *TransRepo {
	return &TransRepo{
		transFactory: transFactory,
	}
}

// Execute executes the specified trans command
func (repo *TransRepo) Execute(command domain.TransCommand) (domain.TransResponse, error) {
	response := domain.TransResponse{
		Params: make(map[string]string),
	}
	resp, err := repo.transaction(command.Command, command.Params)
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

func (repo *TransRepo) transaction(method string, transParams []domain.TransParams) (map[string]string, error) {
	trans := repo.transFactory.MakeTransHandler()
	for _, transParam := range transParams {
		if reflect.TypeOf(transParam.Value).Kind() == reflect.Int {
			transParam.Value = strconv.Itoa(transParam.Value.(int))
		}
	}
	return trans.SendCommand(method, transParams)
}
