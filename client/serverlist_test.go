package client

import "testing"

func TestServerList(t *testing.T) {
	s := newServerList()

	// New lists should be empty
	if e := s.get(); e != nil {
		t.Fatalf("expected empty list to return nil, but received: %v", e)
	}
	if e := s.all(); len(e) != 0 {
		t.Fatalf("expected empty list to return an empty list, but received: %+q", e)
	}

	mklist := func() endpoints {
		return endpoints{
			&endpoint{"b", nil, 1},
			&endpoint{"c", nil, 1},
			&endpoint{"g", nil, 2},
			&endpoint{"d", nil, 1},
			&endpoint{"e", nil, 1},
			&endpoint{"f", nil, 1},
			&endpoint{"h", nil, 2},
			&endpoint{"a", nil, 0},
		}
	}
	s.set(mklist())

	orig := mklist()
	all := s.all()
	if len(all) != len(orig) {
		t.Fatalf("expected %d endpoints but only have %d", len(orig), len(all))
	}

	// Assert list is properly randomized+sorted
	for i, pri := range []int{0, 1, 1, 1, 1, 1, 2, 2} {
		if all[i].priority != pri {
			t.Errorf("expected endpoint %d (%+q) to be priority %d", i, all[i], pri)
		}
	}

	// Subsequent sets should reshuffle (try multiple times as they may
	// shuffle in the same order)
	tries := 0
	max := 3
	for ; tries < max; tries++ {
		s.set(mklist())
		// First entry should always be the same
		if e := s.get(); *e != *all[0] {
			t.Fatalf("on try %d get returned the wrong endpoint: %+q", tries, e)
		}

		all2 := s.all()
		if all.String() == all2.String() {
			// eek, matched; try again in case we just got unlucky
			continue
		}
		break
	}
	if tries == max {
		t.Fatalf("after %d attempts servers were still not random reshuffled", tries)
	}

	// Mark should rotate list items in place
	s.mark(&endpoint{"a", nil, 0})
	all3 := s.all()
	if s.get().name == "a" || all3[len(all3)-1].name != "a" {
		t.Fatalf("endpoint a shold have been rotated to end")
	}
	if len(all3) != len(all) {
		t.Fatalf("marking should not have changed list length")
	}

	// Marking a non-existant endpoint should do nothing
	s.mark(&endpoint{})
	if s.all().String() != all3.String() {
		t.Fatalf("marking a non-existant endpoint alterd the list")
	}
}
