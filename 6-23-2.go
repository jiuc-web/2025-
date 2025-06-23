package main
import "fmt"

func Deduplicate(nums []int) []int{
	num := []int{}// 创建一个空切片用于存储去重后的元素
	if len(nums) == 0 {
		return num // 如果输入切片为空，直接返回空切片
	}
	num = append(num, nums[0]) // 将第一个元素添加到num中
	// 遍历输入切片中的每个元素
	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(num); j++ {
			if nums[i] == num[j] {
				break // 如果当前元素已经存在于num中，跳过
			}
			if j == len(num)-1 { // 如果遍历到最后一个元素还没有找到相同的元素
				num = append(num, nums[i]) // 将当前元素添加到num中
			}
		}		
	}
	return num // 返回去重后的切片
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5}
	fmt.Println("原始切片：", nums)
	num := Deduplicate(nums)
	fmt.Println("去重后的切片：", num)
}