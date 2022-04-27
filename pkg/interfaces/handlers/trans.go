package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
	"gitlab.com/yapo_team/legacy/commons/trans/pkg/domain"
	"gitlab.com/yapo_team/legacy/commons/trans/pkg/usecases"
)

// TransHandler implements the handler interface and responds to /execute
// requests with a message. Expected response format:
// { status: string, response: json }
type TransHandler struct {
	Interactor usecases.ExecuteTransUsecase
}

// TransHandlerInput struct that represents the input
type TransHandlerInput struct {
	Command string                 `get:"command"`
	Params  map[string]interface{} `json:"params"`
}

// TransRequestOutput struct that represents the output
type TransRequestOutput struct {
	Status   string            `json:"status"`
	Response map[string]string `json:"response"`
}

// Input returns a fresh, empty instance of transHandlerInput
func (t *TransHandler) Input() HandlerInput {
	return &TransHandlerInput{}
}

// Execute executes the given trans request and returns the response
// of the execution.
// Expected response format:
//   { Status: string - "TRANS_OK" or error }
func (t *TransHandler) Execute(ig InputGetter) *goutils.Response {
	input, response := ig()
	if response != nil {
		return response
	}
	in := input.(*TransHandlerInput)
	command := parseInput(in)
	var val domain.TransResponse
	val, err := t.Interactor.ExecuteCommand(command)
	// handle trans errors, database errors, or general reported errors by trans
	if _, ok := val.Params["error"]; ok ||
		val.Status == usecases.TransError ||
		val.Status == usecases.TransDatabaseError {
		response = &goutils.Response{
			Code: http.StatusBadRequest,
			Body: TransRequestOutput{
				Status:   val.Status,
				Response: val.Params,
			},
		}
		return response
	}

	// handle errors given by the interactor
	if err != nil {
		response = &goutils.Response{
			Code: http.StatusInternalServerError,
			Body: &goutils.GenericError{
				ErrorMessage: err.Error(),
			},
		}
		return response
	}

	response = &goutils.Response{
		Code: http.StatusOK,
		Body: TransRequestOutput{
			Status:   val.Status,
			Response: val.Params,
		},
	}
	return response
}

func parseInput(input *TransHandlerInput) domain.TransCommand {
	command := domain.TransCommand{
		Command: input.Command,
	}

	params := make([]domain.TransParams, 0)

	for key, value := range input.Params {
		if _, ok := value.([]interface{}); ok {
			for _, val := range value.([]interface{}) {
				if _, ok := val.(map[string]interface{}); ok {
					for k, v := range val.(map[string]interface{}) {
						param := domain.TransParams{
							Key:   k,
							Value: v,
							Blob:  key == "blobs",
						}
						params = append(params, param)
					}
				} else if _, ok := val.(string); ok {
					param := domain.TransParams{
						Key:   key,
						Value: val,
					}
					params = append(params, param)
				}
			}
		} else {
			param := domain.TransParams{
				Key:   key,
				Value: value,
			}
			params = append(params, param)
		}
	}
	command.Params = params
	return command
}
