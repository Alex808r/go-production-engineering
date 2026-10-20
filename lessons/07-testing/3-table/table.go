// Package table демонстрирует подход табличного тестирования (Table-Driven Tests) в Go.
package table

// reverse разворачивает строку с корректной поддержкой многобайтовых символов Unicode (UTF-8).
//
// Для разворота строки обязательно преобразование в []rune (руны),
// так как один символ в UTF-8 (например, кириллица или эмодзи) может занимать от 1 до 4 байт.
// Разворот байтового среза []byte привел бы к повреждению (corrupted bytes) не-ASCII символов.
func reverse(s string) string {
	runes := []rune(s)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

