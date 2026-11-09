package main

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in   string
		host string
		path string
		url  string
	}{
		{"alpha-omega-security/scrutineer", "github.com", "alpha-omega-security/scrutineer", "https://github.com/alpha-omega-security/scrutineer.git"},
		{"https://github.com/alpha-omega-security/scrutineer/", "github.com", "alpha-omega-security/scrutineer", "https://github.com/alpha-omega-security/scrutineer.git"},
		{"https://github.com/alpha-omega-security/scrutineer.git", "github.com", "alpha-omega-security/scrutineer", "https://github.com/alpha-omega-security/scrutineer.git"},
		{"https://github.com/andrew/clone/tree/main/cmd", "github.com", "andrew/clone", "https://github.com/andrew/clone.git"},
		{"http://www.github.com/andrew/clone?tab=readme", "github.com", "andrew/clone", "https://github.com/andrew/clone.git"},
		{"github.com/andrew/clone", "github.com", "andrew/clone", "https://github.com/andrew/clone.git"},
		{"git@github.com:andrew/clone.git", "github.com", "andrew/clone", "git@github.com:andrew/clone.git"},
		{"ssh://git@github.com/andrew/clone.git", "github.com", "andrew/clone", "ssh://git@github.com/andrew/clone.git"},
		{"git@codeberg.org:foo/bar.git", "codeberg.org", "foo/bar", "git@codeberg.org:foo/bar.git"},
		{"https://codeberg.org/foo/bar", "codeberg.org", "foo/bar", "https://codeberg.org/foo/bar"},
		{"codeberg.org/foo/bar", "codeberg.org", "foo/bar", "https://codeberg.org/foo/bar"},
		{"https://gitlab.com/group/sub/repo.git", "gitlab.com", "group/sub/repo", "https://gitlab.com/group/sub/repo"},
		{"https://gitlab.com/group/sub/repo/-/tree/main", "gitlab.com", "group/sub/repo", "https://gitlab.com/group/sub/repo"},
		{"https://git.example.com:3000/foo/bar", "git.example.com", "foo/bar", "https://git.example.com:3000/foo/bar"},
		{"ssh://git@git.example.com:2222/foo/bar.git", "git.example.com", "foo/bar", "ssh://git@git.example.com:2222/foo/bar.git"},
	}
	for _, tt := range tests {
		r, err := parse(tt.in)
		if err != nil {
			t.Errorf("parse(%q) error: %v", tt.in, err)
			continue
		}
		if r.Host != tt.host || r.Path != tt.path || r.URL != tt.url {
			t.Errorf("parse(%q) = %+v, want %s %s %s", tt.in, r, tt.host, tt.path, tt.url)
		}
	}
}

func TestParseErrors(t *testing.T) {
	for _, in := range []string{"", "scrutineer", "codeberg.org/foo", "ftp://github.com/a/b", "../x", "a/..", "git@codeberg.org:foo"} {
		if _, err := parse(in); err == nil {
			t.Errorf("parse(%q) expected error", in)
		}
	}
}
