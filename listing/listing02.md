Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)

func test() (x int) {
	defer func() {
		x++
	}()
	x = 1 // x=0 -> x=1
	// defer развернется здесь x=1 -> x = 2
	return // return x=2
}

func anotherTest() int {
	var x int // 0
	defer func() {
		x++
	}()
	x = 1
	return x // возращаемое значение зафиксировано return 1
}

func main() {
	fmt.Println(test())        // 2
	fmt.Println(anotherTest()) //1
}

```

Ответ:
```
в коде

```
