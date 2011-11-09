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
                 **args):
        
        if args['mouseSensitivity'] > 0:
            self.MOUSESENSITIVITY = args['mouseSensitivity']
        if args['maxSpeed'] > 0:
            self.MAXSPEED = args['maxSpeed']
        if args['acceleration'] > 0:
            self.ACCELERATION = args['acceleration']
        
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
        self.base.accept(args['forward'], self.move, ["y", 1])
        self.base.accept(args['forward'] + '-up', self.move, ["y", 0])
        self.base.accept(args['backward'], self.move, ["y", -1])
        self.base.accept(args['backward'] + "-up", self.move, ["y", 0])
        self.base.accept(args['left'], self.move, ["x", -1])
        self.base.accept(args['left'] + "-up", self.move, ["x", 0])
        self.base.accept(args['right'], self.move, ["x", 1])
        self.base.accept(args['right'] + "-up", self.move, ["x", 0])
        self.base.accept(args['up'], self.move, ["z", 1])
        self.base.accept(args['up'] + "-up", self.move, ["z", 0])
        self.base.accept(args['down'], self.move, ["z", -1])
        self.base.accept(args['down'] + "-up", self.move, ["z", 0])
        
        self.base.accept(args['activate'], self.setControllerActiveState, [True])
        self.base.accept(args['activate'] + '-up', self.setControllerActiveState, [False])
        
        self.base.taskMgr.add(self.updateCamera, 'updateCamera') 
        
    def updateCamera(self, task):
        """ handles camera movement"""
        if not self.vcActive:
            return task.cont
        
        
        pointer = self.base.win.getPointer(0)
        x = pointer.getX()
        y = pointer.getY()
        
        
        #Reset pointer position
        self.base.win.movePointer(0, self.mouseX, self.mouseY)
        #get amount cursor moved
        x = (x - self.mouseX) * -1 
        y = (y - self.mouseY)
        
        if not self.INVERTMOUSE:
            y = y * -1
        
        x = x * self.MOUSESENSITIVITY
        y = y * self.MOUSESENSITIVITY
        
        
        self.base.camera.setHpr(self.base.camera, x, y, 0)
        
        
        elapsed = task.time - self.lastTask
        #Move player
        for k in self.direction.keys():
            if self.direction[k] != 0:
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
        #print ('direction', self.direction['x'])
        
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
        
        wp = WindowProperties(self.base.win.getProperties())
        wp.setForeground(active) 
        
        wp.setCursorHidden(active)
        self.base.win.requestProperties(wp)
        
    
class RotateViewController():
    
    MOUSESENSITIVITY = 0.1
    INVERTMOUSE = True
    CAMERAFOV = 90
    MAXSPEED = 0.05
    ACCELERATION = .25
    
    def __init__(self):
        pass