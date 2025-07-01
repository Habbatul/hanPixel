
# hanPixel

**🕹️ Play Now (WebAssembly Build):** [Click Here](https://habbatul.github.io/hanPixel/)

## 📺 About the Project

**hanPixel** is a pixel-art portfolio game created using [Ebiten](https://ebiten.org/), a 2D game library for Go. This project was built from scratch to showcase fundamental game mechanics. **While some implementations are not fully modular** (maybe it's hard to read and change), the project serves as a demonstration of basic foundational game dev concepts. Now all players can do multiplayer via peer-to-peer, but since STUN is used, there is a possibility that some players may not be able to establish a direct connection due to strict NAT types.

### ✨ Features

* **Online Multiplayer (Peer-to-Peer)**
  Utilizes WebRTC with STUN for peer-to-peer multiplayer. A signaling server is required, accessible at [hanPixel_SignalingServer](https://github.com/Habbatul/hanPixel_SignalingServer).

* **Collision Detection.**
  Detects collisions between the player, NPCs, and obstacles.

* **Pixel-Perfect Obstacle Collisions.**
  Ensures accurate interaction with obstacles using pixel-based collision logic.

* **Render Order Logic.**
  Determines front and back object rendering dynamically for depth simulation.

* **Frame-Based Animation.**
  NPCs animate based on predefined frame sequences.

* **Zoom Factor Support.**
  Zoom level can be adjusted via code (note: current implementation could be improved).

* **Textboxes on Interaction.**
  Displays dialogue or information upon collision with certain objects or NPCs.

* **Touchscreen Support (Mobile).**
  Touch input enabled for better mobile gameplay experience.

* **Keyboard & Mouse Support (Desktop).**
  Full control support for desktop environments.


### ✨ Tech Stack

- pion/webrtc
- coder/websocket
- hajimehoshi/ebiten

---

## 🧪 Run Locally on Your PC

### 📦 Install Dependencies

Before running the project, make sure all Go module dependencies are fetched:

```bash
go mod tidy
```

### 🌐 Run in Web Browser Locally

To serve the project locally using WebAssembly:

```bash
go run github.com/hajimehoshi/wasmserve@latest ./path/to/thisProject
```

---

Feel free to explore, fork, or contribute. This project is a sandbox for experimenting and learning about 2D game development with Go and Ebiten!

---

