# clone

Clone a git repository into `~/code/<owner>/<repo>`, unless it's already there.

```sh
clone alpha-omega-security/scrutineer
clone https://github.com/alpha-omega-security/scrutineer/
```

Both of these end up in `~/code/alpha-omega-security/scrutineer`. If that directory already exists, nothing is cloned and the path is printed.

## Install

```sh
go install github.com/andrew/clone@latest
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`.

## Accepted inputs

- `owner/repo` (GitHub)
- `github.com/owner/repo`
- `https://github.com/owner/repo`, with or without `.git`, a trailing slash, or deeper paths like `/tree/main/cmd`
- `git@github.com:owner/repo.git` and `ssh://git@github.com/owner/repo.git` (cloned over SSH)
- Full URLs for any other host, e.g. `git@codeberg.org:foo/bar.git` or `https://gitlab.com/group/sub/repo`

The host isn't part of the destination, so `git@codeberg.org:foo/bar.git` goes to `~/code/foo/bar`. GitLab-style subgroups keep their full path (`~/code/group/sub/repo`).

## Clashes

If the destination already exists but none of its remotes point at the repo you asked for (say `~/code/foo/bar` is a clone of `github.com/foo/bar` and you ask for `codeberg.org/foo/bar`), `clone` exits with an error. Pick a different location with `--path`, relative to `~/code` unless absolute:

```sh
clone --path foo/bar-codeberg git@codeberg.org:foo/bar.git
```

## Jumping into the repo

The destination path is printed to stdout, and git's output goes to stderr, so you can wrap it in a shell function to clone and `cd` in one step:

```sh
c() { local dir; dir="$(command clone "$@")" && cd "$dir"; }
```

## Configuration

Set `CLONE_ROOT` to use a base directory other than `~/code`:

```sh
CLONE_ROOT=~/src clone owner/repo
```

## License

MIT
