package main
import "fmt"

func Prime(n int) bool {
	if n < 2 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func main(){
	n := 0;
	fmt.Print("请输入一个整数：")
	fmt.Scanf("%d", &n)
	if Prime(n) {
		fmt.Printf("%d 是素数\n", n)
	} else {
		fmt.Printf("%d 不是素数\n", n)
	}
}