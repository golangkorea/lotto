package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	file, err := os.Open("names.txt")
	if err != nil {
		fmt.Println("파일을 열 수 없습니다:", err)
		return
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		if name != "" {
			names = append(names, name)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("파일을 읽는 도중 에러 발생:", err)
		return
	}

	names = removeDuplicates(names)
	fmt.Printf("총 %d명의 추첨 후보자가 있습니다.\n", len(names))
	rand.Seed(time.Now().UnixNano())

	for len(names) > 0 {
		fmt.Println("엔터를 누르면 랜덤으로 한 명을 추첨합니다.")
		fmt.Scanln()

		index := rand.Intn(len(names))

		for i := 0; i < 3; i++ {
			fmt.Print("두구")
			time.Sleep(500 * time.Millisecond)
		}
		fmt.Println("!")

		printGopherArt()
		time.Sleep(1 * time.Second)

		fmt.Println("추첨된 이름:", names[index])
		names = append(names[:index], names[index+1:]...)

		fmt.Printf("남은 추첨 후보자: %d명\n", len(names))
	}

	fmt.Println("더 이상 추첨할 이름이 없습니다.")
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
		} else {
			encountered[elements[v]] = true
			result = append(result, elements[v])
		}
	}
	return result
}

func printGopherArt() {
	gopherArt := `
           ,_---~~~~~----._                     
    _,,_,*^____      _____'''*g*\"*, 
 / __/ /'     ^.  /      \ ^@q   f 
[  @f | @))    |  | @))   l  0 _/  
 \'/   \~____ / __ \_____/    \   
         ''*---*''
`
	fmt.Println(gopherArt)
}

