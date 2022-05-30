#!/bin/bash

# Author: Sergio Joselli
# Created: 30th May 2022
# Last Modified: 30th May 2022

# Description:
# Starts application using go run

# Usage:
# devstart.bash

source $(dirname ${BASH_SOURCE[0]})/rootpath.bash &&
go run ${rootpath}/cmd/app/main.go
