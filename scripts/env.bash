#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 3th June 2022

# Description:
# Exports environment variables defined in a file named "env" in project root
# Each line of the env file must be as follows
#   ENV_VARIABLE_NAME=ENV_VARIABLE_VALUE with an optional NEWLINE in the end (\n)
#   ENV_VARIABLE_NAME accepts upper case letters (A-Z) and underscores (_)
#   ENV_VARIABLE_VALUE accepts upper and lower case letters (A-Za-z), numbers (0-9),
#   underscores (_), forward slashes (/), colons (:), full stops (.), commas (,)
#   and dashes (-)

# Usage:
# source env.bash

echo "setting env variables..."

root_path=$($(dirname "${BASH_SOURCE[0]}")/root_path.bash)
env_path=$root_path/env

[ -f $env_path ] || {
	echo "error: env file not found"
	echo "there must be a file named 'env' in the project root directory"
	exit 1
}

has_equal_sign='^.*=.*$'
has_whitespace='^.*[ 	].*$'
empty_var_name='^=.*$'
allowed_name='^[A-Z_]+$'
empty_var_value='^.*=$'
allowed_value='^[0-9A-Za-z_/:.,\-]+$'

line_counter=0
while IFS= read -r line || [ -n "$line" ]; do
	line_counter=$(( $line_counter + 1 ))

	[[ $line =~ $has_equal_sign ]] || {
		echo "error: no '=' sign on line $line_counter"
		echo "lines must be formated like VARIABLE_NAME=VARIABLE_VALUE"
		exit 1
	}

	[[ $line =~ $has_whitespace ]] && {
		echo "error: whitespace character on line $line_counter"
		exit 1
	}

	[[ $line =~ $empty_var_name ]] && {
		echo "error: empty variable name on line $line_counter"
		exit 1
	}

	line_arr=($(echo $line | tr "=" "\n"))
	var_name=$(echo ${line_arr[0]} | tr -d "\n")
	var_value=$(echo ${line_arr[1]} | tr -d "\n")

	[[ $var_name =~ $allowed_name ]] || {
		echo "error: invalid variable name $var_name"
		echo "allowed characters are A-Z_"
		exit 1
	}

	[[ $line =~ $empty_var_value ]] && {
		echo "error: empty variable value for $var_name"
		exit 1
	}

	[[ $var_value =~ $allowed_value ]] || {
		echo "error: invalid variable value $var_value for $var_name"
		echo "allowed characters are A-Za-z0-9/:.,_-"
		exit 1
	}

	export $var_name=$var_value
	echo "$var_name=$var_value"
done < "$env_path"
