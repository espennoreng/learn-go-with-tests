package arrayslices

func SumAllTails(tailsToSum ...[]int) []int {
	res := []int{}

	for _, slice := range tailsToSum {
		if len(slice) == 0 {
			res = append(res, 0)
		} else {
			tail := slice[1:]
			res = append(res, Sum(tail))
		}
	}

	return res
}
