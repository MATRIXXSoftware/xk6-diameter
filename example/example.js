import diam from 'k6/x/diameter'
import avp from 'k6/x/diameter/avp'
import dict from 'k6/x/diameter/dict'
import { cmd, cmdFlag, app, code, flag, vendor } from './diam/const.js'
import { check } from 'k6'

export let options = {
    iterations: 1,
    vus: 1,
}

// Load additional custom AVP definition
// dict.load("dict/extra.xml")

let data = diam.DataType()

let client = diam.Client({
    requestTimeout: "50ms",
    enableWatchdog: false,
    authApplicationId: [app.ChargingControl],
    vendorSpecificApplicationId: [
        {
            authApplicationId: app.ChargingControl,
            vendorId: vendor.TGPP,
        }
    ],
    capabilityExchange: {
        vendorId: 35838,
    },
})

export default function () {
    client.connect("localhost:3868")

    let ccr = diam.newMessage(cmd.CreditControl, app.ChargingControl, cmdFlag.Request);
    ccr.add(avp.New(code.OriginHost,         0,     0,       data.DiameterIdentity("origin.host")))
    ccr.add(avp.New(code.OriginRealm,        0,     0,       data.DiameterIdentity("origin.realm")))
    ccr.add(avp.New(code.DestinationHost,    0,     0,       data.DiameterIdentity("dest.host")))
    ccr.add(avp.New(code.DestinationRealm,   0,     0,       data.DiameterIdentity("dest.realm")))
    ccr.add(avp.New(code.SessionId,          0,     flag.M,  data.UTF8String("Session-8888")))
    ccr.add(avp.New(code.CCRequestType,      0,     flag.M,  data.Enumerated(1)))
    ccr.add(avp.New(code.CCRequestNumber,    0,     flag.M,  data.Unsigned32(1000)))
    ccr.add(avp.New(code.SubscriptionId,     0,     flag.M,  data.Grouped([
        avp.New(code.SubscriptionIdData,     0,     flag.M,  data.UTF8String("subs-data")),
        avp.New(code.SubscriptionIdType,     0,     flag.M,  data.Enumerated(1))
    ])))
    ccr.add(avp.New(code.ServiceInformation, 10415, flag.M,  data.Grouped([
        avp.New(code.PSInformation,          10415, flag.M,  data.Grouped([
            avp.New(code.CalledStationId,    0,     flag.M,  data.UTF8String("10099"))
        ]))
    ])))

    const cca = client.send(ccr)
    console.log(`CCA: ${cca}`)

    const resultCode = cca.findAVP(code.ResultCode, 0)
    check(resultCode, {'Result-Code == 2001': r => r == 2001,})
}
