package storer

import "time"

type Song struct {
	ID          int        `json:"id" db:"id"`
	GroupName   string     `json:"group_name" db:"group_name"`
	Song        string     `json:"song" db:"song"`
	ReleaseDate string     `json:"release_date" db:"release_date"`
	Link        string     `json:"link" db:"link"`
	CreateAt    time.Time  `json:"created_at" db:"created_at"`
	UpdateAt    *time.Time `json:"updated_at" db:"updated_at"`
}

type Verse struct {
	ID           int        `json:"id" db:"id"`
	SongID       int        `json:"song_id" db:"song_id"`
	VerseTextEng string     `json:"verse_text_en" db:"verse_text_en"`
	VerseTextRu  string     `json:"verse_text_ru" db:"verse_text_ru"`
	CreateAt     time.Time  `json:"created_at" db:"created_at"`
	UpdateAt     *time.Time `json:"updated_at" db:"updated_at"`
}
