# k6 Extension for Diameter

## Overview

This extension adds support for the Diameter protocol to [k6](https://k6.io/).

## Build

```bash
make
```
The Makefile will automatically download [xk6](https://github.com/grafana/xk6), which is required to compile this project.

## Example

```js
import diam from 'k6/x/diameter'
import avp from 'k6/x/diameter/avp'
import dict from 'k6/x/diameter/dict'
import { cmd, app, code, flag, vendor } from './diam/const.js'

import { check } from 'k6'

let client = diam.Client()
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
```

Use your custom k6 binary to run an example k6 script.
```bash
./bin/k6 run example/example.js
```

## Generator

There are thousands of AVPs, each with a unique avp-code and vendor-id. To aid readability and enhance the developer experience, we recommend defining them as constants in a separate file, for example, using `diam/const.js`.

You can either create the constant yourself or use the bin/dict_generator CLI tool to generate a full list of AVPs for you. Use the following command:
```bash
./bin/dict_generator -output example/diam/const.js
```

The CLI also supports generating additional AVPs that are not defined in the default list. Simply add the -dictionary flag to include the additional AVP definition:
```bash
./bin/dict_generator -output example/diam/const.js -dictionary dict/extra.xml
```
