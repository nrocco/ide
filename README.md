ide
===

ide provides a powerful tool set that gets out of your way

[![Actions Status](https://github.com/nrocco/ide/actions/workflows/master.yml/badge.svg)](https://github.com/nrocco/ide/actions/workflows/master.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/nrocco/ide.svg)](https://pkg.go.dev/github.com/nrocco/ide)
[![Go Report Card](https://goreportcard.com/badge/github.com/nrocco/ide)](https://goreportcard.com/report/github.com/nrocco/ide)

> a tool that glues vim/git/ctags together to provide a powerful integrated
> development environment



Usage
-----

After installing `ide` you can invoke it without any arguments to get help:

    % ide
    ide provides a powerful ide that gets out of your way

    Usage:
      ide [command]

    Available Commands:
      binary      Manage binaries for an ide project
      completion  Output shell completion code for the specified shell
      destroy     Remove all ide configuration for a repository
      help        Help about any command
      hook        Manage git hooks for an ide project
      init        Initialize a git repository as an ide project
      status      Get the current status of your ide project
      version     Display version and build information

    Flags:
          --config string   config file (default is .ide.yaml in $PWD, $HOME, /etc)
      -h, --help            help for ide

    Use "ide [command] --help" for more information about a command.


To setup an existing git repository as an `ide` project run:

    % ide init
    2017/07/04 20:15:12 Setting up the repository as a ide project...
    2017/07/04 20:15:12 Setting the project language to go


To remove any traces of `ide` run:

    % ide destroy
    2017/07/04 20:16:18 Repository is no longer an ide project


You can also view the current status of your ide project:

    $ ide status
    Ide
      Name: my-project
      Branch: master
      Language: go
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
the `post-checkout` hook (which runs ctags) run:

    % ide hook enable post-checkout
    2017/07/04 20:17:37 Hook post-checkout enabled


You can see the hook is enabled:

    % ls -ilah .git/hooks/post-checkout
    29546377 lrwxr-xr-x 1 nrocco staff 52 Jul  4 20:17 .git/hooks/post-checkout -> /usr/local/bin/ide



Alternatives
------------

Execute a command defined in `.git/config` can alternatively achieved with a
simple bash script:

    #!/bin/sh
    set -xe
    NAME="ide.binaries.$(basename $0)"
    exec $(git config --local --get "${NAME}" || echo echo No configuration found for "${NAME}") "$@"



Contributing
------------

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Make sure that tests pass (`make test`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request


Contributors
------------

- Nico Di Rocco (https://github.com/nrocco)
