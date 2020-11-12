package querylog

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"testing"

	aglog "github.com/AdguardTeam/golibs/log"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestDecode_decodeQueryLog(t *testing.T) {
	stdWriter := log.Writer()
	stdLevel := aglog.GetLevel()
	t.Cleanup(func() {
		log.SetOutput(stdWriter)
		aglog.SetLevel(stdLevel)
	})

	logOut := &bytes.Buffer{}
	log.SetOutput(logOut)

	aglog.SetLevel(aglog.DEBUG)

	strBytes := &bytes.Buffer{}
	enc := base64.NewEncoder(base64.StdEncoding, strBytes)

	testCases := []struct {
		name    string
		msgFunc func() error
		want    string
	}{{
		name: "back_compatibility_all_right",
		msgFunc: func() error {
			q := &dns.Msg{}
			q.SetQuestion(dns.Fqdn("google.com"), dns.TypeMX)
			bts, err := q.Pack()
			if err != nil {
				return err
			}
			if _, err := enc.Write(bts); err != nil {
				return err
			}
			return nil
		},
		want: "default",
	}, {
		name: "back_compatibility_bad_msg",
		msgFunc: func() error {
			q := &dns.Msg{}
			q.SetQuestion(dns.Fqdn(""), dns.TypeMX)
			bts, err := q.Pack()
			if err != nil {
				return err
			}
			if _, err := enc.Write(bts); err != nil {
				return err
			}
			return nil
		},
		want: "decodeLogEntry err: dns: overflow unpacking uint32\n",
	}, {
		name: "back_compatibility_no_questions",
		msgFunc: func() error {
			q := &dns.Msg{}
			q.SetQuestion(dns.Fqdn("google.com"), dns.TypeMX)
			q.Question = []dns.Question{}
			bts, err := q.Pack()
			if err != nil {
				return err
			}
			if _, err := enc.Write(bts); err != nil {
				return err
			}
			return nil
		},
		want: "default",
	}, {
		name: "back_compatibility_bad_decoding",
		msgFunc: func() error {
			q := &dns.Msg{}
			q.SetQuestion(dns.Fqdn("google.com"), dns.TypeMX)
			q.Question = []dns.Question{{Qclass: dns.ClassINET, Qtype: dns.TypeMX}}
			bts, err := q.Pack()
			enc := base64.NewEncoder(
				base64.NewEncoding("################################################################"),
				strBytes)
			if err != nil {
				return err
			}
			if _, err := enc.Write(bts); err != nil {
				return err
			}
			return nil
		},
		want: "decodeLogEntry err: illegal base64 data at input byte 0\n",
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := logOut.Write([]byte("default"))
			assert.Nil(t, err)

			err = tc.msgFunc()
			assert.Nil(t, err)

			oldTags := fmt.Sprintf("{\"Question\":%q,\"Time\":\"2006-01-02T15:04:05Z\"}", strBytes)

			l := &logEntry{}
			decodeLogEntry(l, oldTags)

			assert.True(t, strings.HasSuffix(logOut.String(), tc.want))

			logOut.Reset()
			strBytes.Reset()
		})
	}
}

func TestJSON(t *testing.T) {
	s := `
	{"keystr":"val","obj":{"keybool":true,"keyint":123456}}
	`
	k, v, jtype := readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTStr))
	assert.Equal(t, "keystr", k)
	assert.Equal(t, "val", v)

	k, _, jtype = readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTObj))
	assert.Equal(t, "obj", k)

	k, v, jtype = readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTBool))
	assert.Equal(t, "keybool", k)
	assert.Equal(t, "true", v)

	k, v, jtype = readJSON(&s)
	assert.Equal(t, jtype, int32(jsonTNum))
	assert.Equal(t, "keyint", k)
	assert.Equal(t, "123456", v)

	_, _, jtype = readJSON(&s)
	assert.True(t, jtype == jsonTErr)
}
