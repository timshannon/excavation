from direct.showbase.ShowBase import ShowBase
 
class PandaTest(ShowBase):
 
    def __init__(self):
        ShowBase.__init__(self)
        
        #disable mouse
        base.disableMouse()
        self.load_level()
        self.add_tasks()
        

    def load_level(self):
        level = self.loader.loadModel("/home/tshannon/workspace/excavation/data/models/levels/leveltest.egg")
        level.reparentTo(self.render)
        
    def add_tasks(self):
        pass
 
 
app = PandaTest()
app.run()

