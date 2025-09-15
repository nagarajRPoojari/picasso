package test

import (
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	dir := t.TempDir()
	src := `
	import io;
	fn main(): int32 {
		say z: int = 190;
		io.printf("hello world");
		return 0;
	}
	`
	output, err := testutils.CompileAndRun(src, dir)
	assert.NoError(t, err)
	assert.Equal(t, output, `"hello world"`)
}
