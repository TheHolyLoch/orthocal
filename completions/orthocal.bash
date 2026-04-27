# Orthocal - Developed by dgm (dgm@tuta.com)
# orthocal/completions/orthocal.bash

_orthocal()
{
	local cur prev commands flags search_categories
	COMPREPLY=()
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[COMP_CWORD-1]}"
	commands="today tomorrow date saints readings hymns events info update serve export-web search version help"
	flags="--db --plain --json --limit --addr --force --help"
	search_categories="saints western primary readings hymns events feasts fasts remembrances"

	case "$prev" in
		--db|--addr|--limit|update|export-web)
			return 0
			;;
		search)
			COMPREPLY=( $(compgen -W "$search_categories" -- "$cur") )
			return 0
			;;
	esac

	if [[ "$cur" == --* ]]; then
		COMPREPLY=( $(compgen -W "$flags" -- "$cur") )
	else
		COMPREPLY=( $(compgen -W "$commands $flags" -- "$cur") )
	fi
}

complete -F _orthocal orthocal
