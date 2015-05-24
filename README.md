TODO: default git commit template that contains INFRA- on the first line


configuration
-------------

Configure git

    $ git config --global ide.ctrlp.ignore '.svn .git __pycache__'

    $ git config --global alias.ide '!/Users/nrocco/Develop/ide/bin/ide'
    $ git config --global alias.ctags '!git ide ctags'
    $ git config --global alias.ctrlp '!git ide ctrlp'
    $ git config --global alias.warmup '!f() { git ide ctags; git ide ctrlp; }; f'

    $ git config --global alias.phpchecker '!f() { git ls-files $@ | /bin/grep -E 'php$' | phpchecker; }; f'


git aliases
-----------
    [ide "ctags"]
        config = /home/nrocco/ide/ctags-php
