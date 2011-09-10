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
from tools.collision import Collision
from utility.globalDef import GlobalDef
from panda3d.core import loadPrcFileData, WindowProperties
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
    ID_LOADMODEL = wx.NewId()

            
    SETTINGSFILE = ".collide"
    
    filename = '' 
    
    
    def __init__(self, *args, **kwargs): 
        wx.Frame.__init__(self, *args, **kwargs) 
        self.Show()
        split = wx.SplitterWindow(self, -1)
        
        self.pandapanel = PandaPanel(split, wx.ID_ANY, size=self.ClientSize) 
        self.pandapanel.initialize()
        self.jEditor = JsonEditor(split)
        split.SetSashGravity(1.0)
        
        #settings
        self.settings = {}
        self.saveDir = GlobalDef.MODELPATH
        self.recentFiles = []   
        self.loadSettings()
        
        self.CreateStatusBar()
        self.createMenus()
        
        split.SplitVertically(self.pandapanel, self.jEditor, 800)
        
       
        
    def loadSettings(self):
        if not os.access(self.SETTINGSFILE, os.F_OK):
            self.saveSettings()
        
        file = open(self.SETTINGSFILE, "rb")
        self.settings = cPickle.load(file)
        file.close()
        
        self.recentFiles = self.settings["recentFiles"]
        self.saveDir = self.settings["saveDir"]
        
        
    def saveSettings(self):
        self.settings["recentFiles"] = self.recentFiles
        self.settings["saveDir"] = self.saveDir
        
        file = open(self.SETTINGSFILE, "wb")
        cPickle.dump(self.settings, file)
        file.close()
    
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
        
        mAdd = wx.Menu()
        addList = [(self.ID_ADDSPHERE, 'Sphere', 'Add a sphere', self.add), \
                    (self.ID_ADDPLANE, 'Plane', 'Add a plane', self.add)
                    ]
        buildMenu(mAdd, addList)
        menuBar.Append(mAdd, '&Add')
        
        self.SetMenuBar(menuBar)
        
    def new(self,event):
        self.actionManager.execute('new')
        self.filename = ''
    
    def add(self, event):
        """Adds an item, determined by the event's sender id"""
        
        if event.id == self.ID_ADDMODEL:
            self.actionManager.execute('addModel')
    
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
        
      
class JsonEditor(wx.TextCtrl):
    collision = ''
    def __init__(self, *args, **kwargs):
        kwargs['style'] = wx.TE_MULTILINE
        super(JsonEditor, self).__init__(*args, **kwargs)
        
    def UpdateJson(self):
        self.Clear()
        self.AppendText(self.collision.toJson())
        
                    
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
        self.jEditor = self.frame.jEditor
        self.modelNode = None
        
        self.registerActions()
        self.collision = Collision()
        self.jEditor.collision = self.collision
        
        self.wxStep()    
    
    
    def registerActions(self):
        '''Register all of the actions to the editor functions so
            they can be used with the action manager'''
        self.actionManager.registerAction('open', Action(self.open))
        self.actionManager.registerAction('new', Action(self.new))
        self.actionManager.registerAction('save', Action(self.save))
        self.actionManager.registerAction('loadModel', Action(self.loadModel))
        
    
    def open(self, parms):
        self.new(parms)
        
    def new(self, parms):
        self.actionManager.reset()
        self.collision = Collision()
        self.jEditor.UpdateJson()
        self.modelNode = None
        for node in render.getChildren():
            node.removeNode()
            
        
    def loadModel(self, parms):
        if self.modelNode:
            self.modelNode.removeNode()
        
        self.modelNode = self.loader.loadModel(parms['model'])
        self.modelNode.reparentTo(self.render)
        
        
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

app = Collide() 
run() 
