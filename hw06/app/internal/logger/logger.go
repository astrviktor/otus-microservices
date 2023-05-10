package logger

import (
	"encoding/json"
	"go.uber.org/zap"
)

func New() (*zap.Logger, error) {
	sampleJSON := []byte(`{
       "level" : "info",
       "encoding": "json",
       "outputPaths":["stdout"],
       "errorOutputPaths":["stderr"],
       "encoderConfig": {
           "messageKey":"message",
           "levelKey":"level",
           "levelEncoder":"lowercase"
       }
   }`)

	var cfg zap.Config

	if err := json.Unmarshal(sampleJSON, &cfg); err != nil {
		panic(err)
	}

	//logger, err := cfg.Build()
	return cfg.Build()
}
