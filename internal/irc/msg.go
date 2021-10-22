package irc

import (
	"strconv"
	"strings"
	"time"
)

const (
	MsgTypePrivMsg         = "PRIVMSG"
	MsgTypeLeave           = "PART"
	MsgTypeJoin            = "JOIN"
	MsgTypeClearChat       = "CLEARCHAT"
	MsgTypeClearMsg        = "CLEARMSG"
	MsgTypeHostTarget      = "HOSTTARGET"
	MsgTypeNotice          = "NOTICE"
	MsgTypeReconnect       = "RECONNECT"
	MsgTypeRoomState       = "ROOMSTATE"
	MsgTypeUserNotice      = "USERNOTICE"
	MsgTypeUserState       = "USERSTATE"
	MsgTypeGlobalUserState = "GLOBALUSERSTATE"

	hostName    = "tmi.twitch.tv"
	hostNameLen = len(hostName)
)

func ParseMsg(rawMsg string) *Msg {
	var msg Msg

	msg.Raw = rawMsg

	msg.Tags, rawMsg = parseTags(rawMsg)

	msg.User, rawMsg = parseUser(rawMsg)

	msg.Type, rawMsg = parseType(rawMsg)

	msg.Channel, rawMsg = parseChannel(rawMsg)

	msg.Text = parseText(rawMsg)

	return &msg
}

func parseTags(rawMsg string) (map[string]string, string) {
	if strings.IndexByte(rawMsg, '@') != 0 {
		return nil, rawMsg
	}

	var rawTags, tail string
	if endOff := strings.IndexByte(rawMsg, ' '); endOff != -1 {
		rawTags = rawMsg[1:endOff]
		tail = rawMsg[endOff+1:]
	} else {
		rawTags = rawMsg[1:]
	}

	tags := map[string]string{}

	pairs := strings.Split(rawTags, ";")

	for _, pair := range pairs {
		sepOff := strings.IndexByte(pair, '=')
		if sepOff == -1 {
			continue
		}

		tagName := string(pair[:sepOff])
		tagValue := string(pair[sepOff+1:])

		tags[tagName] = strings.ReplaceAll(tagValue, "\\s", " ")
	}

	return tags, tail
}

func parseUser(rawMsg string) (string, string) {
	if strings.IndexByte(rawMsg, ':') != 0 {
		return "", rawMsg
	}

	hostOff := strings.Index(rawMsg, hostName)
	if hostOff == -1 {
		return "", rawMsg
	}

	var tail string
	if endOff := strings.IndexByte(rawMsg, ' '); endOff != -1 {
		tail = rawMsg[endOff+1:]
	}

	if hostOff == 1 {
		return "", tail
	}

	atOff := strings.IndexByte(rawMsg, '@')

	var user string
	if atOff != -1 && atOff < hostOff-1 {
		user = string(rawMsg[atOff+1 : hostOff-1])
	}

	return user, tail
}

func parseType(rawMsg string) (string, string) {
	endOff := strings.IndexByte(rawMsg, ' ')

	var tp, tail string
	if endOff != -1 {
		tp = rawMsg[:endOff]
		tail = rawMsg[endOff+1:]
	} else {
		tp = rawMsg
	}

	return tp, tail
}

func parseChannel(rawMsg string) (string, string) {
	if strings.IndexByte(rawMsg, '#') != 0 {
		return "", rawMsg
	}

	endOff := strings.IndexByte(rawMsg, ' ')

	var channel, tail string
	if endOff != -1 {
		channel = rawMsg[1:endOff]
		tail = rawMsg[endOff+1:]
	} else {
		channel = rawMsg[1:]
	}

	return channel, tail
}

func parseText(rawMsg string) string {
	off := strings.IndexByte(rawMsg, ':')
	if off == -1 {
		return ""
	}

	return rawMsg[off+1:]
}

type Msg struct {
	Tags    map[string]string
	Type    string
	Channel string
	Text    string
	User    string
	Raw     string
}

// dead code below

const (
	timeFormat = "15:04"
)

func (msg *Msg) String() string {
	switch msg.Type {
	case MsgTypePrivMsg:
		return renderPrivMsg(msg.Channel, msg.Text, msg.Tags)
	}
	return ""
}

func renderPrivMsg(channel, text string, tags map[string]string) string {
	var sb strings.Builder

	sb.WriteByte('<')
	sb.WriteString(channel)
	sb.WriteString(">\n")

	millis, _ := strconv.ParseInt(tags["tmi-sent-ts"], 10, 64)
	if t := time.Unix(millis/1000, 0); !t.IsZero() {
		sb.WriteString(t.Format(timeFormat))
	} else {
		sb.WriteString("--:--")
	}

	sb.WriteByte(' ')

	if strings.Contains(tags["badges"], "predictions/blue") {
		sb.WriteString("[blue]")
	}
	if strings.Contains(tags["badges"], "predictions/pink") {
		sb.WriteString("[pink]")
	}
	if strings.Contains(tags["badges"], "moderator") {
		sb.WriteString("[mod]")
	}
	if strings.Contains(tags["badges"], "vip") {
		sb.WriteString("[vip]")
	}
	if strings.Contains(tags["badges"], "subscriber") {
		sb.WriteString("[sub]")
	}
	if strings.Contains(tags["badges"], "premium") {
		sb.WriteString("[prime]")
	}
	if strings.Contains(tags["badges"], "glhf-pledge") {
		sb.WriteString("[glhf]")
	}

	sb.WriteString(tags["display-name"])
	sb.WriteString(":\n")
	sb.WriteString(fmtMsg(text))
	sb.WriteString("\n")

	return sb.String()
}
