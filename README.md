# k6 Extension for Diameter

## Overview

This extension adds support for the Diameter protocol to [k6](https://k6.io/).

## Build

```bash
make
```
The Makefile will automatically download [xk6](https://github.com/grafana/xk6), which is required to compile this project.

## Generator
`bin/dict_generator` is used to generate `diam/const.js`, which contains a list of AVP codes and vendor IDs.

```bash
./bin/dict_generator -output example/diam/const.js
```

## Example

```js
import diam from 'k6/x/diameter'
import avp from 'k6/x/diameter/avp'
import { avpCode, flags, vendorId } from './diam/const.js'

let client = diam.Client()
let dataType = diam.DataType()

export default function () {
    client.connect("localhost:3868")

    let msg = diam.newMessage("CCR");

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

    const response = diam.send(client, msg)

    check(response, {'Result-Code == 2001': r => r == 2001,})
}
```

Use your custom k6 binary to run an example k6 script.
```bash
./bin/k6 run example/example.js
```