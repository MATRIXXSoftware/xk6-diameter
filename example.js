import diameter from 'k6/x/diameter';
import { check } from 'k6';

export let options = {
    iterations: 5,
    vus: 1,
}

let client = diameter.newClient();

export default function () {
    let msg = diameter.newMessage("CCR");
    msg.addAVP(263).UTF8String("Session-8888");         // Session ID
    msg.addAVP(264).DiameterIdentity("origin.host")     // Origin-Host
    msg.addAVP(296).DiameterIdentity("origin.realm")    // Origin-Realm
    msg.addAVP(283).DiameterIdentity("dest.host")       // Destination-Host
    msg.addAVP(293).DiameterIdentity("dest.realm")      // Destination-Realm
    msg.addAVP(1).UTF8String("ValueFooBar");            // User-Name

    const response = diameter.send(client, msg);
    console.log("Result Code:", response);

    check(response, {'Result-Code == 2001': r => r == 2001,});
}
