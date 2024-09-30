package server

import (
	"context"

	"github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/storer"
)

type SongService struct {
	storer *storer.PgStorer
}

func NewSongService(storer *storer.PgStorer) *SongService {
	return &SongService{storer: storer}
}

func (s *SongService) AddSong(ctx context.Context, song *storer.Song) (*storer.Song, error) {
	return s.storer.AddSong(ctx, song)
}

func (s *SongService) GetSongByID(ctx context.Context, id int) (*storer.Song, error) {
	return s.storer.GetSongByID(ctx, id)
}

func (s *SongService) GetAllSongs(ctx context.Context, limit int, offset int, filter string) ([]storer.Song, error) {
	return s.storer.GetAllSongs(ctx, limit, offset, filter)
}

func (s *SongService) UpdateSongByID(ctx context.Context, song *storer.Song) (*storer.Song, error) {
	return s.storer.UpdateSongByID(ctx, song)
}

func (s *SongService) DeleteSongByID(ctx context.Context, id int) error {
	return s.storer.DeleteSongById(ctx, id)
}

func (s *SongService) AddVerseBySongID(ctx context.Context, songID int, verse *storer.Verse) (*storer.Verse, error) {
	return s.storer.AddVerseBySongID(ctx, songID, verse)
}

func (s *SongService) GetVersesBySongID(ctx context.Context, songID int) (*storer.Verse, error) {
	return s.storer.GetVersesBySongID(ctx, songID)
}

func (s *SongService) DeleteVerseBySongID(ctx context.Context, verseID int) error {
	return s.storer.DeleteVerseBySongID(ctx, verseID)
}

func (s *SongService) UpdateVerseBySongID(ctx context.Context, songID int, verse *storer.Verse) (*storer.Verse, error) {
	return s.storer.UpdateVerseBySongID(ctx, songID, verse)
}
