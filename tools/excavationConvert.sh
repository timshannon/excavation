#!/bin/sh
# Author : shannon.timothy@gmail.com
# Date : 11/19/2012
# version 0.1
# This script converts collada files to horde format for excavation
# Copy this script in your ~/.gnome2/nautilus-scripts directory

#cp $NAUTILUS_SCRIPT_SELECTED_FILE_PATHS /home/myuser/mydirectory/
/home/tshannon/workspace/go/src/excavation/tools/ColladaConv $NAUTILUS_SCRIPT_SELECTED_FILE_PATHS -base /home/tshannon/workspace/go/src/excavation/data/ 
exit 0
