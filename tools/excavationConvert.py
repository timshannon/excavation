#!/usr/bin/env python
# Author : shannon.timothy@gmail.com
# Date : 11/19/2012
# version 0.1
# This script converts collada files to horde format for excavation
# Copy this script in your ~/.local/share/nautilus/scripts directory

#cp $NAUTILUS_SCRIPT_SELECTED_FILE_PATHS /home/myuser/mydirectory/
import subprocess
import string
import os

base = "/home/tshannon/workspace/go/src/excavation/data/"
toolpath = "/home/tshannon/workspace/go/src/excavation/tools/ColladaConv"
selected = os.environ.get('NAUTILUS_SCRIPT_SELECTED_FILE_PATHS', '')
files = string.strip(string.replace(selected, base, ""))
subprocess.call([toolpath, files, "-base" , base, "-dest" , base])



