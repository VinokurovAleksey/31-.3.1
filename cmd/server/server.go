package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/memdb"
	"GoNews/pkg/storage/mongo"
	"GoNews/pkg/storage/postgres"
	"log"
	"net/http"
)

// Сервер GoNews.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера.
	var srv server

	// Создаём объекты баз данных.
	//
	// БД в памяти.
	db := memdb.New()

	// Реляционная БД PostgreSQL.
	connectionString := "user=postgres dbname=GoNews password=123 host=localhost port=5432 sslmode=disable"
	db2, err := postgres.NewDatabase(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer db2.Close()

	// Документная БД MongoDB.
	connectionString = "mongodb://localhost:27017"
	dbName := "mydb"
	collectionName := "posts"

	db3, err := mongo.NewDatabase(connectionString, dbName, collectionName)
	if err != nil {
		log.Fatal(err)
	}

	defer db3.Close()

	_, _, _ = db, db2, db3

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db2

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов,
	// поэтому сервер будет все запросы отправлять на маршрутизатор.
	// Маршрутизатор будет выбирать нужный обработчик.
	http.ListenAndServe(":8080", srv.api.Router())
}
