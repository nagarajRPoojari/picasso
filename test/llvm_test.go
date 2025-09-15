package test

import (
	"fmt"
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestClass(t *testing.T) {
	dir := t.TempDir()
	src := `
	fn main(): int32 {
		say z: int = 190;
		return 0;
	}
	`

	output, err := testutils.CompileAndRun(src, dir)
	assert.NoError(t, err)

	fmt.Println(output)
}
