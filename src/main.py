#Copyright (c) 2009-2010 Tim Shannon
#
#Permission is hereby granted, free of charge, to any person obtaining a copy
#of this software and associated documentation files (the "Software"), to deal
#in the Software without restriction, including without limitation the rights
#to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
#copies of the Software, and to permit persons to whom the Software is
#furnished to do so, subject to the following conditions:
#
#The above copyright notice and this permission notice shall be included in
#all copies or substantial portions of the Software.
#
#THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
#IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
#FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
#AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
#LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
#OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
#THE SOFTWARE.


from direct.showbase.ShowBase import ShowBase
from panda3d.core import WindowProperties
from panda3d.core import ConfigVariableString
from tools.scene import Scene
from utility.globalDef import GlobalDef
import os
import sys
import copy
 
class Excavation(ShowBase):
    def __init__(self):
        ShowBase.__init__(self)
        
        #load config file
        #set panda core settings
        #load keyconfig file
            #set keys
            
               
        if "-scene" in sys.argv:
            sceneFile = os.path.join(GlobalDef.SCENEPATH, sys.argv[sys.argv.index("-scene") + 1])
            
            self.loadScene(sceneFile)
            
            
            
       
    def loadScene(self, fileName):
        """Loads the models, entities, lights, etc from the scene file."""
        scene = Scene(fileName)
        
        #TODO:  Show Loading Screen maybe on a different camera
        scene.loadScene(self.render)
        #TODO: Remove loading screen / camera
        
    
            
        
        
        
 
main = Excavation()
main.run()



        
        
        
        
        
        
        
        
        
        
        
        