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

__all__ = ['Collision', 'Sphere', 'Plane', 'Box', 'Cylinder', 'Capsule', 'Cone', 'ConvexHull', 'TriangleMesh']

'''
    Collisions will be defined in a separate file from the model.
    When a model file is specified for loading, the scene loader will
    automatically check for a .collision file of the same name as the model
    and in the same directory.  If one exists, it will load the collision
    shapes defined in the collision file first, then attach the model as a 
    child to that node
'''

class Collision():
    '''
        Main collision class for reading writing and loading collision files
        Also defines the relative position to it's parent model, and it's mass.
        if there are no shapes defined, then it will default to Triangle Mesh shape with
        every geom in the egg file being automatically added.
    '''
    
    shapes = []
    relX = 0
    relY = 0
    relZ = 0
    relH = 0
    relP = 0
    relR = 0
    mass = 0.0
    
    def __init__(self,
                 file=''):
        if file:
            self.load(file)
             
    def load(self, file):
        pass
    
    def write(self, file):
        file = open(file, 'wb')
        jObject = []
        
        for s in self.shapes:
            jObject.append(s.toJson())
        
        json.dump(jObject, file, sort_keys=True, indent=4)
                
    

class Sphere():
    radius = 0
    
    def __init__(self,
                 radius):
        self.radius = radius
    
    def setRadius(self, radius):
        self.radius = radius
        
    def toJson(self):
        return {'type':'Sphere', 
                'radius':0}
    
        
class Plane():
    normal = Vec3(0,0,1)
    distance = 0
    
    def __init__(self,
                 normal,
                 distance):
    
        self.normal = normal
        self.distance = distance
                 
                 
    def toJson(self):
        return {'type':'Plane',
                'normal':[self.normal.getX(), self.normal.getY(), self.normal.getZ()], 
                'distance':self.distance}
    
class Box():
    x = 0
    y = 0
    z = 0
        
    def __init__(self,
                 x,
                 y,
                 z):
        self.x = x
        self.y = y
        self.z = z
    def toJson(self):
        return {'type':'Box',
                'x':self.x,
                'y':self.y,
                'z':self.z}
        
class Cylinder():
    radius = 0
    height = 0
    axis = 0    #bullet enum
    
    def __init__(self,
                 radius,
                 height,
                 axis):
        self.radius = radius
        self.height = height
        self.axis = axis
        
    def toJson(self):
        return {'type':'Cylinder',
                'radius':self.radius,
                'height':self.height,
                'axis':self.axis}
    
    
class Capsule():
    radius = 0
    height = 0
    
    def __init__(self,
                 radius,
                 height):
        self.height = height
        self.radius = radius
        
    def toJson(self):
        return {'type':'Capsule',
                'radius':self.radius,
                'height':self.height}
    
class Cone():
    radius = 0
    height = 0
    
    def __init__(self,
                 radius,
                 height):
        self.radius = radius
        self.height = height
        
    def toJson(self):
        return {'type':'Cone',
                'radius':self.radius,
                'height': self.height}
        
    
class ConvexHull():
    points = []
    
    def __init__(self,
                 points):
        self.points = points
        
    def addPoint(self, point3):
        self.points.append(point3)
        
    def toJson(self):
        points = []
        
        for p in self.points:
            points.append([p.getX(),p.getY(),p.getZ()])
        
        return {'type':'ConvexHull',
                'points':points}
    
