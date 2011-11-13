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

import sys
import os 
import wx
import cPickle
import time

from tools.actionManager import Action, ActionManager
from tools import collision
from tools.collision import Collision
from tools.viewController import FreeViewController, RotateViewController 
from utility.globalDef import GlobalDef
from panda3d.core import WindowProperties, ModifierButtons
from panda3d.core import loadPrcFile, loadPrcFileData
from panda3d.core import ConfigVariableString, ConfigVariableInt, ConfigVariableDouble 
from panda3d.core import Vec3
from panda3d.core import Point3
from panda3d.core import TransformState

from panda3d.bullet import BulletWorld
from panda3d.bullet import BulletRigidBodyNode
from panda3d.bullet import BulletDebugNode
from direct.showbase.ShowBase import ShowBase


loadPrcFileData('startup', 'window-type none')

class PandaPanel(wx.Panel): 
    def __init__(self, *args, **kwargs): 
        wx.Panel.__init__(self, *args, **kwargs) 
    
    def initialize(self): 
        assert self.GetHandle() != 0 
        wp = WindowProperties() 
        wp.setOrigin(0,0) 
        wp.setSize(self.ClientSize.GetWidth(), self.ClientSize.GetHeight()) 
        wp.setParentWindow(self.GetHandle()) 
        base.openDefaultWindow(props = wp, gsg = None) 
        self.Bind(wx.EVT_SIZE, self.onResize) 
    
    def onResize(self, event): 
        frame_size = event.GetSize() 
        wp = WindowProperties() 
        wp.setOrigin(0,0) 
        wp.setSize(frame_size.GetWidth(), frame_size.GetHeight()) 
        base.win.requestProperties(wp) 

