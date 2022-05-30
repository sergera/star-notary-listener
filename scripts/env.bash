#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 30th May 2022

# Description:
# Exports environment variables defined in a file named "env"
# Receives one positional parameter with path to directory containing env file
# Each line of the env file must be as follows
#   ENV_VARIABLE_NAME=ENV_VARIABLE_VALUE with an optional NEWLINE in the end (\n)
#		ENV_VARIABLE_NAME accepts upper case letters (A-Z) and underscores (_)
#		ENV_VARIABLE_VALUE accepts upper and lower case letters (A-Za-z), numbers (0-9), 
#		underscores (_), forward slashes (/), colons (:), full stops (.), commas (,)
#		and dashes (-)

# Usage:
# source env.bash DIRECTORYPATH

echo "setting env variables"

envpath=$1

case $envpath in
	*/) envpath="${envpath}env";;
	*) envpath="${envpath}/env";;
esac

if [[ !(-f $envpath) ]]; then
	echo "env file not found"
	exit 1
fi

allowedname="^[A-Z_]+$"
allowedvalue="^[0-9A-Za-z_/:.,\-]+$"

while IFS= read -r line || [ -n "$line" ]; do
	linearr=($(echo $line | tr "=" "\n"))
	varname=$(echo ${linearr[0]} | tr -d "\n")
	varvalue=$(echo ${linearr[1]} | tr -d "\n")

	if [[ -z $varname ]]; then
		echo "empty variable name"
		exit 1
	fi

	if [[ -z $varvalue ]]; then
		echo "empty variable value for $varname"
		exit 1
	fi

	if [[ !($varname =~ $allowedname) ]]; then
		echo "invalid variable name: $varname"
		exit 1
	fi

	if [[ !($varvalue =~ $allowedvalue) ]]; then
		echo "invalid variable value: $varvalue"
		exit 1
	fi

	export $varname=$varvalue
	echo "$varname=$varvalue"
done < "$envpath"
