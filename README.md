<div align="center">
<img alt="TinShop" src="./logo.png" width="50%"><br><br>
Your own personal shop right into tinfoil!<br><br>

[![golangci-lint](https://github.com/DblK/tinshop/actions/workflows/golangci-lint.yml/badge.svg?branch=master)](https://github.com/DblK/tinshop/actions/workflows/golangci-lint.yml)
[![test](https://github.com/DblK/tinshop/actions/workflows/ginkgo.yml/badge.svg?branch=master)](https://github.com/DblK/tinshop/actions/workflows/ginkgo.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/DblK/tinshop.svg)](https://github.com/DblK/tinshop)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/DblK/tinshop)
[![GoReportCard](https://goreportcard.com/badge/github.com/DblK/tinshop)](https://goreportcard.com/report/github.com/DblK/tinshop)
[![GitHub release](https://img.shields.io/github/release/DblK/tinshop.svg)](https://GitHub.com/DblK/tinshop/releases/)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)
</div>

# Disclaimer

This program **DOES NOT** encourage piracy at all!  
It was designed to reduce the time to download/install a game from the Nintendo eShop.  
In case you have a ADSL connection, to install latest `Zelda` ([14.4Gb](https://www.nintendo.com/games/detail/the-legend-of-zelda-breath-of-the-wild-switch/)) it can take ages!

On top of that, if you have bought a game on eShop like [Jump Force](https://www.bandainamcoent.com/news/jump-force-sunsetting-announcement), once it is not anymore on the shop how can you install it again?  
Using your personal NSP dump, with `tinfoil` and `tinshop` everything should be fine and fast!

# Use

To proper use this software, here is the checklist:
- [ ] _Optional:_ A proper configured `config.yaml`
    - [ ] Copy/Paste [`config.example.yaml`](https://raw.githubusercontent.com/DblK/tinshop/master/config.example.yaml) to `config.yaml`
    - [ ] Comment/Uncomment parts in the config according to your needs
- [ ] Games should have in their name `[ID][v0]` to be recognized
- [ ] Games extension should be `nsp` or `nsz`
- [ ] Retrieve binary from [latest release](https://github.com/DblK/tinshop/releases) or build from source (See `Dev` section below)

Now simply run it and add a shop inside tinfoil with the address setup in `config` (or `http://localIp:3000` if not specified).

# Features

Here is the list of all main features so far:
- [X] Automatically download `titles.US.en.json` if missing at startup
- [X] Basic protection from forged queries (should allow only tinfoil to use the shop)
- [X] Serve from several mounted directories
- [X] Serve from several network directories (Using NFS)
- [X] Display a webpage for forbidden devices
- [X] Auto-refresh configuration on file change
- [X] Add the possibility to whitelist or blacklist a switch
- [X] Add the possibility to ban theme
- [X] You can specify custom titledb to be merged with official one
- [X] Auto-watch for mounted directories
- [X] Add filters path for shop
- [X] Simple ticket check in NSP/NSZ (based on titledb file)
- [X] Collect basic statistics
- [X] An API to query information about your shop
- [X] Handle Basic Auth from Tinfoil through Forward Auth Endpoint

## Filtering

When you setup your shop inside `tinfoil` you can now add the following path:
- `multi` : Filter only multiplayer games
- `fr`, `en`, ... : Filter by languages
- `world` : All games without any filter (equivalent without path)

# Dev or build from source

I suggest to use a tiny executable [gow](https://github.com/mitranim/gow) to help you during the process (hot reload, etc..).  
For example I use the following command to develop `gow -c run .`.

If you want to build `TinShop` from source, please run `go build`.

And then, simply run `./tinshop`.

## Want to do cross-build generation?

Wanting to generate all possible os binaries (macOS, linux, windows) with all architectures (arm, amd64)?  
Here is the command `goreleaser release --snapshot --skip-publish --rm-dist`.

Dead simple, thanks to Golang!

## Changing the structure of an interface?

If you change an interface (or add a new one), do not forget to execute `./update_mocks.sh` to generate up-to-date mocks for tests.

Do not forget to install `mockgen` first:
```sh
go install github.com/golang/mock/mockgen@v1.6.0
```

## What to launch tests?

You can run `ginkgo -r` for one shot or `ginkgo watch -r` during development.  
Note: you can add `-cover` to have an idea of the code coverage.
# Roadmap

You can see the [roadmap here](https://github.com/DblK/tinshop/projects/1).

If you have any suggestions, do not hesitate to participate!

# Q & A

## Why use this instead of `X` (NUT or others software)?

It's dead simple, and no dependencies! It's just a single small executable.  
Easier to install games without connecting switch or by updating SD card (Nightmare if you are on macOS).

The upcoming features will also be a huge advantage against others software.

## Where do I put my games?

By default, `TinShop` will look into the `games` directory relative to `tinshop` executable.

However in the `config.yaml` file, you can change this.  
In the `sources` section, you can have the following:
- `directories`: List of directories where you put your games
- `nfs`: List of NFS shares that contains your games


## Can I set up a `https` endpoint?

Yes, you can!  
Use a reverse proxy (like [traefik](https://github.com/traefik/traefik), [caddy](https://github.com/caddyserver/caddy), nginx...) to do tls termination and forward to your instance on port `3000`.

### Example for caddy

To work with [`caddy`](https://caddyserver.com/), you need to put in your `Caddyfile` something similar to this:

```Caddyfile
tinshop.example.com:80 {
	reverse_proxy 192.168.1.2:3000
}
```

and your `config.yaml` as follow:

```yaml
host: tinshop.example.com
protocol: http
port: 3000
reverseProxy: true
```

If you want to have HTTPS, ensure `caddy` handle it (it will with Let's Encrypt) and change `https` in the config and remove `:80` in the `Caddyfile` example.

## How can I add a `basic auth` to protect my shop?

TinShop **does** handle basic auth but not by itself.  
You should look for `forwardAuth` in the `config.yaml` to set the endpoint that will handle the authentication in behalf of TinShop.

In the future, a proper user management will be incorporated into TinShop to handle it.

In addition, for other type of protection, you can whitelist/blacklist your own switch and this will do the trick.

## I have tons of missing title displayed in `tinfoil`, what should I do?

First, download and replace the latest [`titles.US.en.json`](https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json) available (or delete it, it will be automatically downloaded at startup).  
If this does not solve your issue, then you should use custom titledb entry to describe those which are missing.

## Why I still get `NCA signature verification failed` error in `tinfoil` and nothing in `tinshop`?

The current implementation to verify the NSP/NSZ are basic and based on the Ticket information.  
So you might still get some error about signature failed even with `checkVerified` enabled.

Maybe later, this feature will be enhanced to add additional checks on game files (PR Welcome!).

## `tinfoil` does not display the name of the game but the file name, what should I do?

You must follow the naming convention for the games as follow:  
`[gameId][version].(nsp/nsz)`

`gameId` should be a 16 characters long string.

For example, those are invalid:
- `0000000000000000 [v0].nsp`
- `[0000000000000000].nsp`
- `[0000000000000000][v0].xxx`

Those are valid:
- `[0000000000000000] [v0].nsp`
- `[0000000000000000][v131072].nsz`
- `My Saved Game [0000000000000000] [v0].nsp`
- `Awesome title [0000000000000000][v0] (15Gb).nsz`

# Credits

I would like to give back thanks to the people who helped me with or without knowing!
- [Bogdan Rosu Creative](https://www.iconfinder.com/icons/353439/basket_purse_shopping_cart_ecommerce_shop_buy_online_icon) for his shop icon.
- [Dono](https://github.com/Donorhan) for his support and tests.
- [AdamK2003](https://github.com/AdamK2003/titledb) for his up-to-date [`titles.US.en.json`](https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json) and his answers on discord.
- [nxdumptool](https://github.com/DarkMatterCore/nxdumptool) for the information taken of NSP format
