#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 6th June 2022

# Description:
# Prints project root directory canonical path to stdout
# Note that the script is dependant on project structure

# Usage:
# root_path.bash
#
# Typically the stdout is assigned to a variable
# root_path=$(root_path.bash)

script_path="${BASH_SOURCE[0]}"
while [ -h "$script_path" ]; do script_path="$(readlink "$script_path")"; done
echo "$( cd -P "$( dirname "$script_path" )/.." && pwd )"
