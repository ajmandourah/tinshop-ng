<div align="center">
<img alt="TinShop" src="./logo.png" width="50%"><br><br>
Your own personal shop right into tinfoil!<br><br>

</div>

# Why A Next Generation

This is a take on the original Tinshop repo originally by DblK. 
Current solutions have some issues:
- Language like python will have some issues when handling massive influx of data. 
- Can't handle large libraries
- lack of concurrent processing. 
- require certain naming schemes and no fallback to alternatives.
- will need a separate tool for renaming\identifying contents.
- support only nsp and nsz with no xci.
- slow as shit.

this will try to solve many of these issues.

# What is New Here? 

Some of the new features implemented so far:
- XCI extention support has been added.
- Content identification using your own Keys for content that does not meet the naming schemes. this is optional as instead of skipping unindentified content tinshop will now try to decrypt the data in the files.
- Tinshop-ng uses `fastwalk` instead of `filepath.walk`. Traversing directories uses multiple goroutines which work concurrently leading to faster processing of your data.
- Updated titlesdb now using Tinfoil's own data. 
- You can now rename unidentified content to one that meet the naming scheme's requiremnets for faster next time processing. 
- Implemented some features and functions from the popular `Switch-library-manager`
- Tested on rclone mount with more than 10K titles. processing with titles matching the name schemes in almost 10 sec.

# ⚠️ Disclaimer

