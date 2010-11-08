import wx 
import time 

import direct 
from pandac.PandaModules import * 
loadPrcFileData('startup', 'window-type none') 
from direct.directbase.DirectStart import * 
from direct.showbase import DirectObject 

class App(wx.PySimpleApp, DirectObject.DirectObject): 
    def __init__(self): 
        wx.PySimpleApp.__init__(self) 

        #Create a new event loop (to overide default wxEventLoop) 
        self.evtloop = wx.EventLoop() 
        self.old = wx.EventLoop.GetActive() 
        wx.EventLoop.SetActive(self.evtloop) 
        taskMgr.add(self.wx, "Custom wx Event Loop") 

    # wxWindows calls this method to initialize the application 
    def OnInit(self): 
        self.SetAppName('My wx app') 
        self.SetClassName('MyAppClass') 

        self.parent = wx.MDIParentFrame(None, -1, 'My wx app') 
        self.child = wx.MDIChildFrame(self.parent, -1, 'Panda window') 

        self.parent.SetClientSize((600, 400)) 
        self.parent.Show(True) 
        self.child.SetClientSize((400, 300)) 
        self.child.Show(True) 

        base.windowType = 'onscreen' 
        props = WindowProperties.getDefault() 
        props.setParentWindow(self.parent.GetHandle()) 
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

app = App() 
run() 