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

from panda3d.bullet import BulletWorld
from panda3d.bullet import BulletSphereShape
from panda3d.bullet import BulletPlaneShape
from panda3d.core import Vec3
from panda3d.core import Point3

import json
import os
from utility.globalDef import GlobalDef

__all__ = ['Collision', 'Shape', 'Sphere']

'''*******
    Collisions will be defined in a separate file from the model.
    When a model file is specified for loading, the scene loader will
    automatically check for a .collision file of the same name as the model
    and in the same directory.  If one exists, it will load the collision
    shapes defined in the collision file first, then attach the model as a 
    child to that node
'''

class Collision():
    '''Main collision class for reading writing and loading collision files'''
    shapes = []
    
    def __init__(self,
                 file=''):
        if file:
            self.load(file)
             
    def load(self, file):
        pass
    
    def write(self, file):
        file = open(file, 'wb')
        json.dump(self.shapes, file, sort_keys=True, indent=4)
        
    
class Shape():
    '''Collision shape base class
        Position is based relative to it's model's origin
        The collision shape will be loaded first, then the model
    '''
    
    def __init__(self,
                 relX=0,
                 relY=0,
                 relZ=0,
                 relH=0,
                 relP=0,
                 relR=0,
                 mass=0):
        self.relX = relX
        self.relY = relY
        self.relZ = relZ
        self.relH = relH
        self.relP = relP
        self.relR = relR


class Sphere(Shape):
    radius = 0
    
    def setRadius(self, radius):
        self.radius = radius
        
class Plain(Shape):
    normal = Vec3(0,0,1)
    distance = 0
    
    
class Box(Shape):
    x = 0
    y = 0
    z = 0
    
    def setShape(self,
                 x,
                 y,
                 z):
        self.x = x
        self.y = y
        self.z = z
        
class Cylinder(Shape):
    radius = 0
    height = 0
    axis = 0    #bullet enum
    
class Capsule(Shape):
    radius = 0
    height = 0
    
class Cone(Shape):
    radius = 0
    height = 0
    
class ConvexHull(Shape):
    points = []
    
    def addPoint(self, point3):
        self.points.append(point3)
    
class TriangleMesh(Shape):
    '''Should be used for static level structure only.
        Will auto add geoms for the specified egg
    '''
    geoms = []