# Microbot

Chat bot for twitch.

## Run

You need to provide `config` and `cred` files. By default these are `./config.yml` and `creds.yml`.

### Config

Config is expected to have following structure:

```yaml
debug: true # bool
channels: # list of channels
  - name: <channel-1> # name of channel
    chat: # object that contains settings for specific channel
      middlewares: # optional field that provides ability to filter/process messages
        - type: <middleware-type> # middlewares are called starting from the first
          settings:
            first: value-1
            second: value-2
        - type: <middleware-type-2> # this middleware will be called after the first one
          settings:
            third: value-3
            fourth: value-4
      rewards:
        - key: <reward-uuid>
          middlewares: # we can add middlewares for specific rewards and commands
          action:
            type: <action-type>
            settings: # collection of action-specific settings
              fifth: value-5
              sixth: value-6
      commands:
        - key: help # specifies what goes after "!", in this case it will be triggered on "!help"
          middlewares:
          action:
            type: <action-type>
            settings:
              seventh: value-7
              eighth: value-8
```

### Creds

Example for creds:

```yaml
twitchuser: <bot_username>
twitchpass: oauth:<token> 
riotapikey: RGAPI-<key> # this field is optional and required only for interacting with riot API
```

Easiest way to generate token is to use this [tool](https://twitchapps.com/tmi).

To get more info visit [twitch docs](https://dev.twitch.tv/docs/irc).

### Actions

TBD

### Middlewares

TBD
