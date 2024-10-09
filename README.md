# AudioSwitch v1.0.0
## Quickly toggle audio devices

AudioSwitch allows you to simply create any keyboard shortcut in order to toggle any audio device
with a convenient overlay to always see the state of said device.
For now the only eligible devices are microphones.

### Project structure

The project is created with:
+ **[GoLang](https://go.dev)** for the backend logic
+ **C++** and **CGO** to interact with the Windows Core Audio APIs
+ **[SolidJS](https://www.solidjs.com/)** for the frontend
+ **[Wails3](https://v3alpha.wails.io/)** to link the backend and the frontend

### Project build

Dependencies:
+ **[GoLang](https://go.dev/doc/install)**, at least `v1.23.0`
+ **[NodeJS and NPM](https://nodejs.org/en/download/package-manager)**
+ A working **C/C++ Compiler**

Install steps:
+ First clone the repo:
  ```
  git clone https://github.com/nixpare/AudioSwitch
  ```
+ Install **Wails3**:
  ```
  go install -v github.com/wailsapp/wails/v3/cmd/wails3@latest
  ```
+ Build the project:
  ```
  wails3 task build:windows:prod:amd64
  ```

In order to run in DevMode:
```
wails3 dev
```
