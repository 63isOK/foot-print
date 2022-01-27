package hello_test

import (
	"testing"

	"github.com/63isOK/foot-print/mock/hello"
	"github.com/golang/mock/gomock"
)

func TestSayHi(t *testing.T) {
	ctrl := gomock.NewController(t)

	hi := NewMockHi(ctrl)
	hi.
		EXPECT().
		Hello("gopher").
		Return("Hi, gopher!")
	hi.
		EXPECT().
		Hello(gomock.Eq("aa")).
		DoAndReturn(func(s string) string {
			return "Hi, " + s + "!"
		}).
		MaxTimes(1)

	user := hello.UserInfo{
		Name: "gopher",
		Age:  3,
	}

	hi.
		EXPECT().
		HelloAgain("gopher", "63isOK").
		Return(user)

	got := hello.SayHi(hi, "gopher")
	if got != "Hi, gopher!" {
		t.Errorf("got %q, want %q", got, "Hi, gopher!")
	}

	got = hello.SayHi(hi, "aa")
	if got != "Hi, aa!" {
		t.Errorf("got %q, want %q", got, "Hi, gopher!")
	}

	info := hello.SayHiAgain(hi, "gopher", "63isOK")
	if info != user {
		t.Errorf("got %q, want %q", info, user)
	}
}
