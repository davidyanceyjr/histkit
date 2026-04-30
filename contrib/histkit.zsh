histkit_pick_zsh() {
  local selected
  selected="$(histkit pick "$@")" || return $?

  if [[ -z "$selected" ]]; then
    zle redisplay
    return 0
  fi

  BUFFER="$selected"
  CURSOR=${#BUFFER}
  zle redisplay
}

histkit_bind_zsh_pick() {
  zle -N histkit_pick_zsh
  bindkey '^R' histkit_pick_zsh
}
