const ws = new WebSocket("ws://" + window.location.host + "/ws")
const playerCustom = new Plyr("#player", {
  controls: [
    "play-large",
    "play",
    "progress",
    "current-time",
    "mute",
    "volume",
    "fullscreen"
  ],
  settings: []
});
const player = document.getElementById("player")

const arbuzPath = "/static/assets/watermelon-3593669936.jpg"

let videoFiles = null
fetch("http://" + window.location.host + "/api/videofiles")
  .then(response => response.json())
  .then(data => {
    videoFiles = data.videofiles
  })
  .catch(error => console.error("ошибка запроса: ", error))

ws.addEventListener("message", (event) => {
  const data = JSON.parse(event.data)
  switch (data.action) {
    case "updateUsers":
      updateUsers(data.usernames)
      break
    case "pause":
      player.pause()
      break
    case "play":
      player.play()
  }
})

player.addEventListener("pause", sendPauseSignal)
player.addEventListener("play", sendPlaySignal)

function updateUsers(usernames) {
  const list = document.getElementById("users-list")
  list.innerHTML = ""
  usernames.forEach(function(username) {
    const li = document.createElement("li")
    const liText = document.createElement("h2")
    const avatar = document.createElement("img")
    liText.textContent = username
    avatar.src = arbuzPath
    avatar.className = "users-list__item-image"
    li.className = "users-list__item"

    li.appendChild(avatar)
    li.appendChild(liText)

    list.appendChild(li)
  })
}

function sendPauseSignal() {
  const msg = {
    action: "pause"
  }
  const data = JSON.stringify(msg)
  ws.send(data)
}

function sendPlaySignal() {
  const msg = {
    action: "play"
  } 
  const data = JSON.stringify(msg)
  ws.send(data)
}

