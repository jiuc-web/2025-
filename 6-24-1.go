package main
import(
	"fmt"
	"sync"
)

var (
	count = 0
	mu sync.Mutex
)

func f(wg *sync.WaitGroup) {
	for i := 0; i < 100 ; i++{
		mu.Lock()
		count++
		mu.Unlock()
	}
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 50 ;i++{
		wg.Add(1)
		go f(&wg)
	}
	wg.Wait()
	fmt.Println("Final count: ",count)
}