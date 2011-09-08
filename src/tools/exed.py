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

from tools.actionManager import Action, ActionManager
from tools.scene import Scene
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
    ID_ADDMODEL = wx.NewId()
    ID_ADDENTITY = wx.NewId()
    ID_ADDPOINTLIGHT = wx.NewId()
    ID_ADDSPOTLIGHT = wx.NewId()
    ID_ADDDIRECTIONALLIGHT = wx.NewId()
    ID_RUNEXCAVATION = wx.NewId()
    ID_RUNDISCOURSE = wx.NewId()
    ID_SCENEPROP = wx.NewId()
    ID_RUNCOLLIDE = wx.NewId()
        
    SETTINGSFILE = "settings.exed"
    
    filename = '' 
    
    
    def __init__(self, *args, **kwargs): 
        wx.Frame.__init__(self, *args, **kwargs) 
        self.Show()
        self.pandapanel = PandaPanel(self, wx.ID_ANY, size=self.ClientSize) 
        self.pandapanel.initialize()
        
        #settings
        self.settings = {}
        self.saveDir = GlobalDef.SCENEPATH
        self.recentFiles = []   
        self.loadSettings()
        
        self.CreateStatusBar()
        self.createMenus()
        vSizer = wx.BoxSizer(wx.VERTICAL)
        hSizer = wx.BoxSizer(wx.HORIZONTAL)
                
        self.sceneTree = SceneTree(self)
        self.propList = PropList(self, style=wx.SIMPLE_BORDER)
        
        vSizer.Add(self.sceneTree, 1, wx.EXPAND)
        vSizer.Add(self.propList, 1, wx.EXPAND)  
        hSizer.Add(self.pandapanel, 4, wx.EXPAND)
        hSizer.Add(vSizer, 1, wx.EXPAND)
              
        self.SetSizer(hSizer)
        self.SetAutoLayout(1)
        hSizer.Fit(self)
        
        
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
        fileList = [(wx.ID_NEW, '&New', 'Create a new scene', self.newScene), \
                    (wx.ID_OPEN, '&Open', 'Open an existing scene', self.openScene), \
                    (wx.ID_SAVE, '&Save', 'Save the current scene', self.saveScene), \
                    (wx.ID_SAVEAS, 'Save As', 'Save the current scene as a new file', self.saveSceneAs), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (self.ID_RECENTFILES, None, None, None), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (wx.ID_EXIT, 'Exit', 'Exit ExEd', self.exit)]
        buildMenu(mFile, fileList)
        menuBar.Append(mFile, '&File')
        
        mEdit = wx.Menu()
        editList = [(wx.ID_UNDO, 'Undo', 'Undo the previous action', self.undo), \
                    (wx.ID_REDO, 'Redo', 'Redo the previous undone action', self.redo), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (wx.ID_CUT, 'Cut', 'Cut the selected item', self.cut), \
                    (wx.ID_COPY, '&Copy', 'Copy the selected item', self.copy), \
                    (wx.ID_PASTE, 'Paste', 'Paste the contents of the clipboard', self.paste), \
                    (wx.ID_DELETE, 'Delete', 'Delete the selected item', self.delete), \
                    (wx.ID_SEPARATOR, None, None, None), \
                    (self.ID_SCENEPROP, 'Scene Properties', 'Edit the scene''s properties', self.sceneProperties)]
        buildMenu(mEdit, editList)
        menuBar.Append(mEdit, '&Edit')
        
        mAdd = wx.Menu()
        addList = [(self.ID_ADDMODEL, 'Add Model', 'Add model to scene', self.addItem), \
                   (self.ID_ADDENTITY, 'Add Entity', 'Add entity to scene', self.addItem), \
                   (wx.ID_SEPARATOR, None, None, None), \
                   (self.ID_ADDPOINTLIGHT, 'Add Point Light', 'Add point light to scene', self.addItem), \
                   (self.ID_ADDSPOTLIGHT, 'Add Spot Light', 'Add spot light to scene', self.addItem), \
                   (self.ID_ADDDIRECTIONALLIGHT, 'Add Directional Light', 'Add directional light to scene', self.addItem)]
        buildMenu(mAdd, addList)
        menuBar.Append(mAdd, '&Add')
        
        mRun = wx.Menu()
        runList = [(self.ID_RUNEXCAVATION, 'Run Excavation', 'Run Excavation', self.runExternal), \
                   (self.ID_RUNCOLLIDE, 'Run Collide', 'Run Collide', self.runExternal), \
                   (self.ID_RUNDISCOURSE, 'Run Discourse', 'Run Discourse', self.runExternal)]
        buildMenu(mRun, runList)
        menuBar.Append(mRun, '&Run')
        
        self.SetMenuBar(menuBar)
        
    def newScene(self,event):
        self.actionManager.execute('newScene')
        self.filename = ''

    def openScene(self, event):
        dlg = wx.FileDialog(self, 
                            message='Open a Scene',
                            defaultDir=self.saveDir,
                            defaultFile='Scene File (*.scene)|*.scene',
                            style=wx.FD_OPEN)
        if dlg.ShowModal() == wx.ID_OK:
            self.filename = os.path.join(dlg.GetDirectory(), 
                                           dlg.GetFilename())
            self.actionManager.execute('openScene', 
                                         filename=self.filename)
    
    def saveScene(self, event):
        if self.filename <> '':
            self.actionManager.execute('saveScene', filename=self.filename)
        else:
            self.saveSceneAs(None)
            
    
    def saveSceneAs(self, event):
        dlg = wx.FileDialog(self, 
                            message='Save a Scene',
                            defaultDir=self.saveDir,
                            defaultFile='.scene',
                            style=wx.FD_SAVE)
        if dlg.ShowModal() == wx.ID_OK:
            self.filename = os.path.join(dlg.GetDirectory(),dlg.GetFilename())
            self.actionManager.execute('saveScene', 
                                         filename=self.filename)
    
    def exit(self, event):
        self.Close()
        
        
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
    def sceneProperties(self, event):
        pass
    
    def addItem(self, event):
        """Adds an item, determined by the event's sender id"""
        
        if event.id == self.ID_ADDMODEL:
            self.actionManager.execute('addModel')
    
    def runExternal(self, event):
        """Runs an external application"""
        pass
    
