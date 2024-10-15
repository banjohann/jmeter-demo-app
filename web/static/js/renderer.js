
const onNewMessage = (message) => {
    if (message.type === 1) {
        renderClientName(message.clientName)
    } else {
        renderNewMessage(message)
    }
}

const renderNewMessage = (message) => {
    // get element by id
    var messagesContainer = document.getElementById("messages_container")

    

    // get childs
    // append new child
}

const renderClientName = (clientName) => {
  document.getElementById('user__name').textContent = clientName;
}

export { onNewMessage };