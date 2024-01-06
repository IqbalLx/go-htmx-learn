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
	FirstName string
	LastName string
	Email string
}

type NewUser struct {
	FirstName string
	LastName string
	Email string
}

func findUser(users [3]User, id int) int {
	for i := 0; i < 3; i++ {
		if (users[i].ID == id) {
			return i
		}
	}

	return -1
}

func getUserById(users [3]User, id int) templ.Component {
	userIndex := findUser(users, id)
	if (userIndex == -1) {
		return showError("user not found", "400")
	}

	return showUser(users[userIndex])
}

func updateUserByIdView(users [3]User, id int) templ.Component {
	userIndex := findUser(users, id)
	if (userIndex == -1) {
		return showError("user not found", "400")
	}

	return updateUser(users[userIndex])
}

func updateUserById(users *[3]User, id int, newUser NewUser) templ.Component {
	userIndex := findUser(*users, id)
	if (userIndex == -1) {
		return showError("user not found", "400")
	}

	(*users)[userIndex] = User{
		ID: id,
		Email: newUser.Email,
		FirstName: newUser.FirstName,
		LastName: newUser.LastName,
	}

	return showUser(users[userIndex])
}

func main() {

	users := [3]User{
		{
			ID: 1,
			FirstName: "Iqbal",
			LastName: "Maulana",
			Email: "iqbal@email.com",
		},
		{
			ID: 2,
			FirstName: "Maulana",
			LastName: "Iqbal",
			Email: "maulana@email.com",
		},
		{
			ID: 3,
			FirstName: "Iqbalz",
			LastName: "Maulanaz",
			Email: "iqbalz@email.com",
		},
	}
	
	http.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		id, error := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/user/")); if error != nil {
			root(showError(error.Error(), "500")).Render(r.Context(), w)
			return
		}

		switch r.Method {
		case "GET":
			root(getUserById(users, id)).Render(r.Context(), w)
		}
	})

	http.HandleFunc("/user/edit/", func(w http.ResponseWriter, r *http.Request) {
		id, error := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/user/edit/")); if error != nil {
			root(showError(error.Error(), "500")).Render(r.Context(), w)
			return
		}

		switch r.Method {
		case "GET":
			root(updateUserByIdView(users, id)).Render(r.Context(), w)
		
		case "PUT":
			newUser := NewUser{
				Email: r.FormValue("email"),
				FirstName: r.FormValue("firstName"),
				LastName: r.FormValue("lastName"),
			}

			root(updateUserById(&users, id, newUser)).Render(r.Context(), w)
		}
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}