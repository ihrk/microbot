package irc

import (
	"reflect"
	"testing"
)

func assertEq(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("\nexpected: '%v',\nactual: '%v'", expected, actual)
	}
}

type testCase struct {
	rawMsg   string
	expected Msg
}

func (tc *testCase) run(t *testing.T) {
	actual := ParseMsg(tc.rawMsg)

	assertEq(t, tc.expected.Channel, actual.Channel)
	assertEq(t, tc.expected.Type, actual.Type)
	assertEq(t, tc.expected.User, actual.User)
	assertEq(t, tc.expected.Text, actual.Text)
	assertEq(t, tc.expected.Tags, actual.Tags)
}

var cases = []*testCase{
	{
		rawMsg: "PING :tmi.twitch.tv",
		expected: Msg{
			Type: "PING",
			Text: "tmi.twitch.tv",
		},
	},
	{
		rawMsg: ":tmi.twitch.tv HOSTTARGET #generic_streamer :- 0",
		expected: Msg{
			Type:    MsgTypeHostTarget,
			Channel: "generic_streamer",
			User:    "",
			Text:    "- 0",
		},
	},
	{
		rawMsg: "@msg-id=host_on :tmi.twitch.tv NOTICE #pro_channel :Now hosting generic_streamer.",
		expected: Msg{
			Type:    MsgTypeNotice,
			Channel: "pro_channel",
			User:    "",
			Text:    "Now hosting generic_streamer.",
			Tags: map[string]string{
				"msg-id": "host_on",
			},
		},
	},
	{
		rawMsg: `@badge-info=;badges=;color=#FFFFFF;display-name=Raider;emotes=;flags=;id=some-uuid;login=raider;mod=0;msg-id=raid;msg-param-displayName=Raider;msg-param-login=raider;msg-param-profileImageURL=https://static-cdn.jtvnw.net/jtv_user_pictures/uuid-profile_image-70x70.jpg;msg-param-viewerCount=80;room-id=999999999;subscriber=0;system-msg=80\sraiders\sfrom\sRaider\shave\sjoined!;tmi-sent-ts=1627703100614;user-id=121212121;user-type= :tmi.twitch.tv USERNOTICE #pro_channel`,
		expected: Msg{
			Type:    MsgTypeUserNotice,
			Channel: "pro_channel",
			User:    "",
			Text:    "",
			Tags: map[string]string{
				"badge-info":                "",
				"badges":                    "",
				"color":                     "#FFFFFF",
				"display-name":              "Raider",
				"emotes":                    "",
				"flags":                     "",
				"id":                        "some-uuid",
				"login":                     "raider",
				"mod":                       "0",
				"msg-id":                    "raid",
				"msg-param-displayName":     "Raider",
				"msg-param-login":           "raider",
				"msg-param-profileImageURL": "https://static-cdn.jtvnw.net/jtv_user_pictures/uuid-profile_image-70x70.jpg",
				"msg-param-viewerCount":     "80",
				"room-id":                   "999999999",
				"subscriber":                "0",
				"system-msg":                "80 raiders from Raider have joined!",
				"tmi-sent-ts":               "1627703100614",
				"user-id":                   "121212121",
				"user-type":                 "",
			},
		},
	},
	{
		rawMsg: "@badge-info=;badges=;client-nonce=nonce-value;color=#FFFFFF;display-name=Twitch_Viewer;emotes=;flags=some-flag;id=some-uuid;mod=0;room-id=7777777;subscriber=0;tmi-sent-ts=1627670607601;turbo=0;user-id=3333333;user-type= :twitch_viewer!twitch_viewer@twitch_viewer.tmi.twitch.tv PRIVMSG #twitch_streamer :hello",
		expected: Msg{
			Type:    MsgTypePrivMsg,
			Channel: "twitch_streamer",
			Text:    "hello",
			User:    "twitch_viewer",
			Tags: map[string]string{
				"badge-info":   "",
				"badges":       "",
				"client-nonce": "nonce-value",
				"color":        "#FFFFFF",
				"display-name": "Twitch_Viewer",
				"emotes":       "",
				"flags":        "some-flag",
				"id":           "some-uuid",
				"mod":          "0",
				"room-id":      "7777777",
				"subscriber":   "0",
				"tmi-sent-ts":  "1627670607601",
				"turbo":        "0",
				"user-id":      "3333333",
				"user-type":    "",
			},
		},
	},
	{
		rawMsg: "@badge-info=;badges=moderator/1;color=;display-name=twitch_viewer;emote-sets=0;mod=1;subscriber=0;user-type=mod :tmi.twitch.tv USERSTATE #random_channel",
		expected: Msg{
			Type:    MsgTypeUserState,
			Channel: "random_channel",
			Text:    "",
			User:    "",
			Tags: map[string]string{
				"badge-info":   "",
				"badges":       "moderator/1",
				"color":        "",
				"display-name": "twitch_viewer",
				"emote-sets":   "0",
				"mod":          "1",
				"subscriber":   "0",
				"user-type":    "mod",
			},
		},
	},
	{
		rawMsg: "@msg-id=vips_success :tmi.twitch.tv NOTICE #little_streamer :The VIPs of this channel are: vip_user_1, vip_user_2, vip_user_3.",
		expected: Msg{
			Type:    MsgTypeNotice,
			Channel: "little_streamer",
			Text:    "The VIPs of this channel are: vip_user_1, vip_user_2, vip_user_3.",
			User:    "",
			Tags: map[string]string{
				"msg-id": "vips_success",
			},
		},
	},
	{
		rawMsg: ":tmi.twitch.tv CAP * ACK :twitch.tv/tags",
		expected: Msg{
			Type: "CAP",
			Text: "twitch.tv/tags",
		},
	},
	{
		rawMsg: ":tmi.twitch.tv NOTICE * :Login authentication failed",
		expected: Msg{
			Type: MsgTypeNotice,
			Text: "Login authentication failed",
		},
	},
	{
		rawMsg: ":tmi.twitch.tv 001 twitch_viewer :Welcome, GLHF!",
		expected: Msg{
			Type: "001",
			Text: "Welcome, GLHF!",
		},
	},
	{
		rawMsg: "@badge-info=;badges=;color=;display-name=twitch_viewer;emote-sets=0;user-id=123;user-type= :tmi.twitch.tv GLOBALUSERSTATE",
		expected: Msg{
			Type: MsgTypeGlobalUserState,
			Tags: map[string]string{
				"badge-info":   "",
				"badges":       "",
				"color":        "",
				"display-name": "twitch_viewer",
				"emote-sets":   "0",
				"user-id":      "123",
				"user-type":    "",
			},
		},
	},
}

func TestParseMsg(t *testing.T) {
	for _, tc := range cases {
		tc.run(t)
	}
}

func BenchmarkParseMsg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseMsg(cases[4].rawMsg)
	}
}
