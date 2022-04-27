package loggers

import (
	"gitlab.com/yapo_team/legacy/commons/trans/pkg/domain"
	"gitlab.com/yapo_team/legacy/commons/trans/pkg/usecases"
)

// TransInteractorDefaultLogger logger user in the TransInteractor
type TransInteractorDefaultLogger struct {
	logger Logger
}

// LogBadInput logs a bad input error
func (t *TransInteractorDefaultLogger) LogBadInput(command domain.TransCommand) {
	t.logger.Debug("Invalid trans command. Input: %+v", command)
}

// LogRepositoryError logs a repository error
func (t *TransInteractorDefaultLogger) LogRepositoryError(command domain.TransCommand, err error) {
	t.logger.Error("Error executing trans command %+v: %s", command, err)
}

// MakeTransInteractorLogger sets up a TransInteractorLogger instrumented
// via the provided logger
func MakeTransInteractorLogger(logger Logger) usecases.TransInteractorLogger {
	return &TransInteractorDefaultLogger{
		logger: logger,
	}
}
