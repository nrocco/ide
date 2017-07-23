#compdef ide
_ide() {
  local -a commands
  commands=(
    'destroy:Remove all ide configuration for a repository'
    'help:Help about any command'
    'hook:Manage git hooks for an ide project'
    'init:Initialize a git repository as an ide project'
    'link:Manage link to executables for this ide project'
    'status:Get the current status of your ide project'
    'version:Get the version of ide'
  )

  local -a hook_commands
  hook_commands=(
    'disable:Disable a git hook for this ide project'
    'enable:Enable a git hook for this ide project'
    'run:Run a git hook against an ide project'
  )

  local -a link_commands
  link_commands=(
    'add:Link to an executable and add it to this ide project'
    'exec:Execute a linked program in this ide projects environment'
    'rm:Remove a linked program from this ide project'
  )

  local -a hook_run_commands
  hook_run_commands=(
    'commit-msg:Run the commit-msg hook for the ide project'
    'post-checkout:Run the post-checkout hook for the ide project'
    'prepare-commit-msg:Run the prepare-commit-msg hook for the ide project'
  )

  if [[ CURRENT -eq 2 ]]
  then
    _describe -t commands 'commands' commands
  elif [[ CURRENT -eq 3 && $words[2] == 'hook' ]]
  then
    _describe -t hook_commands 'hook subcommands' hook_commands
  elif [[ CURRENT -eq 3 && $words[2] == 'link' ]]
  then
    _describe -t link_commands 'link subcommands' link_commands
  elif [[ CURRENT -eq 4 && $words[2] == 'hook' ]]
  then
    _describe -t hook_run_commands 'hook run subcommands' hook_run_commands
  elif [[ CURRENT -eq 4 && $words[2] == 'link' && $words[3] != 'add' ]]
  then
    _values 'linked executables' $(ls .git/bin)
  fi

  return 0
}

_ide
