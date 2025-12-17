package testutils

import (
	"fmt"
)

func CompileAndRunSafe(src string, dir string, libs ...string) (out string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic occurred: %v", r)
		}
	}()
	out, err = CompileAndRun(src, dir, libs...)
	return
}

func CompileAndRun(src string, dir string, libs ...string) (string, error) {
	return "", nil
}
