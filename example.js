import diameter from 'k6/x/diameter';
import { check } from 'k6';

export let options = {
  iterations: 5,
  vus: 1,
}

let client = diameter.newClient();

export default function () {

  // Send CCR
  let msg = diameter.newMessage("CCR");
  msg.addAVP(1, "ValueFooBar");
  const response = diameter.send(client, msg);
  console.log("Result Code:", response);
  check(response, {'Result-Code == 2001': r => r == 2001,});
}
