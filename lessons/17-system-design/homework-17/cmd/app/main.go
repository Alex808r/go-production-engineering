// Исполняемый файл приложения Homework-17 (Clean Architecture Demonstration).
// Демонстрирует Composition Root: сборку слоев Чистой Архитектуры
// (Entity -> UseCase -> Adapter -> Framework/HTTP).
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

	"github.com/gorilla/mux"

	"go-core-4/lessons/17-system-design/homework-17/cleanarch/adapter/controller"
	"go-core-4/lessons/17-system-design/homework-17/cleanarch/adapter/repository"
	"go-core-4/lessons/17-system-design/homework-17/cleanarch/usecase"
)

func main() {
	// 1. Создание вторичного адаптера (хранилища)
	repo := repository.NewMemoryMovieRepository()

	// 2. Создание сценариев использования (Use Cases) с инверсией зависимостей
	uc := usecase.NewMovieUseCase(repo)

	// 3. Предзаполнение демонстрационными данными
	seed(uc)

	// 4. Создание первичного адаптера (HTTP-контроллера)
	ctrl := controller.NewMovieHTTPController(uc)

	// 5. Внешний слой: роутер и сервер
	router := mux.NewRouter()
	ctrl.RegisterRoutes(router)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		fmt.Println("🚀 Clean Architecture Cinema Service запущен на http://localhost:8080")
		fmt.Println("   GET    /api/v1/movies?studio_id=X — список фильмов")
		fmt.Println("   GET    /api/v1/movies/{id}         — фильм по ID")
		fmt.Println("   POST   /api/v1/movies              — создать фильм")
		fmt.Println("   POST   /api/v1/movies/batch        — создать пакет фильмов")
		fmt.Println("   DELETE /api/v1/movies/{id}         — удалить фильм")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ошибка сервера: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\n🛑 Остановка сервиса...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("принудительная остановка: %v", err)
	}

	fmt.Println("✅ Сервис успешно остановлен.")
}

func seed(uc usecase.MovieUseCase) {
	ctx := context.Background()
	_ = uc.CreateMoviesBatch(ctx, []usecase.CreateMovieDTO{
		{Title: "Inception", ReleaseYear: 2010, BoxOffice: 836000000, Rating: "PG-13", StudioID: 1},
		{Title: "Interstellar", ReleaseYear: 2014, BoxOffice: 701000000, Rating: "PG-13", StudioID: 2},
		{Title: "The Dark Knight", ReleaseYear: 2008, BoxOffice: 1006000000, Rating: "PG-13", StudioID: 1},
		{Title: "Pulp Fiction", ReleaseYear: 1994, BoxOffice: 213000000, Rating: "PG-18", StudioID: 3},
	})
}
