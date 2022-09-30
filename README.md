<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/level27/lvl/deaa8fca41ac2c53f2fb44fe5abe9d04ce4229a4/docs/static/img/level_27_logo_white.svg">
  <img alt="The Level27 Logo" height="175" src="https://raw.githubusercontent.com/level27/lvl/deaa8fca41ac2c53f2fb44fe5abe9d04ce4229a4/docs/static/img/level_27_logo.svg">
</picture>
</p>


`lvl` is a command line tool that empowers Level27 customers to manage and automate their infrastructure. It allows managing of systems, apps, domains, and more.

![GitHub release (latest by date)](https://img.shields.io/github/v/release/level27/lvl) ![](https://img.shields.io/badge/docs-cli.docs.level27.eu-green)

## Installation

You can get compiled builds of `lvl` from the [Releases](https://github.com/level27/lvl/releases/latest) page. We do not currently provide a convenient automatic method of installation, so it has to be done manually:

1. Download the appropriate version from the link above (download `-amd64` if you're unsure about your CPU architecture).
2. Rename it to `lvl`
3. Put it in a good location on your system.
4. Make sure the containing directory is added to your `PATH` environment variable. (in your `~/.bashrc` or such)

## Usage

`lvl` follows a simple scheme for the various commands: simply run `lvl` or `lvl help` to see what it can do.

We also have a [web-based documentation](https://cli.docs.level27.eu/) of the various commands available.

## Building & contributing

`lvl` is built using [Go](https://go.dev/), currently requiring at least Go **1.18**.

To build and run `lvl` yourself, simply run:

```
go run .
```

**Note: `lvl` is developed in-sync with [`l27-go`](git@github.com:level27/l27-go.git). At some times, `lvl` may need to be compiled against the `main` branch of `l27-go`. To do this, you can edit `go.mod` and add a line like so:**

```go.mod
replace github.com/level27/l27-go => /Users/pjb/Projects/l27go-api
```

Edit the path to point to your local repo of `l27-go`.

## License

`lvl` is licensed under the Apache 2.0 license. See [LICENSE](/LICENSE).