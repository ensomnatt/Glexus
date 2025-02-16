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
const usersListBtn = document.getElementById("users")
const videoListBtn = document.getElementById("videos")
const usersList = document.getElementById("users-list")
const videoList = document.getElementById("video-list")

const arbuzPath = "/static/assets/watermelon-3593669936.jpg"

let videoDir = null 
fetch("http://" + window.location.host + "/api/videodir")
  .then(response => response.json())
  .then(data => {
    videoDir = data.videodir
  })

let videoFiles = null
fetch("http://" + window.location.host + "/api/videofiles")
  .then(response => response.json())
  .then(data => {
    videoFiles = data.videofiles
    
    videoList.innerHTML = ""
    videoFiles.forEach(function(file) {
      const li = document.createElement("li")
      const a = document.createElement("a")
      a.textContent = cutFileName(file)
      a.className = "video-list__item-link"
      li.className = "video-list__item"
      a.id = file

      a.addEventListener("click", function(event) {
        event.preventDefault()
        const videoName = this.id
        const splitted = videoName.split(".")
        const type = splitted[1]
        player.src = videoDir + videoName
        player.type = "video/" + type
        console.log(player.src)
      })

      li.appendChild(a)

      videoList.appendChild(li)
    })
  })
  .catch(error => console.error("request error: ", error))


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
usersListBtn.addEventListener("click", changeSidePanelToUsers)
videoListBtn.addEventListener("click", changeSidePanelToVideos)

function updateUsers(usernames) {
  usersList.innerHTML = ""
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

    usersList.appendChild(li)
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

function changeSidePanelToUsers() {
  usersList.style.display = "flex"
  videoList.style.display = "none"
}

function changeSidePanelToVideos() {
  videoList.style.display = "flex"
  usersList.style.display = "none"
}

function cutFileName(file) {
  const numberMatch = file.match(/\d+/g)
  
  if (!numberMatch) return file
  
  const number = numberMatch[0]
  const indexOfNumber = file.lastIndexOf(number)    

  if (file.length <= 20) file 

  const truncatedString = file.slice(0, 20 - number.length - 3) + "..." + number
  
  return truncatedString
}
