// Серверное приложение для Homework-16.
// Запускает REST API для информационного сайта о фильмах на базе in-memory хранилища (memsql).
// Для переключения на PostgreSQL достаточно заменить memsql.New() на pgsql.New(ctx, connStr).
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-core-4/lessons/16-db-apps/homework-16/pkg/api"
	"go-core-4/lessons/16-db-apps/homework-16/pkg/db"
	"go-core-4/lessons/16-db-apps/homework-16/pkg/db/memsql"
)

func main() {
	// Инициализация хранилища (in-memory для демонстрации).
	store := memsql.New()

	// Предзаполняем демонстрационными данными о фильмах.
	seedData(store)

	// Создание API.
	a := api.New(store)

	// Настройка HTTP-сервера с промышленными таймаутами.
	server := &http.Server{
		Addr:         ":8080",
		Handler:      a.Router(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Graceful Shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("🎬 Movies API запущен на http://localhost:8080")
		fmt.Println("   GET    /api/v1/movies?studio_id=X — фильмы (все при 0 или фильтр по студии)")
		fmt.Println("   GET    /api/v1/movies/{id}         — фильм по ID")
		fmt.Println("   POST   /api/v1/movies              — создать фильм или массив фильмов")
		fmt.Println("   PUT    /api/v1/movies/{id}         — обновить фильм")
		fmt.Println("   DELETE /api/v1/movies/{id}         — удалить фильм")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка сервера: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\n🛑 Получен сигнал завершения. Остановка сервера...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("принудительная остановка: %v", err)
	}

	fmt.Println("✅ Сервер остановлен корректно.")
}

// seedData наполняет хранилище демонстрационными данными о фильмах.
func seedData(store db.Interface) {
	ctx := context.Background()
	movies := []db.Movie{
		{
			Title:       "Inception",
			ReleaseYear: 2010,
			BoxOffice:   836800000,
			Rating:      "PG-13",
			StudioID:    1,
			StudioName:  "Warner Bros.",
		},
		{
			Title:       "Interstellar",
			ReleaseYear: 2014,
			BoxOffice:   701729206,
			Rating:      "PG-13",
			StudioID:    2,
			StudioName:  "Paramount Pictures",
		},
		{
			Title:       "The Dark Knight",
			ReleaseYear: 2008,
			BoxOffice:   1006000000,
			Rating:      "PG-13",
			StudioID:    1,
			StudioName:  "Warner Bros.",
		},
		{
			Title:       "Pulp Fiction",
			ReleaseYear: 1994,
			BoxOffice:   213928762,
			Rating:      "PG-18",
			StudioID:    3,
			StudioName:  "Miramax",
		},
		{
			Title:       "Toy Story",
			ReleaseYear: 1995,
			BoxOffice:   373554033,
			Rating:      "PG-10",
			StudioID:    4,
			StudioName:  "Pixar Animation Studios",
		},
	}

	if err := store.AddMovies(ctx, movies); err != nil {
		log.Printf("ошибка начального заполнения: %v", err)
	}
}
