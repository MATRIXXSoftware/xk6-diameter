package diameter

import (
	"encoding/json"
	"errors"
	"time"
)

type DiameterConfig struct {
	RequestTimeout              *Duration                 `json:"requestTimeout,omitempty"`
	MaxRetransmits              *uint                     `json:"maxRetransmits,omitempty"`
	RetransmitInterval          *Duration                 `json:"retransmitInterval,omitempty"`
	EnableWatchdog              *bool                     `json:"enableWatchdog,omitempty"`
	WatchdogInterval            *Duration                 `json:"watchdogInterval,omitempty"`
	WatchdogStream              *uint                     `json:"watchdogStream,omitempty"`
	SupportedVendorID           *[]uint32                 `json:"supportedVendorID,omitempty"`
	AcctApplicationID           *[]uint32                 `json:"acctApplicationId,omitempty"`
	AuthApplicationId           *[]uint32                 `json:"authApplicationId,omitempty"`
	VendorSpecificApplicationID *[]uint32                 `json:"vendorSpecificApplicationId,omitempty"`
	TransportProtocol           *string                   `josn:"transportProtocol,omitempty"`
	CapabilityExchange          *CapabilityExchangeConfig `json:"capabilityExchange,omitempty"`
}

type CapabilityExchangeConfig struct {
	VendorID         *uint32   `json:"vendorID"`
	ProductName      *string   `json:"productName,omitempty"`
	OriginHost       *string   `json:"originHost,omitempty"`
	OriginRealm      *string   `json:"originRealm,omitempty"`
	FirmwareRevision *uint32   `json:"firmwareRevision,omitempty"`
	HostIPAddresses  *[]string `json:"hostIPAddresses,omitempty"`
}

func parseConfig(arg map[string]interface{}) (*DiameterConfig, error) {

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
	var defaultFirmwareRevision uint32 = 1
	var defaultHostIPAddresses = []string{"127.0.0.1"}
	var defaultSupportedVendorID = []uint32{}
	var defaultAcctApplicationID = []uint32{}
	var defaultAuthApplicationID = []uint32{}
	var defaultVendorSpecificApplicationID = []uint32{}
	var defaultTransportProtocol = "tcp"

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
	if config.SupportedVendorID == nil {
		config.SupportedVendorID = &defaultSupportedVendorID
	}
	if config.AcctApplicationID == nil {
		config.AcctApplicationID = &defaultAcctApplicationID
	}
	if config.AuthApplicationId == nil {
		config.AuthApplicationId = &defaultAuthApplicationID
	}
	if config.VendorSpecificApplicationID == nil {
		config.VendorSpecificApplicationID = &defaultVendorSpecificApplicationID
	}
	if config.TransportProtocol == nil {
		config.TransportProtocol = &defaultTransportProtocol
	}

	// Set defaults for CapabilityExchangeConfig
	if config.CapabilityExchange == nil {
		config.CapabilityExchange = &CapabilityExchangeConfig{}
	}
	if config.CapabilityExchange.VendorID == nil {
		config.CapabilityExchange.VendorID = &defaultVendorID
	}
	if config.CapabilityExchange.ProductName == nil {
		config.CapabilityExchange.ProductName = &defaultProductName
	}
	if config.CapabilityExchange.OriginHost == nil {
		config.CapabilityExchange.OriginHost = &defaultOriginHost
	}
	if config.CapabilityExchange.OriginRealm == nil {
		config.CapabilityExchange.OriginRealm = &defaultOriginRealm
	}
	if config.CapabilityExchange.FirmwareRevision == nil {
		config.CapabilityExchange.FirmwareRevision = &defaultFirmwareRevision
	}
	if config.CapabilityExchange.HostIPAddresses == nil {
		config.CapabilityExchange.HostIPAddresses = &defaultHostIPAddresses
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
