# i3-recents

[![Build Status](https://travis-ci.org/Logiraptor/i3-recents.svg?branch=master)](https://travis-ci.org/Logiraptor/i3-recents)

This is a simple server which will listen for events when window focus changes and allow you to focus the last focused window. This replaces the majority of *my* use of `Alt+Tab` on other window managers; ie quickly switching between two windows.

## Setup

Clone this repo, and start the server

```
go run main.go
```

Add a hotkey to your i3 config to bind a key to make a request to the locally running server.

```
bindsym Mod1+Tab exec curl http://localhost:7364
```

