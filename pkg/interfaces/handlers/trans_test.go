package handlers

/*
import (
	"errors"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/domain"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/usecases"
)

const (
	fakeEmail = "user@test.com"
)

type MockTransInteractor struct {
	mock.Mock
}

func (m *MockTransInteractor) ExecuteCommand(command domain.TransCommand) (domain.TransResponse, error) {
	ret := m.Called(command)
	return ret.Get(0).(domain.TransResponse), ret.Error(1)
}

func MakeMockInputTransGetter(input HandlerInput, response *goutils.Response) InputGetter {
	return func() (HandlerInput, *goutils.Response) {
		return input, response
	}
}

func TestTransHandlerInput(t *testing.T) {
	m := MockTransInteractor{}
	h := TransHandler{Interactor: &m}
	input := h.Input()
	var expected *TransHandlerInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestTransHandlerExecuteOK(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "trans-proxyinfo"}
	command := domain.TransCommand{
		Command: "trans-proxyinfo",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{
		Status: usecases.TransOK,
	}
	m.On("ExecuteCommand", command).Return(response, nil).Once()
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: TransRequestOutput{
			Status: usecases.TransOK,
		},
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerParseInput(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{
		Command: "get_account",
		Params:  make(map[string]interface{}),
	}
	input.Params["email"] = fakeEmail
	command := domain.TransCommand{
		Command: "get_account",
		Params:  make([]domain.TransParams, 0),
	}
	param := domain.TransParams{
		Key:   "email",
		Value: fakeEmail,
	}
	command.Params = append(command.Params, param)

	response := domain.TransResponse{
		Status: usecases.TransOK,
		Params: make(map[string]string),
	}
	response.Params["account_id"] = "1"
	response.Params["email"] = fakeEmail
	response.Params["is_company"] = "true"
	m.On("ExecuteCommand", command).Return(response, nil).Once()
	h := TransHandler{Interactor: &m}

	requestOutput := TransRequestOutput{
		Status:   usecases.TransOK,
		Response: make(map[string]string),
	}
	requestOutput.Response["account_id"] = "1"
	requestOutput.Response["email"] = fakeEmail
	requestOutput.Response["is_company"] = "true"
	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: requestOutput,
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerExecuteError(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "get_account"}
	command := domain.TransCommand{
		Command: "get_account",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{
		Status: usecases.TransError,
	}
	m.On("ExecuteCommand", command).Return(response, nil).Once()
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: TransRequestOutput{
			Status: usecases.TransError,
		},
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerExecuteInternalError(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "get_account"}
	command := domain.TransCommand{
		Command: "get_account",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{}
	m.On("ExecuteCommand", command).Return(response, errors.New("Error")).Once()
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusInternalServerError,
		Body: &goutils.GenericError{
			ErrorMessage: "Error",
		},
	}

	getter := MakeMockInputTransGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestTransHandlerInputError(t *testing.T) {
	m := MockTransInteractor{}
	h := TransHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: "Error",
	}

	getter := MakeMockInputTransGetter(nil, expectedResponse)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}
*/
