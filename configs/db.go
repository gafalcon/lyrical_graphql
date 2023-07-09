package configs

import (
	"context"
	"log"
	"time"

	"github.com/gafalcon/lyrical_graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func ConnectDB() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &DB{client: client}
}
func colHelper(db *DB, collectionName string) *mongo.Collection {
	return db.client.Database("lyrical_graphql").Collection(collectionName)
}

func (db *DB) AddSong(title string) (*model.Song, error) {
	collection := colHelper(db, "song")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	song := model.NewSong{
		Title: title,
	}

	res, err := collection.InsertOne(ctx, song)

	if err != nil {
		return nil, err
	}

	newSong := &model.Song{
		ID:    res.InsertedID.(primitive.ObjectID).Hex(),
		Title: &title,
	}
	return newSong, nil
}

func (db *DB) AddLyric(songId string, content string) (*model.Lyric, error) {
	collection := colHelper(db, "lyric")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(songId)
	lyric := &model.NewLyric{
		SongID:  objId,
		Content: content,
	}
	res, err := collection.InsertOne(ctx, lyric)

	if err != nil {
		return nil, err
	}

	newLyric := &model.Lyric{
		ID:      res.InsertedID.(primitive.ObjectID).Hex(),
		Content: &lyric.Content,
		SongID:  lyric.SongID.Hex(),
	}
	return newLyric, nil
}

func (db *DB) GetSongs() ([]*model.Song, error) {
	collection := colHelper(db, "song")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var songs []*model.Song
	defer cancel()

	res, err := collection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	defer res.Close(ctx)
	for res.Next(ctx) {
		var singleSong *model.Song
		if err = res.Decode(&singleSong); err != nil {
			log.Fatal(err)
		}
		songs = append(songs, singleSong)
	}

	return songs, err
}

func (db *DB) GetSong(id string) (*model.Song, error) {
	collection := colHelper(db, "song")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var song *model.Song
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&song)

	return song, err
}

func (db *DB) GetLyric(id string) (*model.Lyric, error) {
	collection := colHelper(db, "lyric")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var lyric *model.Lyric
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)

	err := collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&lyric)

	return lyric, err
}

func (db *DB) DeleteSong(id string) (int64, error) {
	collection := colHelper(db, "song")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)
	res, err := collection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (db *DB) LikeLyric(id string) (*model.Lyric, error) {
	collection := colHelper(db, "lyric")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"$inc", bson.D{{"likes", 1}}}}
	_, err := collection.UpdateOne(ctx, bson.M{"_id": objId}, filter)
	if err != nil {
		return nil, err
	}

	var lyric *model.Lyric
	err = collection.FindOne(ctx, bson.M{"_id": objId}).Decode(&lyric)
	if err != nil {
		return nil, err
	}
	return lyric, err
}

func (db *DB) GetSongLyrics(songId string) ([]*model.Lyric, error) {
	collection := colHelper(db, "lyric")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var lyrics []*model.Lyric
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(songId)
	res, err := collection.Find(ctx, bson.M{"songId": objId})

	if err != nil {
		return nil, err
	}
	defer res.Close(ctx)
	for res.Next(ctx) {
		var singleLyric *model.Lyric
		if err = res.Decode(&singleLyric); err != nil {
			log.Fatal(err)
		}
		lyrics = append(lyrics, singleLyric)
	}

	return lyrics, err
}
