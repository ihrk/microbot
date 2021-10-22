package bot

import (
	"context"
	"fmt"
	"time"

	"github.com/ihrk/microbot/internal/irc"
)

type response struct {
	channel     string
	text        string
	parentMsgID string
}

type Handler interface {
	Serve(cs *Sender)
}

type Middleware func(Handler) Handler

type Server struct {
	c *irc.Client
	h Handler
}

func NewServer(c *irc.Client, h Handler) *Server {
	return &Server{
		c: c,
		h: h,
	}
}

func (srv *Server) serve(respCh <-chan response) {
	for resp := range respCh {
		if resp.parentMsgID == "" {
			srv.c.PrivMsg(resp.channel, resp.text)
		} else {
			srv.c.PrivMsgReply(resp.channel, resp.text, resp.parentMsgID)
		}
	}
}

const msgBuf = 10

func (srv *Server) ListenAndServe(ctx context.Context) error {
	respCh := make(chan response, msgBuf)

	go srv.serve(respCh)

	defer close(respCh)

	for {
		msg, err := srv.c.ReadMsg(ctx)
		if err != nil {
			return err
		}

		go srv.h.Serve(NewSender(msg, respCh))
	}
}

type Sender struct {
	Msg    *irc.Msg
	respCh chan<- response
}

func NewSender(msg *irc.Msg, respCh chan<- response) *Sender {
	return &Sender{
		Msg:    msg,
		respCh: respCh,
	}
}

func (s *Sender) RewardID() string {
	return s.Msg.Tags["custom-reward-id"]
}

func (s *Sender) Send(text string) {
	s.respCh <- response{
		channel: s.Msg.Channel,
		text:    text,
	}
}

func (s *Sender) Reply(text string) {
	s.respCh <- response{
		channel:     s.Msg.Channel,
		text:        text,
		parentMsgID: s.Msg.Tags["id"],
	}
}

func (s *Sender) Timeout(d time.Duration, reason string) {
	s.TimeoutUser(s.Msg.User, d, reason)
}

func (s *Sender) TimeoutUser(
	username string,
	d time.Duration,
	reason string,
) {
	seconds := int64(d / time.Second)
	text := fmt.Sprintf("/timeout %s %d %s", username, seconds, reason)
	s.Send(text)
}

func (s *Sender) Ban(reason string) {
	s.BanUser(s.Msg.User, reason)
}

func (s *Sender) BanUser(username, reason string) {
	text := fmt.Sprintf("/ban %s %s", username, reason)
	s.Send(text)
}
