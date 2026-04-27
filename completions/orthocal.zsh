#compdef orthocal
# Orthocal - Developed by dgm (dgm@tuta.com)
# orthocal/completions/orthocal.zsh

_orthocal()
{
	local -a commands search_categories
	commands=(
		'today:show the local system date'
		'tomorrow:show the local system date plus one'
		'date:show a specific Gregorian date'
		'saints:show saints for a Gregorian date'
		'readings:show scripture readings for a Gregorian date'
		'hymns:show hymns for a Gregorian date'
		'events:show events for a Gregorian date'
		'info:show database metadata and counts'
		'update:replace the configured database'
		'serve:start the local read-only web server'
		'export-web:export a static read-only website'
		'search:search calendar data'
		'version:show version and build information'
		'help:show help'
	)
	search_categories=(saints western primary readings hymns events feasts fasts remembrances)

	_arguments -C \
		'--db[use a specific SQLite database]:database:_files' \
		'--plain[disable terminal styling]' \
		'--json[print JSON for supported commands]' \
		'--limit[limit search results]:limit:' \
		'--addr[web server address]:address:' \
		'--force[allow update to install a newer schema version]' \
		'--help[print help text]' \
		'1:command:->command' \
		'2:argument:->argument'

	case "$state" in
		command)
			_describe 'command' commands
			;;
		argument)
			if [[ "$words[2]" == search ]]; then
				_values 'search category' $search_categories
			fi
			;;
	esac
}

_orthocal "$@"
