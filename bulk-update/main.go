package main

import (
	"net/http"
	"strconv"

	"github.com/a-h/templ"
)

type User struct {
	Name string
	IsActive bool
}

type UserWithIndicator struct {
	Name string
	IsActive bool
	IsChanged bool
}

func main() {
	users := []User{
		{
			Name: "Iqbal",
			IsActive: true,
		},
		{
			Name: "Maulana",
			IsActive: true,
		},
		{
			Name: "Balqi",
			IsActive: false,
		},
		{
			Name: "qiBal",
			IsActive: false,
		},
		{
			Name: "Labqi",
			IsActive: true,
		},
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", templ.Handler(root(userList(&users))))

	http.HandleFunc("/activate", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ids := r.Form["ids"]

		usersWithIndicator := []UserWithIndicator{}
		for _, user := range users {
			usersWithIndicator = append(usersWithIndicator, UserWithIndicator{
				Name: user.Name,
				IsActive: user.IsActive,
				IsChanged: false,
			})
		}

		for _, idx := range ids {
			i, err := strconv.Atoi(idx); if err != nil {
				showError(err.Error(), http.StatusText(500)).Render(r.Context(), w)
				return
			}

			users[i].IsActive = true
			usersWithIndicator[i].IsActive = true
			usersWithIndicator[i].IsChanged = true
		}

		loopUserWithIndicator(&usersWithIndicator).Render(r.Context(), w)

	})

	http.HandleFunc("/deactivate", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ids := r.Form["ids"]

		usersWithIndicator := []UserWithIndicator{}
		for _, user := range users {
			usersWithIndicator = append(usersWithIndicator, UserWithIndicator{
				Name: user.Name,
				IsActive: user.IsActive,
				IsChanged: false,
			})
		}

		for _, idx := range ids {
			i, err := strconv.Atoi(idx); if err != nil {
				showError(err.Error(), http.StatusText(500)).Render(r.Context(), w)
				return
			}

			users[i].IsActive = false
			usersWithIndicator[i].IsActive = false
			usersWithIndicator[i].IsChanged = true
		}

		loopUserWithIndicator(&usersWithIndicator).Render(r.Context(), w)

	})

	http.ListenAndServe(":3000", nil)
}