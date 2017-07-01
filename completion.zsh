#compdef ide
_ide() {
  local -a commands
  commands=(
    'destroy:Remove all ide configuration for a repository'
    'exec:Initialize a git repository as an ide project'
    'help:Help about any command'
    'init:Initialize a git repository as an ide project'
    'status:Get the current status of your ide project'
    'version:Get the version of ide'
  )

  local -a exec_commands
  exec_commands=(
    'help:Help about any command'
    'list:List passwords'
    'show:Show a single password'
    'permissions:Show permissions for a single password'
    'generate:Generate a strong, random password'
    'delete:Delete a single password'
    'lock:Lock a single password'
    'unlock:Unlock a single password'
  )

  local -a hook_commands
  hook_commands=(
    'help:Help about any command'
  )

  if (( CURRENT == 2 )); then
    _describe -t commands 'commands' commands
  elif (( CURRENT == 3)); then
    if [[ $words[2] == 'exec' ]]; then
        _describe -t exec_commands 'exec_commands' exec_commands
    elif [[ $words[2] == 'hook' ]]; then
        _describe -t hook_commands 'hook_commands' hook_commands
    fi
  fi

  return 0
}

_ide
