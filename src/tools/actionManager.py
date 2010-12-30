class ActionManager():
    MAXUNDO = 1000
    lstUndo = []
    lstRedo = []
    actions = {}
    
    
    def registerAction(self, name, action):
        """Registers a method for use with the undo / redo queues"""
        if name in self.actions.keys():
            raise Exception("Action already registered")
        else:
            self.actions[name] = action
        
    def executeAction(self, name, **params):
        """Executes the given Action"""
        actionItem = ActionItem(self.actions[name], params)
        self.lstUndo.append(actionItem)
        self.lstRedo = []
        
        while len(self.lstUndo) > self.MAXUNDO:
            self.lstUndo.delete(0)
        
        actionItem.execute()
        
        
    def undo(self):
        """Reverses the previously executed action"""
        if len(self.lstUndo) > 0:
            actionItem = self.lstUndo.pop()
            self.lstRedo.append(actionItem)
            
            actionItem.undo()
            
        
    def redo(self):
        """Reverses the previously undone action if one exists"""
        if len(self.lstRedo) > 0:
            actionItem = self.lstRedo.pop()
            self.lstUndo.append(actionItem)
            
            actionItem.execute()
            
class Action():
    params = {}
    
    def __init__(self, method, undoMethod):
        self.method = method
        self.undoMethod = undoMethod
    
class ActionItem():
    
    def __init__(self, action, **params):
        self.action = action
        self.params = params
    
    def execute(self):
        method = self.action.method
        method(self.params)
    
    def undo(self):
        method = self.action.undoMethod
        method(self.params)
        
        