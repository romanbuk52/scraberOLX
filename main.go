package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/antchfx/htmlquery"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBrequests struct {
	gorm.Model
	id           int `gorm:"primaryKey"`
	Name         string
	URL          string
	statusReaded bool
}

func main() {
	db, err := gorm.Open(sqlite.Open("URL.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&DBrequests{})
	// scrubbLinks(db)
	scrubbSaleAds(db)
}

func scrubbLinks(db *gorm.DB) {
	// https://www.olx.ua/uk/elektronika/?page=2
	doc, err := htmlquery.LoadURL("https://www.olx.ua/uk/elektronika/")
	if err != nil {
		log.Fatal(err)
	}

	list := htmlquery.Find(doc, "//div/h3/a")
	for i, n := range list {
		a := htmlquery.FindOne(n, "//a")
		link := htmlquery.SelectAttr(a, "href")
		name := htmlquery.InnerText(a)
		fmt.Printf("%d %s(%s)\n", i, name, link)
		data := &DBrequests{Name: name, URL: link, statusReaded: false}
		db.Create(data)

	}
	return
}

func scrubbSaleAds(db *gorm.DB) {
	allRequests := make([]DBrequests, 0, 2)
	db.Find(&allRequests)
	println("len allRequests:", len(allRequests))
	// fmt.Printf("%+v\n", allRequests)

	for _, request := range allRequests {
		if request.statusReaded == false {
			doc, err := htmlquery.LoadURL(request.URL)
			println("name:", request.Name)
			println("requestURL:", request.URL)
			list := htmlquery.Find(doc, "//div/img/src")
			for i, r := range list {
				img := htmlquery.FindOne(r, "//img")
				link := htmlquery.SelectAttr(img, "src")
				println("linkIMG:", link)
				err = downloadFile(request.Name, link)
				// name := htmlquery.InnerText(a)
				// fmt.Printf("%d %s(%s)\n", i, link)
				println("порядковий номер в list:", i)
			}

			// set status readed if no errors
			if err == nil {
				// request.statusReaded = true
			}

		}

	}

}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
