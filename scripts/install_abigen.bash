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
#   1- abigen binary with the same version as the geth version used in project is not in GOBIN
#   2- geth module with the same version as the geth version used in project is in go modules cache

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

dependency_version_result=$($this_dir/get_dependency_ver.bash go-ethereum) && {
	project_geth_version=$dependency_version_result
	echo "project using geth version $project_geth_version"
} || {
	echo $dependency_version_result
	exit 1
}

check_geth_module_version() {
	geth_version_result=$($this_dir/check_mod_ver.bash github.com/ethereum/go-ethereum $project_geth_version) && {
		echo "geth version $geth_version_result found in modules cache"
	} || {
		echo "error: geth version $project_geth_version not found in go modules cache"
		echo "please run: go mod tidy"
		exit 1
	}
}

install() {
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
		check_geth_module_version
		echo "overwriting abigen binary..."
		install && exit 0
	fi

else
	echo "abigen not detected in go binaries"
	check_geth_module_version
	echo "adding abigen binary..."
	install && exit 0
fi
