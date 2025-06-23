package arrayslices

import "testing"

func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}
		got := Sum(numbers)
		want := 15
		assertCorrectMessage(t, got, want)
	})

	t.Run("collection of any size of numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		got := Sum(numbers)
		want := 55
		assertCorrectMessage(t, got, want)
	})

}

func assertCorrectMessage(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func BenchmarkSum(b *testing.B) {
	numbers := []int{1, 2, 3, 4, 5}

	for b.Loop() {
		Sum(numbers)
	}
}
