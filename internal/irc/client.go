package irc

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nhooyr.io/websocket"
)

const (
	chatURL = "wss://irc-ws.chat.twitch.tv:443"

	CapMembership = ":twitch.tv/membership"
	CapTags       = ":twitch.tv/tags"
	CapCommands   = ":twitch.tv/commands"
)

type Client struct {
	ctx  context.Context
	conn *websocket.Conn
	rd   *bufio.Reader
}

func Dial(ctx context.Context, timeout time.Duration) (*Client, error) {
	opts := &websocket.DialOptions{
		HTTPClient: &http.Client{Timeout: timeout},
	}

	conn, _, err := websocket.Dial(ctx, chatURL, opts)
	if err != nil {
		return nil, err
	}

	return &Client{ctx, conn, bufio.NewReader(nil)}, nil
}

const linebreak = "\r\n"

func fmtMsg(msg string) string {
	words := strings.Fields(msg)
	return strings.Join(words, " ")
}

func fmtChannel(channel string) string {
	return strings.ToLower(channel)
}

func (c *Client) pong() error {
	return c.printf("PONG :tmi.twitch.tv")
}

func (c *Client) Login(nick, pass string) error {
	if err := c.printf("PASS %s", pass); err != nil {
		return err
	}

	return c.printf("NICK %s", nick)
}

func (c *Client) RegCaps(caps ...string) error {
	for _, cap := range caps {
		err := c.printf("CAP REQ %s", cap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) Join(channel string) error {
	return c.printf("JOIN #%s", fmtChannel(channel))
}

func (c *Client) PrivMsg(channel, msg string) error {
	return c.printf("PRIVMSG #%s :%s", fmtChannel(channel), fmtMsg(msg))
}

func (c *Client) PrivMsgReply(channel, msg, parentMsgID string) error {
	return c.printf("@reply-parent-msg-id=%s PRIVMSG #%s :%s", parentMsgID, fmtChannel(channel), fmtMsg(msg))
}
func (c *Client) Disconnect() error {
	return c.conn.Close(websocket.StatusNormalClosure, "client disconnect")
}

func (c *Client) printf(format string, args ...interface{}) error {
	return c.conn.Write(c.ctx, websocket.MessageText,
		[]byte(fmt.Sprintf(format, args...)))
}

func (c *Client) ReadMsg(ctx context.Context) (*Msg, error) {
	for {
		if c.rd.Buffered() == 0 {
			_, r, err := c.conn.Reader(ctx)
			if err != nil {
				return nil, err
			}

			c.rd.Reset(r)
		}

		line, err := c.rd.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if strings.HasPrefix(line, "PING") {
			err = c.pong()
			if err != nil {
				return nil, err
			}

			continue
		}

		line = strings.TrimSuffix(line, linebreak)

		return ParseMsg(line), nil
	}
}
