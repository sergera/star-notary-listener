#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 1th June 2022

# Description:
# Installs abigen according to the geth version used in the project
#
# Abigen is necessary for go contract generation
# This script will add its binary in GOBIN if it's not there
# And will overwrite the binary if its version does not match the project's geth dependency
# It runs "go install" in the "abigen" directory of the geth main package
#
# The installation will only happen if
#   1- go root binary directory is in the $PATH environment variable
#   2- the "go env GOPATH" command returns a non-empty string
#   3- abigen binary with the same version as the geth version used in project is not in GOBIN
#   4- geth module with the same version as the geth version used in project is in go modules cache

# Usage:
# install_abigen.bash

echo "installing abigen..."

this_dir=$(dirname "${BASH_SOURCE[0]}")

go_paths_result=($($this_dir/get_go_paths.bash)) && {
	go_root=${go_paths_result[0]}
	go_path=${go_paths_result[1]}
} || {
	echo ${go_paths_result[*]}
	exit 1
}

geth_version_result=$($this_dir/get_dependency_ver.bash go-ethereum) && {
	project_geth_version=$geth_version_result
	echo "project using geth version $project_geth_version"
} || {
	echo $geth_version_result
	exit 1
}

function check_geth_module_versions() {
	if [[ !($(compgen -G $go_path"/pkg/mod/github.com/ethereum/go-ethereum*")) ]]; then
		echo "error: geth not found in go modules cache"
		echo "please run: go mod tidy"
		exit 1
	fi

	gopath_has_project_version=false
	for file in "$go_path"/pkg/mod/github.com/ethereum/go-ethereum*; do
		currently_read_geth_version=${file#*@v}
		echo "found geth version $currently_read_geth_version in go modules cache"

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
	(cd $geth_path && $(go install ./cmd/abigen))
}

if [[ -x $go_path/bin/abigen ]]; then
	abigen_version=$($go_path/bin/abigen -v)
	abigen_version=${abigen_version##* }
	abigen_version=${abigen_version%-*}
	echo "abigen version $abigen_version"

	if [[ $project_geth_version = $abigen_version ]]; then
		echo "abigen already installed with correct version"
		exit 0
		else
		echo "abigen version differs from geth version used in project"
		check_geth_module_versions
		echo "overwriting abigen binary..."
		install && exit 0
	fi

	else
	echo "abigen not detected in go binaries"
	check_geth_module_versions
	echo "adding abigen binary..."
	install && exit 0
fi
