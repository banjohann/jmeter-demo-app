const reconnectInterval = 5000;
const maxReconnectAttempts = 3;
const isSecureConn = window.location.protocol === "https:";

let socket;
let reconnectAttempts = 0;

const newWebsocketConnection = () => {
  let protocol = isSecureConn ? "wss" : "ws";
  socket = new WebSocket(`${protocol}://${window.location.host}/ws`);

  socket.onopen = (event) => {
    console.log(event.data)

  };

  socket.onclose = (event) => {};

  socket.onmessage = (event) => {
    let data = JSON.parse(event.data)
    console.log(data)

    if (data.type == 1) {
      console.log(data.client_name)
      renderClientName(data.client_name)
    }
  };
};

const renderClientName = (clientName) => {
  document.getElementById('user__name').textContent = clientName;
}


export { socket, newWebsocketConnection };
