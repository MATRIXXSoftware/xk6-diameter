import diam from 'k6/x/diameter'
import avp from 'k6/x/diameter/avp'
import dict from 'k6/x/diameter/dict'
import { cmd, app, code, flag, vendor } from './diam/const.js'

import { check } from 'k6'

export let options = {
    iterations: 1,
    vus: 1,
}

// Load additional custom AVP definition
// dict.load("dict/extra.xml")

// Init Client
let client = diam.Client({
    requestTimeout: "50ms",
    enableWatchdog: false,
})
let dataType = diam.DataType()

export default function () {
    client.connect("localhost:3868")

    let ccr = diam.newMessage(cmd.CreditControl, app.ChargingControl);
    ccr.add(avp.New(code.OriginHost,         0,     0,       dataType.DiameterIdentity("origin.host")))
    ccr.add(avp.New(code.OriginRealm,        0,     0,       dataType.DiameterIdentity("origin.realm")))
    ccr.add(avp.New(code.DestinationHost,    0,     0,       dataType.DiameterIdentity("dest.host")))
    ccr.add(avp.New(code.DestinationRealm,   0,     0,       dataType.DiameterIdentity("dest.realm")))
    ccr.add(avp.New(code.SessionId,          0,     flag.M,  dataType.UTF8String("Session-8888")))
    ccr.add(avp.New(code.CCRequestType,      0,     flag.M,  dataType.Enumerated(1)))
    ccr.add(avp.New(code.CCRequestNumber,    0,     flag.M,  dataType.Unsigned32(1000)))
    ccr.add(avp.New(code.SubscriptionId,     0,     flag.M,  dataType.Grouped([
        avp.New(code.SubscriptionIdData,     0,     flag.M,  dataType.UTF8String("subs-data")),
        avp.New(code.SubscriptionIdType,     0,     flag.M,  dataType.Enumerated(1))
    ])))

    const cca = client.send(ccr)
    console.log("CCA: ", cca)

    const resultCode = cca.findAVP(code.ResultCode, 0)
    check(resultCode, {'Result-Code == 2001': r => r == 2001,})
}
