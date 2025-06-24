package main
import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    chNum := make(chan int)
    chChar := make(chan byte)

    // 打印字母
    go func() {
        defer wg.Done()
        for i := 0; i < 26; i++ {
            chChar <- 'A' + byte(i)      // 发送字母
            num := <-chNum               // 等待数字
            fmt.Printf("%c%d", 'A'+byte(i), num)
        }
    }()

    // 打印数字
    go func() {
        defer wg.Done()
        for i := 1; i <= 26; i++ {
            _ = <-chChar               // 等待字母
            chNum <- i                   // 发送数字
            // 字母和数字的打印已在字母协程完成
        }
    }()

    wg.Wait()
    fmt.Println()
}