# Autoposting bot

<p align="center"> 
   <img src="shitpost.png" width="60%">
</p>

Autoposting is a Telegram bot designed to post regularly on a Telegram channel. It is currently being employed on the [Shitpost - Shitposting.io](https://t.me/shitpost) channel.

The bot is designed to support multiple users, assuming they're unaware of what others are sending.

## Key features

- Feature-based image and video duplicate detection thanks to [our custom Analysis API](https://github.com/shitpostingio/analysis-api).

- Queue length aware posting algorithm.

- Based on Telegram's TDlib and not on the Bot API for a richer feature set.

- Easy addition of commands through the `command` package.

- Ability to add custom scheduling algorithms via the `posting/algorithm` package.

- Multiple language support

## Languages

Currently supported languages are:

- English
- Italian
- Russian (thanks to [Rutori](https://github.com/Rutori))
- Brazilian Portuguese (thanks to [Hellstrike12](https://github.com/hellstrike12))
- Spanish (thanks to [Alvaro P.](https://t.me/Tag_if_magic_stones_dont_drop))
- Arabic (thanks to [TMDR](https://github.com/TMDR))

## How to build

In order to build (and run) the bot you will need to have both [Go](https://golang.org/dl/) and [TDlib](https://tdlib.github.io/td/build.html) installed. To simplify this process, you're free to use our custom Docker image that you can find in [this repository](https://github.com/shitpostingio/golang).

Once you have all the prerequisites, you will only need to run

```bash
make build
```

Pre-built artifacts are available in the GitHub Actions of the project: [https://github.com/shitpostingio/autopostingbot/actions](https://github.com/shitpostingio/autopostingbot/actions)

## Contributions

Contributions are welcome: suggest new features, add them yourself, translate the bot into new languages!

