package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/storer"
)

func fetchSongDetails(group, song string) (songDetail, error) {

	groupURL := url.QueryEscape(group)
	songURL := url.QueryEscape(song)
	apiURL := fmt.Sprintf("http://127.0.0.1:8080/info?group=%s&song=%s", groupURL, songURL)

	log.Printf("Fetching song details from: %s", apiURL)
	response, err := http.Get(apiURL)
	if err != nil {
		log.Printf("Error fetching song details from external API: %v", err)
		return songDetail{}, fmt.Errorf("Error fetching song details from external API: %w", err)

	}
	var songDetails songDetail
	err = json.NewDecoder(response.Body).Decode(&songDetails)
	if err != nil {
		log.Printf("Error decoding song details JSON: %v", err)
		return songDetail{}, fmt.Errorf("Error decoding song details JSON: %w", err)
	}
	return songDetails, nil
}

func patchSong(song updateSongRequest, s storer.Song) *storer.Song {
	if song.GroupName != "" {
		s.GroupName = song.GroupName
	}
	if song.Song != "" {
		s.Song = song.Song
	}
	if song.ReleaseDate != "" {
		s.ReleaseDate = song.ReleaseDate
	}
	if song.Link != "" {
		s.Link = song.Link
	}

	return &s
}

func patchVerse(verse updateVerseRequest, v storer.Verse) *storer.Verse {
	if verse.VerseTextEn != "" {
		v.VerseTextEng = verse.VerseTextEn
	}
	if verse.VerseTextRu != "" {
		v.VerseTextRu = verse.VerseTextRu
	}

	return &v
}
