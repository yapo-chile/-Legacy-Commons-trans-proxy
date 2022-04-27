package domain

// TransParams is a struct with Trans format params
type TransParams struct {
	Key   string
	Value interface{}
	Blob  bool
}

// TransCommand represents a trans-proxy command with params to be executed on a trans-proxy server
type TransCommand struct {
	// the command to be executed
	Command string
	// Params the params of the command
	Params []TransParams
}

// TransResponse represents the response given to the execution of a TransCommand
type TransResponse struct {
	// Status the status of the response (normally TRANS_OK or TRANS_ERROR)
	Status string
	// Params additional params returned
	Params map[string]string
}

// TransRepository defines a storage for the trans-proxy commands
type TransRepository interface {
	// Execute executes the command on a trans-proxy server
	Execute(command TransCommand) (TransResponse, error)
}
