import copy
from panda3d.core import WindowProperties

class FreeViewController():
    
    MOUSESENSITIVITY = 0.1
    INVERTMOUSE = True
    MAXSPEED = 0.05
    ACCELERATION = .25
    CAMERAFOV = 90
    
    def __init__(self,
                 base,
                 kMoveForward,
                 kMoveBackward,
                 kMoveLeft,
                 kMoveRight,
                 kMoveUp,
                 kMoveDown,
                 mouseSensitivity,
                 maxSpeed,
                 acceleration):
        
        if mouseSensitivity > 0:
            self.MOUSESENSITIVITY = mouseSensitivity
        if maxSpeed > 0:
            self.MAXSPEED = maxSpeed
        if acceleration > 0:
            self.ACCELERATION = acceleration
        
        self.base = base
        base.camLens.setFov(self.CAMERAFOV)
        
        self.direction = {"x":0, "y":0, "z":0}
        self.speed = {"x":0, "y":0, "z":0}
        self.lastTask = 0

        #Add keys
        base.accept(kMoveForward, self.move, ["y", 1])
        base.accept(kMoveForward + '-up', self.move, ["y", 0])
        base.accept(kMoveBackward, self.move, ["y", -1])
        base.accept(kMoveBackward + "-up", self.move, ["y", 0])
        base.accept(kMoveLeft, self.move, ["x", -1])
        base.accept(kMoveLeft + "-up", self.move, ["x", 0])
        base.accept(kMoveRight, self.move, ["x", 1])
        base.accept(kMoveRight + "-up", self.move, ["x", 0])
        base.accept(kMoveUp, self.move, ["z", 1])
        base.accept(kMoveUp + "-up", self.move, ["z", 0])
        base.accept(kMoveDown, self.move, ["z", -1])
        base.accept(kMoveDown + "-up", self.move, ["z", 0])
        
        self.base.taskMgr.add(self.updateCamera, 'updateCamera') 
        
    def updateCamera(self, task):
        """ handles player movement"""

        
        if self.base.mouseWatcherNode.hasMouse():
            x = self.base.mouseWatcherNode.getMouseX()
            y = self.base.mouseWatcherNode.getMouseY()
        else:
            pointer = self.base.win.getPointer(0)
            x = pointer.getX()
            y = pointer.getY()
            
        #Reset pointer position
        self.base.win.movePointer(0, 300, 300)
        #get amount cursor moved
        x = (x - 300) * -1 
        y = (y - 300)
        
        if not self.INVERTMOUSE:
            y = y * -1
        
        
        quat = self.base.camera.getQuat()
        upQ = copy.copy(quat)
        rightQ = copy.copy(quat)
        #forwardQ = copy.copy(quat)
        forward = self.base.camera.getQuat().getForward()
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
        self.base.camera.setQuat(newQuat)
        
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
            
        
            
        self.base.camera.setFluidPos(self.base.camera, 
                                       self.speed["x"], 
                                       self.speed["y"], 
                                       self.speed["z"])
        
        self.lastTask = task.time
        return task.cont
        
    def move(self, dir, value):
        """Moves player"""
        self.direction[dir] = value
        
    
    def activateController(self):
        #disable mouse and hide cursor
        self.base.disableMouse()
        props = WindowProperties()
        props.setCursorHidden(True)
        self.base.win.requestProperties(props)
    
class RotateViewController():
    
    MOUSESENSITIVITY = 0.1
    INVERTMOUSE = True
    CAMERAFOV = 90
    MAXSPEED = 0.05
    ACCELERATION = .25
    
    def __init__(self):
        pass