package diameter

import (
	"encoding/json"
	"errors"
	"time"
)

type CapacityExchangeConfig struct {
	VendorID    *uint32 `json:"vendorID"`
	ProductName *string `json:"productName,omitempty"`
	OriginHost  *string `json:"originHost,omitempty"`
	OriginRealm *string `json:"originRealm,omitempty"`
}

type DiameterConfig struct {
	RequestTimeout     *Duration               `json:"requestTimeout,omitempty"`
	MaxRetransmits     *uint                   `json:"maxRetransmits,omitempty"`
	RetransmitInterval *Duration               `json:"retransmitInterval,omitempty"`
	EnableWatchdog     *bool                   `json:"enableWatchdog,omitempty"`
	WatchdogInterval   *Duration               `json:"watchdogInterval,omitempty"`
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

	return &config, nil
}

func setDiameterConfigDefaults(config *DiameterConfig) {
	// Default values
	var defaultRequestTimeout = Duration{1 * time.Second}
	var defaultMaxRetransmits uint = 1
	var defaultRetransmitInterval = Duration{1 * time.Second}
	var defaultEnableWatchdog = true
	var defaultWatchdogInterval = Duration{5 * time.Second}
	var defaultWatchdogStream uint = 0

	var defaultVendorID uint32 = 13
	var defaultProductName = "xk6-diameter"
	var defaultOriginHost = "origin.host"
	var defaultOriginRealm = "origin.realm"

	// Set defaults for DiameterConfig
	if config.RequestTimeout == nil {
		config.RequestTimeout = &defaultRequestTimeout
	}
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

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid duration")
	}
}
