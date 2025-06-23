package arrayslices

func SumAll(sliceOfNums ...[]int) []int {
	res := []int{}

	for _, slice := range sliceOfNums {
		res = append(res, Sum(slice))
	}

	return res
}
