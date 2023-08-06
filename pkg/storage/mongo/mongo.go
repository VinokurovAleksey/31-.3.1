package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"GoNews/pkg/storage"
)

// База данных Mongo
type Database struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// Создает новый экземпляр структуры базы данных
func NewDatabase(connection string, dbName string, collectionName string) (*Database, error) {
	clientOptions := options.Client().ApplyURI(connection)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	db := client.Database(dbName)
	collection := db.Collection(collectionName)

	return &Database{client: client, collection: collection}, nil
}

// Закрывает соединение
func (d *Database) Close() error {
	err := d.client.Disconnect(context.Background())
	if err != nil {
		return fmt.Errorf("failed to close the database connection: %v", err)
	}
	return nil
}

// Извлекает все сообщения
func (d *Database) Posts() ([]storage.Post, error) {
	cursor, err := d.collection.Find(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve posts: %v", err)
	}
	defer cursor.Close(context.Background())

	var posts []storage.Post
	for cursor.Next(context.Background()) {
		var post storage.Post
		err := cursor.Decode(&post)
		if err != nil {
			return nil, fmt.Errorf("failed to decode post: %v", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// Добавляет новый пост
func (d *Database) AddPost(post storage.Post) error {
	_, err := d.collection.InsertOne(context.Background(), post)
	if err != nil {
		return fmt.Errorf("failed to add post: %v", err)
	}
	return nil
}

// Обновляет существующий пост
func (d *Database) UpdatePost(post storage.Post) error {
	_, err := d.collection.ReplaceOne(context.Background(), bson.M{"_id": post.ID}, post)
	if err != nil {
		return fmt.Errorf("failed to update post: %v", err)
	}
	return nil
}

// Удаляет запись
func (d *Database) DeletePost(post storage.Post) error {
	_, err := d.collection.DeleteOne(context.Background(), bson.M{"_id": post.ID})
	if err != nil {
		return fmt.Errorf("failed to delete post: %v", err)
	}
	return nil
}
