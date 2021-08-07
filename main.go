package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Page struct {
	Title string
	Body  []byte
}

type PassReq struct {
	Length  int
	Number  bool
	Special bool
	Upper   bool
}

var lowerCase = []rune("qwertyuiopasdfghjklzxcvbnm")
var upperCase = []rune("QWERTYUIOPASDFGHJKLZXCVBNM")
var number = []rune("1234567890")
var special = []rune("`~!@#$%^&*()_-+={}[];:/?,<>'")

func generateLower() rune {
	rand.Seed(time.Now().UnixNano())

	v := lowerCase[rand.Intn(len(lowerCase))]
	return v
}

func generateUpper() rune {
	rand.Seed(time.Now().UnixNano())

	v := upperCase[rand.Intn(len(upperCase))]
	return v
}

func generateSpecial() rune {
	v := special[rand.Intn(len(special))]
	return v
}

func generateNumber() rune {
	v := number[rand.Intn(len(number))]
	return v
}

func Generate(upper, number, special bool, characterCount int) string {
	v := make([]rune, characterCount)
	for i := range v {
		rand.Seed(time.Now().UnixNano())
		t := rand.Intn(2)
		tempCount := rand.Intn(3)
		if t == 1 {
			if upper == true && tempCount <= 0 {
				v[i] = generateUpper()
			}
			if number == true && tempCount <= 1 {
				v[i] = generateNumber()
			}
			if special == true && tempCount <= 2 {
				v[i] = generateSpecial()
			}
		} else if t == 0 {
			v[i] = generateLower()
		}
	}
	return string(v)
}

func main() {
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
	fmt.Println(Generate(true, false, true, 30))
}

func loadPage(title string) *Page {
	filename := "tamplates/" + title + ".html"
	body, _ := ioutil.ReadFile(filename)
	return &Page{Title: title, Body: body}
}

func convertCheckbox(value string) bool {
	if value == "on" {
		return true
	}
	return false
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	// title := "Random Password Generator"
	// p := loadPage(title)
	// p = &Page{Title: title}
	t, _ := template.ParseFiles("templates/index.html")

	if r.Method != http.MethodPost {
		t.Execute(w, nil)
		return
	}
	l, _ := strconv.Atoi(r.FormValue("length"))
	details := PassReq{
		Length:  l,
		Number:  convertCheckbox(r.FormValue("number")),
		Special: convertCheckbox(r.FormValue("special")),
		Upper:   convertCheckbox(r.FormValue("upper")),
	}

	// do something with details
	fmt.Printf("%#v", details)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.Execute(w, struct {
		Success  bool
		Password string
	}{true, Generate(details.Upper, details.Number, details.Special, details.Length)})

}
