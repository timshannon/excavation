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
import wx 

import direct 
from pandac.PandaModules import * 
from tools.actionManager import ActionManager
loadPrcFileData('startup', 'window-type none') 
from direct.directbase.DirectStart import * 
from direct.showbase import DirectObject
from tools import actionManager 

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
    
    ID_RECENTFILES = wx.NewId()
    ID_ADDMODEL = wx.NewId()
    ID_ADDENTITY = wx.NewId()
    ID_ADDPOINTLIGHT = wx.NewId()
    ID_ADDSPOTLIGHT = wx.NewId()
    ID_ADDDIRECTIONALLIGHT = wx.NewId()
    
    ID_RUNEXCAVATION = wx.NewId()
    ID_RUNDISCOURSE = wx.NewId()
    
    
    def __init__(self, *args, **kwargs): 
        wx.Frame.__init__(self, *args, **kwargs) 
        self.Show()
        self.pandapanel = PandaPanel(self, wx.ID_ANY, size=self.ClientSize) 
        self.pandapanel.initialize()
        
        self.recentFiles = []   #load from config file with other items?
                
        self.CreateStatusBar()
        self.createMenus()
        
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
        
        mFile = wx.Menu()
        fileList = [(wx.ID_NEW, '&New', 'Create a new scene', self.newScene), \
                    (wx.ID_OPEN, '&Open', 'Open an existing scene', self.openScene), \
                    (wx.ID_SAVE, '&Save', 'Save the current scene', self.saveScene), \
                    (wx.ID_SAVEAS, 'Save As', 'Save the current scene as a new file', self.saveSceneAs), \
                    (wx.ID_SEPARATOR, None, None), \
                    (self.ID_RECENTFILES, None, None), \
                    (wx.ID_SEPARATOR, None, None), \
                    (wx.ID_EXIT, 'Exit', 'Exit ExEd', self.exit)]
        buildMenu(mFile, fileList)
        
        mEdit = wx.Menu()
        editList = [(wx.ID_UNDO, 'Undo', self.undo), \
                    (wx.ID_SEPARATOR, None, None), \
                    (wx.ID_REDO, 'Redo', self.redo), \
                    (wx.ID_CUT, 'Cut', self.cut), \
                    (wx.ID_COPY, '&Copy', self.copy), \
                    (wx.ID_PASTE, 'Paste', self.paste), \
                    (wx.ID_SEPARATOR, None, None), \
                    (wx.ID_DELETE, 'Delete', self.delete)]
        buildMenu(mEdit, editList)
        
        mAdd = wx.Menu()
        addList = [(self.ID_ADDMODEL, 'Add Model', self.addItem), \
                   (self.ID_ADDENTITY, 'Add Entity', self.addItem), \
                   (wx.ID_SEPARATOR, None, None), \
                   (self.ID_ADDPOINTLIGHT, 'Add Point Light', self.addItem), \
                   (self.ID_ADDSPOTLIGHT, 'Add Spot Light', self.addItem), \
                   (self.ID_ADDDIRECTIONALLIGHT, 'Add Directional Light', self.addItem)]
        buildMenu(mAdd, addList)
        
        mRun = wx.Menu()
        runList = [(self.ID_RUNEXCAVATION, 'Run Excavation', self.runExternal), \
                   (self.ID_RUNDISCOURSE, 'Run Discourse', self.runExternal)]
        buildMenu(mRun, runList)
        
        
        
    def newScene(self,event):
        pass
    def openScene(self, event):
        pass
    def saveScene(self, event):
        pass
    def saveSceneAs(self, event):
        pass
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
    
    def addItem(self, event):
        """Adds an item, determined by the event's sender id"""
        pass
    
    def runExternal(self, event):
        """Runs an external application"""
        pass
    
                
class ExEd(wx.App, DirectObject.DirectObject):
    """Panda object for handling all panda related tasks and events""" 
    def __init__(self): 
        wx.App.__init__(self) 
        self.replaceEventLoop()
        self.frame = PandaFrame(None, wx.ID_ANY, 'ExEd') 
        self.frame.Bind(wx.EVT_CLOSE, self.quit) 
        self.frame.actionManager = ActionManager()
        
        self.registerActions()
        
        self.wxStep()    
    
    
    def registerActions(self):
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
        if task != None: return task.cont 

app = ExEd() 
run() 