#!/bin/bash

# Author: Sergio Joselli
# Created: 6th June 2022
# Last Modified: 6th June 2022

# Description:
# Generates Go contract file into internal/gocontracts/CONTRACT_PACKAGE_NAME
# Receives positional parameters
#   1- contract name
#   2- contract package name
#
# Fails if 
#   1- parameters are not provided
#   2- contract name has any character other than A-Za-z
#   3- contract package name has any character other than A-Za-z
# In case of failure it prints the error message to stdout

# Usage:
# generate_bytecode.bash CONTRACT_NAME CONTRACT_PACKAGE_NAME
#
# Typically stdout is assigned to a variable
# result=$(generate_bytecode.bash 0.8.11 SomeContract somecontract ~/sometruffleproject) || {
#	  error=$result
# }

echo "generating go contract file..."

contract_name="$1"
contract_package_name="$2"
[ -z "$contract_package_name" ] && {
	echo "error: missing positional parameter(s)"
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

this_dir=$(dirname "${BASH_SOURCE[0]}")
go_paths_result=($("$this_dir"/../get_go_paths.bash)) && {
	go_root="${go_paths_result[0]}"
	go_path="${go_paths_result[1]}"
} || {
	echo "${go_paths_result[*]}"
	exit 1
}

mkdir -p "internal/gocontracts/$contract_package_name" && 
"$go_path"/bin/abigen -bin="build/contracts/$contract_package_name/bin/$contract_name.bin" --abi="build/contracts/$contract_package_name/abi/$contract_name.abi" --pkg="$contract_package_name" --out="internal/gocontracts/$contract_package_name/$contract_package_name.go" &&
echo "contract file generated successfully to 'internal/gocontracts/$contract_package_name/$contract_package_name.go'"
