#Stolen from treeform at the panda forums (http://www.panda3d.org/forums/viewtopic.php?t=4068)
# to use as an example reference

import direct.directbase.DirectStart 
from direct.showbase.DirectObject import DirectObject 
from pandac.PandaModules import * 

class FPS(object,DirectObject): 
    def __init__(self): 
        self.initCollision() 
        self.loadLevel() 
        self.initPlayer() 
        
    def initCollision(self): 
        #initialize traverser 
        base.cTrav = CollisionTraverser() 
        #initialize pusher 
        self.pusher = CollisionHandlerPusher() 
        
    def loadLevel(self): 
        
        #load level 
        # must have 
        #<Group> *something* { 
        #  <Collide> { Polyset keep descend } in the egg file 
        level = loader.loadModel('level.egg') 
        level.reparentTo(render) 
        level.setPos(0,0,0) 
        level.setTwoSided(True) 
        level.setColor(1,1,1,.5) 
                
    def initPlayer(self): 
        
        #load man 
        self.man = loader.loadModel('teapot') 
        self.man.reparentTo(render) 
        self.man.setPos(0,0,2) 
        self.man.setScale(.05) 
        base.camera.reparentTo(self.man) 
        base.camera.setPos(0,0,0) 
        base.disableMouse() 
        #create a collsion solid for the man 
        cNode = CollisionNode('man') 
        cNode.addSolid(CollisionSphere(0,0,0,3)) 
        manC = self.man.attachNewNode(cNode) 
        base.cTrav.addCollider(manC,self.pusher) 
        self.pusher.addCollider(manC,self.man, base.drive.node()) 
        
        speed = 50 
        Forward = Vec3(0,speed*2,0) 
        Back = Vec3(0,-speed,0) 
        Left = Vec3(-speed,0,0) 
        Right = Vec3(speed,0,0) 
        Stop = Vec3(0) 
        self.walk = Stop 
        self.strife = Stop 
        self.jump = 0 
        taskMgr.add(self.move, 'move-task') 
        self.accept( "space" , self.__setattr__,["jump",1.]) 
        self.accept( "s" , self.__setattr__,["walk",Back] ) 
        self.accept( "w" , self.__setattr__,["walk",Forward]) 
        self.accept( "s" , self.__setattr__,["walk",Back] ) 
        self.accept( "s-up" , self.__setattr__,["walk",Stop] ) 
        self.accept( "w-up" , self.__setattr__,["walk",Stop] ) 
        self.accept( "a" , self.__setattr__,["strife",Left]) 
        self.accept( "d" , self.__setattr__,["strife",Right] ) 
        self.accept( "a-up" , self.__setattr__,["strife",Stop] ) 
        self.accept( "d-up" , self.__setattr__,["strife",Stop] ) 
        
        self.manGroundRay = CollisionRay() 
        self.manGroundRay.setOrigin(0,0,-.2) 
        self.manGroundRay.setDirection(0,0,-1) 
        
        self.manGroundCol = CollisionNode('manRay') 
        self.manGroundCol.addSolid(self.manGroundRay) 
        self.manGroundCol.setFromCollideMask(BitMask32.bit(0)) 
        self.manGroundCol.setIntoCollideMask(BitMask32.allOff()) 
        
        self.manGroundColNp = self.man.attachNewNode(self.manGroundCol) 
        self.manGroundColNp.show() 
        self.manGroundHandler = CollisionHandlerQueue() 
        
        base.cTrav.addCollider(self.manGroundColNp, self.manGroundHandler) 
        
    def move(self,task): 
        
        # mouse 
        md = base.win.getPointer(0) 
        x = md.getX() 
        y = md.getY() 
        if base.win.movePointer(0, base.win.getXSize()/2, base.win.getYSize()/2): 
            self.man.setH(self.man.getH() -  (x - base.win.getXSize()/2)*0.1) 
            base.camera.setP(base.camera.getP() - (y - base.win.getYSize()/2)*0.1) 
        # move where the keys set it 
        self.man.setPos(self.man,self.walk*globalClock.getDt()) 
        self.man.setPos(self.man,self.strife*globalClock.getDt()) 
        
        highestZ = -100 
        for i in range(self.manGroundHandler.getNumEntries()): 
            entry = self.manGroundHandler.getEntry(i) 
            if entry.getIntoNode().getName() == "Cube": 
                z = entry.getSurfacePoint(render).getZ() 
                if z > highestZ: 
                    highestZ = z 
        # gravity 
        self.man.setZ(self.man.getZ()+self.jump*globalClock.getDt()) 
        self.jump -= 1*globalClock.getDt() 
        
        if highestZ > self.man.getZ()-.3: 
            self.jump = 0 
            self.man.setZ(highestZ+.3) 
        
        return task.cont 
FPS() 
render.setShaderAuto() 
run() 

