// Command clone clones a git repository into ~/code/<owner>/<repo>.
package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type repo struct {
	Host string
	Path string
	URL  string
}

// parse accepts owner/repo (GitHub), host/owner/repo, http(s), ssh and scp-style URLs.
func parse(arg string) (repo, error) {
	s := strings.TrimSpace(arg)
	s, _, _ = strings.Cut(s, "?")
	s, _, _ = strings.Cut(s, "#")

	var r repo
	var u *url.URL
	ssh := false

	switch {
	case strings.Contains(s, "://"):
		var err error
		if u, err = url.Parse(s); err != nil {
			return r, err
		}
		switch u.Scheme {
		case "https", "http":
		case "ssh", "git+ssh":
			ssh = true
		default:
			return r, fmt.Errorf("unsupported URL: %s", arg)
		}
		r.Host, r.Path = u.Hostname(), u.Path
	case isSCP(s):
		host, path, _ := strings.Cut(s, ":")
		if _, h, ok := strings.Cut(host, "@"); ok {
			host = h
		}
		r.Host, r.Path = host, path
		ssh = true
	default:
		first, rest, _ := strings.Cut(s, "/")
		if strings.Contains(first, ".") {
			r.Host, r.Path = first, rest
		} else {
			r.Host, r.Path = "github.com", s
		}
	}

	r.Host = strings.TrimPrefix(strings.ToLower(r.Host), "www.")
	if r.Host == "" {
		return r, fmt.Errorf("missing host: %s", arg)
	}

	path := strings.Trim(r.Path, "/")
	path, _, _ = strings.Cut(path, "/-/")
	parts := strings.Split(path, "/")
	if r.Host == "github.com" && len(parts) > 2 {
		parts = parts[:2]
	}
	parts[len(parts)-1] = strings.TrimSuffix(parts[len(parts)-1], ".git")
	if len(parts) < 2 {
		return r, fmt.Errorf("expected owner/repo, got %q", arg)
	}
	for _, p := range parts {
		if p == "" || p == "." || p == ".." {
			return r, fmt.Errorf("expected owner/repo, got %q", arg)
		}
	}
	r.Path = strings.Join(parts, "/")

	switch {
	case ssh:
		r.URL = s
	case r.Host == "github.com":
		r.URL = "https://github.com/" + r.Path + ".git"
	case u != nil:
		u.Path = "/" + r.Path
		r.URL = u.String()
	default:
		r.URL = "https://" + r.Host + "/" + r.Path
	}
	return r, nil
}

// isSCP reports whether s looks like user@host:path.
func isSCP(s string) bool {
	colon := strings.IndexByte(s, ':')
	slash := strings.IndexByte(s, '/')
	return colon > 0 && (slash == -1 || colon < slash)
}

// clonedFrom reports whether dir is a git checkout with a remote pointing at r.
func clonedFrom(dir string, r repo) bool {
	out, err := exec.Command("git", "-C", dir, "remote", "-v").Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(string(out), "\n") {
		f := strings.Fields(line)
		if len(f) < 2 {
			continue
		}
		other, err := parse(f[1])
		if err == nil && other.Host == r.Host && strings.EqualFold(other.Path, r.Path) {
			return true
		}
	}
	return false
}

func root() (string, error) {
	if dir := os.Getenv("CLONE_ROOT"); dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "code"), nil
}

func run(arg, path string) error {
	r, err := parse(arg)
	if err != nil {
		return err
	}
	base, err := root()
	if err != nil {
		return err
	}
	if path == "" {
		path = r.Path
	}
	dest := path
	if !filepath.IsAbs(dest) {
		dest = filepath.Join(base, dest)
	}

	if _, err := os.Stat(dest); err == nil {
		if !clonedFrom(dest, r) {
			return fmt.Errorf("%s already exists and isn't a clone of %s/%s; choose another location with --path", dest, r.Host, r.Path)
		}
		fmt.Fprintf(os.Stderr, "already cloned: %s\n", dest)
		fmt.Println(dest)
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	cmd := exec.Command("git", "clone", r.URL, dest)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}
	fmt.Println(dest)
	return nil
}

func usage(code int) {
	fmt.Fprintln(os.Stderr, "usage: clone [--path <dir>] <owner/repo | git URL>")
	fmt.Fprintln(os.Stderr, "clones into $CLONE_ROOT/<owner>/<repo> (default ~/code)")
	fmt.Fprintln(os.Stderr, "--path overrides the destination, relative to $CLONE_ROOT unless absolute")
	os.Exit(code)
}

func main() {
	var arg, path string
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch a := args[i]; {
		case a == "-h" || a == "--help":
			usage(0)
		case a == "-p" || a == "--path":
			i++
			if i == len(args) {
				usage(2)
			}
			path = args[i]
		case strings.HasPrefix(a, "--path="):
			path = strings.TrimPrefix(a, "--path=")
		case arg == "":
			arg = a
		default:
			usage(2)
		}
	}
	if arg == "" {
		usage(2)
	}
	if err := run(arg, path); err != nil {
		fmt.Fprintln(os.Stderr, "clone:", err)
		os.Exit(1)
	}
}
