#!/bin/bash

# Author: Sergio Joselli
# Created: 1th June 2022
# Last Modified: 1th June 2022

# Description:
# Prints version of a project dependency present in go.mod to stdout
# It receives a positional parameter with the dependency module path/name/substring to match
# Only the first match is considered, so it's best to provide the full module path
#
# Fails if
#   1- parameter is not provided
#   2- parameter contains illegal characters to a go module path
#   3- dependency not found in go.mod file
# In case of failure it prints the error message to stdout
#
# IMPORTANT:
# - go.mod file must be formatted according to gofmt
# - replace and exclude statements are not considered

# Usage:
# get_dependency_ver.bash DEPENDENCY_MODULE_PATH
#
# Typically stdout is assigned to a variable
# result=$(get_dependency_ver.bash github.com/someorg/somemodule) && {
#	  version=$result
# } || {
#	  error=$result
# }

this_dir=$(dirname "${BASH_SOURCE[0]}")
project_go_mod_path=$($this_dir/root_path.bash)/go.mod

module_path=$1
[[ -z $module_path ]] && {
	echo "error: missing dependency module path positional parameter"
	exit 1
}

illegal_module_path_char='^.*[][,;\!@#$%¨&*=+`´|(){}<>"'\''[:space:]]+.*$'
[[ $module_path =~ $illegal_module_path_char ]] && {
	echo "error: illegal go mod path characters in parameter $module_path"
	exit 1
}

dependency="^ *require *.*$"
dependency_group_begin="^ *require *\( *$"
dependency_group_end="^ *\) *$"

is_one_liner=false
is_dependency_group_start=false
is_dependency_group_middle=false
while IFS= read -r line || [ -n "$line" ]; do
	[[ $line =~ $dependency_group_begin ]] && is_dependency_group_start=true
	[[ $line =~ $dependency_group_end ]] && is_dependency_group_middle=false
	[[ $line =~ $dependency ]] && [[ $is_dependency_group_start = false ]] && is_one_liner=true

	if [[ $is_one_liner = true  ]] || [[ $is_dependency_group_middle = true ]]; then
		case $line in
			*$module_path*)
				dependency_version=${line##* v}
				dependency_version=${dependency_version%% //*};;
			*);;
		esac
	fi

	[[ -n $dependency_version ]] && break

	[[ $is_one_liner = true ]] && is_one_liner=false
	[[ $is_dependency_group_start = true ]] && {
		is_dependency_group_start=false
		is_dependency_group_middle=true
	}
done < "$project_go_mod_path"

[[ -z $dependency_version ]] && {
	echo "error: $module_path dependency not found"
	exit 1
}

echo $dependency_version
exit 0
