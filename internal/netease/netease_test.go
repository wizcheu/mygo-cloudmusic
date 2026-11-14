package netease

import "testing"

// Marks say what a track needs to play in full, whatever the listener may
// play.
func TestTrackMarks(t *testing.T) {
	cases := []struct {
		name            string
		tr              Track
		vip, paid, down bool
	}{
		{"VIP", Track{Fee: 1}, true, false, false},
		{"low quality free", Track{Fee: 8}, true, false, false},
		{"paid album", Track{Fee: 4}, false, true, false},
		{"free", Track{Fee: 0}, false, false, false},
		{"taken down", Track{Privilege: &Privilege{St: -200}}, false, false, true},
		{"fee from privilege", Track{Privilege: &Privilege{Fee: 1}}, true, false, false},
	}
	for _, c := range cases {
		if c.tr.VIP() != c.vip || c.tr.PaidAlbum() != c.paid || c.tr.TakenDown() != c.down {
			t.Errorf("%s: VIP %v, paid %v, taken down %v", c.name, c.tr.VIP(), c.tr.PaidAlbum(), c.tr.TakenDown())
		}
	}
}

func TestPlayableFollowsPrivilege(t *testing.T) {
	tracks := []Track{{ID: 1, Fee: 1}, {ID: 2, Fee: 8}, {ID: 3, Fee: 8}, {ID: 4, Fee: 1}}
	if tracks[0].Playable() || !tracks[1].Playable() {
		t.Fatal("without privileges, VIP songs do not play and free ones do")
	}
	withPrivileges(tracks, []Privilege{
		{ID: 1, Pl: 320000},           // VIP song, the listener subscribes
		{ID: 2, Pl: 0},                // nothing to play
		{ID: 3, St: -200, Pl: 128000}, // taken down
	})
	want := []bool{true, false, false, false}
	for i, tr := range tracks {
		if tr.Playable() != want[i] {
			t.Errorf("track %d: playable %v", tr.ID, tr.Playable())
		}
	}
}
