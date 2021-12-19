<p align="center">
<p align="center">
    <img alt="TinShop" src="./logo.png" width="50%">  
</p>
<p align="center">
    Your own personal shop right into tinfoil!
</p>

[![golangci-lint](https://github.com/DblK/tinshop/actions/workflows/golangci-lint.yml/badge.svg?event=release)](https://github.com/DblK/tinshop/actions/workflows/golangci-lint.yml)
[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/dblk/tinshop.svg)](https://github.com/dblk/tinshop)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/dblk/tinshop/v0.0.1)
[![GoReportCard example](https://goreportcard.com/badge/github.com/dblk/tinshop)](https://goreportcard.com/report/github.com/dblk/tinshop)
[![GitHub release](https://img.shields.io/github/release/dblk/tinshop.svg)](https://GitHub.com/dblk/tinshop/releases/)

# Disclaimer

This program **DOES NOT** encourage piracy at all!  
It was designed to reduce the time to download/install a game from the Nintendo eShop.  
In case you have a ADSL connection, to install latest `Zelda` ([14.4Gb](https://www.nintendo.com/games/detail/the-legend-of-zelda-breath-of-the-wild-switch/)) it can take ages!

On top of that, if you have bought a game on eShop like [Jump Force](https://www.bandainamcoent.com/news/jump-force-sunsetting-announcement), once it is not anymore on the shop how can you install it again?  
Using your personal NSP dump, with `tinfoil` and `tinshop` everything should be fine and fast!

# Use

To proper use this software, here is the checklist:
- [ ] _Optional:_ A proper configured `config.yaml`
    - [ ] Copy/Paste [`config.example.yaml`](https://github.com/DblK/tinshop/blob/master/config.example.yaml) to `config.yaml`
    - [ ] Comment/Uncomment parts in the config according to your needs
- [ ] Games should have in their name `[ID][v0]` to be recognized
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

## How can I add a `basic auth` to protect my shop?

TinShop **does not** implement basic auth by itself.  
You should configure it inside your reverse proxy.

For other type of protection, you can whitelist your own switch and this will do the trick.

## I have tons of missing title displayed in `tinfoil`, what should I do?

First, download and replace the latest [`titles.US.en.json`](https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json) available (or delete it, it will be automatically downloaded at startup).  
If this does not solve your issue, then you should use custom titledb entry (__*__) to describe those which are missing.

__*__ Feature not yet implemented!

# Credits

I would like to give back thanks to the people who helped me with or without knowing!
- [Bogdan Rosu Creative](https://www.iconfinder.com/icons/353439/basket_purse_shopping_cart_ecommerce_shop_buy_online_icon) for his shop icon.
- [Dono](https://github.com/Donorhan) for his support and tests.
- [AdamK2003](https://github.com/AdamK2003/titledb) for his up-to-date [`titles.US.en.json`](https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json) and his answers on discord.
