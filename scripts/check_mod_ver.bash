#!/bin/bash

# Author: Sergio Joselli
# Created: 1th June 2022
# Last Modified: 1th June 2022

# Description:
# Checks if go module is installed in GOPATH with the requested version
# Prints version to stdout if a module with the version is found in GOPATH
#
# It receives two positional parameters
#   1- the full go module path
#   2- the requested version
#
# Fails if
#   1- parameters are not provided
#   2- module path parameter contains illegal characters
#   3- version parameter not compatible with semver
#   e.g. version parameter should be x.x.x with optional concatenated tags -y 
#        where x is numeric and y is alphanumberic
#   4- module with requested version not found in GOPATH
# In case of failure it prints the error message to stdout

# Usage:
# check_mod_ver.bash MODULE_PATH VERSION
#
# Typically stdout is assigned to a variable
# result=$(check_mod_ver.bash github.com/someorg/somemodule 1.0.0) && {
#	  version=$result
# } || {
#	  error=$result
# }

module_path=$1
[[ -z $module_path ]] && {
	echo "error: missing module path positional parameter"
	exit 1
}

requested_version=$2
[[ -z $requested_version ]] && {
	echo "error: missing version positional parameter"
	exit 1
}

illegal_module_path_char='^.*[][,;\!@#$%¨&*=+`´|(){}<>"'\''[:space:]]+.*$'
[[ $module_path =~ $illegal_module_path_char ]] && {
	echo "error: illegal go module path characters in parameter '$module_path'"
	exit 1
}

# translate uppercase to lowercase with a leading exclamation mark to match go modules cache
module_path_cache_format=$(echo $module_path | sed --posix 's/[[:upper:]]/!&/g' | tr '[:upper:]' '[:lower:]')

semver_format='^[0-9]+\.[0-9]+\.[0-9]+(?:-[0-9a-z]+)*$'
[[ $requested_version =~ $semver_format ]] && {
	echo "error: invalid semver format in parameter '$requested_version'"
	exit 1
}

this_dir=$(dirname "${BASH_SOURCE[0]}")
go_paths_result=($($this_dir/get_go_paths.bash)) && {
	go_root=${go_paths_result[0]}
	go_path=${go_paths_result[1]}
} || {
	echo ${go_paths_result[*]}
	exit 1
}

if [[ !($(compgen -G $go_path"/pkg/mod/$module_path_cache_format*")) ]]; then
	echo "error: '$module_path' not found in go modules cache"
	exit 1
fi

for file in "$go_path"/pkg/mod/$module_path_cache_format*; do
	currently_read_version=${file##*@v}

	[[ $currently_read_version = $requested_version ]] && {
		echo $currently_read_version
		exit 0
	}
done

echo "error: module '$module_path' with version '$requested_version' not found"
exit 1
