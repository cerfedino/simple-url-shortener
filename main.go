package main

import "log"

func main() {
	var model Storage = &postgresstorage{}
	model.init("")

	model.storeShortenedURL("https://www.google.com", "tg")

	long_url, err := model.getLongerURL("tg")
	if err != nil {
		log.Println(err)
	}
	log.Println(long_url)

	model.removeShortenedURL("tg")
	_, err = model.getLongerURL("tg")
	if err != nil {
		log.Println(err)
	}
}
