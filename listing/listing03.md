Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```go
package main

import (
	"fmt"
	"os"
)

func Foo() error {
	var err *os.PathError = nil
	return err
}

func main() {
	err := Foo()
	// конкретное значение равно nil
	fmt.Println(err)

	fmt.Println(err == nil)
}

```

Ответ:
```
в golang интерфейс (в данном случае interface error равен nil),
если tab и data равный nil. В данном слачае метаданные itab не nil потому что
используется конкретный тип os.PathError
```
