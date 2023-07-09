package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Lyric struct {
	ID      string  `json:"id" bson:"_id"`
	Likes   *int    `json:"likes,omitempty"`
	Content *string `json:"content,omitempty"`
	Song    *Song   `json:"song" bson:"song,omitempty"`
	SongID  string  `json:"songId" bson:"songId"`
}

type NewLyric struct {
	Likes   int
	Content string
	SongID  primitive.ObjectID `bson:"songId"`
}
