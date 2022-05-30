#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 30th May 2022

# Description:
# Installs geth tools (namely abigen) in go modules cache according to project geth version
#
# The installation won't happen if
#   1- go binary is not in the $PATH environment variable
#   2- the "go env GOPATH" command returns an empty string
#   3- go-ethereum is not installed in go modules cache
#   4- geth tools detected and only one go-ethereum version detected in go modules cache
#
# The installation will happen if none of the above and
#   1- geth tools not detected
#   2- geth tools detected but more than one go-ethereum version detected in go modules cache
#      which is a safe measure in case the abigen binary is from a different geth version than
#      the project

# Usage:
# install_geth_tools.bash

function install_geth_tools() {
	geth_path=$go_path/pkg/mod/github.com/ethereum/go-ethereum@v${geth_project_version}/
	sudo -E env "PATH=$$PATH" &&
	(cd $geth_path && make) &&
	(cd $geth_path && make devtools)
}

echo "installing geth tools..."

go_root=$(which go)
if [[ -z $go_root ]]; then
	echo "error: go not installed or not in PATH"
	exit 1
fi

go_path=$(go env GOPATH)
if [[ -z $go_path ]]; then
	echo "error: could not find GOPATH"
	exit 1
fi

if [[ $(compgen -G $go_path"/pkg/mod/github.com/ethereum/go-ethereum*") ]]; then
	echo "geth detected"
	else
	echo "error: geth not installed"
	echo "please run: go mod tidy"
	exit 1
fi

source $(dirname ${BASH_SOURCE[0]})/root_path.bash
go_mod_path=$root_path/go.mod
while IFS= read -r line || [ -n "$line" ]; do
	case $line in
		*/go-ethereum*) geth_project_version=${line#*go-ethereum v};;
		*) ;;
	esac

	if [[ -n $geth_project_version ]]; then
		break
	fi
done < "$go_mod_path"
echo "project using geth version $geth_project_version"

geth_versions_in_gopath=0
unset -v latest_geth_version
for file in "$go_path"/pkg/mod/github.com/ethereum/go-ethereum*; do
	geth_versions_in_gopath=$((geth_versions_in_gopath + 1))
  [[ $file -nt $latest_geth_version ]] && latest_geth_version=${file#*go-ethereum@v}
done
echo "latest geth gopath version $latest_geth_version"

if [[ -x $go_path/bin/abigen ]]; then
	echo "geth tools detected"
	else 
	echo "geth tools not detected"
	echo "installing..."
	install_geth_tools && exit 0
fi

if [[ geth_versions_in_gopath -gt 1 ]]; then
	echo "$geth_versions_in_gopath geth versions encountered in gopath"
	echo "geth tools binaries might be from a different geth version"
	echo "overwriting binaries..."
	install_geth_tools && exit 0
	else
	echo "only one geth version encontered in gopath"
	echo "installation not necessary"
	exit 0
fi
