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
dict.load("dict/extra.xml")

// Init Client
let client = diam.Client({
    requestTimeout: "50ms",
})
let dataType = diam.DataType()

export default function () {
    client.connect("localhost:3868")

    let msg = diam.newMessage(cmd.CreditControl, app.ChargingControl);

    msg.AVP(code.OriginHost,         0,     0,       dataType.DiameterIdentity("origin.host"))
    msg.AVP(code.OriginRealm,        0,     0,       dataType.DiameterIdentity("origin.realm"))
    msg.AVP(code.DestinationHost,    0,     0,       dataType.DiameterIdentity("dest.host"))
    msg.AVP(code.DestinationRealm,   0,     0,       dataType.DiameterIdentity("dest.realm"))
    msg.AVP(code.SessionId,          0,     flag.M,  dataType.UTF8String("Session-8888"))
    msg.AVP(code.CCRequestType,      0,     flag.M,  dataType.Enumerated(1))
    msg.AVP(code.CCRequestNumber,    0,     flag.M,  dataType.Unsigned32(1000))
    msg.AVP(code.SubscriptionId,     0,     flag.M,  dataType.Grouped([
        avp.New(code.SubscriptionIdData,     0,     flag.M,  dataType.UTF8String("subs-data")),
        avp.New(code.SubscriptionIdType,     0,     flag.M,  dataType.Enumerated(1))
    ]))             

    const response = client.send(msg)
    console.log("Response: ", response.dump())

    const resultCode = response.findAVP(code.ResultCode, 0)
    check(resultCode, {'Result-Code == 2001': r => r == 2001,})
}
