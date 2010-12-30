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
        self.lstUndo.append(ActionItem(self.actions[name], params))
        
        method = self.actions[name]
        method(params)
        
    def undo(self):
        """Reverses the previously executed action"""
        pass
        
    def redo(self):
        """Reverses the previously undone action if one exists"""
        pass
        
class Action():
    params = {}
    
    def __init__(self, method, undoMethod):
        self.method = method
        self.undoMethod = undoMethod
    
class ActionItem():
    
    def __init__(self, action, **params):
        self.action = action
        self.params = params