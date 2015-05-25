# ide

ide provides a powerful ide that gets out of your way.

> a tool that glues vim/git/ctags together to provide a powerful integrated
> development environment


## installation

TODO: write this section


## configuration

Configure git

    $ git config --global ide.ctrlp.ignore '.svn .git __pycache__'

    $ git config --global alias.ide '!/Users/nrocco/Develop/ide/bin/ide'
    $ git config --global alias.ctags '!git ide ctags'
    $ git config --global alias.ctrlp '!git ide ctrlp'
    $ git config --global alias.warmup '!f() { git ide ctags; git ide ctrlp; }; f'

    $ git config --global alias.phpchecker '!f() { git ls-files $@ | /bin/grep -E 'php$' | phpchecker; }; f'


git aliases

    [ide "ctags"]
        config = /home/nrocco/ide/ctags-php


## usage

To get help:

    $ git ide help
    usage: git ide [ctags ctrlp help symfony2]

    ide provides a powerful ide that gets out of your way.

    available sub commands:

        help            print this help and exit
        ctags           (re)generates the tags file
        ctrlp           (re)generates the ctrlp cache file
        symfony         provides checks for symfony2 based projects


## todos

TODO: default git commit template that contains INFRA- on the first line