class PandaFrame(wx.Frame):
    """wx object for handling wx events and actions"""
    actionManager = None
    
    #Menu ID's
    ID_RECENTFILES = wx.NewId()
    ID_ADDSPHERE = wx.NewId()
    ID_ADDPLANE = wx.NewId()
    ID_ADDBOX = wx.NewId()
    ID_ADDCYLINDER = wx.NewId()
    ID_ADDCAPSULE = wx.NewId()
    ID_ADDCONE = wx.NewId()
    ID_LOADMODEL = wx.NewId()

            
    SETTINGSFILE = ".collide"
    
    filename = '' 
    
    
    def __init__(self, *args, **kwargs): 
        wx.Frame.__init__(self, *args, **kwargs) 
        self.Show()
        split = wx.SplitterWindow(self, -1)
        
        self.pandapanel = PandaPanel(split, wx.ID_ANY, size=self.ClientSize) 
        self.pandapanel.initialize()
        
        sidebar = wx.SplitterWindow(split, -1)
        
        self.shapeList = ShapeList(sidebar)
        self.shapeProp = ShapeProperties(sidebar)
        
        sidebar.SetSashGravity(0.0)
        split.SetSashGravity(1.0)
        
        #settings
        self.settings = {}
        self.saveDir = GlobalDef.MODELPATH
        self.recentFiles = []   
        self.loadSettings()
        
        self.CreateStatusBar()
        self.createMenus()
        
        sidebar.SplitHorizontally(self.shapeList, self.shapeProp, 0)
        split.SplitVertically(self.pandapanel, sidebar, 900)
        
               
    def loadSettings(self):
        if not os.access(self.SETTINGSFILE, os.F_OK):
            self.saveSettings()
        
        sFile = open(self.SETTINGSFILE, "rb")
        self.settings = cPickle.load(sFile)
        sFile.close()
        
        self.recentFiles = self.settings["recentFiles"]
        self.saveDir = self.settings["saveDir"]
        
        
    def saveSettings(self):
        self.settings["recentFiles"] = self.recentFiles
        self.settings["saveDir"] = self.saveDir
        
        sFile = open(self.SETTINGSFILE, "wb")
        cPickle.dump(self.settings, sFile)
        sFile.close()
    
    def createMenus(self):
        def buildMenu(menu, valueList):
            for id, label, hintText, handle in valueList:
                if id == wx.ID_SEPARATOR:
                    menu.AppendSeparator()
                elif id == self.ID_RECENTFILES:
                    idFiles = []
                    for f in self.recentFiles:
                        id = wx.NewId()
                        idFiles.append((id, f))
                        mItem = menu.Append(id, f, 'Open recent file', self.openScene)
                    self.recentFiles = idFiles
                else:
                    mItem = menu.Append(id, label, hintText)
                    self.Bind(wx.EVT_MENU, handle, mItem)
        
        menuBar = wx.MenuBar()
        
        #file menu
        mFile = wx.Menu()
        fileList = [(wx.ID_NEW, '&New', 'Create a new collision', self.new), \
                    (wx.ID_OPEN, '&Open', 'Open an existing collision', self.open), \
                    (self.ID_LOADMODEL, '&Load Model', 'Load a model', self.loadModel), \
                    (wx.ID_SAVE, '&Save', 'Save the current collision', self.save), \
                    (wx.ID_SAVEAS, 'Save As', 'Save the current collision as a new file', self.saveAs), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (self.ID_RECENTFILES, None, None, None), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (wx.ID_EXIT, 'Exit', 'Exit ExEd', self.exit)]
        buildMenu(mFile, fileList)
        menuBar.Append(mFile, '&File')
        
        #edit menu
        mEdit = wx.Menu()
        editList = [(wx.ID_UNDO, 'Undo', 'Undo the previous action', self.undo), \
                    (wx.ID_REDO, 'Redo', 'Redo the previous undone action', self.redo), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (wx.ID_CUT, 'Cut', 'Cut the selected item', self.cut), \
                    (wx.ID_COPY, '&Copy', 'Copy the selected item', self.copy), \
                    (wx.ID_PASTE, 'Paste', 'Paste the contents of the clipboard', self.paste), \
                    (wx.ID_DELETE, 'Delete', 'Delete the selected item', self.delete), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (wx.ID_REFRESH, 'Reload File', 'Reload the current collision file', self.reloadFile)]
        buildMenu(mEdit, editList)
        menuBar.Append(mEdit, '&Edit')
        
        #add menu
        mAdd = wx.Menu()
        addList = [(self.ID_ADDSPHERE, 'Sphere', 'Add a sphere', self.add), \
                    (self.ID_ADDPLANE, 'Plane', 'Add a plane', self.add), \
                    (self.ID_ADDBOX, 'Box', 'Add a box', self.add), \
                    (self.ID_ADDCYLINDER, 'Cylinder', 'Add a cylinder', self.add), \
                    (self.ID_ADDCAPSULE, 'Capsule', 'Add a capsule', self.add), \
                    (self.ID_ADDCONE, 'Cone', 'Add a cone', self.add)
                    ]
        buildMenu(mAdd, addList)
        menuBar.Append(mAdd, '&Add')
        
        self.SetMenuBar(menuBar)
        
    def new(self,event):
        self.actionManager.execute('new')
        self.filename = ''
    
    def add(self, event):
        """Adds an item, determined by the event's sender id"""
        if event.Id == self.ID_ADDBOX:
            self.actionManager.execute('addBox',
                                       shape='box',
                                       x=1,y=1,z=1)
    
    def loadModel(self, event):
        '''Loads a model for viewing with the collision file, or for pulling
            vertexes from
        '''
        dlg = wx.FileDialog(self, 
                            message='Load a model',
                            defaultDir=self.saveDir,
                            defaultFile='"Egg files (*.egg)|*.egg"',
                            style=wx.FD_OPEN)
        if dlg.ShowModal() == wx.ID_OK:
            model = os.path.join(dlg.GetDirectory(), 
                                 dlg.GetFilename())
            self.actionManager.execute('loadModel', 
                                         model=model)
            
        
        
    def undo(self, event):
        pass
    
    def redo(self, event):
        pass
    
    def cut(self, event):
        pass
    
    def copy(self, event):
        pass
    
    def paste(self, event):
        pass
    
    def delete(self, event):
        pass
    
    def reloadFile(self, event):
        if self.filename:
            self.actionManager.execute('open', 
                                       filename=self.filename)
            
    
    def open(self, event):
        dlg = wx.FileDialog(self, 
                            message='Open a Collision',
                            defaultDir=self.saveDir,
                            defaultFile='Collision File (*.collision)|*.collision',
                            style=wx.FD_OPEN)
        if dlg.ShowModal() == wx.ID_OK:
            self.filename = os.path.join(dlg.GetDirectory(), 
                                           dlg.GetFilename())
            self.actionManager.execute('open', 
                                         filename=self.filename)
    
    def save(self, event):
        if self.filename <> '':
            self.actionManager.execute('save', filename=self.filename)
        else:
            self.saveSceneAs(None)
            
    
    def saveAs(self, event):
        dlg = wx.FileDialog(self, 
                            message='Save a Collision',
                            defaultDir=self.saveDir,
                            defaultFile='.collision',
                            style=wx.FD_SAVE)
        if dlg.ShowModal() == wx.ID_OK:
            self.filename = os.path.join(dlg.GetDirectory(),dlg.GetFilename())
            self.actionManager.execute('save', 
                                         filename=self.filename)
    
    def exit(self, event):
        self.Close()
        
      
class ShapeList(wx.ListCtrl):
    
    def __init__(self, *args, **kwargs):
        super(ShapeList, self).__init__(*args, **kwargs)
        
    def addShape(self, shape):
        pass
        
class ShapeProperties(wx.ListCtrl):
    
    def __init__(self, *args, **kwargs):
        super(ShapeProperties, self).__init__(*args, **kwargs)
        
                     
