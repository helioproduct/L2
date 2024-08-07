Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

**Ответ:**

При присваивании возвращаемого значения test переменной err -
значение будет обвернуто в интерфейс. Интерфейс будет иметь динамический тип
customError и динамическое значение nil.

Но это не тоже самое, что интерфейсное значение nil, так как в последнем случае интерфейс не имеет динамического типа и значения. Поэтому ```err != nil``` - верное утверждение, и будет выполнена ветка кода, печатающая "error".

Вывод программы:
error