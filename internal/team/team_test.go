package team

import (
	"errors"
	"testing"
	"time"

	"github.com/haritabh17/theirtime/internal/config"
	"github.com/haritabh17/theirtime/internal/slack"
)

func defaultCfg() *config.Config {
	c := &config.Config{}
	config.ApplyDefaults(c)
	return c
}

type fakeInfoPresenceClient struct {
	info        slack.UserInfo
	presence    slack.Presence
	presenceErr error
}

func (c fakeInfoPresenceClient) GetUserInfo(userID string) (slack.UserInfo, error) {
	info := c.info
	info.ID = userID
	return info, nil
}

func (c fakeInfoPresenceClient) GetUserPresence(userID string) (slack.Presence, error) {
	if c.presenceErr != nil {
		return "", c.presenceErr
	}
	return c.presence, nil
}

func TestFormatMenubarTitleDefault(t *testing.T) {
	cfg := defaultCfg()
	members := []MemberTime{
		{Label: "sugu", Time: "10.46 AM"},
		{Label: "rafa", Time: "3.15 PM"},
	}
	got := FormatMenubarTitle(cfg, members)
	want := "10.46 AM | 3.15 PM"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestListWithTimesPopulatesPresence(t *testing.T) {
	cfg := defaultCfg()
	cfg.Team = []config.TeamMember{{Label: "sugu", ID: "U12345678"}}
	client := fakeInfoPresenceClient{
		info: slack.UserInfo{
			DisplayName: "Sugu",
			TZ:          "America/Los_Angeles",
			AvatarURL:   "https://example.com/avatar.png",
		},
		presence: slack.PresenceActive,
	}

	got, err := ListWithTimes(client, cfg, time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d members want 1", len(got))
	}
	if got[0].Presence != slack.PresenceActive {
		t.Fatalf("presence got %q want %q", got[0].Presence, slack.PresenceActive)
	}
	if got[0].AvatarURL == "" || got[0].Time == "" || got[0].Time == "—" {
		t.Fatalf("expected profile and time fields to be populated: %#v", got[0])
	}
}

func TestListWithTimesIgnoresPresenceError(t *testing.T) {
	cfg := defaultCfg()
	cfg.Team = []config.TeamMember{{Label: "sugu", ID: "U12345678"}}
	client := fakeInfoPresenceClient{
		info: slack.UserInfo{
			DisplayName: "Sugu",
			TZ:          "America/Los_Angeles",
			AvatarURL:   "https://example.com/avatar.png",
		},
		presenceErr: errors.New("missing scope"),
	}

	got, err := ListWithTimes(client, cfg, time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("got %d members want 1", len(got))
	}
	if got[0].Presence != "" {
		t.Fatalf("presence got %q want empty", got[0].Presence)
	}
	if got[0].AvatarURL == "" || got[0].Time == "" || got[0].Time == "—" {
		t.Fatalf("presence error should not drop profile or time fields: %#v", got[0])
	}
}

func TestFormatMenubarTitleWithNames(t *testing.T) {
	cfg := defaultCfg()
	show := true
	cfg.ShowName = &show
	members := []MemberTime{
		{Label: "sugu", Time: "10.46 AM"},
		{Label: "rafa", Time: "3.15 PM"},
	}
	got := FormatMenubarTitle(cfg, members)
	want := "sugu - 10.46 AM | rafa - 3.15 PM"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatMenubarTitleGroupsSameTimezone(t *testing.T) {
	cfg := defaultCfg()
	members := []MemberTime{
		{Label: "sugu", TZ: "America/Los_Angeles", Time: "10.46 AM"},
		{Label: "rafa", TZ: "America/New_York", Time: "1.46 PM"},
		{Label: "ann", TZ: "America/Los_Angeles", Time: "10.46 AM"},
	}
	got := FormatMenubarTitle(cfg, members)
	want := "10.46 AM | 1.46 PM"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatMenubarTitleGroupsNamesWithSharedTime(t *testing.T) {
	cfg := defaultCfg()
	show := true
	cfg.ShowName = &show
	members := []MemberTime{
		{Label: "sugu", TZ: "America/Los_Angeles", Time: "10.46 AM"},
		{Label: "rafa", TZ: "America/Los_Angeles", Time: "10.46 AM"},
		{Label: "ann", TZ: "America/New_York", Time: "1.46 PM"},
	}
	got := FormatMenubarTitle(cfg, members)
	want := "sugu, rafa - 10.46 AM | ann - 1.46 PM"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatMenubarTitleGroupsSameCurrentOffset(t *testing.T) {
	cfg := defaultCfg()
	show := true
	cfg.ShowName = &show
	members := []MemberTime{
		{Label: "manan", TZ: "Europe/Belgrade", Time: "3.02 PM"},
		{Label: "nico", TZ: "Europe/Amsterdam", Time: "3.02 PM"},
		{Label: "mats", TZ: "Europe/Amsterdam", Time: "3.02 PM"},
	}
	got := FormatMenubarTitle(cfg, members)
	want := "manan, nico, mats - 3.02 PM"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatMenubarTitleShowTimeFalseGroupsNames(t *testing.T) {
	cfg := defaultCfg()
	showName := true
	showTime := false
	cfg.ShowName = &showName
	cfg.ShowTime = &showTime
	members := []MemberTime{
		{Label: "sugu", TZ: "America/Los_Angeles", Time: "10.46 AM"},
		{Label: "rafa", TZ: "America/Los_Angeles", Time: "10.46 AM"},
	}
	got := FormatMenubarTitle(cfg, members)
	want := "sugu, rafa"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestGroupMemberTimesAtSortsByUTCOffset(t *testing.T) {
	at := time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC)
	members := []MemberTime{
		{Label: "first-ny", TZ: "America/New_York"},
		{Label: "first-belgrade", TZ: "Europe/Belgrade"},
		{Label: "first-la", TZ: "America/Los_Angeles"},
		{Label: "first-amsterdam", TZ: "Europe/Amsterdam"},
	}
	groups := GroupMemberTimesAt(members, at)
	if len(groups) != 3 {
		t.Fatalf("got %d groups want 3", len(groups))
	}
	if groups[0].Members[0].Label != "first-la" {
		t.Fatalf("group 0 got %q want first-la", groups[0].Members[0].Label)
	}
	if groups[1].Members[0].Label != "first-ny" {
		t.Fatalf("group 1 got %q want first-ny", groups[1].Members[0].Label)
	}
	gotThird := []string{groups[2].Members[0].Label, groups[2].Members[1].Label}
	wantFirst := []string{"first-belgrade", "first-amsterdam"}
	for i := range wantFirst {
		if gotThird[i] != wantFirst[i] {
			t.Fatalf("group 2 got %v want %v", gotThird, wantFirst)
		}
	}
}

func TestSortMemberTimesByUTCOffsetAtFlattensRenderedOrder(t *testing.T) {
	at := time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC)
	members := []MemberTime{
		{Label: "ny", TZ: "America/New_York"},
		{Label: "belgrade", TZ: "Europe/Belgrade"},
		{Label: "la", TZ: "America/Los_Angeles"},
		{Label: "amsterdam", TZ: "Europe/Amsterdam"},
		{Label: "kolkata", TZ: "Asia/Kolkata"},
	}
	sorted := SortMemberTimesByUTCOffsetAt(members, at)
	got := make([]string, 0, len(sorted))
	for _, m := range sorted {
		got = append(got, m.Label)
	}
	want := []string{"la", "ny", "belgrade", "amsterdam", "kolkata"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got %v want %v", got, want)
		}
	}
}

func TestGroupMemberTimesAtKeepsBlankTimezoneSeparate(t *testing.T) {
	at := time.Date(2026, 6, 23, 12, 0, 0, 0, time.UTC)
	members := []MemberTime{
		{Label: "missing-1", Time: "—"},
		{Label: "missing-2", Time: "—"},
	}
	groups := GroupMemberTimesAt(members, at)
	if len(groups) != 2 {
		t.Fatalf("got %d groups want 2", len(groups))
	}
	if groups[0].Members[0].Label != "missing-1" || groups[1].Members[0].Label != "missing-2" {
		t.Fatalf("blank timezone members should keep separate groups: %#v", groups)
	}
}

func TestTruncateMenubarTitle(t *testing.T) {
	cfg := defaultCfg()
	show := true
	cfg.ShowName = &show
	long := MemberTime{Label: "verylonglabel", Time: "10.46 AM"}
	title := FormatMenubarTitle(cfg, []MemberTime{long, long, long, long, long})
	if len([]rune(title)) > maxMenubarTitleRunes {
		t.Fatalf("title too long: %d runes", len([]rune(title)))
	}
	runes := []rune(title)
	if runes[len(runes)-1] != '…' {
		t.Fatalf("expected ellipsis, got %q", title)
	}
}

func TestFormatMemberTime(t *testing.T) {
	cfg := defaultCfg()
	at := time.Date(2026, 6, 13, 15, 30, 0, 0, time.UTC)
	got := FormatMemberTime(cfg, "America/New_York", at)
	if got == "" || got == "—" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatMemberDisplayAvatarTimeOnly(t *testing.T) {
	cfg := defaultCfg()
	got := FormatMemberDisplay(cfg, MemberTime{Label: "bob", Time: "4.07 PM"}, "")
	if got != "4.07 PM" {
		t.Fatalf("got %q", got)
	}
}
