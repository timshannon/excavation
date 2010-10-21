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

import cPickle
import StringIO
from ex_tag import TagGroup

__all__ = ["Scene", "Node", "Light", "Spotlight", "Entity", "PointLight", "DirectionalLight"]

class Scene(): 
    #Defines the structure of a scene for saving to a JSON file for loading
    
    def __init__(self, fileName=""):
        if fileName == "":
            self.tree = Node("render")
            self.keyValues = {}  #dictionary for scene level values, ambient light level, skybox, name, level description, level load hints, etc
        else:
            self.read(fileName)
            
        
    def write(self, fileName):
        file = open(fileName, "wb")
        cPickle.dump(self, file)
        file.close()
        
    def read(self, fileName):
        file = open(fileName, "rb")
        scene = cPickle.load(file)
        file.close()
        
        self.keyValues = scene.keyValues
        self.tree = scene.tree
    
        
class Node():
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
        self.tags = TagGroup()
        self.x = x
        self.y = y
        self.z = z
        self.h = h
        self.p = p
        self.r = r
        self.parent = parent
        self.children = []
        
        if self.parent <> None:
            self.parent.add_child(self)
            
            
    def add_child(self, child):
        if child not in self.children:
            self.children.append(child)
            if child.parent == None:
                child.parent = self
                
    def get_siblings(self):
        if self.parent == None:
            return None
        else:
            siblings = self.parent.children[:]
            siblings.remove(self)
    
    def set_pos(self, x, y, z):
        """Sets position of the node"""
        self.x = x
        self.y = y
        self.z = z
            
    def set_hpr(self, h, p, r):
        """Sets heading, pitch and rotation"""
        self.h = h
        self.p = p
        self.r = r
                    
class Entity(Node):
    keyValues = {}  #keyvalue dictionary to hold any settings the entity may make use of
    type = ""       #type of entity used to identify to python file to use
    
class Model(Node):
    model = ""
    collision = ""
    scaleX = 1.0
    scaleY = 1.0
    scaleZ = 1.0
    
    def set_scale(self, **scale):
        if "x" in scale:
            self.scaleX = scale["x"]
        
        if "y" in scale:
            self.scaley = scale["y"]
            
        if "z" in scale:
            self.scaleZ = scale["z"]
        
        
            
            
    
class Light(Node):
    color = {"red":1,"green":1,"blue":1,"alpha":1}
    specColor = {"red":1,"green":1,"blue":1,"alpha":1}
    
    def set_color(self, **colors):
        for k in colors.keys():
            if k in self.color.keys():
                self.color[k] = colors[k]
                
    def set_spec_color(self, **colors):
        for k in colors.keys():
            if k in self.specColor.keys():
                self.specColor[k] = color[k]
    
class PointLight(Light):
    attenuation = {"constant":0,"linear":0,"quadratic":0}
    
    def set_attenuation(self, **attenuation):
        for k in attenuation.keys():
            if k in self.attenuation.keys():
                self.attenuation[k] = attenuation[k]
    
class DirectionalLight(Light):
    direction = {"x":0,"y":0,"z":0}
    
    def set_direction(self, **direction):
        for k in direction.keys():
            if k in self.direction.keys():
                self.direction[k] = direction[k]
                
class Spotlight(Light):
    attenuation = {"constant":0,"linear":0,"quadratic":0}
    exponent = 0.0
    
    def set_attenuation(self, **attenuation):
        for k in attenuation.keys():
            if k in self.attenuation.keys():
                self.attenuation[k] = attenuation[k]
        
    