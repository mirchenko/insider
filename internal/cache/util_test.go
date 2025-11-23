package cache

import "testing"

func TestMakeRedisKey(t *testing.T) {
	cases := []struct{ entity, key, want string }{
		{"msg", "123", "msg.123"},
		{"user", "abc", "user.abc"},
		{"", "k", ".k"},
		{"e", "", "e."},
	}
	for _, c := range cases {
		if got := makeRedisKey(c.entity, c.key); got != c.want {
			t.Fatalf("makeRedisKey(%q,%q)=%q; want %q", c.entity, c.key, got, c.want)
		}
	}
}
