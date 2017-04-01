# i3-recents

[![Build Status](https://travis-ci.org/Logiraptor/i3-recents.svg?branch=master)](https://travis-ci.org/Logiraptor/i3-recents)

[![Go Report Card](https://goreportcard.com/badge/github.com/Logiraptor/i3-recents)](https://goreportcard.com/report/github.com/Logiraptor/i3-recents)

This is a simple server which will listen for events when window focus changes and allow you to focus the last focused window. This replaces the majority of *my* use of `Alt+Tab` on other window managers; ie quickly switching between two windows.

## Setup

Download the compiled binary from the releases page, or build it yourself. There are no dependencies outside the Go standard library.

Add a line to your i3 config to start the server automatically.

```
exec_always --no-startup-id path/to/i3-recents
```

Add a line to your i3 config so your chosen key will trigger the client

```
bindsym Mod1+Tab exec path/to/i3-recents --back
```

