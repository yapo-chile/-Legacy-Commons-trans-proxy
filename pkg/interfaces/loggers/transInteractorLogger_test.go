package loggers

import (
	"testing"

	"gitlab.com/yapo_team/legacy/commons/trans/pkg/domain"
)

// There are no return values to assert on, as logger only cause side effects
// to communicate with the outside world. These tests only ensure that the
// loggers don't panic

func TestTransInteractorDefaultLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeTransInteractorLogger(m)
	input := domain.TransCommand{
		Command: "",
	}
	l.LogBadInput(input)
	l.LogRepositoryError(input, nil)
}
