package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Input() HandlerInput {
	args := m.Called()
	return args.Get(0).(HandlerInput)
}
func (m *MockHandler) Execute(getter InputGetter) *goutils.Response {
	args := m.Called(getter)
	_, response := getter()
	if response != nil {
		return response
	}
	return args.Get(0).(*goutils.Response)
}

type MockPanicHandler struct {
	mock.Mock
}

func (m *MockPanicHandler) Input() HandlerInput {
	args := m.Called()
	return args.Get(0).(HandlerInput)
}
func (m *MockPanicHandler) Execute(getter InputGetter) *goutils.Response {
	m.Called(getter)
	panic("dead")
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogRequestStart(r *http.Request) {
	m.Called(r)
}
func (m *MockLogger) LogRequestEnd(r *http.Request, response *goutils.Response) {
	m.Called(r, response)
}
func (m *MockLogger) LogRequestPanic(r *http.Request, response *goutils.Response, err interface{}) {
	m.Called(r, response, err)
}

type DummyInput struct {
	X int
}

type DummyInputGet struct {
	Method string `get:"method"`
	X      int
}

type DummyOutput struct {
	Y string
}

func TestJsonHandlerFuncOK(t *testing.T) {
	h := MockHandler{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 400,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input").Return(input).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response)

	fn := MakeJSONHandlerFunc(&h, &l)
	fn(w, r)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	l.AssertExpectations(t)
}
func TestJsonHandlerFillGet(t *testing.T) {
	h := MockHandler{}
	l := MockLogger{}
	input := &DummyInputGet{}
	response := &goutils.Response{
		Code: 400,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input").Return(input).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{\"x\":12}"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response)

	fn := MakeJSONHandlerFunc(&h, &l)
	fn(w, r)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFillGetInvalidStruct(t *testing.T) {
	h := MockHandler{}
	l := MockLogger{}
	input := 6
	response := &goutils.Response{
		Code: 400,
		Body: goutils.GenericError{
			ErrorMessage: "Is not a valid struct",
		},
	}

	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input").Return(input).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{\"x\":12}"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response)

	fn := MakeJSONHandlerFunc(&h, &l)
	fn(w, r)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"ErrorMessage":"Is not a valid struct"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncParseError(t *testing.T) {
	h := MockHandler{}
	l := MockLogger{}
	input := &DummyInput{}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter)
	h.On("Input").Return(input).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, mock.AnythingOfType("*goutils.Response"))

	fn := MakeJSONHandlerFunc(&h, &l)
	fn(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"ErrorMessage":"unexpected EOF"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncPanic(t *testing.T) {
	h := MockPanicHandler{}
	l := MockLogger{}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{"))

	l.On("LogRequestStart", r)
	l.On("LogRequestPanic", r, mock.AnythingOfType("*goutils.Response"), "dead")

	fn := MakeJSONHandlerFunc(&h, &l)
	fn(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "null\n", w.Body.String())
	h.AssertExpectations(t)
	l.AssertExpectations(t)
}
