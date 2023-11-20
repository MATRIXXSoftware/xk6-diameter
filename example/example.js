import diam from 'k6/x/diameter'
import avp from 'k6/x/diameter/avp'
import dict from 'k6/x/diameter/dict'
import { cmd, avpCode, flags, vendorId } from './diam/const.js'

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

    let msg = diam.newMessage(cmd.CreditControl);

    msg.AVP(avpCode.OriginHost,         0,     0,           dataType.DiameterIdentity("origin.host"))
    msg.AVP(avpCode.OriginRealm,        0,     0,           dataType.DiameterIdentity("origin.realm"))
    msg.AVP(avpCode.DestinationHost,    0,     0,           dataType.DiameterIdentity("dest.host"))
    msg.AVP(avpCode.DestinationRealm,   0,     0,           dataType.DiameterIdentity("dest.realm"))
    msg.AVP(avpCode.SessionId,          0,     flags.Mbit,  dataType.UTF8String("Session-8888"))
    msg.AVP(avpCode.CCRequestType,      0,     flags.Mbit,  dataType.Enumerated(1))
    msg.AVP(avpCode.CCRequestNumber,    0,     flags.Mbit,  dataType.Unsigned32(1000))
    msg.AVP(avpCode.SubscriptionId,     0,     flags.Mbit,  dataType.Grouped([
        avp.New(avpCode.SubscriptionIdData,     0,     flags.Mbit,  dataType.UTF8String("subs-data")),
        avp.New(avpCode.SubscriptionIdType,     0,     flags.Mbit,  dataType.Enumerated(1))
    ]))             

    const response = client.send(msg)
    console.log("Result Code:", response)

    check(response, {'Result-Code == 2001': r => r == 2001,})
}
