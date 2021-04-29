package bin

import "fmt"

func a() string {
	return "a"
}

func main() {
	aa := a() + a()
	fmt.Println(aa)
}
