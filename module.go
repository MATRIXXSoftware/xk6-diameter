package diameter

import (
	"go.k6.io/k6/js/modules"
)

type (
	Diameter struct {
		vu      modules.VU
		metrics DiameterMetrics
	}
	RootModule struct{}
)

func init() {
	modules.Register("k6/x/diameter", New())
	modules.Register("k6/x/diameter/avp", &AVP{})
	modules.Register("k6/x/diameter/dict", &Dict{})
}

func New() *RootModule {
	return &RootModule{}
}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Diameter{
		vu:      vu,
		metrics: registerMetrics(vu),
	}
}

func (d *Diameter) Exports() modules.Exports {
	return modules.Exports{Default: d}
}
