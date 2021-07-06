package main

import (
	"fmt"
	"log"

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
		gsd := &DBrequests{Name: name, URL: link}
		db.Create(gsd)
		// doc2, err := htmlquery.LoadURL(link)
		// if err != nil {
		// 	log.Fatal(err)
		// }

	}

}
