package runtime

import (
	"encoding/json"
	"github.com/layotto/layotto/pkg/services/pubsub"

	"github.com/layotto/layotto/pkg/services/configstores"
	"github.com/layotto/layotto/pkg/services/hello"
)

type AppConfig struct {
	AppId            string `json:"app_id"`
	GrpcCallbackPort int    `json:"grpc_callback_port"`
}

type MosnRuntimeConfig struct {
	AppManagement          AppConfig                           `json:"app"`
	HelloServiceManagement map[string]hello.HelloConfig        `json:"hellos"`
	ConfigStoreManagement  map[string]configstores.StoreConfig `json:"config_stores"`
	PubSubManagement       map[string]pubsub.Config            `json:"pub_subs"`
}

func ParseRuntimeConfig(data json.RawMessage) (*MosnRuntimeConfig, error) {
	cfg := &MosnRuntimeConfig{}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
