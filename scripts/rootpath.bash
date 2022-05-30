#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 30th May 2022

# Description:
# Exports project root directory canonical path in the "rootpath" variable
# Note that the script is dependant on project structure

# Usage:
# source projectroot.bash

scriptpath="${BASH_SOURCE[0]}"
while [ -h "$scriptpath" ]; do scriptpath="$(readlink "$scriptpath")"; done
export rootpath="$( cd -P "$( dirname "$scriptpath" )/.." && pwd )"
