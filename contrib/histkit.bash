__histkit_pick_bash() {
  local selected
  selected="$(histkit pick "$@")" || return $?

  if [[ -z "$selected" ]]; then
    return 0
  fi

  READLINE_LINE="$selected"
  READLINE_POINT=${#READLINE_LINE}
}

histkit_bind_bash_pick() {
  bind -x '"\C-r":"__histkit_pick_bash"'
}
