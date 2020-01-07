# git-match-run

Utility to run commands when specific files have changed on Git.
It's designed to be used in CI/CD systems to run builds steps only for specific parts of the Git repos.

## Installation

First, you need to install [Go](https://golang.org/dl/) version 1.12 or higher.
After doing so, you can use `go get` to install git-match-run:

```
go get github.com/Lepovirta/git-match-run
```

## Usage

This command accepts the following parameters:

* `-config`: Location of the config file. See the Config section below for more details. Default: `gitmatchrun.yaml`
* `-from`: Git ref to start finding changes from
* `-to`: Git ref to end finding changes. Default: `HEAD`
* `-run`: Enable to actually run the commands. If not enabled, the run is only simulated.

## Config

Configuration is provided in YAML format.
The document is a list of entries, where each entry can contain the following fields:

* `pattern`: The RegEx pattern to match a file on.
* `command`: The command to execute when a matching file is found.
* `args`: The arguments to provide to the command.

## Example

Your config could have the following entries in file `gitmatchrun.yaml`:

```yaml
- pattern: Dockerfile
  command: docker
  args:
  - build
  - -t
  - mydockerimage
  - .
- pattern: README.md
  command: spellcheck
  args:
  - README.md
```

You can then find out what commands would be run between two Git refs:

```
git-match-run -from b6cf08d -to ac2509c
```

You can run the actual commands by including the `-run` flag:

```
git-match-run -from b6cf08d -to ac2509c
```

When integrated in a CI/CD system, the from and to refs would be provided by your CI system.

## License

GNU General Public License v3.0

See LICENSE file for more information.

