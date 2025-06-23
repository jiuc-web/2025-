package main
import "fmt"

func merge(nums1 []int, m int, nums2 []int, n int) []int {
    i, j, k := 0, 0, 0
    nums := make([]int, m+n)
    for i < m && j < n {
        if nums1[i] < nums2[j] {
            nums[k] = nums1[i]
            i++
        } else {
            nums[k] = nums2[j]
            j++
        }
        k++
    }
    for i < m {
        nums[k] = nums1[i]
        i++
        k++
    }
    for j < n {
        nums[k] = nums2[j]
        j++
        k++
    }
    return nums
}

func main() {
	nums1 := []int{1, 3, 5, 7, 9}
	m := len(nums1)
	nums2 := []int{2, 4, 6, 8, 10}
	n := len(nums2)
	nums := merge(nums1, m, nums2, n)
	fmt.Println("合并后的切片：", nums)
}