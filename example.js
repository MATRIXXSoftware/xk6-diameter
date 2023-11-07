import diameter from 'k6/x/diameter';

export let options = {
  iterations: 5,
  vus: 1,
}

let client = diameter.newClient();

export default function () {

  // Send CCR
  let msg = diameter.newMessage("CCR");
  msg.addAVP("Session-Id");
  const response = diameter.send(client, msg);

  console.log("Result Code:", response);
}
