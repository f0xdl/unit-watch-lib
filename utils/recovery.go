package utils

import (
	"fmt"
	"log"
	"runtime/debug"
)

// SafeCall do func with recovery and return error
func SafeCall(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = fmt.Errorf("panic: %s", x)
			case error:
				err = fmt.Errorf("panic: %w", x)
			default:
				err = fmt.Errorf("panic: %v", x)
			}
			log.Printf("SafeCall panic: %v\nStack:\n%s", err, debug.Stack())
		}
	}()

	return fn()
}
