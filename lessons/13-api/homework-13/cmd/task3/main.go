package main

import (
	"fmt"
	"unsafe"

	"github.com/DmitriyVTitov/size"
)

type example struct {
	a []int
	b bool
	c int32
	d string
}

func main() {
	fmt.Println("=== Разбор Задачи №3 (Занятие №13: Модель памяти и size.Of) ===")

	ex := example{
		a: []int{1, 2, 3}, // 24 байта (заголовок слайса) + 24 байта (массив из 3 int по 8 байт) = 48 байт
		b: true,           // 1 байт
		d: "1234",         // 16 байт (заголовок строки) + 4 байта (буфер 4 байт UTF-8) = 20 байт
	} // 48 байт (поверхностный unsafe.Sizeof) / 76 байт (глубокий size.Of)
	fmt.Println("Размер в байтах для ex:", size.Of(ex))

	ex1 := example{
		a: []int{1, 2, 3}, // 24 байта (заголовок слайса) + 24 байта (массив из 3 int по 8 байт) = 48 байт
		b: true,           // 1 байт
		d: "1234",         // 16 байт (заголовок строки) + 4 байта (буфер 4 байт UTF-8) = 20 байт
		c: 100,            // 4 байта внутри структуры
	} // 48 байт (поверхностный unsafe.Sizeof) / 76 байт (глубокий size.Of)
	fmt.Println("Размер в байтах для ex1:", size.Of(ex1))

	fmt.Println("\n--- Детальная анатомия памяти ---")
	fmt.Printf("1. Поверхностный размер структуры unsafe.Sizeof(example): %d байт\n", unsafe.Sizeof(ex))
	fmt.Printf("   - смещение поля a ([]int):  %d..%d (размер: %d, выравнивание: %d)\n",
		unsafe.Offsetof(ex.a), unsafe.Offsetof(ex.a)+unsafe.Sizeof(ex.a)-1, unsafe.Sizeof(ex.a), unsafe.Alignof(ex.a))
	fmt.Printf("   - смещение поля b (bool):   %d..%d (размер: %d, выравнивание: %d)\n",
		unsafe.Offsetof(ex.b), unsafe.Offsetof(ex.b)+unsafe.Sizeof(ex.b)-1, unsafe.Sizeof(ex.b), unsafe.Alignof(ex.b))
	fmt.Printf("   - padding между b и c:      %d..%d (3 байта padding для кратности 4 байтам int32)\n",
		unsafe.Offsetof(ex.b)+unsafe.Sizeof(ex.b), unsafe.Offsetof(ex.c)-1)
	fmt.Printf("   - смещение поля c (int32):  %d..%d (размер: %d, выравнивание: %d)\n",
		unsafe.Offsetof(ex.c), unsafe.Offsetof(ex.c)+unsafe.Sizeof(ex.c)-1, unsafe.Sizeof(ex.c), unsafe.Alignof(ex.c))
	fmt.Printf("   - смещение поля d (string): %d..%d (размер: %d, выравнивание: %d)\n",
		unsafe.Offsetof(ex.d), unsafe.Offsetof(ex.d)+unsafe.Sizeof(ex.d)-1, unsafe.Sizeof(ex.d), unsafe.Alignof(ex.d))
	fmt.Println("2. Дополнительная память в куче/буферах:")
	fmt.Printf("   - Массив элементов среза ex.a: 3 элемента int * 8 байт = %d байт\n", cap(ex.a)*8)
	fmt.Printf("   - Буфер строки ex.d: 4 символа ASCII = %d байта\n", len(ex.d))
	fmt.Printf("3. Итого глубокий размер: 48 (структура) + 24 (массив) + 4 (строка) = %d байт\n", size.Of(ex))
	fmt.Printf("4. Почему size.Of(ex) == size.Of(ex1)?\n")
	fmt.Printf("   Потому что поле c int32 УЖЕ входит в 48 байт структуры на смещении 28..31.\n")
	fmt.Printf("   В ex оно равно 0 (zero value), а в ex1 равно 100. Скаляр int32 не аллоцирует память в куче.\n")
}
