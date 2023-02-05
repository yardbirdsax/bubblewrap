package main

import (
	"log"

	"github.com/yardbirdsax/bubblewrap"
)

func main() {
  userInput, err := bubblewrap.Input("What's your favorite color? > ")
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("favorite color was: %s", userInput)
}