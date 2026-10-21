// -----------------------------------------------------------------------------
// ИМЯ МОДУЛЯ (MODULE PATH):
// Префикс "github.com/Alex808r/wallet" обеспечивает глобальную уникальность
// пространства имен в экосистеме Go Modules.
// При сборке и тестах компилятор Go сопоставляет импорты проекта с этим именем
// и разрешает их ЛОКАЛЬНО в каталоге проекта без сетевых обращений к GitHub.
// -----------------------------------------------------------------------------
module github.com/Alex808r/wallet

go 1.27.1

require (
	github.com/google/uuid v1.6.0
	github.com/jackc/pgx/v5 v5.7.2
	github.com/shopspring/decimal v1.4.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)
