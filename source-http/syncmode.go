package sourcehttp

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bitstrapped/airbyte"
)

func FullRefreshPhoneNumber(stream airbyte.ConfiguredStream, baseURL string, numElements int16, tracker airbyte.MessageTracker) error {
	if stream.Stream.Name == "PhoneNumber" {
		var pns []PhoneNumber

		res, err := http.Get(fmt.Sprintf("%s/api/phone_number/random_phone_number?size=%d", baseURL, numElements))

		if err != nil {
			return err
		}

		json.NewDecoder(res.Body).Decode(&pns)

		for _, pn := range pns {
			tracker.Record(pn, stream.Stream.Name, stream.Stream.Namespace)
		}

	}

	return nil
}

func FullRefreshCode(stream airbyte.ConfiguredStream, baseURL string, numElements int16, tracker airbyte.MessageTracker) error {
	if stream.Stream.Name == "Code" {
		var cs []Code

		res, err := http.Get(fmt.Sprintf("%s/api/code/random_code?size=%d", baseURL, numElements))

		if err != nil {
			return err
		}

		json.NewDecoder(res.Body).Decode(&cs)

		for _, c := range cs {
			tracker.Record(c, stream.Stream.Name, stream.Stream.Namespace)
		}

	}

	return nil
}
