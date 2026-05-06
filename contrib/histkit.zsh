histkit_pick_zsh() {
  local selected
  zle -I
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
  local keyseq
  keyseq="${1:-^R}"
  zle -N histkit_pick_zsh
  bindkey "$keyseq" histkit_pick_zsh
}
