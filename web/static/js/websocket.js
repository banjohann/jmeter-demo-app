import { onNewMessage } from "./renderer";
import { Message } from "./message"

const isSecureConn = window.location.protocol === "https:";

let socket;

const newWebsocketConnection = () => {
  let protocol = isSecureConn ? "wss" : "ws";
  socket = new WebSocket(`${protocol}://${window.location.host}/ws`);

  socket.onopen = (event) => {
    console.log(event.data)

  };

  socket.onclose = (event) => {};

  socket.onmessage = (event) => {
    let jsonData = JSON.parse(event.data) 
    onNewMessage(new Message(jsonData.clientName, jsonData.text, jsonData.type))
  };
};

export { socket, newWebsocketConnection };