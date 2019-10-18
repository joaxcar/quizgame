package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

type player struct {
	UserID string
	solvedProbs []problem
}

type webGame struct{
	Db *sql.DB
	DoneID []int
	CurID int
	ProbSet []int
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
func (g *webGame) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ans := r.FormValue("answer")
	prob := getProblem(g.Db, g.CurID)
	if ans == prob.Answer {
		g.DoneID = append(g.DoneID, g.CurID)
		g.CurID++
		if g.CurID >= len(g.ProbSet) {
			http.Redirect(w,r,"/result",http.StatusFound)
		}
		prob = getProblem(g.Db, g.ProbSet[g.CurID])
	}
	t,_ := template.ParseFiles("site.html")
	err := t.Execute(w, prob)
	if err != nil {
		fmt.Println("Fel med sidan")
	}

}
func initiateServer(db *sql.DB) {
	wg := webGame{Db:db}
	q := `SELECT id FROM quiz`
	res, err := db.Query(q)
	if err != nil {
		fmt.Println("Database problems while selectiong all probs")
	}
	for res.Next() {
		var probID int
		res.Scan(&probID)
		wg.ProbSet = append(wg.ProbSet, probID)
	}
	rand.Shuffle(len(wg.ProbSet), func(i, j int) {wg.ProbSet[i],wg.ProbSet[j] = wg.ProbSet[j],wg.ProbSet[i]})
	fmt.Println(len(wg.ProbSet))
	wg.CurID = 0
	http.Handle("/game", &wg)
	_ = http.ListenAndServe(":8881", nil)
	fmt.Println("hej")
}
