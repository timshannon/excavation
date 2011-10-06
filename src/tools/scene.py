#Copyright (c) 2009 Tim Shannon
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

import cPickle
import os
import json
from utility.globalDef import GlobalDef
from panda3d.core import VBase4, Point3

__all__ = ["Scene", "Node", "Light", "Spotlight", "Entity", "PointLight", "DirectionalLight"]

class Scene(): 
    #Defines the structure of a scene for saving to a file for loading

    def __init__(self, fileName=""):
        if fileName == "":
            self.tree = Node("render")
            self.keyValues = {}  #dictionary for scene level values, ambient light level, skybox, name, level description, level load hints, etc
        else:
            self.read(fileName)
    
    def write(self, fileName):
        file = open(fileName, "wb")
        #TODO: Replace pickle with json file format
        json.dump(self.toJson(), file, indent=4, sort_keys=True) 
        file.close()
        
    def load(self, fileName):
        file = open(fileName, "rb")
        scene = cPickle.load(file)
        file.close()
        
        self.keyValues = scene.keyValues
        self.tree = scene.tree
       
    def toJson(self):
        return self.tree.toJson()

    def loadScene(self, render):
        self.loadNode(self.tree, render)
        
    def loadNode(self, node, parentNode):
        """recursively loads the nodes in the tree"""
        nodeP = None
        
        if type(node).__name__ == "Model":
            nodeP = self.loader.loadModel(os.path.join(GlobalDef.MODELPATH, node.model))
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
        elif type(node).__name__ == "DirectionalLight":
            dLight = DirectionalLight(node.name)
            
            dLight.setColor(VBase4(node.color["red"], 
                                  node.color["green"], 
                                  node.color["blue"], 
                                  node.color["alpha"]))
            dLight.setSpecularColorColor(VBase4(node.color["red"], 
                                                node.specColor["green"], 
                                                node.specColor["blue"],
                                                node.specColor["alpha"]))
            dLight.setDirection(Point3(node.direction["x"], 
                                       node.direction["y"],
                                       node.direction["z"]))
                            
            nodeP = parentNode.attachNewNode(dLight)
            nodeP.setPosHpr(node.x,
                            node.y,
                            node.z,
                            node.h,
                            node.p,
                            node.r)
            #For now each light will light everything under it's parent
            parentNode.setLight(nodeP)
        elif type(node).__name__ == "Entity":
            #lookup entity class and instantiate it
            #TODO: Finish this
            pass


        for c in node.children:
            self.loadNode(c, nodeP)


class Node(object):
    def __init__(self, 
                 name, 
                 parent=None, 
                 x=-1, 
                 y=-1,
                 z=-1,
                 h=-1,
                 p=-1,
                 r=-1):
        self.name = name
        self.x = x
        self.y = y
        self.z = z
        self.h = h
        self.p = p
        self.r = r
        self.parent = parent
        self.children = []
        self.layer = 0
        
        if self.parent <> None:
            self.parent.addChild(self)
            
            
    def addChild(self, child):
        if child not in self.children:
            self.children.append(child)
            if child.parent == None:
                child.parent = self
                
    def getSiblings(self):
        if self.parent == None:
            return None
        else:
            siblings = self.parent.children[:]
            siblings.remove(self)
            return siblings
    
    def setPos(self, x, y, z):
        """Sets position of the node"""
        self.x = x
        self.y = y
        self.z = z
            
    def setHpr(self, h, p, r):
        """Sets heading, pitch and rotation"""
        self.h = h
        self.p = p
        self.r = r

    def toJson(self):
        jChildren = []
        for child in self.children:
            jChildren.append(child.toJson())

        return {'name':self.name,
                'x':self.x,
                'y':self.y,
                'z':self.z,
                'h':self.h,
                'p':self.p,
                'r':self.r,
                'layer':self.layer,
                'children':jChildren}

                    
class Entity(Node):
    #keyvalue dictionary to hold any settings the entity may make use of
    keyValues = {} 
    #type of entity used to identify to python file to use?
    type = ""       

    def __init__(self,
                 keyValues,
                 type):
        self.keyValues = keyValues
        self.type = type

    def toJson(self):
        dict = super(Entity, self).toJson()
        dict['keyValues'] = self.keyValues
        dict['type'] = self.type

        return dict
    
class Model(Node):
    model = ""
    collision = ""  
    scaleX = 1.0
    scaleY = 1.0
    scaleZ = 1.0

    def __init__(self,
                 model,
                 collision,
                 scaleX,
                 scaleY,
                 scaleZ):
        self.model = model
        self.collision = collision
        self.scaleX = scaleX
        self.scaleY = scaleY
        self.scaleZ = scaleZ
    
    def setScale(self, **scale):
        if "x" in scale:
            self.scaleX = scale["x"]
        
        if "y" in scale:
            self.scaley = scale["y"]
            
        if "z" in scale:
            self.scaleZ = scale["z"]

    def toJson(self):
        dict = super(Model, self).toJson()
        dict['model'] = self.model
        dict['collision'] = self.collision
        dict['scaleX'] = self.scaleX
        dict['scaleY'] = self.scaleY
        dict['scaleZ'] = self.scaleZ

        return dict
                  
    
class Light(Node):
    color = {"red":1,"green":1,"blue":1,"alpha":1}
    specColor = {"red":1,"green":1,"blue":1,"alpha":1}
    
    def __init__(self,
                 color,
                 specColor):
        self.color = color
        self.specColor = specColor

    def setColor(self, **colors):
        for k in colors.keys():
            if k in self.color.keys():
                self.color[k] = colors[k]
                
    def setSpecColor(self, **colors):
        for k in colors.keys():
            if k in self.specColor.keys():
                self.specColor[k] = colors[k]

    def toJson(self):
        dict = super(Light, self).toJson()
        dict['color'] = self.color
        dict['specColor'] = self.specColor

        return dict

class PointLight(Light):
    attenuation = {"constant":0,"linear":0,"quadratic":0}
   
    def __init__(self,
                 attenuation):
        self.attenuation = attenuation

    def setAttenuation(self, **attenuation):
        for k in attenuation.keys():
            if k in self.attenuation.keys():
                self.attenuation[k] = attenuation[k]

    def toJson(self):
        dict = super(PointLight, self).toJson()
        dict['attenuation'] = self.attenuation

        return dict
    
class DirectionalLight(Light):
    direction = {"x":0,"y":0,"z":0}
    castShadows = False
   
    def __init__(self,
                 direction,
                 castShadows):
        self.direction = direction
        self.castShadows = castShadows
        
    def setDirection(self, **direction):
        for k in direction.keys():
            if k in self.direction.keys():
                self.direction[k] = direction[k]

    def toJson(self):
        dict = super(DirectionalLight, self).toJson()
        dict['direction'] = self.direction
        dict['castShadows'] = self.castShadows
        
        return dict

class Spotlight(Light):
    attenuation = {"constant":0,"linear":0,"quadratic":0}
    exponent = 0.0
    castShadows = False
   
    def __init_(self,
                attenuation,
                exponent,
                castShadows):
        self.attenuation = attenuation
        self.exponent = exponent
        self.castShadows = castShadows
        
    def setAttenuation(self, **attenuation):
        for k in attenuation.keys():
            if k in self.attenuation.keys():
                self.attenuation[k] = attenuation[k]
    
    def toJson(self):
        dict = super(Spotlight,self).toJson()
        dict['attenuation'] = self.attenuation
        dict['exponent'] = self.exponent
        dict['castShadows'] = self.castShadows
        
        return dict
     
