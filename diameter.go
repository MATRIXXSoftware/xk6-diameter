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

type Diameter struct{}

type DiameterClient struct {
	client *sm.Client
	conn   diam.Conn
	hopIds map[uint32]chan *diam.Message
}

type DiameterMessage struct {
	//header
	name string // test

	diamMsg *diam.Message
	avps    []*DiameterAVP
}

type DiameterAVP struct {
	code   uint32
	flags  uint8
	vendor uint32
	data   datatype.Type
}

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
	// TODO need to support other diameter CMD

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
	return func(_ diam.Conn, m *diam.Message) {
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

	diamMsg := diam.NewRequest(diam.CreditControl, 4, dict.Default)

	return &DiameterMessage{
		name:    name,
		diamMsg: diamMsg,
		avps:    []*DiameterAVP{},
	}
}

func (m *DiameterMessage) AddAVP() *DiameterAVP {
	// populate later
	avp := DiameterAVP{
		code:   0,
		flags:  0,
		vendor: 0,
		data:   nil,
	}
	m.avps = append(m.avps, &avp)
	return &avp
}

func (a *DiameterAVP) XCode(code uint32) *DiameterAVP {
	a.code = code
	return a
}

func (a *DiameterAVP) XMbit() *DiameterAVP {
	a.flags = a.flags | avp.Mbit
	return a
}

func (a *DiameterAVP) XPbit() *DiameterAVP {
	a.flags = a.flags | avp.Pbit
	return a
}

func (a *DiameterAVP) XVbit() *DiameterAVP {
	a.flags = a.flags | avp.Vbit
	return a
}

func (a *DiameterAVP) XVendor(vendor uint32) *DiameterAVP {
	a.vendor = vendor
	return a
}

func (a *DiameterAVP) XUTF8String(value string) *DiameterAVP {
	a.data = datatype.UTF8String(value)
	return a
}

func (a *DiameterAVP) XDiameterIdentity(value string) *DiameterAVP {
	a.data = datatype.DiameterIdentity(value)
	return a
}

// TODO add more data type

func (d *Diameter) Send(client *DiameterClient, msg *DiameterMessage) (uint32, error) {

	req := msg.diamMsg

	for _, avp := range msg.avps {
		req.NewAVP(avp.code, avp.flags, avp.vendor, avp.data)
	}

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

func init() {
	diameter := &Diameter{}
	modules.Register("k6/x/diameter", diameter)
}
