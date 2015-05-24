TODO: default git commit template that contains INFRA- on the first line


git aliases
-----------

    [alias]
        ctags = "!/home/nrocco/ide/general/ctags-warmup-cache.sh"
        ctrlp = "!/home/nrocco/ide/general/ctrlp-warmup-cache.sh"
        warmup = "!f() { git ctags; git ctrlp; }; f"
        phpchecker = "!f() { git ls-files $@ | /bin/grep -E 'php$' | phpchecker; }; f"

    [ide "ctrlp"]
        ignore = .svn .git

    [ide "ctags"]
        config = /home/nrocco/ide/ctags-php
