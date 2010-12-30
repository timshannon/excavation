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
    
    def __init__(self, *args, **kwargs): 
        wx.Frame.__init__(self, *args, **kwargs) 
        self.Show()
        self.pandapanel = PandaPanel(self, wx.ID_ANY, size=self.ClientSize) 
        self.pandapanel.initialize()
                
        self.CreateStatusBar()
        
        
    def undo(self):
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