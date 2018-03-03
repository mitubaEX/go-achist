package main

import (
	"fmt"
	"net/http"
	"github.com/go-achist/Models"
	"github.com/go-achist/Services"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"strings"
	"time"
	"bytes"
)

func GetResultString(rateData *Models.RateData) string {
	var transition = "Highest"
	if rateData.Subset < 0 {
		transition = "Lowest"
	}
	return fmt.Sprintf("%d->%d(%d) %s\nPerformance %d\n",
		rateData.OldRate, rateData.NewRate, rateData.Subset, transition, rateData.Performance)
}

func handler(w rest.ResponseWriter, r *rest.Request) {
	name := r.PathParam("name")
	fmt.Println(name)
	result := GetResultString(Services.CrawRateData(name))
	w.WriteJson(map[string]string{"Body": result})
}

func sendMessage (message string) {
	name := "ACNotification"
	text := message
	channel := "bots"

	jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + text + `"}`

	req, err := http.NewRequest(
		"POST",
		"https://hooks.slack.com/services/hogehoge",
		bytes.NewBuffer([]byte(jsonStr)),
	)

	if err != nil {
		fmt.Print(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Print(err)
	}

	defer resp.Body.Close()
}

func notification (w rest.ResponseWriter, r *rest.Request) {
	result := Services.CrawContestData()
	splitedString := strings.Split(result[1].Date, " ")

	t := time.Now()
	const layout = "2006/01/02"

	if t.Format(layout) == splitedString[0] {
		sendMessage(result[1].Date + " " + result[1].Name)
	}

	w.WriteJson(map[string]string{"challenge": "hello"})
}

func main() {
	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)
	router, err := rest.MakeRouter(
		rest.Get("/rate/#name", handler),
		rest.Get("/", func(w rest.ResponseWriter, r *rest.Request){
			w.WriteJson(map[string]string{"Body": "Hello, World"})
		}),
		rest.Get("/notification", notification),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(router)
	log.Fatal(http.ListenAndServe(":8080", api.MakeHandler()))
}
