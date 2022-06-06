#!/bin/bash

# Author: Sergio Joselli
# Created: 6th June 2022
# Last Modified: 6th June 2022

# Description:
# Compiles contract ABI using solc docker image and write to build/contracts/CONTRACT_PACKAGE_NAME/abi
# Receives positional parameters
#   1- solidity version for contract
#   2- contract name
#   3- contract package name
#   4- truffle project canonical path where contract is located
#
# Fails if 
#   1- parameters are not provided
#   2- solidity version has illegal semver characters
#   3- contract name has any characters other than A-Za-z
#   4- contract package name has any characters other than A-Za-z
#   5- truffle project path is not a directory
#   6- docker is not in PATH
# In case of failure it prints the error message to stdout

# Usage:
# generate_abi.bash SOLIDITY_VERSION CONTRACT_NAME CONTRACT_PACKAGE_NAME TRUFFLE_PROJECT_PATH
#
# Typically stdout is assigned to a variable
# result=$(generate_abi.bash 0.8.11 SomeContract somecontract ~/sometruffleproject) || {
#	  error=$result
# }

echo "generating contract abi..."

solidity_version="$1"
contract_name="$2"
contract_package_name="$3"
truffle_root_path="$4"
[ -z "$truffle_root_path" ] && {
	echo "error: missing positional parameter(s)"
	exit 1
}

semver_format='^[0-9]+\.[0-9]+\.[0-9]+(-[.A-Za-z0-9]+)*$'
[[ $solidity_version =~ $semver_format ]] || {
	echo "error: invalid semver format in solidity version parameter '$solidity_version'"
	exit 1
}

allowed_contract_name='^[A-Za-z]+$'
[[ $contract_name =~ $allowed_contract_name ]] || {
	echo "error: illegal characters in contract name parameter '$contract_name'"
	echo "allowed characters are A-Za-z"
	exit 1
}

[[ $contract_package_name =~ $allowed_contract_name ]] || {
	echo "error: illegal characters in contract package name parameter '$contract_package_name'"
	echo "allowed characters are A-Za-z"
	exit 1
}

[ -d "$truffle_root_path" ] || {
	echo "error: truffle root path parameter '$truffle_root_path' is not a directory"
	exit 1
}

docker_root=$(which docker)
[ -z "$docker_root" ] && {
	echo "error: docker not installed or not in PATH"
	exit 1
}

root_path=$($(dirname "${BASH_SOURCE[0]}")/../get_root_path.bash)

mkdir -p "build/contracts/$contract_package_name/abi" && 
docker pull "ethereum/solc:$solidity_version" && 
docker run --rm -v "$root_path:/root" -v "$truffle_root_path:/truffle" "ethereum/solc:$solidity_version" "openzeppelin-solidity=/truffle/node_modules/openzeppelin-solidity" --abi "/truffle/contracts/$contract_name.sol" --overwrite -o "/root/build/contracts/$contract_package_name/abi" &&
echo "contract abi generated successfully"
