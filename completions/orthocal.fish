# Orthocal - Developed by dgm (dgm@tuta.com)
# orthocal/completions/orthocal.fish

set -l commands today tomorrow date saints readings hymns events info update serve export-web search version help
set -l search_categories saints western primary readings hymns events feasts fasts remembrances

complete -c orthocal -f
complete -c orthocal -l db -r -d 'Use a specific SQLite database'
complete -c orthocal -l plain -d 'Disable terminal styling'
complete -c orthocal -l json -d 'Print JSON for supported commands'
complete -c orthocal -l limit -r -d 'Limit search results'
complete -c orthocal -l addr -r -d 'Web server address'
complete -c orthocal -l force -d 'Allow update to install a newer schema version'
complete -c orthocal -l help -d 'Print help text'

for command in $commands
	complete -c orthocal -n "__fish_use_subcommand" -a "$command"
end

for category in $search_categories
	complete -c orthocal -n "__fish_seen_subcommand_from search" -a "$category"
end
