package diameter

import (
	"go.k6.io/k6/js/modules"
)

func init() {
	modules.Register("k6/x/diameter", New())
}

func New() *Diameter {
	return &Diameter{}
}

type Diameter struct{}

type DiameterMessage struct {
	//header
	//avps
}

func (*Diameter) XNew() *DiameterClient {
	return &DiameterClient{}
}

type DiameterClient struct {
}

func (d *Diameter) XSend(client *DiameterClient, message string) string {
	return "Send " + message
}
