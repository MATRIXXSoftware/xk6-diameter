package diameter

import (
	"github.com/dop251/goja"
	"go.k6.io/k6/js/modules"
)

type (
	Diameter struct {
		vu      modules.VU
		exports *goja.Object
	}
	RootModule struct{}
	Module     struct {
		*Diameter
	}
)

func init() {
	modules.Register("k6/x/diameter", &Diameter{})
	modules.Register("k6/x/diameter/avp", &AVP{})
	modules.Register("k6/x/diameter/dict", &Dict{})
}
