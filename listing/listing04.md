Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func(ch chan int) {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}(ch)

	for n := range ch {
		println(n)
	}
}
```

Ответ:
```
первая goroutine и main подсоединятся к каналу и последовательно запишут/прочитают числа от 0 до 9.
main будет бесконечно ждать в незакрытый канал, однако runtime обнаружит состояние гонки и завершит main с ошибкой 
```
