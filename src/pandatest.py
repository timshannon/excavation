from direct.showbase.ShowBase import ShowBase
from panda3d.core import WindowProperties
import os, sys, copy
 
class PandaTest(ShowBase):
 
    def __init__(self):
        ShowBase.__init__(self)
        
        self.RUNNINGDIR = os.path.abspath(sys.path[0])
        self.MODELPATH = "../data/models/"
        self.MOUSESENSITIVITY = 0.1
        self.INVERTMOUSE = True
        
               
        #disable mouse and hide cursor
        base.disableMouse()
        props = WindowProperties()
        props.setCursorHidden(True)
        base.win.requestProperties(props)
        
        
        self.load_level()
        self.add_keys()
        
        taskMgr.add(self.update_player, 'update_player') 
        
      

    def add_keys(self):
        self.accept("escape", sys.exit, [0])
        self.accept("w", self.)        
        
        
    def load_level(self):
        level = self.loader.loadModel(os.path.join(self.RUNNINGDIR, self.MODELPATH + "levels/leveltest.egg"))
        level.reparentTo(self.render)
               
                
    def update_player(self, task):
        """ handles player movement"""
        pointer = base.win.getPointer(0)
        x = pointer.getX()
        y = pointer.getY()
        
        #Reset pointer position
        if base.win.movePointer(0, 300, 300):
            #get amount cursor moved
            x = (x - 300) * -1 
            y = (y - 300)
            
            if not self.INVERTMOUSE:
                y = y * -1
            
            
            quat = base.camera.getQuat()
            upQ = copy.copy(quat)
            rightQ = copy.copy(quat)
            #forwardQ = copy.copy(quat)
            up = quat.getUp()
            right = quat.getRight()
            forward = quat.getForward()
            up.normalize()
            right.normalize()
            forward.normalize()
                       
            upQ.setFromAxisAngle(x * self.MOUSESENSITIVITY, up)
            rightQ.setFromAxisAngle(y * self.MOUSESENSITIVITY, right)
            #forwardQ.setFromAxisAngle(45, right)
                        
            
            base.camera.setQuat(quat.multiply(upQ.multiply(rightQ))) 
            
        
        return task.cont
    
        
 
app = PandaTest()
app.run()

