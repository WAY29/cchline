package tui

import "testing"

func TestParseSemVerToken(t *testing.T) {
	cases := []struct {
		in      string
		wantOK  bool
		wantVer semVer
	}{
		{"1.2.3", true, semVer{major: 1, minor: 2, patch: 3}},
		{"v1.2.3", true, semVer{major: 1, minor: 2, patch: 3}},
		{"1.2.3-beta.1", true, semVer{major: 1, minor: 2, patch: 3}},
		{"1.2", false, semVer{}},
		{"dev", false, semVer{}},
		{"", false, semVer{}},
	}

	for _, tc := range cases {
		got, ok := parseSemVerToken(tc.in)
		if ok != tc.wantOK {
			t.Fatalf("parseSemVerToken(%q) ok=%v want=%v", tc.in, ok, tc.wantOK)
		}
		if ok && got != tc.wantVer {
			t.Fatalf("parseSemVerToken(%q)=%v want=%v", tc.in, got, tc.wantVer)
		}
	}
}

func TestParseCCHLineVersionOutput(t *testing.T) {
	cases := []struct {
		in       string
		wantDev  bool
		wantOK   bool
		wantSem  semVer
		semverOK bool
	}{
		{"cchline dev\n", true, true, semVer{}, false},
		{"cchline 1.2.3\n", false, true, semVer{major: 1, minor: 2, patch: 3}, true},
		{"something cchline 2.0.1", false, true, semVer{major: 2, minor: 0, patch: 1}, true},
		{"", false, false, semVer{}, false},
		{"cchline ???", false, false, semVer{}, false},
	}

	for _, tc := range cases {
		isDev, v, ok := parseCCHLineVersionOutput([]byte(tc.in))
		if ok != tc.wantOK || isDev != tc.wantDev {
			t.Fatalf("parseCCHLineVersionOutput(%q)=(dev=%v, ok=%v) want (dev=%v, ok=%v)", tc.in, isDev, ok, tc.wantDev, tc.wantOK)
		}
		if ok && !isDev {
			if v != tc.wantSem {
				t.Fatalf("parseCCHLineVersionOutput(%q) sem=%v want=%v", tc.in, v, tc.wantSem)
			}
		}
	}
}

func TestCompareSemVer(t *testing.T) {
	if compareSemVer(semVer{1, 0, 0}, semVer{1, 0, 0}) != 0 {
		t.Fatalf("expected equal")
	}
	if compareSemVer(semVer{1, 0, 0}, semVer{1, 0, 1}) >= 0 {
		t.Fatalf("expected less")
	}
	if compareSemVer(semVer{2, 0, 0}, semVer{1, 9, 9}) <= 0 {
		t.Fatalf("expected greater")
	}
}
