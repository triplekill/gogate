package discovery

import (
	"github.com/hashicorp/consul/api"
	"github.com/wanghongfei/gogate/conf"
	. "github.com/wanghongfei/gogate/conf"
	"github.com/wanghongfei/gogate/utils"
	"strconv"
	"strings"
)

type ConsulClient struct {
	client *api.Client
}


func NewConsulClient() (Client, error) {
	cfg := &api.Config{}
	cfg.Address = conf.App.ConsulConfig.Address
	cfg.Scheme = "http"

	c, err := api.NewClient(cfg)
	if nil != err {
		return nil, utils.Errorf("failed to init consule client => %w", err)
	}

	return &ConsulClient{client:c}, nil
}

func (c *ConsulClient) QueryServices() ([]*InstanceInfo, error) {
	servMap, err := c.client.Agent().Services()
	if nil != err {
		return nil, err
	}

	// 查出所有健康实例
	healthList, _, err := c.client.Health().State("passing", &api.QueryOptions{})
	if nil != err {
		return nil, utils.Errorf("failed to query consul => %w", err)
	}

	instances := make([]*InstanceInfo, 0, 10)
	for _, servInfo := range servMap {
		servName := servInfo.Service
		servId := servInfo.ID

		// 查查在healthList中有没有
		isHealth := false
		for _, healthInfo := range healthList {
			if healthInfo.ServiceName == servName && healthInfo.ServiceID == servId {
				isHealth = true
				break
			}
		}

		if !isHealth {
			Log.Warn("following instance is not health, skip; service name: %v, service id: %v", servName, servId)
			continue
		}

		instances = append(
			instances,
			&InstanceInfo{
				ServiceName: strings.ToUpper(servInfo.Service),
				Addr: servInfo.Address + ":" + strconv.Itoa(servInfo.Port),
				Meta: servInfo.Meta,
			},
		)
	}

	return instances, nil
}

func (c *ConsulClient) Register() error {
	return utils.Errorf("not implement yet")
}

func (c *ConsulClient) UnRegister() error {
	return utils.Errorf("not implement yet")
}
