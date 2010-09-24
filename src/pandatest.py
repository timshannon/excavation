from direct.showbase.ShowBase import ShowBase
import os, sys
 
class PandaTest(ShowBase):
 
    def __init__(self):
        ShowBase.__init__(self)
        
        self.RUNNINGDIR = os.path.abspath(sys.path[0])
        self.MODELPATH = "../data/models/"
        
        
        #disable mouse
        base.disableMouse()
        self.load_level()
        self.add_tasks()
        

    def load_level(self):
        level = self.loader.loadModel(os.path.join(self.RUNNINGDIR, self.MODELPATH + "levels/leveltest.egg"))
        level.reparentTo(self.render)
        
    def add_tasks(self):
        taskMgr.add(self.player_move, 'player_move') 
        
        
        
    def player_move(self, task):
        """ handles player movement"""
        
        
        
        return task.cont
        
 
app = PandaTest()
app.run()

