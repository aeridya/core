package logit

import (
	"fmt"
)

type errors struct {
	status Status
	caller string
	msg    string
}

func (e errors) Error() string {
	return e.String()
}

func (e errors) String() string {
	return fmt.Sprintf("%s | Function: %s\n  %s", e.status.String(), e.caller, e.msg)
}

func Error(stat Status, call string, m string) *errors {
	return &errors{status: stat, caller: call, msg: m}
}

func Errorf(stat Status, call string, format string, a ...interface{}) *errors {
	return &errors{status: stat, caller: call, msg: fmt.Sprintf(format, a...)}
}
