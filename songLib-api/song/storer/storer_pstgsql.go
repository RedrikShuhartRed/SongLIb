package storer

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type PgStorer struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewPgStorer(db *sqlx.DB, logger *zap.Logger) *PgStorer {
	return &PgStorer{
		db:     db,
		logger: logger,
	}
}

func (ps *PgStorer) GetAllSongs(ctx context.Context, limit int, offset int, filter string) ([]Song, error) {
	ps.logger.Debug("Getting all songs with filter", zap.String("filter", filter), zap.Int("limit", limit), zap.Int("offset", offset))
	var songs []Song
	query := "SELECT * FROM songs WHERE song ILIKE $1 OR group_name ILIKE $1 LIMIT $2 OFFSET $3"
	err := ps.db.SelectContext(ctx, &songs, query, "%"+filter+"%", limit, offset)
	if err != nil {
		ps.logger.Error("Error getting songs", zap.Error(err))
		return nil, fmt.Errorf("Error getting songs: %w", err)
	}

	return songs, nil
}

func (ps *PgStorer) GetSongByID(ctx context.Context, id int) (*Song, error) {
	ps.logger.Debug("Getting song by ID", zap.Int("id", id))
	var song Song
	query := "SELECT * FROM songs WHERE id = $1"
	err := ps.db.GetContext(ctx, &song, query, id)
	if err != nil {
		ps.logger.Error("Error getting song by ID", zap.Error(err))
		return nil, fmt.Errorf("Error getting song by ID: %w", err)
	}
	return &song, nil
}

func (ps *PgStorer) DeleteSongById(ctx context.Context, id int) error {
	ps.logger.Debug("Deleting song by ID", zap.Int("id", id))
	query := "DELETE FROM verses WHERE song_id = $1"
	_, err := ps.db.ExecContext(ctx, query, id)
	if err != nil {
		ps.logger.Error("Error deleting song by ID", zap.Error(err))
		return fmt.Errorf("Error deleting song by ID: %w", err)
	}
	return nil
}

func (ps *PgStorer) AddSong(ctx context.Context, s *Song) (*Song, error) {
	ps.logger.Debug("Adding song", zap.String("song", s.Song))
	query := "INSERT INTO songs (group_name, song, release_date, link, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id, created_at"
	err := ps.db.QueryRowContext(ctx, query, s.GroupName, s.Song, s.ReleaseDate, s.Link).Scan(&s.ID, &s.CreateAt)
	if err != nil {
		ps.logger.Error("Error inserting song", zap.Error(err))
		return nil, fmt.Errorf("Error adding song: %w", err)
	}
	return s, nil
}

func (ps *PgStorer) UpdateSongByID(ctx context.Context, s *Song) (*Song, error) {
	ps.logger.Debug("Updating song by ID", zap.Int("id", s.ID))
	query := "UPDATE songs SET group_name=$1, song=$2, release_date=$3, link=$4, updated_at=NOW() WHERE id=$5 RETURNING updated_at"
	err := ps.db.QueryRowContext(ctx, query, s.GroupName, s.Song, s.ReleaseDate, s.Link, s.ID).Scan(&s.UpdateAt)
	if err != nil {
		ps.logger.Error("Error updating song", zap.Error(err))
		return nil, fmt.Errorf("Error updating song: %w", err)
	}
	return s, nil
}

func (ps *PgStorer) AddVerseBySongID(ctx context.Context, songID int, v *Verse) (*Verse, error) {
	ps.logger.Debug("Adding verse to song", zap.Int("song_id", songID))
	query := "INSERT INTO verses (song_id, verse_text_en, verse_text_ru, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, created_at"
	err := ps.db.QueryRowContext(ctx, query, songID, v.VerseTextEng, v.VerseTextRu).Scan(&v.ID, &v.CreateAt)
	if err != nil {
		ps.logger.Error("Error inserting verse", zap.Error(err))
		return nil, fmt.Errorf("Error adding verse: %w", err)
	}
	return v, nil
}
func (ps *PgStorer) GetVersesBySongID(ctx context.Context, songID int) (*Verse, error) {
	ps.logger.Debug("Getting verses by song ID", zap.Int("song_id", songID))
	var verse Verse
	query := "SELECT verse_text_en, verse_text_ru FROM verses WHERE song_id = $1"
	err := ps.db.GetContext(ctx, &verse, query, songID)
	if err != nil {
		ps.logger.Error("Error getting verses by song ID", zap.Error(err))
		return nil, fmt.Errorf("Error getting verses by song ID: %w", err)
	}

	return &verse, nil
}
func (ps *PgStorer) DeleteVerseBySongID(ctx context.Context, songID int) error {
	ps.logger.Debug("Deleting verse by ID", zap.Int("song_id", songID))
	query := "DELETE FROM verses WHERE song_id = $1"
	_, err := ps.db.ExecContext(ctx, query, songID)
	if err != nil {
		ps.logger.Error("Error deleting verse", zap.Error(err))
		return fmt.Errorf("Error deleting verse by ID: %w", err)
	}
	return nil
}

func (ps *PgStorer) UpdateVerseBySongID(ctx context.Context, songID int, v *Verse) (*Verse, error) {
	ps.logger.Debug("Updating verse by ID", zap.Int("verse_id", v.ID))
	query := "UPDATE verses SET verse_text_en=$1, verse_text_ru=$2, updated_at=NOW() WHERE song_id=$3 RETURNING updated_at, id"
	err := ps.db.QueryRowContext(ctx, query, v.VerseTextEng, v.VerseTextRu, songID).Scan(&v.UpdateAt, &v.ID)
	if err != nil {
		ps.logger.Error("Error updating verse", zap.Error(err))
		return nil, fmt.Errorf("Error updating verse: %w", err)
	}
	return v, nil
}
