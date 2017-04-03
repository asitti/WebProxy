package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	proxyHosts        []string          = make([]string, 10)
	proxyDestinations map[string]string = make(map[string]string)
)

func main() {

	// Print the version:
	log.Println(`WebProxy v1.0.0`)

	// Parse all configurations:
	log.Printf("The configuration: %s\n", os.Args[1:])
	for _, arg := range os.Args[1:] {
		elements := strings.Split(arg, ` => `)

		if len(elements) != 2 {
			log.Printf("The configuration '%s' uses a invalid format. Necessary format is: 'my-domain => http://www.other-domain.com'.\n", arg)
			continue
		}

		proxyHosts = append(proxyHosts, elements[0])
		proxyDestinations[elements[0]] = elements[1]
		log.Printf("Add configuration: host='%s', destination='%s'.\n", elements[0], proxyDestinations[elements[0]])
	}

	// Wire the handlers:
	serverMUX := http.NewServeMux()
	serverMUX.HandleFunc("/", proxy)

	// Set up the server:
	server := &http.Server{}
	server.Addr = ":80"
	server.Handler = serverMUX
	server.SetKeepAlivesEnabled(true)
	server.ReadTimeout = 60 * time.Second
	server.WriteTimeout = 60 * time.Second

	log.Println(`Proxy ready.`)
	server.ListenAndServe()
}

func proxy(response http.ResponseWriter, request *http.Request) {

	// Which host was requested?
	requestedHost := request.Host

	// Get all parameters:
	uri := request.RequestURI

	// Loop over all configured hosts:
	for _, host := range proxyHosts {

		// Known?
		if strings.Contains(requestedHost, host) {
			log.Printf("host='%s'\n", host)
			// Read the destination:
			destination := proxyDestinations[host]

			// Create a client in order to connect both:
			client := &http.Client{}

			// Create the client request:
			if clientRequest, clientErr := http.NewRequest(request.Method, destination+uri, request.Body); clientErr != nil {
				log.Printf("Was not able to create connection to destination '%s' for host '%s': %s\n", destination, requestedHost, clientErr.Error())
				http.NotFound(response, request)
				return
			} else {
				// Perform the request:
				if clientResponse, clientErrDo := client.Do(clientRequest); clientErrDo != nil {
					log.Printf("Was not able to perform the proxy request to destination '%s' for host '%s': %s\n", destination, requestedHost, clientErrDo.Error())
					http.NotFound(response, request)
					return
				} else {
					defer clientResponse.Body.Close()
					io.Copy(response, clientResponse.Body)
					return
				}
			}

			// Break the loop:
			return
		}
	}

	// The requested host is unknown:
	log.Printf("Host '%s' was not configured.\n", requestedHost)
	http.NotFound(response, request)
	return
}
