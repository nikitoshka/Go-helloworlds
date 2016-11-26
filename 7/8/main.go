package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const localURL = "localhost:8000"
const htmlTempl = `<html><title>Tracks table</title>
<table>
    <tr style="text-align: center">
        <th><a href='http://{{.LocalURL}}/sort?field=title'>Title</a></th>
        <th><a href='http://{{.LocalURL}}/sort?field=artist'>Artist</a></th>
        <th><a href='http://{{.LocalURL}}/sort?field=album'>Album</a></th>
        <th><a href='http://{{.LocalURL}}/sort?field=year'>Year</a></th>
        <th><a href='http://{{.LocalURL}}/sort?field=length'>Length</a></th>
    </tr>
    {{range .Tracks}}
    <tr>
        <td>{{.Title}}</td>
        <td>{{.Artist}}</td>
        <td>{{.Album}}</td>
        <td>{{.Year}}</td>
        <td>{{.Length}}</td>
    </tr>
    {{end}}
</table>
</html>`

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

type TrackSort struct {
	tracks []*Track
	less   func(i, j *Track) bool
}

func (t TrackSort) Len() int {
	return len(t.tracks)
}

func (t TrackSort) Swap(i, j int) {
	t.tracks[i], t.tracks[j] = t.tracks[j], t.tracks[i]
}

func (t TrackSort) Less(i, j int) bool {
	return t.less(t.tracks[i], t.tracks[j])
}

func duration(s string) time.Duration {
	d, err := time.ParseDuration(s)

	if err != nil {
		return time.Duration(0)
	}

	return d
}

func sortByYear(i, j *Track) bool {
	return i.Year < j.Year
}

func sortByTitle(i, j *Track) bool {
	return strings.ToLower(i.Title) < strings.ToLower(j.Title)
}

func sortByArtist(i, j *Track) bool {
	return strings.ToLower(i.Artist) < strings.ToLower(j.Artist)
}

func sortByAlbum(i, j *Track) bool {
	return strings.ToLower(i.Album) < strings.ToLower(j.Album)
}

func sortByLength(i, j *Track) bool {
	return i.Length < j.Length
}

var tracks = []*Track{
	{"Still into you", "Paramore", "Paramore", 2013, duration("3m12s")},
	{"Take no prisoners", "Megadeth", "Peace sells", 1989, duration("4m10s")},
	{"Hope", "Magenta", "Revolution", 2007, duration("10m6s")},
	{"Be right there", "Dipplo", "Album", 2016, duration("2m59s")},
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	printResponse(tracks, w)
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("field")
	f := (func(i, j *Track) bool)(nil)

	switch key {
	case "title":
		f = sortByTitle
	case "artist":
		f = sortByArtist
	case "album":
		f = sortByAlbum
	case "year":
		f = sortByYear
	case "length":
		f = sortByLength
	}

	if f == nil {
		http.Error(w, "wrong sorting field", http.StatusBadRequest)
		return
	}

	sorting := TrackSort{tracks, f}

	if sort.IsSorted(sorting) {
		sort.Sort(sort.Reverse(sorting))
	} else {
		sort.Sort(sorting)
	}

	printResponse(sorting.tracks, w)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")
	artist := r.URL.Query().Get("artist")
	album := r.URL.Query().Get("album")

	year, err := strconv.Atoi(r.URL.Query().Get("year"))
	if err != nil {
		year = 0
		return
	}

	length := duration(r.URL.Query().Get("length"))

	tracks = append(tracks, &Track{title, artist, album, year, length})
	printResponse(tracks, w)
}

func printResponse(tracks []*Track, w http.ResponseWriter) {
	t := template.Must(template.New("tracks").Parse(htmlTempl))
	htmlStruct := struct {
		LocalURL string
		Tracks   []*Track
	}{localURL, tracks}

	if err := t.Execute(w, htmlStruct); err != nil {
		fmt.Fprintf(w, "template executin error: %v", err)
		return
	}
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/sort", sortHandler)
	http.HandleFunc("/add", addHandler)

	log.Fatal(http.ListenAndServe(localURL, nil))
}
