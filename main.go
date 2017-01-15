package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

// XXX remove; not necessary anymore; has disintegrated into getters and setters
type Question interface {
	Qno()              int
	String()           string
	Radio()            bool
	Choices()          []Choice
	Add(answer, voter string)
	Remove(voter string)
}

type question struct {
	no      int
	text    string
	radio   bool
	voters  []string
	choices []Choice         // no predefined choices for text questions
}

func (q *question) Qno() int          { return q.no }
func (q *question) String() string    { return q.text }
func (q *question) Radio() bool       { return q.radio }
func (q *question) Choices() []Choice { return q.choices }

func (q *question) Add(answer, voter string) {
	// If a user answered a question once, they may change it; first remove
	// the current answer, though.
	for _, v := range q.voters {
		if v == voter {
			q.Remove(voter)
		}
	}

	for i, ch := range q.choices {
		// The identifier answer, used as value in radio button decls, is the only
		// one relevant here. The normal Answer is solely for display.
		if ch.IdentAnswer() == answer {
			q.voters = append(q.voters, voter)
			q.choices[i].Voters = append(q.choices[i].Voters, voter)
			return
		}
	}

	if !q.radio {
		q.voters = append(q.voters, voter)
		q.choices = append(q.choices, Choice{answer, []string{voter}, 0.0})
	}
}

func (q *question) Remove(voter string) {
	for i := range q.choices {
		for j, v := range q.choices[i].Voters {
			if v == voter {
				q.choices[i].Voters = append(q.choices[i].Voters[:j], q.choices[i].Voters[j+1:]...)

				// Remove text question choice entirely if none adhere to it.
				if !q.radio && len(q.choices[i].Voters) == 0 {
					q.choices = append(q.choices[:i], q.choices[i+1:]...)
				}

				return
			}
		}
	}
}

// XXX temp. for testing; implement sessions
var user = "oki"
var loggedIn = false

func rootHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/q/1", http.StatusFound)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		user = r.Form["username"][0]
		loggedIn = true
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// qnoFromPath: "/q/565" → 565, and so on.
func qnoFromPath(path string) int {
	a := strings.FieldsFunc(path, func(r rune) bool {return r == '/'})
	qno, _ := strconv.Atoi(a[len(a)-1])  // get last segment, ignore error XXX
	return qno
}

func questionHandler(w http.ResponseWriter, r *http.Request) {
	qno := qnoFromPath(r.URL.Path)
	if qno <= 0 || qno >= len(questions) {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		var t *template.Template

		if questions[qno].Radio() {
			t, _ = template.ParseFiles("radio.gtpl")
		} else {
			t, _ = template.ParseFiles("text.gtpl")
		}

		t.Execute(w, questions[qno])
	} else {
		r.ParseForm()

		questions[qno].Add(r.Form["answer"][0], user)

		if qno+1 < len(questions) {
			http.Redirect(w, r, fmt.Sprint("/q/", qno+1), http.StatusFound)
		} else {
			fmt.Fprintln(w, "Danke für deine Teilnahme!")
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
		choices[i].Percentage = 100.0 * float64(len(ch.Voters)) / float64(total)
	}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("stats.gtpl")

	for i := range questions {
		sortAndCalcPercentage(questions[i].Choices())
	}

	t.Execute(w, questions[1:])
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/q/", questionHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
