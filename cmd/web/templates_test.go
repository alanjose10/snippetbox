package main

import (
	"testing"
	"time"

	"github.com/alanjose10/snippetbox/internal/assert"
)

func TestHumanDate(t *testing.T) {
	// tm := time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC)

	// hd := humanDate(tm)

	// expectedHt := "17 Mar 2024 at 10:15 AM"

	// if hd != expectedHt {
	// 	t.Errorf("got %q; want %q\n", hd, expectedHt)
	// }

	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.UTC),
			want: "17 Mar 2024 at 10:15",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 3, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			want: "17 Mar 2024 at 09:15",
		},
	}

	for _, item := range tests {
		t.Run(item.name, func(t *testing.T) {
			res := humanDate(item.tm)
			assert.Equal(t, item.want, res)
		})
	}

}
