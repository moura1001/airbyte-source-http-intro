package main

import (
	sourcehttp "airbyte/custom-source/source-http"
	"log"

	"github.com/bitstrapped/airbyte"
)

func main() {
	httpsrc := sourcehttp.NewHTTPSRC("https://random-data-api.com/")
	runner := airbyte.NewSourceRunner(httpsrc)
	err := runner.Start()
	if err != nil {
		log.Fatal(err)
	}
}
