#compdef ide
_ide() {
  local -a commands
  commands=(
    'destroy:Remove all ide configuration for a repository'
    'help:Help about any command'
    'hook:Manage git hooks for an ide project'
    'init:Initialize a git repository as an ide project'
    'status:Get the current status of your ide project'
    'version:Get the version of ide'
  )

  local -a hook_commands
  hook_commands=(
    'help:Help about any command'
    'disable:Disable a git hook for this ide project'
    'enable:Enable a git hook for this ide project'
    'run:Run a git hook against an ide project'
    'list:List all available hooks'
  )

  local -a hook_run_commands
  hook_run_commands=(
    'commit-msg:Run the commit-msg hook for the ide project'
    'post-checkout:Run the post-checkout hook for the ide project'
    'prepare-commit-msg:Run the prepare-commit-msg hook for the ide project'
  )

  if (( CURRENT == 2 ))
  then
    _describe -t commands 'commands' commands
  elif (( CURRENT == 3))
  then
    if [[ $words[2] == 'hook' ]]
    then
        _describe -t hook_commands 'hook_commands' hook_commands
    fi
  elif (( CURRENT == 4))
  then
    if [[ $words[2] == 'hook' && $words[3] == 'run' ]]
    then
        _describe -t hook_run_commands 'hook_run_commands' hook_run_commands
    fi
  fi

  return 0
}

_ide
