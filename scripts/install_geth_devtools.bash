#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 31th May 2022

# Description:
# Installs geth dev tools according to the geth version used in the project
#
# Go contract generation is handled by a tool called "abigen",
# this script will add it, among other tools, to go binaries
# it uses the "devtools" Makefile rule from go-ethereum
#
# The installation will only happen if
#   1- go root binary directory is in the $PATH environment variable
#   2- the "go env GOPATH" command returns a non-empty string
#   3- abigen binary not present in go binaries or has a different version than project geth
#   4- project geth version is in go modules cache

# Usage:
# install_geth_devtools.bash

echo "installing geth tools..."

go_root=$(which go)
if [[ -z $go_root ]]; then
	echo "error: go not installed or not in PATH"
	exit 1
fi

go_path=$(go env GOPATH)
if [[ -z $go_path ]]; then
	echo "error: could not get path to go workspace"
	exit 1
fi

source $(dirname ${BASH_SOURCE[0]})/root_path.bash
project_go_mod_path=$root_path/go.mod
while IFS= read -r line || [ -n "$line" ]; do
	case $line in
		*/go-ethereum*) project_geth_version=${line#* v};;
		*) ;;
	esac

	if [[ -n $project_geth_version ]]; then
		break
	fi
done < "$project_go_mod_path"
echo "project using geth version $project_geth_version"

function check_geth_module_versions() {
	if [[ !($(compgen -G $go_path"/pkg/mod/github.com/ethereum/go-ethereum*")) ]]; then
		echo "error: geth not found in go modules cache"
		echo "please run: go mod tidy"
		exit 1
	fi

	gopath_has_project_version=false
	for file in "$go_path"/pkg/mod/github.com/ethereum/go-ethereum*; do
		currently_read_geth_version=${file#*@v}
		echo "go modules cache has geth version $currently_read_geth_version"

		if [[ $project_geth_version = $currently_read_geth_version ]]; then
			gopath_has_project_version=true
		fi
	done

	if [[ !($gopath_has_project_version = true) ]]; then
		echo "error: geth version used in project not detected in go modules cache"
		echo "please run: go mod tidy"
		exit 1
	fi
}

function install() {
	geth_path=$go_path/pkg/mod/github.com/ethereum/go-ethereum@v$project_geth_version
	(cd $geth_path && sudo -E env "PATH=$PATH" make devtools)
}

if [[ -x $go_path/bin/abigen ]]; then
	abigen_version=$($go_path/bin/abigen -v)
	abigen_version=${abigen_version##* }
	abigen_version=${abigen_version%-*}
	echo "geth tools version $abigen_version"

	if [[ $project_geth_version = $abigen_version ]]; then
		echo "geth tools already installed with correct version"
		exit 0
		else
		echo "geth tools version differs from geth version used in project"
		check_geth_module_versions
		echo "overwriting binaries..."
		install && exit 0
	fi

	else
	echo "geth tools not detected in go binaries"
	check_geth_module_versions
	echo "adding binaries..."
	install && exit 0
fi
