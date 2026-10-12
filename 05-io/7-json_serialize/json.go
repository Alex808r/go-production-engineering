package main

import (
	"encoding/json"
	"fmt"
)

// Point представляет точку с координатами.
// Важные правила работы с encoding/json:
// 1. Поля структуры ОБЯЗАНЫ быть экспортируемыми (начинаться с заглавной буквы), иначе рефлексия json их проигнорирует.
// 2. Теги структур `json:"..."` переопределяют имена ключей в итоговом JSON-объекте.
//    (например, для экономии трафика можно сократить имена до `json:"1"` и `json:"2"` или использовать `json:"lon"` и `json:"lat"`).
type Point struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

func main() {
	// Создаем срез из нескольких точек для наглядности
	var points []Point
	for i := 0; i < 3; i++ {
		points = append(points, Point{Lon: float64(i * 10), Lat: float64(i * 20)})
	}

	// 1. Сериализация (MarshalIndent): форматирует JSON с отступами для удобного чтения человеком.
	// Для компактной передачи по сети используют обычный json.Marshal(points) без лишних пробелов.
	b, err := json.MarshalIndent(points, "", "  ")
	if err != nil {
		fmt.Println("Ошибка сериализации:", err)
		return
	}
	fmt.Println("Сериализованный JSON:")
	fmt.Println(string(b))

	// 2. Десериализация (Unmarshal): парсинг JSON-байтов обратно в срез структур Go.
	var decoded []Point
	err = json.Unmarshal(b, &decoded)
	if err != nil {
		fmt.Println("Ошибка десериализации:", err)
		return
	}
	fmt.Printf("\nУспешно десериализовано объектов: %d\nПервая точка: Lon=%.1f, Lat=%.1f\n", len(decoded), decoded[0].Lon, decoded[0].Lat)
}
