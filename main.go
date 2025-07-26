package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/go-chi/chi/v5"
)

type User struct {
	ID    int64  `datastore:"-"`
	Name  string
	Email string
	Age   int
}

var (
	projectID     = "aero3-467012" // Remplace par ton ID GCP
	userListTmpl  *template.Template
	editFormTmpl  *template.Template
	dsClient      *datastore.Client
)

func main() {
	ctx := context.Background()

	var err error
	dsClient, err = datastore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Erreur client Datastore: %v", err)
	}

	// Templates
	userListTmpl = template.Must(template.New("list").Parse(userListHTML))
	editFormTmpl = template.Must(template.New("edit").Parse(editFormHTML))

	r := chi.NewRouter()

	r.Get("/", listUsersHandler)
	r.Get("/init", initDataHandler)
	r.Get("/edit/{id}", editUserFormHandler)
	r.Post("/edit/{id}", editUserHandler)
	r.Get("/delete/{id}", deleteUserHandler)

	fmt.Println("Serveur sur http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

// === HANDLERS ===

func initDataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	for i := 1; i <= 10; i++ {
		u := &User{
			Name:  fmt.Sprintf("User%d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
			Age:   20 + i,
		}
		key := datastore.IncompleteKey("User", nil)
		if _, err := dsClient.Put(ctx, key, u); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var users []User
	query := datastore.NewQuery("User")
	keys, err := dsClient.GetAll(ctx, query, &users)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	for i, key := range keys {
		users[i].ID = key.ID
	}
	userListTmpl.Execute(w, users)
}

func editUserFormHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	key := datastore.IDKey("User", id, nil)
	var u User
	if err := dsClient.Get(ctx, key, &u); err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	u.ID = id
	editFormTmpl.Execute(w, u)
}

func editUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	key := datastore.IDKey("User", id, nil)

	name := r.FormValue("name")
	email := r.FormValue("email")
	age, _ := strconv.Atoi(r.FormValue("age"))

	u := &User{
		Name:  name,
		Email: email,
		Age:   age,
	}
	if _, err := dsClient.Put(ctx, key, u); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	key := datastore.IDKey("User", id, nil)
	if err := dsClient.Delete(ctx, key); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// === HTML TEMPLATES ===

var userListHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Liste des utilisateurs</title>
</head>
<body>
	<h1>Utilisateurs</h1>
	<a href="/init">Insérer les 10 utilisateurs</a>
	<table border="1">
		<tr><th>ID</th><th>Nom</th><th>Email</th><th>Age</th><th>Actions</th></tr>
		{{range .}}
		<tr>
			<td>{{.ID}}</td>
			<td>{{.Name}}</td>
			<td>{{.Email}}</td>
			<td>{{.Age}}</td>
			<td>
				<a href="/edit/{{.ID}}">Modifier</a>
				<a href="/delete/{{.ID}}" onclick="return confirm('Supprimer ?')">Supprimer</a>
			</td>
		</tr>
		{{end}}
	</table>
</body>
</html>
`

var editFormHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Modifier Utilisateur</title>
</head>
<body>
	<h1>Modifier l'utilisateur</h1>
	<form method="POST">
		<label>Nom: <input name="name" value="{{.Name}}"></label><br>
		<label>Email: <input name="email" value="{{.Email}}"></label><br>
		<label>Age: <input name="age" value="{{.Age}}"></label><br>
		<button type="submit">Enregistrer</button>
	</form>
	<a href="/">Retour</a>
</body>
</html>
`
