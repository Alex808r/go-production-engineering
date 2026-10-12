// Package ocp демонстрирует принцип открытости/закрытости (Open-Closed Principle) в Go.
// Модули должны быть открыты для расширения, но закрыты для изменения.
// Новое поведение добавляется созданием новых типов, реализующих интерфейс,
// без правки существующего проверенного кода.
package ocp

// Movie — сущность фильма.
type Movie struct {
	ID          int64
	Title       string
	ReleaseYear int
	Rating      string // PG-10, PG-13, PG-18
	BoxOffice   int64
}

// FilterSpecification определяет контракт критерия фильтрации (паттерн Specification).
type FilterSpecification interface {
	IsSatisfiedBy(m Movie) bool
}

// YearFilter фильтрует фильмы по минимальному году выпуска.
type YearFilter struct {
	MinYear int
}

func (f YearFilter) IsSatisfiedBy(m Movie) bool {
	return m.ReleaseYear >= f.MinYear
}

// RatingFilter фильтрует фильмы по возрастному рейтингу.
type RatingFilter struct {
	Rating string
}

func (f RatingFilter) IsSatisfiedBy(m Movie) bool {
	return m.Rating == f.Rating
}

// BoxOfficeFilter фильтрует фильмы по минимальным сборам.
type BoxOfficeFilter struct {
	MinBoxOffice int64
}

func (f BoxOfficeFilter) IsSatisfiedBy(m Movie) bool {
	return m.BoxOffice >= f.MinBoxOffice
}

// AndSpecification объединяет несколько спецификаций (логическое И).
type AndSpecification struct {
	Specs []FilterSpecification
}

func (a AndSpecification) IsSatisfiedBy(m Movie) bool {
	for _, spec := range a.Specs {
		if !spec.IsSatisfiedBy(m) {
			return false
		}
	}
	return true
}

// MovieFilter фильтрует список фильмов согласно любой переданной спецификации.
// Код функции закрыт для изменений: при появлении новых фильтров (например, по режиссеру или актеру)
// мы не модифицируем MovieFilter, а просто реализуем интерфейс FilterSpecification!
func FilterMovies(movies []Movie, spec FilterSpecification) []Movie {
	var result []Movie
	for _, m := range movies {
		if spec.IsSatisfiedBy(m) {
			result = append(result, m)
		}
	}
	return result
}
