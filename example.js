import diameter from 'k6/x/diameter';

export default function () {
  let client = diameter.New();
  const response = diameter.Send(client, "Hello, Diameter!");

  console.log(response);
}
