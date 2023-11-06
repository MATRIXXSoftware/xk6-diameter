import diameter from 'k6/x/diameter';

export let options = {
  iterations: 3,
  vus: 1,
}

let client = diameter.NewClient();

export default function () {

  // Send CCR
  let msg = diameter.NewMessage("CCR");
  msg.AddAVP("Session-Id");
  const response = diameter.Send(client, msg);

  console.log(response);
}
