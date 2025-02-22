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
  .catch(error => console.error(error))

let videoFiles = null
fetch("http://" + window.location.host + "/api/videofiles")
  .then(response => response.json())
  .then(data => {
    videoFiles = data.videofiles
    player.src = videoDir + videoFiles[0]
    player.type = "video/" + getVideoType(videoFiles[0])
    
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
        let type = getVideoType(this.id)
        let videoName = this.id 
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
      console.log(data.time)
      player.currentTime = data.time
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
    time: player.currentTime,
    action: "play"
  } 
  const data = JSON.stringify(msg)
  ws.send(data)
}

function changeSidePanelToUsers() {
  usersList.style.display = "flex"
  usersListBtn.style.backgroundColor = "var(--fifth)"
  videoList.style.display = "none"
  videoListBtn.style.backgroundColor = "var(--third)"
}

function changeSidePanelToVideos() {
  videoList.style.display = "flex"
  videoListBtn.style.backgroundColor = "var(--fifth)"
  usersList.style.display = "none"
  usersListBtn.style.backgroundColor = "var(--third)"
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

function getVideoType(name) {
    const splitted = name.split(".")
    return splitted[1]
}
