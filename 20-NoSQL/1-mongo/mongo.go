package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// Официальный драйвер MongoDB для Go.
	// Поддерживает пулы соединений, распределенные транзакции и спецификацию BSON.
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	databaseName   = "data"      // имя учебной БД
	collectionName = "languages" // имя коллекции в учебной БД
)

// lang представляет документ в коллекции MongoDB.
// В отличие от реляционных таблиц, MongoDB хранит BSON (Binary JSON) документы.
// Поля структур сопоставляются с ключами BSON через тэги `bson:"имя"`.
type lang struct {
	ID   int    `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func main() {
	// Конфигурация клиента MongoDB с пулом соединений (Connection Pool).
	// Driver автоматически поддерживает пул TCP-сокетов, восстанавливает соединения при обрывах
	// и балансирует запросы между узлами репликасета (Replica Set).
	mongoOpts := options.Client().ApplyURI("mongodb://localhost:27017/")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		log.Fatal(err)
	}

	// Освобождаем ресурсы и корректно закрываем все соединения пула при завершении программы
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("ошибка при отключении от MongoDB: %v", err)
		}
	}()

	// Проверка доступности СУБД (Ping)
	if err := client.Ping(ctx, nil); err != nil {
		log.Printf("Внимание: сервер MongoDB недоступен (%v). Демонстрация кода продолжается в учебных целях.", err)
		return
	}

	// Вставка документов в БД
	data := []lang{
		{ID: 1, Name: "Go"},
		{ID: 2, Name: "JavaScript"},
	}
	err = insertDocs(client, data)
	if err != nil {
		log.Fatal(err)
	}

	// Получение документов из БД
	docs, err := docs(client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Документы из MongoDB:", docs)
}

// insertDocs вставляет в БД массив документов с использованием InsertOne.
// Для массовой вставки в production-коде рекомендуется InsertMany.
func insertDocs(c *mongo.Client, data []lang) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := c.Database(databaseName).Collection(collectionName)
	for _, doc := range data {
		_, err := collection.InsertOne(ctx, doc)
		if err != nil {
			return fmt.Errorf("ошибка вставки документа: %w", err)
		}
	}
	return nil
}

// docs возвращает все документы из БД с использованием курсора.
// Курсор (mongo.Cursor) читает результаты порциями (batch size) с сервера,
// экономя оперативную память при выборке миллионов строк.
func docs(c *mongo.Client) ([]lang, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := c.Database(databaseName).Collection(collectionName)
	// bson.D{} — пустой фильтр (аналог SELECT * FROM languages)
	// bson.D представляет упорядоченный срез документов (Ordered BSON)
	filter := bson.D{}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения Find: %w", err)
	}
	defer cur.Close(ctx)

	var data []lang
	// Итерация по курсору с декодированием каждого BSON-документа
	for cur.Next(ctx) {
		var l lang
		err := cur.Decode(&l)
		if err != nil {
			return nil, fmt.Errorf("ошибка декодирования документа: %w", err)
		}
		data = append(data, l)
	}

	// Обязательная проверка ошибки итератора курсора
	if err := cur.Err(); err != nil {
		return nil, fmt.Errorf("ошибка курсора MongoDB: %w", err)
	}

	return data, nil
}
