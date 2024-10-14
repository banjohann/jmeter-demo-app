import { socket, newWebsocketConnection } from './websocket.js'

document.addEventListener("DOMContentLoaded", () => {
    newWebsocketConnection();
});

document.getElementById('send__button').addEventListener('click', () => {
    let message = document.getElementById('input__message').value

    if (message != null) {
        sendMessage(message);
        console.log(message);
    }
});

let sendMessage = (message) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ text: message }));
  }
};
