package mosn

import (
	"context"
	"encoding/json"
	"errors"

	"mosn.io/layotto/components/rpc"
	"mosn.io/layotto/components/rpc/callback"
	"mosn.io/layotto/components/rpc/invoker/mosn/channel"
	_ "mosn.io/mosn/pkg/filter/network/proxy"
	"mosn.io/pkg/log"
)

const (
	Name = "mosn"
)

type mosnInvoker struct {
	channel rpc.Channel
	cb      rpc.Callback
}

type mosnConfig struct {
	Before  []rpc.CallbackFunc      `json:"before_invoke"`
	After   []rpc.CallbackFunc      `json:"after_invoke"`
	Channel []channel.ChannelConfig `json:"channel"`
}

func NewMosnInvoker() rpc.Invoker {
	invoker := &mosnInvoker{cb: callback.NewCallback()}
	return invoker
}

func (m *mosnInvoker) Init(conf rpc.RpcConfig) error {
	var config mosnConfig
	if err := json.Unmarshal(conf.Config, &config); err != nil {
		return err
	}

	for _, before := range config.Before {
		m.cb.AddBeforeInvoke(before)
	}

	for _, after := range config.After {
		m.cb.AddAfterInvoke(after)
	}

	if len(config.Channel) == 0 {
		return errors.New("missing channel config")
	}

	// todo support multiple channel
	channel, err := channel.GetChannel(config.Channel[0])
	if err != nil {
		return err
	}
	m.channel = channel
	return nil
}

func (m *mosnInvoker) Invoke(ctx context.Context, req *rpc.RPCRequest) (*rpc.RPCResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.DefaultLogger.Errorf("[runtime][rpc]mosn invoker panic: %v", r)
		}
	}()

	if req.Timeout == 0 {
		req.Timeout = 3000
	}
	log.DefaultLogger.Debugf("[runtime][rpc]request %+v", req)
	req, err := m.cb.BeforeInvoke(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]before filter error %s", err.Error())
		return nil, err
	}

	resp, err := m.channel.Do(req)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]error %s", err.Error())
		return nil, err
	}

	resp.Ctx = req.Ctx
	resp, err = m.cb.AfterInvoke(resp)
	if err != nil {
		log.DefaultLogger.Errorf("[runtime][rpc]after filter error %s", err.Error())
	}
	return resp, err
}
