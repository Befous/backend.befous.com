package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comics struct {
	ID                 *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ID_Mangadex        string              `json:"id_mangadex" bson:"id_mangadex"`
	Type               string              `json:"type" bson:"type"`
	Title              string              `json:"title" bson:"title"`
	Alternative_Titles []string            `json:"alternative_titles,omitempty" bson:"alternative_titles,omitempty"`
	Description        string              `json:"description" bson:"description"`
	Original_Language  string              `json:"original_language" bson:"original_language"`
	Demographic        string              `json:"demographic" bson:"demographic"`
	Status             string              `json:"status" bson:"status"`
	Year               int                 `json:"year" bson:"year"`
	Content_Rating     string              `json:"content_rating" bson:"content_rating"`
	Tags               struct {
		Name string `json:"name" bson:"name"`
	} `json:"tags" bson:"tags"`
	State string `json:"state" bson:"state"`
	Cover string `json:"cover" bson:"cover"`
	Link  struct {
		Raw          string `json:"raw" bson:"raw"`
		Baka_Updates string `json:"baka_updates" bson:"baka_updates"`
		Anime_Planet string `json:"anime_planet" bson:"anime_planet"`
		AniList      string `json:"aniList" bson:"aniList"`
		Kitsu        string `json:"kitsu" bson:"kitsu"`
		MyAnimeList  string `json:"myAnimeList" bson:"myAnimeList"`
	} `json:"link" bson:"link"`
}
