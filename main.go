package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"strings"
	"os"

)

type entry struct {
	ID               int      `json:"id"`
	UserName         string   `json:"user_Name"`
	Ips              []string `json:"ips"`
	Target           string   `json:"target"`
	EVENT0ACTION     string   `json:"EVENT_0_ACTION"`
	DateTimeAndStuff int      `json:"DateTimeAndStuff"`
}

type normalEntry struct {
	AcmeApiId 		int
	UserName  		string
	SourceIp 		string
	Target   		string
	Action   		string
	EventTime 		string
}

type lastEvent struct {
	EntryCount		int `json:"EntryCount"`
	LastEntryHash	string `json:"LastEntryHash"`	
}
var normalEntries []normalEntry
var getVisualizationData = false

func main() {
	var token = getToken()
	var from = getOldEntryNum()
	if from != 0 {
		from++
		loadEntries()
	}
	var to = getNewEntryNum(token)
	if from == to {
		return
	}
	if to - from > 500 {
		paging(token, from, to)
	} else {
		normalizeEntries(token, from, to)
	}
	if getVisualizationData {
		visualizationData()
	}

	
}


func normalizeName (name string) string {
	words := strings.Fields(name)
	var user string
	for i := 0; i < len(words); i++ {
		if strings.Contains(words[i], "@") {
			user = words[i]
		}
	}
	return strings.ToLower(user)
}

func normalizeAction (action string) string {
	if strings.Contains(strings.ToLower(action), "success") {
		return "Logon-Success"
	}

	return "Logon-Failure"

}

func getToken() string {
	url := "https://challenger.btbsecurity.com/auth"

	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	token, _ := ioutil.ReadAll(resp.Body)

	return string(token)
}

func loadEntries() {
	jsonFile, err := os.Open("entries.json")
	if err != nil {
		fmt.Println(err)
	}
	oldEntries, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(oldEntries, &normalEntries)
}

func getOldEntryNum() int {
	jsonFile, err := os.Open("numLogs.json")

	if err != nil {
		return 0
	}

	var lastEvent = lastEvent{}

	event, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(event, &lastEvent)
	return lastEvent.EntryCount
}

func getNewEntryNum(token string) int {
	url := "https://challenger.btbsecurity.com/get-events"

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var lastEvent = lastEvent{}

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &lastEvent)
	if err != nil {
		fmt.Println(err)
	}

	file, _ := json.MarshalIndent(lastEvent, "", " ")
	_ = ioutil.WriteFile("numLogs.json", file, 0644)

	return int(lastEvent.EntryCount)
}

func paging(token string, from int, to int) {
	var hold int = from + 499
	for to-from >= 500 {
		normalizeEntries(token, from, hold)
		hold += 500
		from += 500
	}
	normalizeEntries(token, from, to)
}


func normalizeEntries(token string, from int, to int) {

	url := fmt.Sprintf("https://challenger.btbsecurity.com/get-events?from=%d&to=%d", from, to)

	req, err := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	var entries []entry

	body, _ := ioutil.ReadAll(resp.Body)


	err = json.Unmarshal(body, &entries)
	if err != nil {
		fmt.Println(err)
	}




	for i := 0; i < len(entries); i++ {
		var normal = normalEntry{}
		normal.AcmeApiId = entries[i].ID
		normal.UserName = normalizeName(entries[i].UserName)
		normal.SourceIp = entries[i].Ips[0]
		normal.Target = entries[i].Target
		normal.Action = normalizeAction(entries[i].EVENT0ACTION)
		unixTime := time.Unix(int64(entries[i].DateTimeAndStuff),0)
		time := unixTime.Format(time.RubyDate)
		normal.EventTime = time
		normalEntries = append(normalEntries, normal)
	}
	file, _ := json.MarshalIndent(normalEntries, "", " ")
	_ = ioutil.WriteFile("entries.json", file, 0644)
}

func visualizationData() {
	mon, tue, wed, thu, fri, sat, sun := 0, 0, 0, 0, 0, 0, 0
	for i := 0; i < len(normalEntries); i++ {
		var day = normalEntries[i].EventTime
		if strings.Contains(day, "Mon") {
			mon++
		} else if strings.Contains(day, "Tue") {
			tue++
		} else if strings.Contains(day, "Wed") {
			wed++
		} else if strings.Contains(day, "Thu") {
			thu++
		} else if strings.Contains(day, "Fri") {
			fri++
		} else if strings.Contains(day, "Sat") {
			sat++
		} else if strings.Contains(day, "Sun") {
			sun++
		}
		var data = fmt.Sprintf("# of attempts on Monday: %d\n# of attempts on Tuesday: %d\n# of attempts on Wednesday: %d\n# of attempts on Thursday: %d\n# of attempts on Friday: %d\n# of attempts on Saturday: %d\n# of attempts on Sunday: %d\n",mon,tue,wed,thu,fri,sat,sun)
	
		ioutil.WriteFile("visualizationData.txt", []byte(data), 0644)
	
	}
	
}

