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

env_path=$1

case $env_path in
	*/) env_path="${env_path}env";;
	*) env_path="${env_path}/env";;
esac

if [[ !(-f $env_path) ]]; then
	echo "env file not found"
	exit 1
fi

allowed_name="^[A-Z_]+$"
allowed_value="^[0-9A-Za-z_/:.,\-]+$"

while IFS= read -r line || [ -n "$line" ]; do
	line_arr=($(echo $line | tr "=" "\n"))
	var_name=$(echo ${line_arr[0]} | tr -d "\n")
	var_value=$(echo ${line_arr[1]} | tr -d "\n")

	if [[ -z $var_name ]]; then
		echo "empty variable name"
		exit 1
	fi

	if [[ -z $var_value ]]; then
		echo "empty variable value for $var_name"
		exit 1
	fi

	if [[ !($var_name =~ $allowed_name) ]]; then
		echo "invalid variable name: $var_name"
		exit 1
	fi

	if [[ !($var_value =~ $allowed_value) ]]; then
		echo "invalid variable value: $var_value"
		exit 1
	fi

	export $var_name=$var_value
	echo "$var_name=$var_value"
done < "$env_path"
