# stormy
Stormy is a bot that can be used to brainstorm and collect ideas in a Discord server.
Ideas can be collected in a channel and can get up- or down-voted.
Admins and authorized users can then repost those message to other channels by clicking one of the repost reactions.

## Configuration
Place a `stormy.yml` in the same folder as the stormy executable or specify a path to the config via the `-c` flag.
Because we use [viper](https://github.com/spf13/viper) for configuration, you are not limited to yaml files, but may also write your configuration
in any major configuration file format, such as json or toml.
Below is a sample config with all fields filled:

```yml
token: Bot your token
status: dnd # default: online
activity: watching # default: playing
activityName: a mindmap # if empty, no activity will be displayed
dateFormat: 01/02/2006 # Default January 2, 2006
timeFormat: 15:04 # default: 3:04 PM
location: Europe/Berlin # defaults to system time zone
channelConfigs: # array of channel configurations
  - channelID: 123 # the id of the watched channel
    autoReactions: # array of emojis that will always trigger a reaction
      - 😁
    # array of emojis that will be reacted with only, if found in the message
    scanReactions:
      - 1️⃣
    repostReactions: # array of repost reactions
      - target: 456 # the channel that shall be posted to
        reaction: 🍇 # the reaction that triggers a repost
        # the message that will be sent in target
        # available variables are
        # - Message - the original message
        # - MessageQuoted - the original message, but quoted
        # - Author - the name of the author without descriptor
        # - AuthorMention - a mention of the author
        # - Crossposter - the name of the user who authorized the crosspost
        # - CrossposterMention - a mention of the user who authorized the crosspost
        # - SourceChannel - a mention of the original channel
        # - Time - the time the original message was sent
        # - Date - the date the original message was sent
        message: "{{.MessageQuoted}}\n\n*by {{.Author}}*"
        # defines users that are authorize to repost, admins can always repost
        rigths:
          userIDs: # array of ids of users that can trigger a repost
            - 123
          roleIDs: # array of ids of roles whose owners can repost
            - 456
```
