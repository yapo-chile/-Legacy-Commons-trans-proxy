package handlers

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

type MockTokenValidator struct {
	mock.Mock
}

func (m *MockTokenValidator) CleanAndMatchToken(token string) error {
	ret := m.Called(token)
	return ret.Error(0)
}

func TestTransHandlerInput(t *testing.T) {
	m := MockTransInteractor{}
	mInputRequest := MockInputRequest{}
	mTargetRequest := MockTargetRequest{}
	mInputRequest.On("Set", mock.Anything).Return(&mTargetRequest)
	mTargetRequest.On("FromHeaders").Return()
	mTargetRequest.On("FromPath").Return()
	mTargetRequest.On("FromJSONBody").Return()

	h := TransHandler{Interactor: &m}
	input := h.Input(&mInputRequest)
	var expected *TransHandlerInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestTransHandlerExecuteOK(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "transinfo"}
	command := domain.TransCommand{
		Command: "transinfo",
		Params:  make([]domain.TransParams, 0),
	}
	response := domain.TransResponse{
		Status: usecases.TransOK,
	}
	m.On("ExecuteCommand", command).Return(response, nil).Once()

	mTokenVal := MockTokenValidator{}
	mTokenVal.On("CleanAndMatchToken", "").Return(nil).Once()

	h := TransHandler{Interactor: &m, TokenValidationInteractor: &mTokenVal}

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
	mTokenVal := MockTokenValidator{}
	mTokenVal.On("CleanAndMatchToken", "").Return(nil).Once()

	h := TransHandler{Interactor: &m, TokenValidationInteractor: &mTokenVal}

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
	mTokenVal := MockTokenValidator{}
	mTokenVal.On("CleanAndMatchToken", "").Return(nil).Once()

	h := TransHandler{Interactor: &m, TokenValidationInteractor: &mTokenVal}

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
	mTokenVal := MockTokenValidator{}
	mTokenVal.On("CleanAndMatchToken", "").Return(nil).Once()

	h := TransHandler{Interactor: &m, TokenValidationInteractor: &mTokenVal}

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

func TestTransHandlerExecuteUnauthorized(t *testing.T) {
	m := MockTransInteractor{}
	input := TransHandlerInput{Command: "get_account"}

	mTokenVal := MockTokenValidator{}
	mTokenVal.On("CleanAndMatchToken", "").Return(errors.New("foo")).Once()

	h := TransHandler{Interactor: &m, TokenValidationInteractor: &mTokenVal}

	expectedResponse := &goutils.Response{
		Code: http.StatusUnauthorized,
		Body: &goutils.GenericError{
			ErrorMessage: "foo",
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

func TestBuildCommand(t *testing.T) {
	input := TransHandlerInput{
		Command: "get_account",
		Params:  make(map[string]interface{}),
	}
	sliceParams := make([]interface{}, 0)
	mapParams := make(map[string]interface{})
	mapParams["subemail"] = fakeEmail
	sliceParams = append(sliceParams, mapParams)
	input.Params["email"] = sliceParams

	secondSliceParam := make([]interface{}, 0)
	secondSliceParam = append(secondSliceParam, "true")
	input.Params["is_company"] = secondSliceParam

	command := domain.TransCommand{
		Command: "get_account",
		Params: []domain.TransParams{
			{
				Key:   "subemail",
				Value: "user@test.com",
				Blob:  false,
			},
			{
				Key:   "is_company",
				Value: "true",
				Blob:  false,
			},
		},
	}

	r := BuildCommand(&input)

	assert.Equal(t, command, r)
}
