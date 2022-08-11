package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	date := r.URL.Query().Get("date")
	out := ""
	if name != "" && date != "" {
		out = schedule(name, date)
	}
	fmt.Fprintf(w, `SST MTR Schedule Fetcher: Enter
        your name to get a list of all matching entries.
        For exact matches, enter as seen in class register.`+"<br><br>"+`
        <form action="" method="get">
                <input type="text" name="name">
                <select name="date" id="date">
                  <option value="Today">Today</option>
                  <option value="11 Aug">11 Aug</option>
                  <option value="12 Aug">12 Aug</option>
                </select>
                <input type="submit" value="Enter">
              </form>`+out+`<br>Only schedules for 11 Aug and 12 Aug are implemented.<br> If you have any problems or find any bugs,
my discord is awpgikxdigj#8231<br><br>Made by Ethan Tse Chun Lam, S407 (objectively better computing class)`)
}

func schedule(name string, date string) string {
	directory := os.Getcwd() + "/data"
	f, err := os.Open(directory + "/datelist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	dates := strings.Split(string(b), "\n")
	f, err = os.Open(directory + "/namelist.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	b, err = ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	fullnames := strings.Split(string(b), "\n")
	if date == "Today" {
		temp := time.Now().Format("02 Jan")
		if temp != "11 Aug" && temp != "12 Aug" {
			return "No schedule for Aug " + temp
		} else {
			date = "/" + temp
		}
	} else {
		date = "/" + date
	}
	name = strings.ToLower(name)
	picknames := []string{}
	for _, b := range fullnames {
		if strings.Contains(b, name) {
			picknames = append(picknames, strings.TrimSpace(b))
		}
	}
	times, err := ioutil.ReadDir(directory + "/dates" + date)
	if err != nil {
		log.Fatal(err)
	}
	out := [][][]string{}
	for _, c := range picknames {
		intout := [][]string{}
		for _, a := range times {
			flag := len(intout)
			f, err := os.Open(directory + "/dates" + date + "/" + a.Name())
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				log.Fatal(err)
			}
			classes := strings.Split(string(b), "\n")
			for _, b := range classes {
				details := strings.Split(b, "█")
				names := strings.Split(strings.TrimSpace(details[2]), "|")
				for _, d := range names {
					if d == c {
						intout = append(intout, []string{details[0], details[1]})
						break
					}
				}
			}
		}
		out = append(out, intout)
	}
	message := ""
	if len(out) == 0 {
		message += "schedule not found"
	} else {
		for b := range out {
			if len(out[b]) == 0 {
				continue
			} else {
				message += "<br><b>"
				fragname := strings.Split(picknames[b], " ")
				for x := range fragname {
					message += strings.Title(fragname[x])
					if x == len(fragname)-1 {
						continue
					} else {
						message += " "
					}
				}
				message += "'s schedule on " + date[1:] + "</b><br><br>"
				for a := range out[b] {
					message += strings.TrimSuffix(times[a].Name(), ".txt") + ": " + out[b][a][0]
					if out[b][a][1] != " " {
						message += " (" + out[b][a][1] + ")<br>"
					} else {
						message += "<br>"
					}
					message += "<br>"
				}
			}
		}
	}
	return message
}
