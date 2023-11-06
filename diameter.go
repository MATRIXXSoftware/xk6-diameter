package diameter

import (
	"net"
	"time"

	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	"go.k6.io/k6/js/modules"

	log "github.com/sirupsen/logrus"
)

func init() {
	modules.Register("k6/x/diameter", New())
}

func New() *Diameter {
	return &Diameter{}
}

type Diameter struct{}

func (*Diameter) XNewClient() *DiameterClient {

	// TODO make all this configurable later
	cfg := &sm.Settings{
		OriginHost:       datatype.DiameterIdentity("diam.host"),
		OriginRealm:      datatype.DiameterIdentity("diam.realm"),
		VendorID:         13,
		ProductName:      "xk6-diameter",
		OriginStateID:    datatype.Unsigned32(time.Now().Unix()),
		FirmwareRevision: 1,
		HostIPAddresses: []datatype.Address{
			datatype.Address(net.ParseIP("127.0.0.1")),
		},
	}
	mux := sm.New(cfg)

	client := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     3,
		RetransmitInterval: time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   5 * time.Second,
		AuthApplicationID: []*diam.AVP{
			// Advertise support for credit control application
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)), // RFC 4006
		},
	}

	conn, err := client.DialNetwork("tcp", "localhost:3868")
	if err != nil {
		log.Errorf("Error connecting to %s, %v\n", "localhost:3868", err)
		panic(err)
	}

	log.Infof("Connected to %s\n", "localhost:3868")

	return &DiameterClient{
		client: client,
		conn:   conn,
	}
}

func (*Diameter) XNewMessage(name string) *DiameterMessage {
	return &DiameterMessage{
		name: name,
	}
}

type DiameterClient struct {
	client *sm.Client
	conn   diam.Conn
}

type DiameterMessage struct {
	//header
	name string   // test
	avps []string // test
}

func (m *DiameterMessage) XAddAVP(avp string) {
	m.avps = append(m.avps, avp)
}

func (d *Diameter) XSend(client *DiameterClient, msg *DiameterMessage) string {
	resp := "Send " + msg.name
	for _, avp := range msg.avps {
		resp = resp + " with avp " + avp
	}
	return resp
}
