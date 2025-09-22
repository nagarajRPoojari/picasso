package io_test

import (
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
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
			name: "basic printf without format specifiers",
			src: `
			import io;
			fn main(): int32 {
				say z: int = 190;
				say s: string = "string";
				say n: float64 = 89.0;
				io.printf("hello world %d, %s, %f", z, s, n);
				return 0;
			}
			`,
			want:    "hello world 190, string, 89.000000",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRun(tt.src, t.TempDir())
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
