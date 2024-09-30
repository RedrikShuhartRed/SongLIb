package handler

type songDetail struct {
	ReleaseDate string `json:"release_date" db:"release_date"`
	VersesEn    string `json:"verses_en" db:"verses_en"`
	VersesRu    string `json:"verses_ru" db:"verses_ru"`
	Link        string `json:"link" db:"link"`
}

type addSongRequest struct {
	GroupName string `json:"group_name" binding:"required"`
	Song      string `json:"song" binding:"required"`
}

type updateSongRequest struct {
	GroupName   string `json:"group_name"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
}

type updateVerseRequest struct {
	VerseTextEn string `json:"verse_text_en" `
	VerseTextRu string `json:"verse_text_ru"`
}
