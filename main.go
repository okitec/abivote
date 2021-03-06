/*
Specialised Go HTTP server for taking surveys.
*/
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
)

type Choice struct {
	Answer     string
	Voters     []string
	Percentage float64
}

// IdentAnswer: replace space with hyphens, remove punctuation, cast lowercase to create
// a very probably valid HTML5 identifier.
func (ch Choice) IdentAnswer() string {
	s := strings.Replace(ch.Answer, " ", "-", -1)
	s = strings.Replace(ch.Answer, ".", "", -1)
	s = strings.Replace(ch.Answer, ":", "", -1)
	s = strings.Replace(ch.Answer, "?", "", -1)
	s = strings.Replace(ch.Answer, "!", "", -1)
	s = strings.ToLower(s)
	return s
}

type ByVotes []Choice

func (v ByVotes) Len() int           { return len(v) }
func (v ByVotes) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v ByVotes) Less(i, j int) bool { return len(v[i].Voters) < len(v[j].Voters) }

type question struct {
	Qno     int      // question number; URL path is '/q/no'
	text    string   // question text
	Radio   bool     // is radio question?
	voters  []string // ids of people who voted for this
	Choices []Choice // choices; no predefined choices for text questions
	User    string   // hack: current user that must be passed to template
}

func (q *question) String() string { return q.text }

// Add the answer a voter chose. For text questions, this is the
// text that was entered; for radio questions, this is the value
// of the chosen radio button. That value is the full displayed
// answer without spaces, punctuation and the like.
func (q *question) Add(answer, voter string) {
	// If a user answered a question once, they may change it; first remove
	// the current answer, though.
	for _, v := range q.voters {
		if v == voter {
			q.Remove(voter)
		}
	}

	for i, ch := range q.Choices {
		// The identifier answer, used as value in radio button decls, is the only
		// one relevant here. The normal Answer is solely for display.
		if ch.IdentAnswer() == answer {
			q.voters = append(q.voters, voter)
			q.Choices[i].Voters = append(q.Choices[i].Voters, voter)
			return
		}
	}

	if !q.Radio {
		q.voters = append(q.voters, voter)
		q.Choices = append(q.Choices, Choice{answer, []string{voter}, 0.0})
	}
}

// Remove the answer the voter has given.
func (q *question) Remove(voter string) {
	for i := range q.Choices {
		for j, v := range q.Choices[i].Voters {
			if v == voter {
				q.Choices[i].Voters = append(q.Choices[i].Voters[:j], q.Choices[i].Voters[j+1:]...)

				// Remove text question choice entirely if none adhere to it.
				if !q.Radio && len(q.Choices[i].Voters) == 0 {
					q.Choices = append(q.Choices[:i], q.Choices[i+1:]...)
				}

				return
			}
		}
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := r.Form.Get("user")
	if user == "" {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/q/1?user="+user, http.StatusFound)
}

type userinfo struct {
	Name     string
	HasVoted bool
	Admin    bool
}

// filled from "users.json" file
var users map[string]userinfo

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		user := r.PostForm.Get("username")
		if ok, err := CanVote(user); !ok {
			http.Error(w, err, http.StatusForbidden) // XXX proper error
			return
		}

		log.Printf("User %s logged in", user)
		http.Redirect(w, r, "/q/1?user="+user, http.StatusFound)
	}
}

func CanVote(user string) (ok bool, err string) {
	if user == "" {
		return false, "no name"
	}

	if _, ok := users[user]; !ok {
		return false, "Diese Kennung ist nicht registriert"
	}

	if users[user].HasVoted {
		return false, "Du hast bereits gewählt"
	}

	return true, ""
}

// qnoFromPath: "/q/565" → 565, and so on.
func qnoFromPath(path string) int {
	a := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
	qno, _ := strconv.Atoi(a[len(a)-1]) // get last segment, ignore error XXX
	return qno
}

// Replica of questions, but the choices are never reordered.
var pristineQuestions []*question;

func init() {
	for i, q := range questions {
		q2 := &question{i, q.text, q.Radio, nil, nil, ""}

		ch2 := make([]Choice, len(q.Choices))
		copy(ch2, q.Choices)
		q2.Choices = ch2

		pristineQuestions = append(pristineQuestions, q2)
		fmt.Println(q2)
	}
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := r.Form.Get("user")
	if ok, _ := CanVote(user); !ok {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	qno := qnoFromPath(r.URL.Path)
	if qno <= 0 || qno >= len(questions) {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		var t *template.Template

		if questions[qno].Radio {
			t, _ = template.ParseFiles("radio.html")
		} else {
			t, _ = template.ParseFiles("text.html")
		}

		pristineQuestions[qno].User = user
		t.Execute(w, pristineQuestions[qno])
		pristineQuestions[qno].User = ""
	} else {
		// Don't complain when the user skipped a question and left no answer.
		if r.Form.Get("answer") != "" {
			questions[qno].Add(r.Form["answer"][0], user)
		}

		saveResults()
		if qno+1 < len(questions) {
			http.Redirect(w, r, fmt.Sprintf("/q/%d?user=%s", qno+1, user), http.StatusFound)
		} else {
			fmt.Fprintln(w, "<p>Danke für deine Teilnahme! Nur Geduld, die Ergebnisse findest du dann in der der Abizeitung.</p>")
			// can't assign to struct in map because map store might change position
			// invalid: users[user].HasVoted = true
			u := users[user]
			u.HasVoted = true
			users[user] = u
			log.Printf("User %s completed the questionnaire", user)
			saveUsers()
		}
	}
}

func sortAndCalcPercentage(choices []Choice) {
	sort.Sort(sort.Reverse(ByVotes(choices)))

	total := 0
	for _, ch := range choices {
		total += len(ch.Voters)
	}

	for i, ch := range choices {
		if total > 0 {
			choices[i].Percentage = 100.0 * float64(len(ch.Voters)) / float64(total)
		} else {
			choices[i].Percentage = 0.0
		}
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	user := r.Form.Get("user")
	if _, ok := users[user]; !ok || !users[user].Admin {
		http.Error(w, "Kein Zugang zu den Statistiken für Nicht-Admins", http.StatusForbidden) // XXX proper error
		return
	}

	t, _ := template.ParseFiles("stats.html")

	for i := range questions {
		sortAndCalcPercentage(questions[i].Choices)
	}

	t.Execute(w, questions[1:])
}

func saveUsers() {
	b, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return
	}

	ioutil.WriteFile("users.json", b, 0644)
	log.Println("Saved users.json")
}

func saveResults() {
	b, err := json.Marshal(questions)
	if err != nil {
		log.Println(err)
		return
	}

	ioutil.WriteFile("results.json", b, 0644)
	log.Println("Saved results.json")
}

func main() {
	b, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(b, &users)
	if err != nil {
		log.Fatal("users.json: ", err)
	}
	log.Printf("Loaded users.json. %d users", len(users))
	defer saveUsers()

	// Only try to decode if the file can be read.
	b, err = ioutil.ReadFile("results.json")
	if err == nil {
		err = json.Unmarshal(b, &questions)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Loaded results.json")
	} else {
		log.Println(err)
	}
	defer saveResults()

	// Also save users and results on SIGINT
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	go func() {
		for _ = range sc {
			saveUsers()
			saveResults()
			os.Exit(0)
		}
	}()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/q/", questionHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
