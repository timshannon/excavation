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
import os, sys, copy
 
class Excavation(ShowBase):
    
    RUNNINGDIR = os.path.abspath(sys.path[0])
    MODELPATH = "../data/models/"
    MOUSESENSITIVITY = 0.2
    INVERTMOUSE = True
    CAMERAFOV = 90
    MAXSPEED = 0.5
    ACCELERATION = 2.5

    def __init__(self):
        ShowBase.__init__(self)
        
        self.direction = {"x":0, "y":0, "z":0}
        self.speed = {"x":0, "y":0, "z":0}
        self.lastTask = 0
        
        #disable mouse and hide cursor
        base.disableMouse()
        props = WindowProperties()
        props.setCursorHidden(True)
        base.win.requestProperties(props)
        
        #normal FPS camera FOV
        base.camLens.setFov(self.CAMERAFOV)
        
        #setup player collider
        #cPlayer = CollisionSphere(0, 0, 0, 1)
        #cPlayerNode = base.camera.attachNewNode(CollisionNode("cpnode"))
        
        self.load_level()
        self.add_keys()
        
        taskMgr.add(self.update_player, 'update_player') 
        
      

    def add_keys(self):
        self.accept("escape", sys.exit, [0])
        self.accept("w", self.move, ["y", 1])
        self.accept("w-up", self.move, ["y", 0])
        self.accept("s", self.move, ["y", -1])
        self.accept("s-up", self.move, ["y", 0])
        self.accept("a", self.move, ["x", -1])
        self.accept("a-up", self.move, ["x", 0])
        self.accept("d", self.move, ["x", 1])
        self.accept("d-up", self.move, ["x", 0])
        self.accept("e", self.move, ["z", 1])
        self.accept("e-up", self.move, ["z", 0])
        self.accept("space", self.move, ["z", -1])
        self.accept("space-up", self.move, ["z", 0])
                    
        
        
    def load_level(self):
        level = self.loader.loadModel(os.path.join(self.RUNNINGDIR, self.MODELPATH + "levels/leveltest.egg"))
        level.reparentTo(self.render)
        
        
        
                
    def update_player(self, task):
        """ handles player movement"""
        pointer = base.win.getPointer(0)
        x = pointer.getX()
        y = pointer.getY()
        
        #Reset pointer position
        base.win.movePointer(0, 300, 300)
        #get amount cursor moved
        x = (x - 300) * -1 
        y = (y - 300)
        
        if not self.INVERTMOUSE:
            y = y * -1
        
        
        quat = base.camera.getQuat()
        upQ = copy.copy(quat)
        rightQ = copy.copy(quat)
        #forwardQ = copy.copy(quat)
        forward = base.camera.getQuat().getForward()
        forward.normalize()
        up = quat.getUp()
        right = quat.getRight()
        forward = quat.getForward()
        up.normalize()
        right.normalize()
        forward.normalize()
                   
        upQ.setFromAxisAngle(x * self.MOUSESENSITIVITY, up)
        rightQ.setFromAxisAngle(y * self.MOUSESENSITIVITY, right)
        #forwardQ.setFromAxisAngle(45, right)
                    
        newQuat = quat.multiply(upQ.multiply(rightQ))
        base.camera.setQuat(newQuat)
        
        elapsed = task.time - self.lastTask
        #Move player
        for k in self.direction.keys():
            if self.direction[k] <> 0:
                self.speed[k] = self.speed[k] + (self.ACCELERATION * self.direction[k] * elapsed)
            else:
                #decelerate
                if self.speed[k] > 0:
                    self.speed[k] = self.speed[k] - (self.ACCELERATION * elapsed)
                    if self.speed[k] < 0:
                        self.speed[k] = 0
                else:
                    self.speed[k] = self.speed[k] + (self.ACCELERATION * elapsed)
                    if self.speed[k] > 0:
                        self.speed[k] = 0
                
            if abs(self.speed[k]) > self.MAXSPEED:
                if self.speed[k] > 0:
                    self.speed[k] = self.MAXSPEED
                else:
                    self.speed[k] = self.MAXSPEED * -1
            
        
            
        base.camera.setFluidPos(base.camera, 
                               self.speed["x"], 
                               self.speed["y"], 
                               self.speed["z"])
        
        self.lastTask = task.time
        return task.cont
    
    def move(self, dir, value):
        """Moves player"""
        self.direction[dir] = value
        
       
          
 
app = Excavation()
app.run()

