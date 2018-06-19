# ircDiscordRelay
![maintained](https://img.shields.io/badge/maintained-yes-brightgreen.svg)
![Programming Language](https://img.shields.io/badge/language-Go-orange.svg)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/iceflowRE/irc-discord-relay/blob/master/LICENSE.md)

[![Build Status](https://travis-ci.org/IceflowRE/irc-discord-relay.svg?branch=master)](https://travis-ci.org/IceflowRE/irc-discord-relay)

---

## Configuration

The configuration is a json formatted file.
One example configuration is part of this repository.

- **discord** - *map*, required - contains all settings related to the Discord part
    - **token** - *string*, required - the bot token
    - **channelId** - *string*, required - the channel id where the messages will be send
    - **sharing** - *list of string* - the actions which will be shared to IRC, an absence of this will imply to share all actions
        - **"message"** - discord message
- **irc** - *map*, required - contains all settings related to the IRC part
    - **channel** - *string*, required - the IRC channel where the bot connects to
    - **nick** - *string*, required - the IRC nick
    - **onConnection** - *list of string* - every string will send as a **raw** message
    - **sharing** - *list of string* - the actions which will be shared to Discord, an absence of this will imply to share all actions
        - **"message"** - discord message
        - **"me"** - `/me` messages
        - **"join"** - joining users
        - **"leaving"** - leaving users
        - **"nick"** - nick change
        - **"quit"** - quiting users (e.g. timeout)
        - **"kick"** - user got kicked
        - **"mode"** - a mode changed by another user

---

## Requirements

- Go (>= 1.9)

## Build

- `go get github.com/IceflowRE/irc-discord-relay`
- `go build -x github.com/IceflowRE/irc-discord-relay`

Alternativ
- `git clone https://github.com/IceflowRE/irc-discord-relay.git`
- `cd irc-discord-relay`
- `make`

## Update

- `go get -u github.com/IceflowRE/irc-discord-relay`
- build again

## Run
- create a `config.json`, with their needed values, look into the example for more
- place the config and the executable in one folder
- execute with `irc-discord-relay -c config.json`

---

## Web
https://github.com/IceflowRE/irc-discord-relay

## Credits
- Developer
    - Iceflower S
        - iceflower@gmx.de

### Third Party
- DiscordGo *by* Bruce Marriner ([bwmarrin](https://github.com/bwmarrin))
    - https://github.com/bwmarrin/discordgo
    - [BSD-3-Clause](https://github.com/bwmarrin/discordgo/blob/master/LICENSE)
- errors *by* Dave Cheney
    - https://github.com/pkg/errors
    - [BSD-2-Clause](https://github.com/pkg/errors/blob/master/LICENSE)
- go-ircevent *by* Thomas Jager ([thoj](https://github.com/thoj))
    - https://github.com/thoj/go-ircevent
    - [BSD-3-Clause](https://github.com/thoj/go-ircevent/blob/master/LICENSE)

Some Code taken from:
- Snowflower *by* Iceflower S ([Iceflower](https://gitlab.com/Iceflower))
    - (unpublished)
    - https://gitlab.com/Iceflower/snowflower
    - [GPL-3.0-or-later](https://gitlab.com/Iceflower/snowflower/blob/master/LICENSE.md)

## License
Copyright 2018 Iceflower S (iceflower@gmx.de)

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
