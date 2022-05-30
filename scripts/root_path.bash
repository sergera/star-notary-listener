#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 30th May 2022

# Description:
# Exports project root directory canonical path in the "rootpath" variable
# Note that the script is dependant on project structure

# Usage:
# source root_path.bash

script_path="${BASH_SOURCE[0]}"
while [ -h "$script_path" ]; do script_path="$(readlink "$script_path")"; done
export root_path="$( cd -P "$( dirname "$script_path" )/.." && pwd )"
