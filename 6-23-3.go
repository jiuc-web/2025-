package main
import "fmt"

type TreeNode struct {
	num  int
	left *TreeNode
	right *TreeNode
}

func BuildTreePreorder(nums []int, idx *int) *TreeNode {
	if *idx >= len(nums) {
        return nil
    }
    node := &TreeNode{num: nums[*idx]}
    *idx++
    node.left = BuildTreePreorder(nums, idx)
    node.right = BuildTreePreorder(nums, idx)
    return node
}

func PreOrder(node *TreeNode) {
	if node == nil {
		return
	}
	fmt.Printf("%d ", node.num)
	PreOrder(node.left)
	PreOrder(node.right)
}

func InOrder(node *TreeNode) {
	if node == nil {
		return
	}
	InOrder(node.left)
	fmt.Printf("%d ", node.num)
	InOrder(node.right)
}

func PostOrder(node *TreeNode) {
	if node == nil {
		return
	}
	PostOrder(node.left)
	PostOrder(node.right)
	fmt.Printf("%d ", node.num)
}

func main() {
    nums := []int{1, 2, 3, 4, 5, 6, 7}
    idx := 0
    root := BuildTreePreorder(nums, &idx)
    fmt.Print("先序遍历: ")
    PreOrder(root)
    fmt.Println()
    fmt.Print("中序遍历: ")
    InOrder(root)
    fmt.Println()
    fmt.Print("后序遍历: ")
    PostOrder(root)
    fmt.Println()
}