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


from direct.showbase.ShowBase import ShowBase
from panda3d.core import WindowProperties
from panda3d.core import ConfigVariableString
from tools.scene import *
import wx
import os
import sys
import copy

class ExEd(wx.Frame, ShowBase): 
    def __init__(self, parent, title): 
        wx.Frame.__init__(self, parent, title=title, size=(800, 600))
        ShowBase.__init__(self)
        self.CreateStatusBar()
        
        
        #file menu
        fileMenu = wx.Menu()
        fileMenu.Append(wx.ID_ABOUT, "&About", " About ExEd")
        fileMenu.Append(wx.ID_EXIT, "E&xit", " Exit ExEd")
        
        menuBar = wx.MenuBar()
        menuBar.Append(fileMenu, "&File")
        self.SetMenuBar(menuBar)
        
        
        #override  wxEventLoop)
        self.evtloop = wx.EventLoop() 
        self.old = wx.EventLoop.GetActive() 
        wx.EventLoop.SetActive(self.evtloop) 
        taskMgr.add(self.wx, "Custom wx Event Loop")
        
        
        self.Show(True) 

    # wxWindows call to initialize the application 
    def OnInit(self): 
        self.SetAppName("ExEd") 
        self.SetClassName("ExEd") 

        base.windowType = 'onscreen' 
        props = WindowProperties.getDefault() 
        props.setParentWindow(self.GetHandle())
        base.openDefaultWindow(props = props) 

        base.setFrameRateMeter(True) 
        
        return True 

    def wx(self, task): 
        while self.evtloop.Pending(): 
            self.evtloop.Dispatch() 
        #time.sleep(0.01) 
        self.ProcessIdle() 
        return task.cont 

    def close(self): 
        wx.EventLoop.SetActive(self.old) 

wxApp = wx.App(False)
exed = ExEd(None, "ExEd - Scene Editor for Excavation") 
exed.run() 