class PropList(wx.ListCtrl):
    def __init__(self, *args, **kwargs):
        super(PropList, self).__init__(*args, **kwargs)
        self.InsertColumn(0, "Property", wx.LIST_FORMAT_LEFT, 40)
        self.InsertColumn(1, "Value", wx.LIST_FORMAT_RIGHT, 50)
        
    def clear(self):
        pass
        

class SceneTree(wx.TreeCtrl):

    def __init__(self, *args, **kwargs):
        super(SceneTree, self).__init__(*args, **kwargs)
        #self.__collapsing = False
        self.root = self.AddRoot('render')
        
    def loadScene(self, scene):
        self.clear()
        def addChildren(parent, children):
            for c in children:
                newParent = self.AppendItem(parent, c.name)
                addChildren(newParent, c.children)
            
        addChildren(self.root, scene.tree.children)
        
    def clear(self):
        self.DeleteAllItems()
        self.root = self.AddRoot('render')
        
                
class ExEd(wx.App, ShowBase):
    """Panda object for handling all panda related tasks and events"""
    
         
    def __init__(self): 
        wx.App.__init__(self)
        ShowBase.__init__(self) 
        self.replaceEventLoop()
        self.frame = PandaFrame(None, wx.ID_ANY, 'ExEd', size=(1024,768)) 
        self.frame.Bind(wx.EVT_CLOSE, self.quit) 
        
        self.actionManager = ActionManager()
        self.frame.actionManager = self.actionManager
        self.sceneTree = self.frame.sceneTree
        self.propList = self.frame.propList
        
        self.registerActions()
        self.scene = Scene()
        
        self.wxStep()    
    
    
    def registerActions(self):
        '''Register all of the actions to the editor functions so
            they can be used with the action manager'''
        self.actionManager.registerAction('openScene', Action(self.openScene))
        self.actionManager.registerAction('newScene', Action(self.newScene))
        self.actionManager.registerAction('saveScene', Action(self.saveScene))
        
    
    def openScene(self, parms):
        self.scene.load(parms['filename'])
        self.scene.loadScene(self.render)
        self.sceneTree.loadScene(self.scene)
        self.actionManager.reset()
        
    def newScene(self, parms):
        self.scene = Scene()
        self.sceneTree.loadScene(self.scene)
        self.actionManager.reset()
        
    def saveScene(self, parms):
        self.scene.write(parms['filename'])
        
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
        if task != None: return task.cont 

app = ExEd() 
run() 