class Collide(wx.App, ShowBase):
    """Panda object for handling all panda related tasks and events"""
    
    def __init__(self): 
        wx.App.__init__(self)
        ShowBase.__init__(self) 
        self.replaceEventLoop()
        self.frame = PandaFrame(None, wx.ID_ANY, 'Collide', size=(1200,768)) 
        self.frame.Bind(wx.EVT_CLOSE, self.quit) 
        
        self.actionManager = ActionManager()
        self.frame.actionManager = self.actionManager
        self.shapeList = self.frame.shapeList
        
        
        #initialize bulletworld
        self.bWorld = BulletWorld()
        self.bWorld.setGravity(Vec3(0, 0, -9.81))
        self.bodyNode = BulletRigidBodyNode('baseNode')
        render.attachNewNode(self.bodyNode)
        self.bWorld.attachRigidBody(self.bodyNode)
        #debug node
        self.debugNode = BulletDebugNode('Debug')
        debugNP = render.attachNewNode(self.debugNode)
        debugNP.show()
        self.bWorld.setDebugNode(debugNP.node())
        
        taskMgr.add(self.update, 'update')
        
        #load collide config file
        loadPrcFile(GlobalDef.RUNNINGDIR + "/collide.prc")
        
        self.modelNode = None
        
        self.registerActions()
        self.collision = Collision()
                        
        #viewControllers
        self.fvc = FreeViewController(base, 
                                      mouseSensitivity=ConfigVariableDouble('mouseSensitivity', 0.1).getValue(), 
                                      maxSpeed=ConfigVariableDouble('maxSpeed', .05).getValue(),
                                      acceleration=ConfigVariableDouble('acceleration',1).getValue(),
                                      activate=ConfigVariableString('activate', 'mouse3').getValue(),
                                      forward=ConfigVariableString('forward', 'w').getValue(), 
                                      backward=ConfigVariableString('backward', 's').getValue(), 
                                      left=ConfigVariableString('left', 'a').getValue(), 
                                      right=ConfigVariableString('right', 'd').getValue(), 
                                      up=ConfigVariableString('up', 'e').getValue(), 
                                      down=ConfigVariableString('down', 'space').getValue())
        base.mouseWatcherNode.setModifierButtons(ModifierButtons()) 
        base.buttonThrowers[0].node().setModifierButtons(ModifierButtons())
        
        self.accept('mouse1', self.click)
        #messenger.toggleVerbose()
        
        self.wxStep()   
        
        
    def setBackground(self):
        wp = WindowProperties(base.win.getProperties())
        if wp.getForeground():
            wp.setForeground(False) 
            base.win.requestProperties(wp)
        
            
    def setForeground(self):
        wp = WindowProperties(base.win.getProperties())
        if not wp.getForeground():
            wp.setForeground(True) 
            base.win.requestProperties(wp) 
            
    def click(self):
        '''mouse 1 click'''
        print render.ls()
        
    
    def registerActions(self):
        '''Register all of the actions to the editor functions so
            they can be used with the action manager'''
        self.actionManager.registerAction('open', Action(self.open))
        self.actionManager.registerAction('new', Action(self.new))
        self.actionManager.registerAction('save', Action(self.save))
        self.actionManager.registerAction('loadModel', Action(self.loadModel))
        self.actionManager.registerAction('addBox', Action(self.addShape, self.removeShape))
        
    
    def open(self, parms):
        self.new(parms)
        
    def new(self, parms):
        self.actionManager.reset()
        self.collision = Collision()
        self.modelNode = None
        
        for shape in self.collision.shapes:
            self.bodyNode.removeShape(shape.bulletShape)
        
        
            
        
    def loadModel(self, parms):
        if self.modelNode:
            self.modelNode.removeNode()
        
        #TODO: Load ordered list of vertexes from the model so that
        #    sizing operations can quick snap to the vertex bounds
        #    Regrabbing the unordered list of vertexes for each collision
        #    shape change operation would take way to long, especially on
        #    complex models.  Grab them once on load instead and order them
        self.modelNode = self.loader.loadModel(parms['model'])
        self.modelNode.reparentTo(render)
        
    def addShape(self, parms):
        '''Adds a given shape based on the parms passed in'''
        if parms['shape'] == 'box':
            colShape = collision.Box(parms['x'], 
                                     parms['y'],
                                     parms['z'])
            self.collision.shapes.append(colShape)
            self.bodyNode.addShape(colShape.createShape(), 
                                          colShape.transformState())
            
            
                
        
        
    def removeShape(self, parms):
        '''removes the passed in shape'''
        pass
        
    def save(self, parms):
        pass
       
       
    def replaceEventLoop(self): 
        self.evtLoop = wx.EventLoop() 
        self.oldLoop = wx.EventLoop.GetActive() 
        wx.EventLoop.SetActive(self.evtLoop) 
        taskMgr.add(self.wxStep, "evtLoopTask") 
    
    def onDestroy(self, event=None): 
        self.wxStep() 
        wx.EventLoop.SetActive(self.oldLoop) 
    
    def quit(self, event=None): 
        self.onDestroy(event) 
        try: 
            base 
        except NameError: 
            sys.exit() 
        base.userExit() 
    
    def wxStep(self, task=None): 
        while self.evtLoop.Pending(): 
            self.evtLoop.Dispatch() 
        self.ProcessIdle() 
        time.sleep(0.01)
        if task != None: return task.cont
        
    def update(self, task):
        dt = globalClock.getDt()
        self.bWorld.doPhysics(dt)
        return task.cont 

app = Collide()
run() 
