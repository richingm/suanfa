package main

import (
	_ "embed"
	"fmt"
	"os"
	"sync"
)

//go:embed lines
var c string

func main() {
	//writeFile()
	testKMP()
}

func produce(data chan string, wg *sync.WaitGroup) {
	n := "abcdefg|"
	data <- n
	wg.Done()
}

func consume(data chan string, done chan bool) {
	f, err := os.Create("lines")
	if err != nil {
		fmt.Println(err)
		return
	}
	for d := range data {
		_, err = fmt.Fprint(f, d)
		if err != nil {
			fmt.Println(err)
			f.Close()
			done <- false
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		done <- false
		return
	}
	done <- true
}

func writeFile() {
	data := make(chan string)
	done := make(chan bool)
	wg := sync.WaitGroup{}
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go produce(data, &wg)
	}
	go consume(data, done)
	go func() {
		wg.Wait()
		close(data)
	}()
	d := <-done
	if d == true {
		fmt.Println("File written successfully")
	} else {
		fmt.Println("File writing failed")
	}
}

func testKMP() {
	fmt.Println(KMP(c, "abcdd"))
	fmt.Println(KMP(c, "cd"))
	fmt.Println(KMP(c, "h,的i"))
}

func KMP(str, substr string) int {
	if substr == "" {
		return 0
	}
	strLen := len(str)
	subLen := len(substr)
	next := make([]int, subLen)
	for i, j := 1, 0; i < subLen; i++ {
		for j > 0 && substr[i] != substr[j] {
			j = next[j-1]
		}
		if substr[i] == substr[j] {
			j++
		}
		next[i] = j
	}
	for i, j := 0, 0; i < strLen; i++ {
		for j > 0 && str[i] != substr[j] {
			j = next[j-1]
		}
		if str[i] == substr[j] {
			j++
		}
		if j == subLen {
			return i - subLen + 1
		}
	}
	return -1
}
