import diameter from 'k6/x/diameter';

export let options = {
  iterations: 2,
  vus: 1,
}

// We want a diameter client per VU
export function setup() {
  console.log("Setting up diameter client for VU: " + __VU);
  let client = diameter.NewClient();
  return { client: client };
}

export default function (data) {
  let client = data.client;

  // Send CCR
  let msg = diameter.NewMessage("CCR");
  msg.AddAVP("Session-Id");
  const response = diameter.Send(client, msg);

  console.log(response);
}
