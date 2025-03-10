# ide

ide provides a tool set that gets out of your way

[![Actions Status](https://github.com/nrocco/ide/actions/workflows/master.yml/badge.svg)](https://github.com/nrocco/ide/actions/workflows/master.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nrocco/ide.svg)](https://pkg.go.dev/github.com/nrocco/ide)
[![Go Report Card](https://goreportcard.com/badge/github.com/nrocco/ide)](https://goreportcard.com/report/github.com/nrocco/ide)

> a tool that glues vim, git, direnv, docker compose and ctags together
> to provide a powerful integrated development environment

## Usage

After installing `ide` you can invoke it without any arguments to get help:

    $ ide
    ide provides a powerful ide that gets out of your way

    Usage:
      ide [command]

    Available Commands:
      compare     Compare a file in the current project with another project
      completion  Output shell completion code for the specified shell
      destroy     Remove all ide configuration for a repository
      fix         Fix source code
      help        Help about any command
      hook        Manage git hooks for an ide project
      lint        Lint source code and report errors
      rgit        Run a git command in multiple git projects
      shim        Manage shims for an ide project
      status      Get the current status of your ide project
      version     Display version and build information

    Flags:
      -h, --help   help for ide

    Use "ide [command] --help" for more information about a command.

To remove any traces of `ide` run:

    $ ide destroy
    2017/07/04 20:16:18 Repository is no longer an ide project

You can also view the current status of your ide project:

    $ ide status
    Ide
      Name: my-project
      Branch: master
      Location: /Users/nrocco/dev/my-project
      Ctags:
        File: /Users/nrocco/dev/my-project/.git/tags
        Age: 5 days ago
        Size: 18 kB
      Hooks: ~
      Binaries:
        go: docker-compose exec --workdir=$PWD backend go
        goimports: docker-compose exec --workdir=$PWD backend goimports
        gofmt: docker-compose exec --workdir=$PWD backend gofmt
        npm: docker-compose exec --workdir=$PWD frontend npm

In the above case no hooks are enabled for this project. In order to enable
the `prepare-commit-msg` hook run:

    $ ide hook enable prepare-commit-msg
    2017/07/04 20:17:37 Hook prepare-commit-msg enabled

You can see the hook is enabled:

    $ ls -ilah .git/hooks/prepare-commit-msg
    29546377 lrwxr-xr-x 1 nrocco staff 52 Jul  4 20:17 .git/hooks/prepare-commit-msg -> /usr/local/bin/ide

## Use docker-compose

In `.git/compose.yaml`:

    ---
    name: "my-name"
    services:
      ruby:
        image: "ruby:3"
        platform: "linux/amd64"
        init: true
        command: ["sleep", "infinity"]
        working_dir: "${PWD}"
        volumes:
          - "${PWD}:${PWD}"
          - "ruby_cache:/usr/local/bundle"
    volumes:
      ruby_cache:


## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Make sure that tests pass (`make test`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request

## Contributors

- Nico Di Rocco (https://github.com/nrocco)
