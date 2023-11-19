package diameter

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
)

type CapacityExchangeConfig struct {
	VendorID    *uint32 `json:"vendorID"`
	ProductName *string `json:"productName,omitempty"`
	OriginHost  *string `json:"originHost,omitempty"`
	OriginRealm *string `json:"originRealm,omitempty"`
}

type DiameterConfig struct {
	MaxRetransmits     *uint                   `json:"maxRetransmits,omitempty"`
	RetransmitInterval *time.Duration          `json:"retransmitInterval,omitempty"`
	EnableWatchdog     *bool                   `json:"enableWatchdog,omitempty"`
	WatchdogInterval   *time.Duration          `json:"watchdogInterval,omitempty"`
	WatchdogStream     *uint                   `json:"watchdogStream,omitempty"`
	CapacityExchange   *CapacityExchangeConfig `json:"capacityExchange,omitempty"`
}

func processConfig(arg map[string]interface{}) (*DiameterConfig, error) {

	var config DiameterConfig
	if b, err := json.Marshal(arg); err != nil {
		return nil, err
	} else {
		if err = json.Unmarshal(b, &config); err != nil {
			return nil, err
		}
	}

	setDiameterConfigDefaults(&config)

	log.Infof("Config %+v\n", config)
	log.Infof("CE Config %+v\n", config.CapacityExchange)

	return &config, nil
}

func setDiameterConfigDefaults(config *DiameterConfig) {
	// Default values
	var defaultMaxRetransmits uint = 1
	var defaultRetransmitInterval = 1 * time.Second
	var defaultEnableWatchdog = true
	var defaultWatchdogInterval = 5 * time.Second
	var defaultWatchdogStream uint = 0

	var defaultVendorID uint32 = 13
	var defaultProductName = "xk6-diameter"
	var defaultOriginHost = "origin.host"
	var defaultOriginRealm = "origin.realm"

	// Set defaults for DiameterConfig
	if config.MaxRetransmits == nil {
		config.MaxRetransmits = &defaultMaxRetransmits
	}
	if config.RetransmitInterval == nil {
		config.RetransmitInterval = &defaultRetransmitInterval
	}
	if config.EnableWatchdog == nil {
		config.EnableWatchdog = &defaultEnableWatchdog
	}
	if config.WatchdogInterval == nil {
		config.WatchdogInterval = &defaultWatchdogInterval
	}
	if config.WatchdogStream == nil {
		config.WatchdogStream = &defaultWatchdogStream
	}

	// Set defaults for CapacityExchangeConfig
	if config.CapacityExchange == nil {
		config.CapacityExchange = &CapacityExchangeConfig{}
	}
	if config.CapacityExchange.VendorID == nil {
		config.CapacityExchange.VendorID = &defaultVendorID
	}
	if config.CapacityExchange.ProductName == nil {
		config.CapacityExchange.ProductName = &defaultProductName
	}
	if config.CapacityExchange.OriginHost == nil {
		config.CapacityExchange.OriginHost = &defaultOriginHost
	}
	if config.CapacityExchange.OriginRealm == nil {
		config.CapacityExchange.OriginRealm = &defaultOriginRealm
	}
}
