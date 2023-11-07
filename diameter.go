package diameter

import (
	"errors"
	"net"
	"time"

	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/fiorix/go-diameter/v4/diam/sm"
	log "github.com/sirupsen/logrus"
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/diameter", New())
}

func New() *Diameter {
	return &Diameter{}
}

type Diameter struct{}

func (*Diameter) NewClient() (*DiameterClient, error) {

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

	hopIds := make(map[uint32]chan *diam.Message)
	mux.Handle("CCA", handleCCA(hopIds))

	client := &sm.Client{
		Dict:               dict.Default,
		Handler:            mux,
		MaxRetransmits:     1,
		RetransmitInterval: time.Second,
		EnableWatchdog:     true,
		WatchdogInterval:   5 * time.Second,
		AuthApplicationID: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)),
		},
	}

	conn, err := client.DialNetwork("tcp", "localhost:3868")
	if err != nil {
		log.Errorf("Error connecting to %s, %v\n", "localhost:3868", err)
		return nil, err
	}

	log.Infof("Connected to %s\n", "localhost:3868")

	return &DiameterClient{
		client: client,
		conn:   conn,
		hopIds: hopIds,
	}, nil
}

func handleCCA(hopIds map[uint32]chan *diam.Message) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		hopByHopID := m.Header.HopByHopID
		v, exists := hopIds[hopByHopID]
		if !exists {
			log.Errorf("Received unexpected CCA with Hop-by-Hop ID %d\n", hopByHopID)
		} else {
			v <- m
		}
	}
}

func (*Diameter) NewMessage(name string) *DiameterMessage {
	return &DiameterMessage{
		name: name,
	}
}

type DiameterClient struct {
	client *sm.Client
	conn   diam.Conn
	hopIds map[uint32]chan *diam.Message
}

type DiameterMessage struct {
	//header
	name string   // test
	avps []string // test
}

func (m *DiameterMessage) AddAVP(avp string) {
	m.avps = append(m.avps, avp)
}

func (d *Diameter) Send(client *DiameterClient, msg *DiameterMessage) (uint32, error) {

	// TODO extract AVPs and Header from DiameterMessage

	req := diam.NewRequest(diam.CreditControl, 4, dict.Default)
	req.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("session-12345"))
	req.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity("origin.host"))
	req.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity("origin.realm"))
	req.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.DiameterIdentity("dest.realm"))
	req.NewAVP(avp.DestinationHost, avp.Mbit, 0, datatype.DiameterIdentity("dest.host"))
	req.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String("foobar"))

	// Keep track of Hop-by-Hop ID
	hopByHopID := req.Header.HopByHopID
	client.hopIds[hopByHopID] = make(chan *diam.Message)

	// Send CCR
	_, err := req.WriteTo(client.conn)
	if err != nil {
		return uint32(0), err
	}

	// Wait for CCA
	resp := <-client.hopIds[hopByHopID]
	//log.Infof("Received CCA \n%s", resp)

	delete(client.hopIds, hopByHopID)

	resultCodeAvp, err := resp.FindAVP(avp.ResultCode, 0)
	if err != nil {
		return uint32(0), errors.New("Result-Code AVP not found")
	}
	resultCode := resultCodeAvp.Data.(datatype.Unsigned32)

	return uint32(resultCode), nil
}
