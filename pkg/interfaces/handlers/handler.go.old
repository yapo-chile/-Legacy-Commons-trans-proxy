package handlers

import (
	"net/http"
	"reflect"

	"github.com/Yapo/goutils"
	mux "gopkg.in/gorilla/mux.v1"
)

// HandlerInput is a placeholder for whatever input a handler may need.
type HandlerInput interface{}

// InputGetter defines a type for all functions that, when called, will attempt
// to retrieve and parse the input of a request and return it. Should any error
// happen, a goutils.Response must be filled with an adequate message and code
type InputGetter func() (HandlerInput, *goutils.Response)

// Handler is the interface for the objects that should process web requests.
// Input() must return a fresh struct to be filled with the request input
// Execute(input) receives a filled input struct to handle the request
type Handler interface {
	// Input should return a pointer to the struct that this handler will need
	// to be filled with the user input for a request
	Input() HandlerInput
	// Execute is the actual handler code. The InputGetter can be used to retrieve
	// the request's input at any time (or not at all).
	Execute(InputGetter) *goutils.Response
}

// MakeJSONHandlerFunc wraps a Handler on a json-over-http context, returning
// a standard http.HandlerFunc
func MakeJSONHandlerFunc(h Handler, l JSONHandlerLogger) http.HandlerFunc {
	jh := jsonHandler{handler: h, logger: l}
	return jh.run
}

// JSONHandlerLogger defines all the events a jsonHandler can report
type JSONHandlerLogger interface {
	LogRequestStart(r *http.Request)
	LogRequestEnd(*http.Request, *goutils.Response)
	LogRequestPanic(*http.Request, *goutils.Response, interface{})
}

// jsonHandler provides an http.HandlerFunc that reads its input and formats
// its output as json
type jsonHandler struct {
	handler Handler
	logger  JSONHandlerLogger
}

// run will prepare the input for the actual handler and format the response
// as json. Also, request information will be logged. It's an instance of
// http.HandlerFunc
func (jh *jsonHandler) run(w http.ResponseWriter, r *http.Request) {
	jh.logger.LogRequestStart(r)
	// Default response
	response := &goutils.Response{
		Code: http.StatusInternalServerError,
	}
	// Function the request can call to retrieve its input
	inputGetter := func() (HandlerInput, *goutils.Response) {
		input := jh.handler.Input()
		// Parse the get params
		getParams := mux.Vars(r)
		response = fillGet(getParams, input)
		if response != nil {
			return input, response
		}

		// Parse the body params
		response = goutils.ParseJSONBody(r, &input)
		return input, response
	}
	// Format the output and send it down the writer
	outputWriter := func() {
		goutils.CreateJSON(response)
		goutils.WriteJSONResponse(w, response)
	}
	// Handle panicking handlers and report errors
	errorHandler := func() {
		if err := recover(); err != nil {
			jh.logger.LogRequestPanic(r, response, err)
		}
	}
	// Setup before calling the actual handler
	defer outputWriter()
	defer errorHandler()
	// Do the Harlem Shake
	response = jh.handler.Execute(inputGetter)
	jh.logger.LogRequestEnd(r, response)
}

// fillGet set variables into the corresponding get param
func fillGet(vars map[string]string, input interface{}) *goutils.Response {
	v := reflect.ValueOf(input)
	reflectedInput := reflect.Indirect(v)
	// Only attempt to set writeable variables
	if reflectedInput.IsValid() && reflectedInput.CanSet() && reflectedInput.Kind() == reflect.Struct {
		// Recursively load inner struct fields
		for i := 0; i < reflectedInput.NumField(); i++ {
			if tag, ok := reflectedInput.Type().Field(i).Tag.Lookup("get"); ok {
				reflectedInput.Field(i).Set(reflect.ValueOf(vars[tag]))
			}
		}
		return nil
	}
	return &goutils.Response{
		Code: http.StatusBadRequest,
		Body: goutils.GenericError{
			ErrorMessage: "Is not a valid struct",
		},
	}
}
