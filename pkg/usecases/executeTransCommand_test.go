package usecases

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/yapo_team/legacy/commons/trans-proxy/pkg/domain"
)

type MockTransRepository struct {
	mock.Mock
}

func (m *MockTransRepository) Execute(command domain.TransCommand) (domain.TransResponse, error) {
	ret := m.Called(command)
	return ret.Get(0).(domain.TransResponse), ret.Error(1)
}

type MockTransInteractorLogger struct {
	mock.Mock
}

func (m *MockTransInteractorLogger) LogBadInput(c domain.TransCommand) {
	m.Called(c)
}

func (m *MockTransInteractorLogger) LogRepositoryError(c domain.TransCommand, err error) {
	m.Called(c, err)
}

func TestTransInteractorInvalidCommand(t *testing.T) {
	logger := &MockTransInteractorLogger{}
	repo := &MockTransRepository{}
	interactor := TransInteractor{
		Logger:     logger,
		Repository: repo,
	}
	command := domain.TransCommand{}
	logger.On("LogBadInput", command)

	_, err := interactor.ExecuteCommand(command)
	assert.Error(t, err)
	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransInteractorRepositoryError(t *testing.T) {
	command := domain.TransCommand{
		Command: "command 1",
	}
	response := domain.TransResponse{}
	err := errors.New("error")
	logger := &MockTransInteractorLogger{}
	repo := &MockTransRepository{}
	repo.On("Execute", command).Return(response, err).Once()
	interactor := TransInteractor{
		Logger:     logger,
		Repository: repo,
	}
	logger.On("LogRepositoryError", command, err).Once()
	expectedErr := fmt.Errorf("error during execution")
	returnResp, returnErr := interactor.ExecuteCommand(command)
	assert.Error(t, returnErr)
	assert.Equal(t, expectedErr, returnErr)
	assert.Equal(t, response, returnResp)
	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransInteractorTransNoCommand(t *testing.T) {
	command := domain.TransCommand{
		Command: "command 1",
	}
	err := errors.New("error command doesn't exists")
	response := domain.TransResponse{
		Status: TransNoCommand,
		Params: make(map[string]string),
	}

	logger := &MockTransInteractorLogger{}
	repo := &MockTransRepository{}
	repo.On("Execute", command).Return(response, err).Once()
	interactor := TransInteractor{
		Logger:     logger,
		Repository: repo,
	}
	logger.On("LogRepositoryError", command, err).Once()
	expectedErr := fmt.Errorf("error command doesn't exists")
	expectedResponse := domain.TransResponse{
		Status: TransError,
		Params: make(map[string]string),
	}
	expectedResponse.Params["error"] = expectedErr.Error()
	returnResp, returnErr := interactor.ExecuteCommand(command)
	assert.Error(t, returnErr)
	assert.Equal(t, expectedErr, returnErr)
	assert.Equal(t, expectedResponse, returnResp)
	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransInteractorTransDatabaseError(t *testing.T) {
	command := domain.TransCommand{
		Command: "command 1",
	}
	errorStringDB := "ERROR EXECUTING QUERY"
	errorString := "Trans Database error"
	err := errors.New(errorString)
	errDB := errors.New(errorStringDB)
	response := domain.TransResponse{
		Status: fmt.Sprintf("%s:%s", TransDatabaseError, errorStringDB),
		Params: make(map[string]string),
	}

	logger := &MockTransInteractorLogger{}
	repo := &MockTransRepository{}
	repo.On("Execute", command).Return(response, err).Once()
	interactor := TransInteractor{
		Logger:     logger,
		Repository: repo,
	}
	logger.On("LogRepositoryError", command, err).Once()
	logger.On("LogRepositoryError", command, errDB).Once()

	expectedResponse := domain.TransResponse{
		Status: TransDatabaseError,
		Params: make(map[string]string),
	}
	expectedResponse.Params["error"] = errorStringDB
	returnResp, returnErr := interactor.ExecuteCommand(command)

	assert.Error(t, returnErr)
	assert.Equal(t, errDB, returnErr)
	assert.Equal(t, expectedResponse, returnResp)
	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestTransInteractorExecuteCommandOK(t *testing.T) {
	command := domain.TransCommand{
		Command: "command 1",
	}
	response := domain.TransResponse{
		Status: TransOK,
	}
	logger := &MockTransInteractorLogger{}
	repo := &MockTransRepository{}
	repo.On("Execute", command).Return(response, nil).Once()
	interactor := TransInteractor{
		Logger:     logger,
		Repository: repo,
	}
	returnResp, returnErr := interactor.ExecuteCommand(command)
	assert.NoError(t, returnErr)
	assert.Equal(t, response, returnResp)
	repo.AssertExpectations(t)
	logger.AssertExpectations(t)
}
