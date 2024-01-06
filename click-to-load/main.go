package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

type User struct {
	ID int
	Name string
	Email string
}

func generateUser(page int) []User {
	users := []User{}
	for i := (page * 5) - 5; i < page * 5; i++ {
		users = append(users, User{
			ID: i + 1,
			Name: fmt.Sprintf("Iqbal %d", i + 1),
			Email: fmt.Sprintf("iqbal%d@email.com", i + 1),
		})
	}

	return users
}

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", templ.Handler(root(showTable(generateUser(1), 2, true))))

	http.HandleFunc("/load-more/", func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/load-more/")); if err != nil {
			root(showError(err.Error(), http.StatusText(500)))
			return
		}
		
		users := generateUser(page)
		showRows(&users, page + 1, page <= 5).Render(r.Context(), w)
	})

	http.ListenAndServe(":3000", nil)
}