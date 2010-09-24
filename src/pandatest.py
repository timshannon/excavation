from direct.showbase.ShowBase import ShowBase
from panda3d.core import WindowProperties
import os, sys
 
class PandaTest(ShowBase):
 
    def __init__(self):
        ShowBase.__init__(self)
        
        self.RUNNINGDIR = os.path.abspath(sys.path[0])
        self.MODELPATH = "../data/models/"
        
        
        
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
        self.accept("w", self.move)
        
        
    def load_level(self):
        level = self.loader.loadModel(os.path.join(self.RUNNINGDIR, self.MODELPATH + "levels/leveltest.egg"))
        level.reparentTo(self.render)
               
                
    def update_player(self, task):
        """ handles player movement"""
        quat = base.camera.getQuat()
        
        #multiply directrion quats with current camera quat
                
        
        return task.cont
    
    def move(self):
        print "test"
        
 
app = PandaTest()
app.run()

