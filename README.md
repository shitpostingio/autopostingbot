# Autoposting bot

Autoposting is a Telegram bot designed to post regularly on a Telegram channel. It is currently being employed on the [Shitpost - Shitposting.io](https://t.me/shitpost) channel.

The bot is designed to support multiple users, assuming they're unaware of what others are sending.

## Key features

- Feature-based image and video duplicate detection thanks to our custom Analysis API (to be released soon).

- Queue length aware posting algorithm in the _Shitpost_ edition.

- Based on Telegram's TDlib and not on the Bot API for a richer feature set.

- Easy addition of commands through the `command` package.

- Ability to add custom scheduling algorithms via the `posting/edition` package.

## How to build

In order to build (and run) the bot you will need to have both [Go](https://golang.org/dl/) and [TDlib](https://tdlib.github.io/td/build.html) installed. To simplify this process, you're free to use our custom Docker image that you can find in [this repository](https://github.com/shitpostingio/golang).

Once you have all the prerequisites, you will only need to run

```bash
make build
```

## Contributions

Contributions are welcome: suggest new features, add them yourself, translate the bot into new languages!

### To be updated
