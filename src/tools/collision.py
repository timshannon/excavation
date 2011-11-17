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
from panda3d.bullet import BulletBoxShape
from panda3d.bullet import BulletCylinderShape
from panda3d.bullet import BulletCapsuleShape
from panda3d.bullet import BulletConeShape
from panda3d.bullet import BulletConvexHullShape
from panda3d.bullet import BulletTriangleMeshShape
from panda3d.bullet import BulletRigidBodyNode

from panda3d.core import Vec3
from panda3d.core import Point3
from panda3d.core import TransformState

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
           
    def __init__(self,
                 file=''):
        self.shapes = []
        self.mass = 0.0
        self.convexOnly = False
        self.name = 'Unnamed Collision'
        self.model = ''
    
        if file:
            self.load(file)
            self.name = file
        
             
    def load(self, file):
        fileObj = open(file, 'rb')
        jObject = json.load(fileObj)
        
        self.name = file
        self.mass = jObject['mass']
        self.convexOnly = jObject['convexOnly']
        self.model = jObject['model']
        shapes = jObject['shapes']
        self.shapes = []
        
        for s in shapes:
            for k,v in s.items():
                if k == 'Sphere':
                    self.shapes.append(Sphere(v['radius']))
                elif k == 'Plane':
                    self.shapes.append(Plane(v['normal'], 
                                             v['distance']))
                elif k == 'Box':
                    self.shapes.append(Box(v['x'], 
                                           v['y'],
                                           v['z']))
                elif k == 'Cylinder':
                    self.shapes.append(Cylinder(v['radius'], 
                                                v['height'], 
                                                v['enmAxis']))
                elif k == 'Capsule':
                    self.shapes.append(Capsule(v['radius'], 
                                               v['height'],
                                               v['enmAxis']))
                elif k == 'Cone':
                    self.shapes.append(Cone(v['radius'],
                                            v['height'],
                                            v['enmAxis']))
                
                self.shapes[-1].relX = v['relX']
                self.shapes[-1].relY = v['relY']
                self.shapes[-1].relZ = v['relZ']
                self.shapes[-1].relH = v['relH']
                self.shapes[-1].relP = v['relP']
                self.shapes[-1].relR = v['relR']
        
    def createNode(self, render, bulletWorld):
        '''Creates the rigid body node, adds it to the world and the renderer
            adds all the shapes, and returns the node'''
        node = BulletRigidBodyNode(self.name)
        nodePath = render.attachNewNode(node)
        bulletWorld.attachRigidBody(node)
        
        for s in self.shapes:
            node.addShape(s.createShape(),
                          s.transformState())
        
        return nodePath
            
    
    
    def write(self, file):
        self.name = file
        fileObj = open(file, 'wb')
        json.dump(self.toJson(),fileObj, indent=4, sort_keys=True)
        
    def toJson(self):
        
        jObject = []
        
        for s in self.shapes:
            jObject.append({s.__class__.__name__:s.toJson()})
        
        return {'mass':self.mass,
                'convexOnly':self.convexOnly,
                'model':self.model,
                'shapes':jObject}

class Shape(object):
    
    def __init__(self, 
                 relX = 0,
                 relY = 0,
                 relZ = 0,
                 relH = 0,
                 relP = 0,
                 relR = 0,
                 bulletShape = None):
        self.relX = relX
        self.relY = relY
        self.relZ = relZ
        self.relH = relH
        self.relP = relP
        self.relR = relR
        self.bulletShape = bulletShape
    
    def transformState(self):
        '''Returns the relative transform state necessary to
            to position the collision shape with its model'''
        return TransformState.makePosHpr(Point3(self.relX,
                                                self.relY,
                                                self.relZ),
                                         Point3(self.relH,
                                                self.relP,
                                                self.relR))
    
    def toJson(self):
        return {'relX':self.relX,
                   'relY':self.relY,
                   'relZ':self.relZ,
                   'relH':self.relH,
                   'relP':self.relP,
                   'relR':self.relR}
    

class Sphere(Shape):
    radius = 0
    
    def __init__(self,
                 radius):
        super(Sphere, self).__init__()
        self.radius = radius
    
    def setRadius(self, radius):
        self.radius = radius
    
    def createShape(self):
        return BulletSphereShape(self.radius)
        
    def toJson(self):
        dict = super(Sphere, self).toJson()
        dict['radius'] = self.radius
        
        return dict
    
        
class Plane(Shape):
    normal = Vec3(0,0,1)
    distance = 0
    
    def __init__(self,
                 normal,
                 distance):
        super(Plane, self).__init__()
        self.normal = Vec3(normal[0],normal[1],normal[2])
        self.distance = distance
                 
    def createShape(self):
        self.bulletShape = BulletPlaneShape(self.normal, self.distance)
        
        return self.bulletShape
    
    def toJson(self):
        dict = super(Plane, self).toJson()
        dict['normal'] = [self.normal.getX(), self.normal.getY(), self.normal.getZ()]
        dict['distance'] = self.distance
        return dict
    
class Box(Shape):
    x = 0
    y = 0
    z = 0
        
    def __init__(self,
                 x,
                 y,
                 z):
        super(Box, self).__init__()
        self.x = x
        self.y = y
        self.z = z
        
    def createShape(self):
        self.bulletShape =  BulletBoxShape(Vec3(self.x, 
                                           self.y, 
                                           self.z))
        return self.bulletShape
    
    def toJson(self):
        dict = super(Box, self).toJson()
        dict['x'] = self.x
        dict['y'] = self.y
        dict['z'] = self.z
        
        return dict 
        
class Cylinder(Shape):
    radius = 0
    height = 0
    enmAxis = 0    #bullet enum
    
    def __init__(self,
                 radius,
                 height,
                 enmAxis):
        super(Cylinder, self).__init__()
        self.radius = radius
        self.height = height
        self.enmAxis = enmAxis
        
    def createShape(self):
        self.bulletShape = BulletCylinderShape(self.radius, self.height, self.enmAxis)
        return self.bulletShape
        
    def toJson(self):
        dict = super(Cylinder, self).toJson()
        dict['radius'] = self.radius
        dict['height'] = self.height
        dict['enmAxis'] = self.enmAxis
        
        return dict
    
    
class Capsule(Shape):
    radius = 0
    height = 0
    enmAxis = 0
    
    def __init__(self,
                 radius,
                 height,
                 enmAxis):
        super(Capsule, self).__init__()
        self.height = height
        self.radius = radius
        self.enmAxis = enmAxis
        
    def createShape(self):
        self.bulletShape = BulletCapsuleShape(self.radius, self.height, self.enmAxis)
        return self.bulletShape
        
    def toJson(self):
        dict = super(Capsule, self).toJson()
        dict['radius'] = self.radius
        dict['height'] = self.height
        dict['enmAxis'] = self.enmAxis
        
        return dict 
    
class Cone(Shape):
    radius = 0
    height = 0
    enmAxis = 0
    
    def __init__(self,
                 radius,
                 height,
                 enmAxis):
        super(Cone, self).__init__()
        self.radius = radius
        self.height = height
        self.enmAxis = enmAxis
        
    def createShape(self):
        self.bulletShape = BulletConeShape(self.radius, self.height, self.enmAxis)
        return self.bulletShape
        
    def toJson(self):
        dict = super(Cone, self).toJson()
        dict['radius'] = self.radius
        dict['height'] = self.height
        dict['enmAxis'] = self.enmAxis
        
        return dict 
        

