package exercises

import (
	"slices"
	"testing"
)

func TestFanOutSquares(t *testing.T) {
	input := Generate(1, 2, 3, 4, 5)
	actual := collect(FanOutSquares(input, 3))
	slices.Sort(actual)

	expected := []int{1, 4, 9, 16, 25}
	if !slices.Equal(actual, expected) {
		t.Fatalf("结果 = %v，期望 %v", actual, expected)
	}
}

func TestFanIn(t *testing.T) {
	actual := collect(FanIn(Generate(1, 2), Generate(3, 4)))
	slices.Sort(actual)

	if !slices.Equal(actual, []int{1, 2, 3, 4}) {
		t.Fatalf("合并结果错误: %v", actual)
	}
}

func TestPipeline(t *testing.T) {
	actual := collect(Square(Positive(Generate(-1, 0, 2, 3))))

	if !slices.Equal(actual, []int{4, 9}) {
		t.Fatalf("流水线结果错误: %v", actual)
	}
}

func collect(input <-chan int) []int {
	var values []int
	for value := range input {
		values = append(values, value)
	}
	return values
}
