package retcode_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/wavesoftware/k8s-aware/pkg/utils/retcode"
)

func TestCalc(t *testing.T) {
	cases := testCases()
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			if got := retcode.Calc(tt.err); got != tt.want {
				t.Errorf("Calc() = %v, want %v", got, tt.want)
			}
		})
	}
}

var errExampleError = errors.New("example error")

func testCases() []testCase {
	return []testCase{{
		name: "nil",
		err:  nil,
		want: 0,
	}, {
		name: "errExampleError",
		err:  errExampleError,
		want: 133,
	}, {
		name: "error of wrap caused by 12345",
		err:  fmt.Errorf("%w: 12345", errExampleError),
		want: 249,
	}}
}

type testCase struct {
	name string
	err  error
	want int
}
