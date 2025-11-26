package testutils

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/nagarajRPoojari/x-lang/generator"
	"github.com/nagarajRPoojari/x-lang/parser"
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
	irFile := path.Join(dir, "output.ll")
	bcFile := path.Join(dir, "output.bc")
	objectFile := path.Join(dir, "output.o")
	runtimeObj := path.Join(dir, "runtime.o")
	execFile := path.Join(dir, "output")

	ast := parser.Parse(src)
	c := generator.NewGenerator()
	c.SetAST(ast)
	c.Compile()
	c.Dump(irFile)

	if out, err := exec.Command("llvm-as", irFile, "-o", bcFile).CombinedOutput(); err != nil {
		fmt.Println("llvm error")
		return "", fmt.Errorf("llvm-as error: %v\n%s", err, out)
	}

	if out, err := exec.Command("llc", "-filetype=obj", bcFile, "-o", objectFile).CombinedOutput(); err != nil {
		return "", fmt.Errorf("llc error: %v\n%s", err, out)
	}

	projectRoot, err := findProjectRoot("c/runtime.c")
	if err != nil {
		return "", err
	}
	runtimeC := filepath.Join(projectRoot, "c", "runtime.c")

	brewPrefixBytes, err := exec.Command("brew", "--prefix", "bdw-gc").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get brew prefix for bdw-gc: %v", err)
	}
	brewPrefix := strings.TrimSpace(string(brewPrefixBytes))

	if out, err := exec.Command(
		"clang",
		"-c", runtimeC,
		"-I"+filepath.Join(brewPrefix, "include"),
		"-o", runtimeObj,
	).CombinedOutput(); err != nil {
		return "", fmt.Errorf("clang compile runtime.c error: %v\n%s", err, out)
	}

	args := []string{
		objectFile, runtimeObj,
		"-L" + filepath.Join(brewPrefix, "lib"),
		"-lgc",
		"-o", execFile,
	}
	args = append(libs, args...)
	if out, err := exec.Command("clang", args...).CombinedOutput(); err != nil {
		return "", fmt.Errorf("clang link error: %v\n%s", err, out)
	}

	absExecFile, err := filepath.Abs(execFile)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(absExecFile)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("exec error: %v\n%s", err, out)
	}

	return string(out), nil
}

func findProjectRoot(target string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, target)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("could not find %s in any parent directory", target)
}
