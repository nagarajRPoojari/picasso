package io_test

import (
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	dir := t.TempDir()
	tests := []struct {
		name    string
		src     string
		want    string
		wantErr bool
	}{
		{
			name: "basic printf without format specifiers",
			src: `
			import io;
			fn main(): int32 {
				say z: int = 190;
				io.printf("hello world");
				return 0;
			}
			`,
			want:    "hello world",
			wantErr: false,
		},
		{
			name: "with basic format specifiers",
			src: `
			import io;
			fn main(): int32 {
				say a: int = 190;
				say b: float64 = 1.23;
				say c: string = "test";
				io.printf("hello world %d %s %f", a, c, b);
				return 0;
			}
			`,
			want:    "hello world 190 test 1.230000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRun(tt.src, dir)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
