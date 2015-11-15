package main

import (
	"fmt"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"net/http"
)

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	bus := c.Value("bus").(chan string) // Acquire global message bus
	fmt.Fprintln(w, "index")
	bus <- "index"
	return nil
}

func main() {
	messageBus := make(chan string)
	// Global message bus is stored in root context for each request
	noodle.Origin = context.WithValue(noodle.Origin, "bus", messageBus)

	go func() {
		// Here we'll do some application-wide message processing
		for msg := range messageBus {
			fmt.Println(msg)
		}
	}()
	n := noodle.Default()
	http.Handle("/", n.Then(index))
	http.ListenAndServe(":8080", nil)
}
