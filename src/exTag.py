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

class TagGroup:
    """
        Object to hold groups of tags and handle the operations on them
        such as calculation and evaluation strings
    """
    def __init__(self, 
                 tags={}):
        if type(tags) is type(dict()):
            self.tags = tags
        else:
            self.tags = {}
            self.addTags(tags)
        
    def addTags(self, tags):
        """Adds the passed in tags to the tag dictionary for the grid
            The tags can be in a comma delimited string, a list or a dictionary"""
        if type(tags) is type(str()):
            if tags.__contains__(","):
                lstTags = tags.split(",")
                 
                self.__addTagsList__(lstTags)
                    
            else:
                self.addTag(tags)
        elif type(tags) is type(list()) or \
             type(tags) is type(tuple()):
            self.__addTagsList__(tags)
            
        elif type(tags) is type(dict()):
            self.__addTagsDict__(tags)
            
        elif type(tags) is type(self):  #if a taggroup is passed in
            self.__addTagsDict__(tags.tags)
            
    def removeTags(self, tags):
        """ removes the passed in tags from the tag group"""
        if type(tags) is type(str()):
            if tags.__contains__(","):
                lstTags = tags.split(",")
                 
                self.__addTagsList__(lstTags, -1)
                    
            else:
                self.removeTag(tags)
        elif type(tags) is type(list()) or \
             type(tags) is type(tuple()):
            self.__addTagsList__(tags, -1)
            
        elif type(tags) is type(dict()):
            self.__addTagsDict__(tags, -1)
            
        elif type(tags) is type(self):  #if a taggroup is passed in
            self.__addTagsDict__(tags.tags, -1)
   
    def removeTag(self, tag):
        """Removes one tag from the tag group"""
        self.addTag(tag, -1)
        
    def addTag(self, tag, value=1):
        """Adds a single tag to the tag group"""
        if self.tags.has_key(tag):
            self.tags[tag] = self.tags[tag] + value
        else:
            self.tags[tag] = value
            
        if self.tags[tag] < 1:
            del self.tags[tag]
            
    def evaluate(self, evalString):
        if evalString <> "":
            try:
                if not eval(evalString):
                    return False
            except:
                return False
        return True
    
    def __eq__(self, other):
        """Checks for equality between tag groups"""
        if self.tags == other.tags:
            return True
        else:
            return False
            
    def __addTagsDict__(self, tags, value=1):
        """ Adds the tags in the dictionary to the existing tag dict"""
        for k, v in tags.items():
            self.addTag(k, (v * value))
    
    def __addTagsList__(self, tags, value=1):
        """ adds the list of tags to the existing tag dict"""
        for t in tags:
            self.addTag(t.lstrip().rstrip(), value)
