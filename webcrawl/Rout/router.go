package Rout

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/PuerkitoBio/goquery"

	"github.com/gorilla/mux"
)

var Resp_Body Dat
var Result_array Response

func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/crawl", Crawl).Methods("POST")
	return r
}

func Crawl(res http.ResponseWriter, req *http.Request) {
	var body Request
	err := json.NewDecoder(req.Body).Decode(&body) //Decode request body
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	wg := new(sync.WaitGroup) //Wait group to avoid panic (wait untill goroutines are complete)
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i, value := range body.Urls {
			fmt.Println(i, "value", value)
			ReqChannel <- value //Using a channel to concurrently process the Urls
		}
		close(ReqChannel)
	}()
	go func() {
		defer wg.Done()
		for i := range ReqChannel {
			Resp_Body.Url = i
			fmt.Println("url is", Resp_Body.Url)
			resp, err := http.Get(Resp_Body.Url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
			}
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			Resp_Body.Data = doc.Find("title").Text()
			fmt.Println("title", Resp_Body.Data)
			Result_array.Result = append(Result_array.Result, Resp_Body)
		}
	}()
	wg.Wait()
	Final_resp, err := json.Marshal(Result_array)
	if err != nil {
		fmt.Println(err)
	}
	_, err = res.Write(Final_resp)
	if err != nil {
		fmt.Println(err)
	}
	return
}
