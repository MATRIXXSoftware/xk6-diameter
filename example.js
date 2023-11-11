import diameter from 'k6/x/diameter';
import { check } from 'k6';

export let options = {
    // iterations: 500000,
    // vus: 32,

    iterations: 2,
    vus: 1,
}

let Vbit = 0x80;
let Mbit = 0x40;
let Pbit = 0x20;
let client = diameter.newClient();
let dataType = diameter.DataType();

export default function () {
    client.connect("localhost:3868")

    let msg = diameter.newMessage("CCR");
    msg.AVP(264, 0,     0,    dataType.DiameterIdentity("origin.host"))     // Origin-Host
    msg.AVP(296, 0,     0,    dataType.DiameterIdentity("origin.realm"))    // Origin-Realm
    msg.AVP(283, 0,     0,    dataType.DiameterIdentity("dest.host"))       // Destination-Host
    msg.AVP(293, 0,     0,    dataType.DiameterIdentity("dest.realm"))      // Destination-Realm

    msg.AVP(263, 0,     Mbit, dataType.UTF8String("Session-8888"))          // Session ID
    msg.AVP(416, 0,     Mbit, dataType.Enumerated(1))                       // CC-Request-Type
    msg.AVP(415, 0,     Mbit, dataType.Unsigned32(1000))                    // CC-Request-Number

    const response = diameter.send(client, msg);
    //console.log("Result Code:", response);

    check(response, {'Result-Code == 2001': r => r == 2001,});
}
