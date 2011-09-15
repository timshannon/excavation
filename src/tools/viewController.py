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
                 mouseSensitivity,
                 maxSpeed,
                 acceleration,
                 **keys):
        
        if mouseSensitivity > 0:
            self.MOUSESENSITIVITY = mouseSensitivity
        if maxSpeed > 0:
            self.MAXSPEED = maxSpeed
        if acceleration > 0:
            self.ACCELERATION = acceleration
        
        self.vcActive = False
        self.base = base
        self.base.camLens.setFov(self.CAMERAFOV)
        
        pointer = self.base.win.getPointer(0)
        self.mouseX = pointer.getX()
        self.mouseY = pointer.getY()

        
        #don't use built in mouse controller
        self.base.disableMouse()
        
        self.direction = {"x":0, "y":0, "z":0}
        self.speed = {"x":0, "y":0, "z":0}
        self.lastTask = 0

        #Add keys
        self.base.accept(keys['forward'], self.move, ["y", 1])
        self.base.accept(keys['forward'] + '-up', self.move, ["y", 0])
        self.base.accept(keys['backward'], self.move, ["y", -1])
        self.base.accept(keys['backward'] + "-up", self.move, ["y", 0])
        self.base.accept(keys['left'], self.move, ["x", -1])
        self.base.accept(keys['left'] + "-up", self.move, ["x", 0])
        self.base.accept(keys['right'], self.move, ["x", 1])
        self.base.accept(keys['right'] + "-up", self.move, ["x", 0])
        self.base.accept(keys['up'], self.move, ["z", 1])
        self.base.accept(keys['up'] + "-up", self.move, ["z", 0])
        self.base.accept(keys['down'], self.move, ["z", -1])
        self.base.accept(keys['down'] + "-up", self.move, ["z", 0])
        
        self.base.accept(keys['activate'], self.setControllerActiveState, [True])
        self.base.accept(keys['activate'] + '-up', self.setControllerActiveState, [False])
        
        self.base.taskMgr.add(self.updateCamera, 'updateCamera') 
        
    def updateCamera(self, task):
        """ handles player movement"""
        if not self.vcActive:
            return task.cont
        
        
        if not self.base.mouseWatcherNode.hasMouse():
            return task.cont
        
        x = self.base.mouseWatcherNode.getMouseX()
        y = self.base.mouseWatcherNode.getMouseY()
        
        #Reset pointer position
        self.base.win.movePointer(0, self.mouseX, self.mouseY)
        #get amount cursor moved
        x = (x - self.mouseX) * -1 
        y = (y - self.mouseY)
        
        if not self.INVERTMOUSE:
            y = y * -1
        
        
        camQuat = self.base.camera.getQuat()    #camQuat  ha ha, get it?
        upQ = copy.copy(camQuat)
        rightQ = copy.copy(camQuat)
        #forwardQ = copy.copy(camQuat)
        forward = self.base.camera.getQuat().getForward()
        forward.normalize()
        up = camQuat.getUp()
        right = camQuat.getRight()
        forward = camQuat.getForward()
        up.normalize()
        right.normalize()
        forward.normalize()
                   
        upQ.setFromAxisAngle(x * self.MOUSESENSITIVITY, up)
        rightQ.setFromAxisAngle(y * self.MOUSESENSITIVITY, right)
        #forwardQ.setFromAxisAngle(45, right)
                    
        newQuat = camQuat.multiply(upQ.multiply(rightQ))
        #self.base.camera.setQuat(newQuat)
        
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
        #print ('camera: ' + str(self.base.camera.getX()))
        print ('direction', self.direction['x'])
        
        self.lastTask = task.time
        return task.cont
        
    def move(self, dir, value):
        """Moves player"""
        self.direction[dir] = value
        
    
    def setControllerActiveState(self, active):
        #disable mouse and hide cursor
        self.vcActive = active
        
        pointer = self.base.win.getPointer(0)
        self.mouseX = pointer.getX()
        self.mouseY = pointer.getY()
        
        props = WindowProperties()
        props.setCursorHidden(active)
        self.base.win.requestProperties(props)
        
    
class RotateViewController():
    
    MOUSESENSITIVITY = 0.1
    INVERTMOUSE = True
    CAMERAFOV = 90
    MAXSPEED = 0.05
    ACCELERATION = .25
    
    def __init__(self):
        pass