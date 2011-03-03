import os
import sys

__all__ = ["GlobalDef"]

class GlobalDef(object):
    '''
        Shared information across all of Excavation code
    '''
    
    RUNNINGDIR = os.path.abspath(sys.path[0])
    MODELPATH = os.path.join(RUNNINGDIR, "../data/models/")
    SCENEPATH = os.path.join(RUNNINGDIR, "../data/scenes/")    
    
