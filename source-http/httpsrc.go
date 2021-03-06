package sourcehttp

import (
	"errors"
	"net/http"

	"github.com/bitstrapped/airbyte"
)

type HTTPSRC struct {
	baseURL string
}

func NewHTTPSRC(baseURL string) airbyte.Source {
	return HTTPSRC{
		baseURL: baseURL,
	}
}

// Spec returns the input "form" spec needed for your source
func (s HTTPSRC) Spec(logTracker airbyte.LogTracker) (*airbyte.ConnectorSpecification, error) {
	logTracker.Log(airbyte.LogLevelInfo, "inside spec")

	return &airbyte.ConnectorSpecification{
		DocumentationURL:      "https://random-data-api.com/",
		ChangeLogURL:          "https://random-data-api.com/",
		SupportsIncremental:   false,
		SupportsNormalization: true,
		SupportsDBT:           true,
		SupportedDestinationSyncModes: []airbyte.DestinationSyncMode{
			airbyte.DestinationSyncModeAppend,
			airbyte.DestinationSyncModeOverwrite,
		},
		ConnectionSpecification: airbyte.ConnectionSpecification{
			Title:       "Random Data API",
			Description: "Random Data Source API",
			Type:        "object",
			Required:    []airbyte.PropertyName{"numElements"},
			Properties: airbyte.Properties{
				Properties: map[airbyte.PropertyName]airbyte.PropertySpec{
					"numElements": {
						Description: "number of elements to pull per instance",
						Examples:    []string{"1", "7", "16"},
						PropertyType: airbyte.PropertyType{
							Type: []airbyte.PropType{
								airbyte.Integer,
							},
						},
					},
				},
			},
		},
	}, nil
}

type InputConfig struct {
	NumElements int16 `json:"numElements"`
}

// Check verifies the source - usually verify creds/connection etc.
func (s HTTPSRC) Check(srcCfgPath string, logTracker airbyte.LogTracker) error {
	logTracker.Log(airbyte.LogLevelInfo, "inside check")

	res, err := http.Get(s.baseURL)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("api internal error")
	}

	var ic InputConfig
	err = airbyte.UnmarshalFromPath(srcCfgPath, &ic)
	if err != nil {
		return err
	}

	if ic.NumElements <= 0 {
		return errors.New("should be a positive value greater than 0")
	}
	if ic.NumElements > 100 {
		return errors.New("can't have more then 100 elements")
	}

	return nil
}

type PhoneNumber struct {
	ID        int64  `json:"id"`
	UID       string `json:"uid"`
	CellPhone string `json:"cell_phone"`
}

type Code struct {
	ID  int64  `json:"id"`
	UID string `json:"uid"`
	NPI string `json:"npi"`
}

// Discover returns the schema of the data you want to sync
func (s HTTPSRC) Discover(srcConfigPath string, logTracker airbyte.LogTracker) (*airbyte.Catalog, error) {
	return &airbyte.Catalog{
		Streams: []airbyte.Stream{
			{
				Name: "PhoneNumber",
				JSONSchema: airbyte.Properties{
					Properties: map[airbyte.PropertyName]airbyte.PropertySpec{
						"id": {
							PropertyType: airbyte.PropertyType{
								Type:        []airbyte.PropType{airbyte.Integer, airbyte.Null},
								AirbyteType: airbyte.BigInteger,
							},
							Description: "",
						},
						"uid": {
							PropertyType: airbyte.PropertyType{
								Type: []airbyte.PropType{airbyte.String, airbyte.Null},
							},
							Description: "",
						},
						"cell_phone": {
							PropertyType: airbyte.PropertyType{
								Type: []airbyte.PropType{airbyte.String, airbyte.Null},
							},
							Description: "",
						},
					},
				},
				SupportedSyncModes: []airbyte.SyncMode{
					airbyte.SyncModeFullRefresh,
				},
				SourceDefinedCursor: false,
				Namespace:           "raw",
			},
			{
				Name: "Code",
				JSONSchema: airbyte.Properties{
					Properties: map[airbyte.PropertyName]airbyte.PropertySpec{
						"id": {
							PropertyType: airbyte.PropertyType{
								Type:        []airbyte.PropType{airbyte.Integer, airbyte.Null},
								AirbyteType: airbyte.BigInteger,
							},
							Description: "",
						},
						"uid": {
							PropertyType: airbyte.PropertyType{
								Type: []airbyte.PropType{airbyte.String, airbyte.Null},
							},
							Description: "",
						},
						"npi": {
							PropertyType: airbyte.PropertyType{
								Type: []airbyte.PropType{airbyte.String, airbyte.Null},
							},
							Description: "",
						},
					},
				},
				SupportedSyncModes: []airbyte.SyncMode{
					airbyte.SyncModeFullRefresh,
				},
				SourceDefinedCursor: false,
				Namespace:           "raw",
			},
		},
	}, nil
}

type State struct {
	LastSyncTime int64 `json:"lastSyncTime"`
}

// Read will read the actual data from your source and use tracker.Record(), tracker.State() and tracker.Log() to sync data with airbyte/destinations
// MessageTracker is thread-safe and so it is completely find to spin off goroutines to sync your data (just don't forget your waitgroups :))
// returning an error from this will cancel the sync and returning a nil from this will successfully end the sync
func (s HTTPSRC) Read(sourceCfgPath string, prevStatePath string, configuredCat *airbyte.ConfiguredCatalog,
	tracker airbyte.MessageTracker) error {

	tracker.Log(airbyte.LogLevelInfo, "inside read")

	var ic InputConfig
	err := airbyte.UnmarshalFromPath(sourceCfgPath, &ic)
	if err != nil {
		return err
	}

	var st State
	err = airbyte.UnmarshalFromPath(prevStatePath, &st)
	if err != nil {
		return err
	}

	for _, stream := range configuredCat.Streams {
		if stream.SyncMode == airbyte.SyncModeFullRefresh {
			errPN := FullRefreshPhoneNumber(stream, s.baseURL, ic.NumElements, tracker)

			errC := FullRefreshCode(stream, s.baseURL, ic.NumElements, tracker)

			if errPN != nil {
				return errPN
			}

			if errC != nil {
				return errC
			}
		}
	}

	return nil
}
