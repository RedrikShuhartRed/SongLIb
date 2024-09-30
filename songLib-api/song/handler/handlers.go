package handler

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	service "github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/service"
	"github.com/RedrikShuhartRed/EfMobSongLib/songLib-api/song/storer"
)

type handler struct {
	ctx     context.Context
	service *service.SongService
	logger  *zap.Logger
}

func NewHandler(service *service.SongService, logger *zap.Logger) *handler {
	return &handler{
		ctx:     context.Background(),
		service: service,
		logger:  logger,
	}
}

// @Summary Add a new song
// @Description Adds a new song along with its verses to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body addSongRequest true "Song details"
// @Success 201 {object} storer.Song
// @Failure 400 {object} string "Error binding JSON"
// @Failure 400 {object} string "Error fetching song details"
// @Failure 500 {object} string "Internal server error"
// @Router /song [post]
func (h *handler) addSong(c *gin.Context) {
	h.logger.Info("querying to add song")
	var AddSongRequest addSongRequest
	var song storer.Song
	var newVerse *storer.Verse
	err := c.ShouldBindJSON(&AddSongRequest)
	if err != nil {
		h.logger.Warn("Error binding JSON", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Debug("Data from JSON:", zap.String("group_name", AddSongRequest.GroupName), zap.String("song", AddSongRequest.Song))
	var songDetails songDetail
	songDetails, err = fetchSongDetails(AddSongRequest.GroupName, AddSongRequest.Song)
	if err != nil {
		h.logger.Warn("Error fetching song details", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	song.ReleaseDate = songDetails.ReleaseDate
	song.Link = songDetails.Link
	song.GroupName = AddSongRequest.GroupName
	song.Song = AddSongRequest.Song

	newSong, err := h.service.AddSong(h.ctx, &song)
	if err != nil {
		h.logger.Error("Error adding song", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	verse := storer.Verse{
		SongID:       newSong.ID,
		VerseTextEng: songDetails.VersesEn,
		VerseTextRu:  songDetails.VersesRu,
	}

	newVerse, err = h.service.AddVerseBySongID(h.ctx, newSong.ID, &verse)
	if err != nil {
		h.logger.Error("Error adding verse", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	newVerse.SongID = newSong.ID
	h.logger.Info("Successfully added song and verses", zap.Int("song_id", newSong.ID))
	c.JSON(201, gin.H{"song": newSong, "verse_ID": newVerse.ID})

}

// @Summary Get all songs
// @Description Retrieves a list of all songs with optional filtering, limit, and offset
// @Tags songs
// @Accept json
// @Produce json
// @Param filter query string false "Filter songs by group name or title"
// @Param limit query int false "Number of songs to return" default(5)
// @Param offset query int false "Number of songs to skip" default(0)
// @Success 200 {object} []storer.Song
// @Failure 400 {object} string "Invalid query parameters"
// @Failure 500 {object} string  "Internal server error"
// @Router /songs [get]
func (h *handler) getAllSongs(c *gin.Context) {
	h.logger.Info("querying to get all songs")
	filter := c.Query("filter")
	limitStr := c.DefaultQuery("limit", "5")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.Warn("Invalid limit", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid limit"})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.Warn("Invalid offset", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid offset"})
	}
	h.logger.Info("Query parameters:", zap.String("filter", filter), zap.Int("limit", limit), zap.Int("offset", offset))

	songs, err := h.service.GetAllSongs(h.ctx, limit, offset, filter)
	if err != nil {
		h.logger.Error("Error getting songs", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully getting songs", zap.Int("count", len(songs)))
	c.JSON(200, gin.H{
		"songs": songs,
	})
}

// @Summary Get a song by ID
// @Description Retrieves a song by its unique identifier
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} storer.Song
// @Failure 400 {object} string  "Invalid ID"
// @Failure 404 {object} string  "Song not found"
// @Failure 500 {object} string  "Internal server error"
// @Router /song/{id} [get]
func (h *handler) getSongByID(c *gin.Context) {
	h.logger.Info("querying to get song by ID")
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	h.logger.Info("Query parameters:", zap.Int("id", id))
	song, err := h.service.GetSongByID(h.ctx, id)
	if err != nil {
		h.logger.Error("Error getting song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully getting song by ID", zap.Int("id", id))
	c.JSON(200, gin.H{"song": song})
}

// @Summary Delete a song by ID
// @Description Deletes a song using its unique identifier
// @Tags songs
// @Param id path int true "Song ID"
// @Success 204 "Successfully deleted song"
// @Failure 400 {object} string "Invalid ID"
// @Failure 404 {object} string "Song not found"
// @Failure 500 {object} string "Internal server error"
// @Router /song/{id} [delete]
func (h *handler) deleteSongByID(c *gin.Context) {
	h.logger.Info("querying to delete song by ID")
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	h.logger.Info("Query parameters:", zap.Int("id", id))
	err = h.service.DeleteSongByID(h.ctx, id)
	if err != nil {
		h.logger.Error("Error deleting song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Successfully deleted song by ID", zap.Int("id", id))
	c.JSON(204, nil)
}

// @Summary Update a song by ID
// @Description Updates a song's details using its unique identifier
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body updateSongRequest true "Song details"
// @Success 200 {object} storer.Song
// @Failure 400 {object} string "Invalid ID or request body"
// @Failure 404 {object} string "Song not found"
// @Failure 500 {object} string "Internal server error"
// @Router /song/{id} [patch]
func (h *handler) updateSongByID(c *gin.Context) {
	h.logger.Info("querying to update song by ID")
	var updateSongRequest updateSongRequest
	//var song storer.Song
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.ShouldBindJSON(&updateSongRequest)
	if err != nil {
		h.logger.Warn("Error binding JSON", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Query parameters:", zap.Int("id", id))

	s, err := h.service.GetSongByID(h.ctx, id)
	if err != nil {
		h.logger.Error("Error getting song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	s = patchSong(updateSongRequest, *s)

	updateSong, err := h.service.UpdateSongByID(h.ctx, s)
	if err != nil {
		h.logger.Error("Error updating song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Successfully updated song by ID", zap.Int("id", id))
	c.JSON(200, gin.H{"song": updateSong})
}

// @Summary Get verses by song ID
// @Description Retrieves verses for a specific song using its unique identifier
// @Tags verses
// @Param id path int true "Song ID"
// @Param lang query string false "Language (en or ru)"
// @Param limit query int false "Number of verses to return" default(10)
// @Param offset query int false "Number of verses to skip" default(0)
// @Success 200 {object} storer.Verse
// @Failure 400 {object} string "Invalid parameters"
// @Failure 404 {object} string "Song not found"
// @Failure 500 {object} string "Internal server error"
// @Router /song/{id}/verse [get]
func (h *handler) getVersesBySongID(c *gin.Context) {
	log.Println("querying to get verses by song ID")
	idStr := c.Param("id")
	lang := c.DefaultQuery("lang", "")
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.logger.Warn("Invalid limit", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid limit"})
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.logger.Warn("Invalid offset", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid offset"})
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}
	song, err := h.service.GetSongByID(h.ctx, id)
	if err != nil {
		h.logger.Error("Error getting song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Query parameters:", zap.Int("id", id), zap.Int("limit", limit), zap.Int("offset", offset))
	verses, err := h.service.GetVersesBySongID(h.ctx, id)
	if err != nil {
		h.logger.Error("Error getting verses by song ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	linesEn := strings.Split(verses.VerseTextEng, "\n")
	linesRu := strings.Split(verses.VerseTextRu, "\n")
	start := offset
	end := offset + limit
	if start >= len(linesEn) {
		c.JSON(200, gin.H{"verses": []string{}})
		return
	}
	if end >= len(linesEn) {
		end = len(linesEn)
	}
	info := struct {
		GroupName string `json:"group_name"`
		Song      string `json:"song"`
		ID        int    `json:"id"`
	}{
		ID:        id,
		GroupName: song.GroupName,
		Song:      song.Song,
	}
	var verseTextsEn, verseTextsRu []string

	switch lang {
	case "en":
		verseTextsEn = linesEn[start:end]
	case "ru":
		verseTextsRu = linesRu[start:end]
	default:
		verseTextsEn = linesEn[start:end]
		verseTextsRu = linesRu[start:end]
	}

	response := gin.H{"info": info}
	if len(verseTextsEn) > 0 {
		response["verses_en"] = verseTextsEn
	}
	if len(verseTextsRu) > 0 {
		response["verses_ru"] = verseTextsRu
	}

	h.logger.Info("Successfully getting verses by song ID", zap.Int("id", id))
	c.JSON(200, response)
}

// @Summary Delete verse by song ID
// @Description Deletes a verse associated with a specific song using its unique identifier
// @Tags verses
// @Param id path int true "Song ID"
// @Success 204 {object} nil "No content"
// @Failure 400 {object} string "Invalid song ID"
// @Failure 404 {object} string "Verse not found"
// @Failure 500 {object} string "Internal server error"
// @Router /song/{id}/verse [delete]
func (h *handler) deleteVerseByID(c *gin.Context) {
	h.logger.Info("querying to delete verse by ID")
	songIDStr := c.Param("id")

	songID, err := strconv.Atoi(songIDStr)
	if err != nil {
		h.logger.Warn("Invalid song ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid song ID"})
		return
	}
	h.logger.Info("Query parameters:", zap.Int("song_id", songID))
	err = h.service.DeleteVerseBySongID(h.ctx, songID)
	if err != nil {
		h.logger.Error("Error deleting verse by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
	}
	h.logger.Info("Successfully deleted verse by songID", zap.Int("song_id", songID))
	c.JSON(204, nil)
}

// @Summary Update verse by song ID
// @Description Updates a verse associated with a specific song using its unique identifier
// @Tags verses
// @Param id path int true "Song ID"
// @Param verse body updateVerseRequest true "Verse data to update"
// @Success 200 {object} storer.Verse
// @Failure 400 {object} string "Invalid parameters"
// @Failure 404 {object} string "Verse not found"
// @Failure 500 {object} string "Internal server error"
// @Router /song/{id}/verse [patch]
func (h *handler) updateVerseBySongID(c *gin.Context) {
	h.logger.Info("querying to update verse by ID")
	var updateVerseRequest updateVerseRequest
	//var verse storer.Verse
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Warn("Invalid ID", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.ShouldBindJSON(&updateVerseRequest)
	if err != nil {
		h.logger.Warn("Error binding JSON", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Query parameters:", zap.Int("id", id))

	v, err := h.service.GetVersesBySongID(h.ctx, id)
	if err != nil {
		h.logger.Error("Error getting song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	v = patchVerse(updateVerseRequest, *v)

	updateVerse, err := h.service.UpdateVerseBySongID(h.ctx, id, v)
	updateVerse.ID = v.ID
	updateVerse.SongID = id
	if err != nil {
		h.logger.Error("Error updating song by ID", zap.Error(err))
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	h.logger.Info("Successfully updated song by ID", zap.Int("id", id))
	c.JSON(200, gin.H{"song": updateVerse})

}
