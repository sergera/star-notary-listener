#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 30th May 2022

# Description:
# Starts application using go run

# Usage:
# dev_start.bash

source $(dirname ${BASH_SOURCE[0]})/root_path.bash &&
go run ${root_path}/cmd/app/main.go
