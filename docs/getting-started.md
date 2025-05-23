# Getting started

Gokord is a Go 1.24+ library.

To use it, you must know Go and have at least the 1.24 installed.

Gokord relies heavily on [discordgo](https://github.com/bwmarrin/discordgo).
It is a simple wrapper of the official [Discord API](https://discord.com/developers/docs/).
If you have any questions, don't forget to check these!

## Installation

You can install Gokord with
```bash
$ go get -u github.com/anhgelus/gokord@latest
```

Replace `latest` by a specific tag or by a commit hash to get a specific version.

## Setuping configs and databases

You must setup configs before doing anything else.

If you want to disable Redis, set `gokord.UseRedis` to `false`.

Use `gokord.SetupConfigs(gokord.BaseConfig, []*cfgInfo) error`.
Check [config](/config) for more information.

Then, you can automigrate your schema with `gokord.DB.AutoMigrate(interface{}...) error`.
Check [databases](/databases/) for more information.

## Update

The bot automatically handles slash commands migration.
You must load an innovation JSON file with `gokord.LoadInnovationFromJson([]byte) (gokord.Innovation, error)`.
Check [innovation](/innovation) for more information.

## Initialize the bot

To create a new Gokord bot, you must use the struct `gokord.Bot`.
- `Token` field is the bot's token
- `Status` contains the bot's [statuses](/statuses)
- `Commands` contains the bot's [slash commands](/slash-commands/)
- `AfterInit` contains a function called after the initialisation of the bot (type `func (*discordgo.Session)`)
- `Version` contains the bot's [version](/innovation)
- `Innovations` contains the bot's [innovation](/innovation)
- `Intents` contains the bot's [intents](https://discord.com/developers/docs/events/gateway#gateway-intents)

Use `bot.start()` to start the bot.
This instruction will be blocked until the program is stopped.

For example, these instructions will start a new simple bot.
```go
innovation // contains bot's innovation
version := gokord.Version{
    Major: 1,
    Minor: 0,
    Patch: 0,
}

bot := gokord.Bot{
    Token: "token",
    Status: []*gokord.Status{},
    Commands: []gokord.CommandBuilder{},
    AfterInit: func(dg *discordgo.Session){},
    Innovations: innovations,
    Version: &version,
    Intents: discordgo.IntentsAllWithoutPrivileged,
}
bot.Start()
```
