#compdef ide
_ide() {
  local -a commands
  commands=(
    'destroy:Remove all ide configuration for a repository'
    'help:Help about any command'
    'hook:Manage git hooks for an ide project'
    'init:Initialize a git repository as an ide project'
    'exec:Manage executables for this ide project'
    'status:Get the current status of your ide project'
    'version:Get the version of ide'
    'server:Run an ide server for processing long running operations'
  )

  local -a hook_commands
  hook_commands=(
    'disable:Disable a git hook for this ide project'
    'enable:Enable a git hook for this ide project'
    'run:Run a git hook against an ide project'
  )

  local -a exec_commands
  exec_commands=(
    'add:Add an executable to this ide project'
    'exec:Execute a program in this ide projects environment'
    'rm:Remove an executable from this ide project'
  )

  local -a hook_run_commands
  hook_run_commands=(
    'commit-msg:Run the commit-msg hook for the ide project'
    'post-checkout:Run the post-checkout hook for the ide project'
    'post-merge:Run the post-merge hook for the ide project'
    'prepare-commit-msg:Run the prepare-commit-msg hook for the ide project'
  )

  if [[ CURRENT -eq 2 ]]
  then
    _describe -t commands 'commands' commands
  elif [[ CURRENT -eq 3 && $words[2] == 'hook' ]]
  then
    _describe -t hook_commands 'hook subcommands' hook_commands
  elif [[ CURRENT -eq 3 && $words[2] == 'exec' ]]
  then
    _describe -t exec_commands 'exec subcommands' exec_commands
  elif [[ CURRENT -eq 4 && $words[2] == 'hook' ]]
  then
    _describe -t hook_run_commands 'hook run subcommands' hook_run_commands
  elif [[ CURRENT -eq 4 && $words[2] == 'exec' && $words[3] != 'add' ]]
  then
    _values 'executables' $(ls .git/bin)
  fi

  return 0
}

_ide
