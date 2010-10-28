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
from tools.scene import *
import os
import sys
import copy
 
class Excavation(ShowBase):
    RUNNINGDIR = os.path.abspath(sys.path[0])
    MODELPATH = os.path.join(RUNNINGDIR, "../data/models/")
    SCENEPATH = os.path.join(RUNNINGDIR, "../data/scenes/")    
    def __init__(self):
        ShowBase.__init__(self)
        
        #load config file
        #set panda core settings
        #load keyconfig file
            #set keys
            
               
        if "-scene" in sys.argv:
            sceneFile = os.path.join(self.SCENEPATH, sys.argv[sys.argv.index("-scene") + 1])
            
            self.load_scene(sceneFile)
            
            
            
       
    def load_scene(self, fileName):
        """Loads the models, entities, lights, etc from the scene file."""
        scene = Scene(fileName)
                
        def load_node(node, parentNode):
            """recursively loads the nodes in the tree"""
            nodeP = None
            
            if type(node).__name__ == "Model":
                nodeP = self.loader.loadModel(os.path.joing(self.MODELPATH, node.model))
                nodeP.reparentTo(parentNode)
                nodeP.setPosHprScale(node.x,
                                     node.y,
                                     node.z,
                                     node.h,
                                     node.p,
                                     node.r,
                                     node.scaleX,
                                     node.scaleY,
                                     node.scaleZ)
            elif type(node).__name__ == "Node":
                nodeP = parentNode.attachNewNode(node.name)
            elif type(node).__name__ == "PointLight":
                light = PointLight(node.name)
                
                light.setColor(VBase4(node.color["red"], 
                                      node.color["green"], 
                                      node.color["blue"], 
                                      node.color["alpha"]))
                light.setSpecularColorColor(VBase4(node.color["red"], 
                                                   node.specColor["green"], 
                                                   node.specColor["blue"],
                                                   node.specColor["alpha"]))
                light.setAttenuation(Point3(node.attenuation["constant"],
                                            node.attenuation["linear"],
                                            node.attenuation["quadratic"]))
                
                nodeP = parentNode.attachNewNode(light)
                nodeP.setPos(node.x,
                             node.y,
                             node.z)
                #For now each light will light everything under it's parent
                parentNode.setLight(nodeP)
            elif type(node).__name__ == "Spotlight":
                sLight = Spotlight(node.name)
                
                sLight.setColor(VBase4(node.color["red"], 
                                      node.color["green"], 
                                      node.color["blue"], 
                                      node.color["alpha"]))
                sLight.setSpecularColorColor(VBase4(node.color["red"], 
                                                    node.specColor["green"], 
                                                    node.specColor["blue"],
                                                    node.specColor["alpha"]))
                sLight.setAttenuation(Point3(node.attenuation["constant"],
                                             node.attenuation["linear"],
                                             node.attenuation["quadratic"]))
                                
                nodeP = parentNode.attachNewNode(sLight)
                nodeP.setPosHpr(node.x,
                                node.y,
                                node.z,
                                node.h,
                                node.p,
                                node.r)
                #For now each light will light everything under it's parent
                parentNode.setLight(nodeP)
            
            
            for c in node.children:
                load_node(c, nodeP)
                    
                    
        load_node(scene.tree, self.render)
            
        
        
        
 
main = Excavation()
main.run()



        
        
        
        
        
        
        
        
        
        
        
        