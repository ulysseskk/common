package goroutineUtil

import (
	"gitlab.ulyssesk.top/common/common/logger/log"
	"gitlab.ulyssesk.top/common/common/model/errors"
	"gitlab.ulyssesk.top/common/common/model/rest"
	"fmt"
	"runtime"
)

func RecoverFunc(hook func(r any)) func() {
	return func() {
		if r := recover(); r != nil {
			if hook != nil {
				hook(r)
			}
			defaultRecoveryFunc(r)
		}
	}
}

func defaultRecoveryFunc(r interface{}) {
	stack := make([]byte, 1<<16)
	stack = stack[:runtime.Stack(stack, false)]
	commonErr := errors.NewError().WithCode(rest.InternalError)
	err, ok := r.(error)
	if ok {
		commonErr = commonErr.WithError(err)
	}
	commonErr = commonErr.WithMessage(fmt.Sprintf("%v", r))
	log.GlobalLogger().Errorf("Panic %v\n%s", commonErr, stack)
}
