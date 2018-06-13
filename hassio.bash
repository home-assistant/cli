# bash completion for hassio

_hassio()
{
    local cur=${COMP_WORDS[COMP_CWORD]} prev=${COMP_WORDS[COMP_CWORD-1]}
    local -a cmds=(homeassistant supervisor host hardware network snapshots
                   addons help)
    local -a opts
    local i cmd action

    # Find out command and action
    for ((i=1; i < COMP_CWORD; i++)); do
        if [[ ${COMP_WORDS[i]} != -* ]]; then
            if [[ -z $cmd ]]; then
                cmd=${COMP_WORDS[i]}
            else
                action=${COMP_WORDS[i]}
                break
            fi
        fi
    done

    # Process top level commands and options
    if [[ -z $cmd ]]; then
        case $cur in
            -*) opts=(--debug --help --version) ;;
            *)  opts=("${cmds[@]}") ;;
        esac
        COMPREPLY=($(compgen -W '${opts[@]}' -- "$cur"))
        return
    fi

    # Handle common command options
    case $prev in --options|-o|--filter|-f|--help) return ;; esac

    # Process command options and actions
    case $cmd in
        homeassistant|ha)
            case $cur in
                -*) opts=(--rawjson --options --filter --help) ;;
                *)  [[ -z $action ]] &&
                        opts=(info logs check restart start stop update) ;;
            esac
            ;;
        supervisor|su)
            case $cur in
                -*) opts=(--rawjson --options --filter --help) ;;
                *)  [[ -z $action ]] && opts=(info logs reload update) ;;
            esac
            ;;
        host|ho)
            case $cur in
                -*) opts=(--rawjson --options --filter --help) ;;
                *)  [[ -z $action ]] && opts=(reboot shutdown update) ;;
            esac
            ;;
        hardware|hw)
            case $cur in
                -*) opts=(--rawjson --filter --help) ;;
                *)  [[ -z $action ]] && opts=(info audio) ;;
            esac
            ;;
        network|ne)
            case $cur in
                -*) opts=(--rawjson --options --filter --help) ;;
                *)  [[ -z $action ]] && opts=(info options) ;;
            esac
            ;;
        snapshots|sn)
            if [[ $prev =~ ^--name|--password$ ]]; then
                :  # nothing
            elif [[ $prev == --slug ]]; then
                opts=($("$1" snapshots list \
                            | jq -r .data.snapshots[].slug 2>/dev/null))
            elif [[ $prev =~ ^info|restore|remove$ ]]; then
                opts=(--slug)
            else
                case $cur in
                    -*) opts=(--rawjson --options --filter --slug --name
                              --password --help) ;;
                    *)  [[ -z $action ]] &&
                            opts=(list info reload new restore remove) ;;
                esac
            fi
            ;;
        addons|ad)
            if [[ $prev == --name ]]; then
                opts=($("$1" addons list | \
                            jq -r .data.addons[].slug 2>/dev/null))
            elif [[ $prev == info ]]; then
                opts=(--name)
            else
                case $cur in
                    -*) opts=(--rawjson --options --filter --name --help) ;;
                    *)  [[ -z $action ]] &&
                            opts=(list info logo changelog logs stats reload
                                  start stop install uninstall update) ;;
                esac
            fi
            ;;
        help|h)
            [[ -z $action ]] && opts=("${cmds[@]}")
            ;;
    esac

    COMPREPLY=($(compgen -W '${opts[@]}' -- "$cur"))
} &&
complete -F _hassio hassio
