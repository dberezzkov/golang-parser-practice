package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	url := "https://oresrekhovo-zuevo.hh.ru/vacancy/94934082?query=Middle+java+developer&hhtmFrom=vacancy_search_list"
	resp, err := http.Get(url)

	// закрыте запроса
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

	//salary := doc.Find("span.bloko-header-section-2 bloko-header-section-2_lite").Text()
	// smth := doc.Find("h1.bloko-header-section-1").Text()
	smth := doc.Find("h1.bloko-header-section-1>div>span.bloko-header-section-2 bloko-header-section-2_lite").Text()
	fmt.Println("1")
	fmt.Println(string(smth))
	fmt.Println("1")
}
