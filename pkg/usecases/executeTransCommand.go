package usecases

import (
	"fmt"
	"strings"

	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/domain"
)

// TransOK Status returned when a trans command executes successfully
const TransOK = "TRANS_OK"

// TransError Error while executing a trans error
const TransError = "TRANS_ERROR"

// TransDatabaseError Error while executing a database request inside a trans
const TransDatabaseError = "TRANS_DATABASE_ERROR"

// TransNoCommand Error when the provided command doesn't exists
const TransNoCommand = "TRANS_ERROR_NO_SUCH_COMMAND:Err no such command"

// ExecuteTransUsecase states:
// As a User, I would like to execute my TransCommand on a Trans server and get the corresponding response
// ExecuteTrans should return a response, or an appropriate error if there was a problem.
type ExecuteTransUsecase interface {
	ExecuteCommand(command domain.TransCommand) (domain.TransResponse, error)
}

// TransInteractorLogger defines all the events a TransInteractor may
// need/like to report as they happen
type TransInteractorLogger interface {
	LogBadInput(domain.TransCommand)
	LogRepositoryError(domain.TransCommand, error)
}

// TransInteractor implements ExecuteTransUsecase by using Repository
// to execute the Trans and to retrieve the response.
type TransInteractor struct {
	Logger     TransInteractorLogger
	Repository domain.TransRepository
}

// ExecuteCommand executes the given TransCommand and returns the corresponding TransResponse.
func (interactor TransInteractor) ExecuteCommand(
	command domain.TransCommand,
) (domain.TransResponse, error) {
	response := domain.TransResponse{
		Status: TransError,
		Params: make(map[string]string),
	}
	// Ensure correct input
	if command.Command == "" {
		interactor.Logger.LogBadInput(command)
		return response, fmt.Errorf("invalid command %+v", command)
	}

	// Execute the command and retrieve the response
	response, err := interactor.Repository.Execute(command)
	if err != nil {
		// Report the error
		interactor.Logger.LogRepositoryError(command, err)
		if transErr, ok := response.Params["error"]; ok {
			err = fmt.Errorf(transErr)
		} else {
			err = fmt.Errorf("error during execution")
		}
	}
	// if the command sent doesn´t exists in the server
	if response.Status == TransNoCommand {
		err = fmt.Errorf("error command doesn't exists")
		response.Status = TransError
		response.Params["error"] = err.Error()
	}
	// if the error is a database error
	if strings.Contains(response.Status, TransDatabaseError) {
		// get the specific error message from the status response
		errorString := strings.Replace(response.Status, TransDatabaseError, "", 1)
		errorString = strings.Replace(errorString, ":", "", 1)
		err = fmt.Errorf(errorString)
		interactor.Logger.LogRepositoryError(command, err)
		response.Status = TransDatabaseError
		response.Params["error"] = err.Error()
	}

	return response, err
}
