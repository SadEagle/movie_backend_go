package handlers

import (
	"context"
	"fmt"
	"io"
	"movie_backend_go/db/sqlc"
	"movie_backend_go/pkg/auth"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// TODO: make proper pathing
// Get this path from volume

const (
	chunkSize     = 10 * 1024 * 1024 // 10MB
	MOVIES_PREFIX = "/movie-data"
)

// NOTE: DownloadMove expect middleware that will handle installation info saving somewhere. We need to know about saving this data

// @Summary     Upload movie
// @Description Upload movie data as []bytes stream
// @Tags        video-manager, admin
// @Accept 		octet-stream
// @Produce     json
// @Security	OAuth2Password
// @Param       movie_id   	path	string 	true  "Movie ID"
// @Param       tequest		body	[]byte 	true  "Streaming Bytes"
// @Success     204
// @Failure     404  {object}  map[string]string
// @Failure     500  {object}  map[string]string
// @Router      /upload/movie/{movie_id} [post]
func (ho *HandlerObj) UploadMovie(rw http.ResponseWriter, r *http.Request) {
	ctx, close := context.WithTimeout(r.Context(), OpTimeContext)
	defer close()

	movieIDStr := r.PathValue("movie_id")
	var movieID pgtype.UUID

	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	userTokenData, err := auth.GetTokenDataContext(ctx)
	if err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Wrong tokend extractor middleware", http.StatusInternalServerError)
	}

	// Verify
	if !userTokenData.IsAdmin {
		ho.Logger.Println("Unauthorized user")
		http.Error(rw, "Unauthorized user", http.StatusUnauthorized)
		return
	}

	moviePath := filepath.Join(MOVIES_PREFIX, movieIDStr+".mp4")

	file, err := os.Open(moviePath)
	if err != nil {
		http.Error(rw, "video not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// TODO: make transaction and write movie_path before operation itself, submit operation after success???
	io.Copy(file, r.Body)

	moviePathAdd := sqlc.AddMoviePathParams{ID: movieID, MoviePath: &moviePath}
	_, err = ho.QuerierDB.AddMoviePath(ctx, moviePathAdd)
	if err != nil {
		ho.Logger.Printf("Can't write downloaded function path into movie table. Delete downloaded file: %v", moviePathAdd)
		err = os.Remove(moviePath)
		if err != nil {
			ho.Logger.Printf("!!CAN't delete downloaded video %v", moviePathAdd)
		}
	}
	rw.WriteHeader(http.StatusNoContent)
}

// @Summary     Stream movie
// @Tags        video-manager
// @Accept      json
// @Produce     video/mp4
// @Param       movie_id 	path	string  true 	"Movie ID"
// @Param 		Range 		header 	string 	false 	"Byte range"
// @Header 		200  	{string} 	Accept-Ranges 	"bytes"
// @Header 		200  	{string} 	Content-Type 	"video/mp4"
// @Header 		200  	{int} 		Content-Lenght 	200
// @Header 		200  	{string} 	Content-Range 	"bytes 1024-10112"
// @Success 	200  	{object} 	[]byte
// @Failure 	404  	{object} 	map[string]string
// @Failure 	500  	{object} 	map[string]string
// @Router     /stream/movie/{movie_id} [get]
func (ho *HandlerObj) StreamMovie(rw http.ResponseWriter, r *http.Request) {
	movieIDStr := r.PathValue("movie_id")
	var movieID pgtype.UUID

	if err := movieID.Scan(r.PathValue("movie_id")); err != nil {
		ho.Logger.Println(err)
		http.Error(rw, "Requested movie id should contain uuid style", http.StatusBadRequest)
		return
	}

	moviePath := filepath.Join(MOVIES_PREFIX, movieIDStr+".mp4")

	file, err := os.Open(moviePath)
	if err != nil {
		http.Error(rw, "video not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(rw, "cannot stat file", http.StatusInternalServerError)
		return
	}
	size := stat.Size()

	rw.Header().Set("Accept-Ranges", "bytes")
	rw.Header().Set("Content-Type", "video/mp4")

	rangeHdr := r.Header.Get("Range")
	if rangeHdr == "" {
		// TODO: handle videos >30s. in my case they are videos with 30+minutes
		// No Range header: serve the whole file (200 OK)
		rw.Header().Set("Content-Length", strconv.FormatInt(size, 10))
		if _, err := io.Copy(rw, file); err != nil {
			ho.Logger.Printf("copy full file: %v", err)
		}
		return
	}

	// Expect formats like: "bytes=START-" or "bytes=START-END"
	if !strings.HasPrefix(strings.ToLower(rangeHdr), "bytes=") {
		http.Error(rw, "invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}
	rangeSpec := strings.TrimPrefix(rangeHdr, "bytes=")
	parts := strings.Split(rangeSpec, "-")
	if len(parts) != 2 || parts[0] == "" {
		http.Error(rw, "invalid range", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	start, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || start < 0 {
		http.Error(rw, "invalid range start", http.StatusRequestedRangeNotSatisfiable)
		return
	}

	// Compute end (inclusive)
	var end int64
	if parts[1] == "" {
		// No end provided: mimic the Python version's behavior (chunked)
		end = start + chunkSize - 1
	} else {
		end, err = strconv.ParseInt(parts[1], 10, 64)
		if err != nil || end < start {
			http.Error(rw, "invalid range end", http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	// Clamp to file size - 1
	if start >= size {
		// Invalid range: start beyond EOF
		rw.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
		http.Error(rw, "range not satisfiable", http.StatusRequestedRangeNotSatisfiable)
		return
	}
	if end >= size {
		end = size - 1
	}

	length := end - start + 1

	// Prepare headers for 206 Partial Content
	rw.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, size))
	rw.Header().Set("Content-Length", strconv.FormatInt(length, 10))
	rw.WriteHeader(http.StatusPartialContent)

	// Seek and write exactly the requested bytes
	if _, err := file.Seek(start, io.SeekStart); err != nil {
		ho.Logger.Printf("seek error: %v", err)
		return
	}
	if _, err := io.CopyN(rw, file, length); err != nil {
		// Client may cancel early; just log
		ho.Logger.Printf("copyN error: %v", err)
	}
}
