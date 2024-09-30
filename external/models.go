package external

type DataAPI struct {
	ReleaseDate string `json:"release_date"`
	VersesEn    string `json:"verses_en" db:"verses_en"`
	VersesRu    string `json:"verses_ru" db:"verses_ru"`
	Link        string `json:"link"`
}
