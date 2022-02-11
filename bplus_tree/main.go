package main

import "log"

func binarySearch(nums []int, n int) int {
	l := 0
	r := len(nums)

	for l < r {
		mid := (l + r) / 2
		if nums[mid] > n {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return r
}

type Stu struct {
	age  int
	name string
}

func main() {
	log.Println(binarySearch([]int{1,3,5,7},1))
}
