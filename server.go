package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello There")
	query := r.URL.Query()
	firstName := query.Get("name")
	fmt.Println("First Name is : ", firstName)
	if firstName == "" {
		firstName = "No-One"
	}

	lastName := query.Get("last")
	fmt.Println("Last Name is : ", lastName)
	if lastName == "" {
		lastName = "No-last-Name"
	}
	w.Write([]byte(fmt.Sprintf("Hello, %s%s%s\n", firstName, " ", lastName)))
}

//curl http://localhost:8081\?name\="Ankit"\&\last\="Chahal"

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	//log.Fatal(http.ListenAndServe(":8080", r)) Without any interrupt
	// Start Server
	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {

	/*

		if somewhere you are listening for output from a channel but in the meanwhile that channel not being given any input,
		the listening place will be blocked until it receives the output.

	*/

	interruptChan := make(chan os.Signal, 1)

	/*

		SIGINT means is Signal Interrupted, send when the user types the INTR character (e.g. Ctrl-C).
		SIGTERM signal is a generic signal used to cause program termination.
		The below statement is to send the interruot signal to the channel "interruptChan".

	*/

	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	/*

		some zombie connections may still be there and use your memory,
		in order to solve this, the safest way is to setup a timeout threshold using the Context.

	*/

	defer cancel()

	/*

		The timer context and the cancel function (the cancel function releases its resources) will be returned,
		 and you can use that context to perform Shutdown(ctx) ,
		inside the Shutdown it check if the timer context Done channel is closed and will not run indefinitely.

	*/

	srv.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}

/* Useful Links

https://www.rudderstack.com/blog/implementing-graceful-shutdown-in-go/
https://medium.com/@pinkudebnath/graceful-shutdown-of-golang-servers-using-context-and-os-signals-cc1fa2c55e97
https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a

*/
