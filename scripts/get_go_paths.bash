#!/bin/bash

# Author: Sergio Joselli
# Created: 1th June 2022
# Last Modified: 1th June 2022

# Description:
# Prints GOROOT and GOPATH to stdout, in this order
#
# Fails if
#   1- go binary is not in PATH
#   2- "go env GOPATH" command returns nothing
# In case of failure it prints the error message to stdout

# Usage:
# get_go_paths.bash
#
# Typically stdout is assigned to an array
# result=($(get_go_paths.bash)) && {
#	  go_root=${result[0]}
#   go_path=${result[1]}
# } || {
#	  error=${result[*]}
# }

go_root=$(which go)
[[ -z $go_root ]] && {
	echo "error: go not installed or not in PATH"
	exit 1
}

go_path=$(go env GOPATH)
[[ -z $go_path ]] && {
	echo "error: could not get GOPATH"
	exit 1
}

echo $go_root
echo $go_path
exit 0