This program **DOES NOT** encourage piracy at all!  
It was designed to reduce the time to download/install a game from the Nintendo eShop.  
In case you have a ADSL connection, to install latest `Zelda` ([14.4Gb](https://www.nintendo.com/games/detail/the-legend-of-zelda-breath-of-the-wild-switch/)) it can take ages!

On top of that, if you have bought a game on eShop like [Jump Force](https://www.bandainamcoent.com/news/jump-force-sunsetting-announcement), once it is not anymore on the shop how can you install it again?  
Using your personal NSP dump, with `tinfoil` and `tinshop` everything should be fine and fast!

# 🎮 Use

To proper use this software, here is the checklist:
- [ ] _Optional:_ A proper configured `config.yaml`
    - [ ] Copy/Paste [`config.example.yaml`](https://raw.githubusercontent.com/DblK/tinshop/master/config.example.yaml) to `config.yaml`
    - [ ] Comment/Uncomment parts in the config according to your needs
- [ ] _Optional:_ Games should have in their name `[ID][v0]` to be recognized. 
- [ ] _Optional:_ You can supply your own keys for content identification.
- [ ] Games extension should be `nsp` or `nsz` or `xci`

Now simply run it and add a shop inside tinfoil with the address setup in `config` (or `http://localIp:3000` if not specified).

# 🎉 Features

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
- [X] Content identification as fallback if naming schemes requirement are not fillfilled . It will try to identify the content and add it to your library.
- [X] Optional renaming of the identified content to an acceptable naming scheme so next time you start the server it will identify it faster

## 🏳️ Filtering

When you setup your shop inside `tinfoil` you can now add the following path:
- `multi` : Filter only multiplayer games
- `fr`, `en`, ... : Filter by languages
- `world` : All games without any filter (equivalent without path)

# Configuration

This is an example of the config.yaml file

```yaml
# Name of the host [optional]
host: tinshop.example.com

# Protocol (Can be http or https) [optional]
# If you use "https" then you should set up a reverse-proxy in front to handle tls
# And forward the port 443 to "yourIp:3000"
protocol: https

# keys [optional]
# This is a fallback in case parsing failed due to bad rename pattern. slower than parsing but more accurate in captureing content information.
# a full path to the key should be provided. [use /data as the path if you are using docker]
#keys: prod.keys

# renameFiles [optional] [keys must be present if true otherwise will do nothing]
# this will rename files that is misrenamed somehow by appending the title id and version to the end of the file. this may result in duplicate ids in the file name
# as it does not check for already present data.
# enabling this option will make processing much faster after the initial rename as there will be no need to decrypt the files any longer
renameFiles: false

# Port [optional]
# This affect the url to download games & the web server will run on that port (default: 3000).
# port: 3000

# Shop name [optional]
# This is used as title when trying to visit the shop with a non switch device
name: TinShop

# Tells if we are using a reverse proxy if front of tinshop [optional]
# This is used to rewrite correct url to download games when reverse proxy used
reverseProxy: false

# Welcome message on Switch [optional]
# The default message will be "Welcome to your own TinShop!"
welcomeMessage: "Welcome to your own TinShop!"
# If you want to disable the welcome message, set it to true [optional]
noWelcomeMessage: false

# All debug flags will be stored here
debug:
  # Display more information when connecting to nfs share
  nfs: false
  # Remove middleware security for retrieving index
  # DO NOT use in production (only for dev purpose)
  noSecurity: false
  # Display more information about ticket's verification
  ticket: false

# All actions related to NSP file will be stored here
nsp:
  # Tells if tinshop should verify the ticket inside NSP to ensure no issue with install
  # This will make processing really slow as every file will be opened for verification 
  checkVerified: false

# All sources where we should look for games
# If this section is commented out, then the directory "games" will be looked at
sources:
  # Local mounted path [optional]
  directories:
    - /my/full/path/to/games
    - ./games

  # NFS Shares [optional]
  nfs:
    - host:sharePath/to/game/files

# All security information will be stored here
security:
  # List of theme to be banned with security
  # Be aware that this should be string (do not forget quotes)
  # You can find the theme of a switch in the log upon access
  bannedTheme:
    - "0000000000000000000000000000000000000000000000000000000000000000"
  # List of switch uid to whitelist
  # If enabled then only switch in this area will be listed
  # You can find the uid of a switch in the log upon access
  whitelist:
    - TESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTEST
  # List of switch uid to blacklist
  # Block access to all switch present in this list
  # You can find the uid of a switch in the log upon access
  blacklist:
    - NOACCESSNOACCESSNOACCESSNOACCESSNOACCESSNOACCESSNOACCESSNOACCESS
  # Endpoint to which a query will be sent to verify user/password/uid to
  # Headers sent :
  # - Authorization: same as sent by switch
  # - Device-Id: Switch fingerprint
  # Response with status code other than 200 will be treated as failure
  forwardAuth: https://auth.tinshop.com/switch

# This section describe all custom title db to show up properly in tinfoil
customTitledb:
  # Id of the entry
  "060000BADDAD0000":
    id: "050000BADDAD0000"
    name: "Tinfoil"
    region: "US"
    releaseDate: 20180801
    description: "Nintendo Switch Title Manager"
    size: 14000000
    iconUrl: ""
```

# 🐋 Docker

_WIP_ not ready yet.

To run with [Docker](https://docs.docker.com/engine/install/), you can use this as a starting `cli` example:

`docker run -d --restart=always -e TINSHOP_SOURCES_DIRECTORIES=/games -e TINSHOP_WELCOMEMESSAGE="Welcome to my Tinshop!" -v /local/game/backups:/games -p 3000:3000 ghcr.io/dblk/tinshop:latest`

This will run Tinshop on  `http://localhost:3000` and persist across reboots!

If `docker compose` is your thing, then start with this example:

```yaml
version: '3.9'
services:
  tinshop:
    container_name: tinshop
    image: ghcr.io/ajmandourah/tinshop-ng:main
    restart: always
    ports:
      - 3000:3000
    environment:
      - TINSHOP_SOURCES_DIRECTORIES=/games
      - TINSHOP_WELCOMEMESSAGE=Welcome to my Tinshop!
    volumes:
      - /media/switch:/games
      - /path/to/config:/data  #this is where config.yaml and titles json file will live. added keys here should also be modifed in config.yaml file as /data/prod.keys 
```
All of the settings in the `config.yaml` file are valid Environment Variables. They must be `UPPERCASE` and prefixed by `TINSHOP_`. Nested properties should be prefixed by `_`. Here are a few examples:

| ENV_VAR                      | `config.yaml` entry | Default Value                  | Example Value                     |
|------------------------------|---------------------|--------------------------------|-----------------------------------|
| TINSHOP_HOST                 | host                | `<empty>`                      | `tinshop.example.com`             |
| TINSHOP_PROTOCOL             | protocol            | `http`                         | `https`                           |
| TINSHOP_NAME                 | name                | `TinShop`                      | `MyShop`                          |
| TINSHOP_REVERSEPROXY         | reverseProxy        | `false`                        | `true`                            |
| TINSHOP_WELCOMEMESSAGE       | welcomeMessage      | `Welcome to your own TinShop!` | `Welcome to my shop!`             |
| TINSHOP_NOWELCOMEMESSAGE     | noWelcomeMessage    | `false`                        | `true`                            |
| TINSHOP_DEBUG_NFS            | debug.nfs           | `false`                        | `true`                            |
| TINSHOP_DEBUG_NOSECURITY     | debug.nosecurity    | `false`                        | `true`                            |
| TINSHOP_DEBUG_TICKET         | debug.ticket        | `false`                        | `true`                            |
| TINSHOP_NSP_CHECKVERIFIED    | nsp.checkVerified   | `false`                        | `true`                            |
| TINSHOP_SOURCES_DIRECTORIES  | sources.directories | `./games`                      | `/games /path/two /path/three`    |
| TINSHOP_SOURCES_NSF          | sources.nfs         | `null`                         | `192.168.1.100:/path/to/games`    |
| TINSHOP_SECURITY_BANNEDTHEME | sources.bannedTheme | `null`                         | `THEME1 THEME2 THEME3`            |
| TINSHOP_SECURITY_WHITELIST   | sources.whitelist   | `null`                         | `NSWID1 NSWID2 NSWID3`            |
| TINSHOP_SECURITY_BLACKLIST   | sources.blacklist   | `null`                         | `NSWID4 NSWID5 NSWID6`            |
| TINSHOP_SECURITY_FORWARDAUTH | sources.forwardAuth | `null`                         | `https://auth.tinshop.com/switch` |

## 🥍 Want to do cross-build generation?

Wanting to generate all possible os binaries (macOS, linux, windows) with all architectures (arm, amd64)?  
Here is the command `goreleaser release --snapshot --skip-publish --rm-dist`.

## Tips for faster processing especially when using cloud shares ie Rclone

If you are going to use rclone or similar cloud storage solutions as your source of content here are some tips:
- make sure all your contents are in one folder without them being in subfolders. As for every subfolder rclone will need to fetch a list of every file in it. this can take long time especially with large number of folders.
- either make sure your content matches the naming format. you can also run switch-library-manager on the folder with renaming option enabled to insure these matches. Tinshop-ng with decryption enabled will need to read every file for decryption which will take more time processing .
- enabling cache will give some boost to processing.
  
# 👂🏻 Q & A

## Why use this instead of `X` (NUT or others software)?

<details>
<summary>Answer</summary>

It's dead simple, and no dependencies! It's just a single small executable.  
Easier to install games without connecting switch or by updating SD card (Nightmare if you are on macOS).

The upcoming features will also be a huge advantage against others software.
</details>

## Where do I put my games?

<details>
<summary>Answer</summary>

By default, `TinShop` will look into the `games` directory relative to `tinshop` executable.

However in the `config.yaml` file, you can change this.  
In the `sources` section, you can have the following:
- `directories`: List of directories where you put your games
- `nfs`: List of NFS shares that contains your games
</details>

## Can I set up a `https` endpoint?

<details>
<summary>Answer</summary>

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

### Example for traefik

To work with [`traefik`](https://traefik.io/), you need to put in your Dynamic Configuration something similar to this:

```yaml
http:
  routers:
    service: tinshop
    rule: Host(`tinshop.example.com`)
    entryPoints: websecure # Could be web if not using https

  services:
    tinshop:
      loadBalancer:
        servers:
          - url: http://192.168.1.2:3000
```

and your `config.yaml` as follow:

```yaml
host: tinshop.example.com
protocol: http
port: 3000
reverseProxy: true
```

If you want to have HTTPS, ensure `traefik` can handle it (it will with Let's Encrypt) and use protocol `https` in the config.

For more details on Traefik + Let's Encrypt, [click here](https://doc.traefik.io/traefik/https/acme/).
</details>

## How can I add a `basic auth` to protect my shop?

<details>
<summary>Answer</summary>

TinShop **does** handle basic auth but not by itself.  
You should look for `forwardAuth` in the `config.yaml` to set the endpoint that will handle the authentication in behalf of TinShop.

In the future, a proper user management will be incorporated into TinShop to handle it.

In addition, for other type of protection, you can whitelist/blacklist your own switch and this will do the trick.
</details>

## I have tons of missing title displayed in `tinfoil`, what should I do?

<details>
<summary>Answer</summary>

First, download and replace the latest [`titles.US.en.json`](https://github.com/AdamK2003/titledb/releases/download/latest/titles.US.en.json) available (or delete it, it will be automatically downloaded at startup).  
If this does not solve your issue, then you should use custom titledb entry to describe those which are missing.
</details>

## Why I still get `NCA signature verification failed` error in `tinfoil` and nothing in `tinshop`?

<details>
<summary>Answer</summary>

The current implementation to verify the NSP/NSZ are basic and based on the Ticket information.  
So you might still get some error about signature failed even with `checkVerified` enabled.

Maybe later, this feature will be enhanced to add additional checks on game files (PR Welcome!).
</details>

## `tinfoil` does not display the name of the game but the file name, what should I do?

<details>
<summary>Answer</summary>

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
</details>

# 🙏 Credits

- [DblK](https://github.com/DblK) for the original effort on the original repo @DblK
- [Trembon](https://github.com/trembon) outstanding work on [switch-library-manager](https://github.com/trembon/switch-library-manager)

# Todo
- ~~A new container with ability to edit config.yaml~~
- ~~workflow edit~~

