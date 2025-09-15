package testutils

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/nagarajRPoojari/x-lang/compiler"
	"github.com/nagarajRPoojari/x-lang/parser"
)

func CompileAndRun(src string, dir string, libs ...string) (string, error) {
	irFile := path.Join(dir, "output.ll")
	bcFile := path.Join(dir, "output.bc")
	objectFile := path.Join(dir, "output.o")
	execFile := path.Join(dir, "output")

	ast := parser.Parse(src)
	c := compiler.NewCompiler()
	c.Compile(ast)
	c.Dump(irFile)

	cmd := exec.Command("llvm-as", irFile, "-o", bcFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("llvm-as error: %v\n%s", err, string(out))
	}

	cmd = exec.Command("llc", "-filetype=obj", bcFile, "-o", objectFile)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("llc error: %v\n%s", err, string(out))
	}

	libs = append(libs, objectFile)

	args := append(libs, "-o", execFile)
	cmd = exec.Command("clang", args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("clang error: %v\n%s", err, string(out))
	}

	cmd = exec.Command(execFile)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec error: %v\n%s", err, string(out))
	}
	return string(out), nil
}
