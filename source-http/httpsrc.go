package sourcehttp

import "github.com/bitstrapped/airbyte"

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

// Check verifies the source - usually verify creds/connection etc.
func (s HTTPSRC) Check(srcCfgPath string, logTracker airbyte.LogTracker) error { return nil }

// Discover returns the schema of the data you want to sync
func (s HTTPSRC) Discover(srcConfigPath string, logTracker airbyte.LogTracker) (*airbyte.Catalog, error) {
	return nil, nil
}

// Read will read the actual data from your source and use tracker.Record(), tracker.State() and tracker.Log() to sync data with airbyte/destinations
// MessageTracker is thread-safe and so it is completely find to spin off goroutines to sync your data (just don't forget your waitgroups :))
// returning an error from this will cancel the sync and returning a nil from this will successfully end the sync
func (s HTTPSRC) Read(sourceCfgPath string, prevStatePath string, configuredCat *airbyte.ConfiguredCatalog,
	tracker airbyte.MessageTracker) error {
	return nil
}
