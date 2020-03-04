package config

import (
	"sync"
	"time"
)

var (
	InitContainerImageName    = "init-image"
	InitContainerImageTagName = "init-image-tag"
	IAMNamespaceName          = "iam-namespace"
	IAMServicePortName        = "iam-service-port"
	IAMServicePort            = "4300"
	DefaultInitImage          = "hyc-cloud-private-edge-docker-local.artifactory.swg-devops.com/ibmcom/icp-initcontainer"
	DefaultInitImageTag       = "1.0.0-build.2"
	DefaultIamNamespace       = "iam-namespace"
	OperatorNS                = "openshift-monitoring"
)

type ControllerConfig struct {
	*sync.Mutex
	Values map[string]interface{}
}

var instance *ControllerConfig
var once sync.Once

func GetControllerConfig() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{
			Mutex:  &sync.Mutex{},
			Values: map[string]interface{}{},
		}
	})
	return instance
}

func (c *ControllerConfig) AddConfigItem(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	if key != "" && value != nil && value != "" {
		c.Values[key] = value
	}
}

func (c *ControllerConfig) RemoveConfigItem(key string) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.Values[key]; ok {
		delete(c.Values, key)
	}
}

func (c *ControllerConfig) GetConfigItem(key string, defaultValue interface{}) interface{} {
	if c.HasConfigItem(key) {
		return c.Values[key]
	}
	return defaultValue
}

func (c *ControllerConfig) GetConfigString(key, defaultValue string) string {
	if c.HasConfigItem(key) {
		return c.Values[key].(string)
	}
	return defaultValue
}

func (c *ControllerConfig) GetConfigBool(key string, defaultValue bool) bool {
	if c.HasConfigItem(key) {
		return c.Values[key].(bool)
	}
	return defaultValue
}

func (c *ControllerConfig) GetConfigTimestamp(key string, defaultValue time.Time) time.Time {
	if c.HasConfigItem(key) {
		return c.Values[key].(time.Time)
	}
	return defaultValue
}

func (c *ControllerConfig) HasConfigItem(key string) bool {
	c.Lock()
	defer c.Unlock()
	_, ok := c.Values[key]
	return ok
}
