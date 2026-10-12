// Package refactoring демонстрирует архитектурный рефакторинг для тестируемости (Testability).
// Разделение чистой бизнес-логики и сетевых зависимостей позволяет писать надежные unit-тесты.
package refactoring

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Rates представляет схему ответа API Центрального Банка РФ с котировками валют.
type Rates struct {
	Valute struct {
		USD struct {
			ID       string  `json:"ID"`
			NumCode  string  `json:"NumCode"`
			CharCode string  `json:"CharCode"`
			Nominal  int     `json:"Nominal"`
			Name     string  `json:"Name"`
			Value    float64 `json:"Value"`
			Previous float64 `json:"Previous"`
		} `json:"USD"`
	}
}

// rublesToUSD — интеграционная функция конвертации рублей в доллары США.
// Связывает получение курса из сети и математический расчет.
func rublesToUSD(rubles float64) (float64, error) {
	ratio, err := data()
	if err != nil {
		return 0, err
	}

	return calc(rubles, ratio), nil
}

// data выполняет HTTP-запрос к внешнему API ЦБ РФ и парсит актуальный курс доллара.
// Функция подвержена сетевым сбоям и задержкам (I/O dependency).
func data() (float64, error) {
	var data Rates
	const url = "https://www.cbr-xml-daily.ru/daily_json.js"
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return 0, err
	}

	if data.Valute.USD.Value == 0 {
		return 0, errors.New("курс USD равен 0, деление невозможно")
	}

	return data.Valute.USD.Value, nil
}

// calc — ЧИСТАЯ ФУНКЦИЯ (pure function) математического пересчета.
// Не зависит от сети, времени, глобального состояния или внешних сервисов.
// Обладает 100% детерминизмом и тестируется тривиально без моков.
func calc(rubles float64, ratio float64) float64 {
	return rubles / ratio
}

