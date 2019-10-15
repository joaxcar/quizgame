package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
)

type webGame struct{
	Db *sql.DB
	DoneID []string
	CurID int
}

type problem struct {
	ID int
	Question string
	Answer string
	Done bool
}

func getProblem(db *sql.DB, ID int) problem{
	q := `SELECT id, question, answer, done FROM quiz WHERE id=$1`
	res, err := db.Query(q,ID)
	if err != nil {
		panic(err)
	}
	defer res.Close()

	prob := problem{}
	if res.Next() {
		err = res.Scan(&prob.ID, &prob.Question, &prob.Answer, &prob.Done)
		if err != nil {
			fmt.Println("Error")
		}
	}

	return prob
}
func (g webGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ans := r.FormValue("answer")
	prob := getProblem(g.Db, g.CurID)
	fmt.Println(g.CurID)
	if ans == prob.Answer {
		g.CurID++
		prob = getProblem(g.Db, g.CurID)
		fmt.Println(g.CurID)
	}
	t,_ := template.ParseFiles("site.html")
	err := t.Execute(w, prob)
	if err != nil {
		fmt.Println("Fel med sidan")
	}

}
func initiateServer(db *sql.DB) {
	wg := webGame{Db:db, CurID:1}
	http.Handle("/", wg)
	http.ListenAndServe(":8888", nil)
}
