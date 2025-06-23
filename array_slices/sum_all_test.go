package arrayslices

import (
	"reflect"
	"testing"
)

func TestSumAll(t *testing.T) {
	x := []int{1, 2}
	y := []int{0, 9}
	got := SumAll(x, y)
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
