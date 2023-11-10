import diameter from 'k6/x/diameter';
import { check } from 'k6';

export let options = {
    iterations: 5,
    vus: 2,
}

let client = diameter.newClient();

export default function () {
    client.connect("localhost:3868")

    let msg = diameter.newMessage("CCR");
    msg.addAVP().Code(263).Mbit().UTF8String("Session-8888");      // Session ID
    msg.addAVP().Code(264).DiameterIdentity("origin.host")         // Origin-Host
    msg.addAVP().Code(296).DiameterIdentity("origin.realm")        // Origin-Realm
    msg.addAVP().Code(283).DiameterIdentity("dest.host")           // Destination-Host
    msg.addAVP().Code(293).DiameterIdentity("dest.realm")          // Destination-Realm
    msg.addAVP().Code(1).Vendor(10415).UTF8String("ValueFooBar");  // User-Name

    const response = diameter.send(client, msg);
    console.log("Result Code:", response);

    check(response, {'Result-Code == 2001': r => r == 2001,});
}
