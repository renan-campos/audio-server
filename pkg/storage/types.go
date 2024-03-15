package storage

type MovieStorageService interface {
	CreateEntry(movieName string) error
	ListMovies() (MovieList, error)
}

type MovieData struct {
	Name string
}

type MovieList struct {
	Total int
	Items []string
}
