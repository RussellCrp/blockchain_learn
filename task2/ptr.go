package main

import (
	"fmt"
)

func main() {
	num := 1
	ptrAdd(&num)
	fmt.Println("num + 10:", num)

	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	ptrSlice(&slice)
	fmt.Println("slice * 2 :", slice)

}

func ptrAdd(num *int) {
	*num = *num + 10
}

func ptrSlice(slice *[]int) {
	for i, val := range *slice {
		(*slice)[i] = val * 2
	}
}
