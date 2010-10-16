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
            self.add_tags(tags)
        
    def add_tags(self, tags):
        """Adds the passed in tags to the tag dictionary for the grid
            The tags can be in a comma delimited string, a list or a dictionary"""
        if type(tags) is type(str()):
            if tags.__contains__(","):
                lstTags = tags.split(",")
                 
                self.__add_tags_list__(lstTags)
                    
            else:
                self.add_tag(tags)
        elif type(tags) is type(list()) or \
             type(tags) is type(tuple()):
            self.__add_tags_list__(tags)
            
        elif type(tags) is type(dict()):
            self.__add_tags_dict__(tags)
            
        elif type(tags) is type(self):  #if a taggroup is passed in
            self.__add_tags_dict__(tags.tags)
            
    def remove_tags(self, tags):
        """ removes the passed in tags from the tag group"""
        if type(tags) is type(str()):
            if tags.__contains__(","):
                lstTags = tags.split(",")
                 
                self.__add_tags_list__(lstTags, -1)
                    
            else:
                self.remove_tag(tags)
        elif type(tags) is type(list()) or \
             type(tags) is type(tuple()):
            self.__add_tags_list__(tags, -1)
            
        elif type(tags) is type(dict()):
            self.__add_tags_dict__(tags, -1)
            
        elif type(tags) is type(self):  #if a taggroup is passed in
            self.__add_tags_dict__(tags.tags, -1)
   
    def remove_tag(self, tag):
        """Removes one tag from the tag group"""
        self.add_tag(tag, -1)
        
    def add_tag(self, tag, value=1):
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
            
    def __add_tags_dict__(self, tags, value=1):
        """ Adds the tags in the dictionary to the existing tag dict"""
        for k, v in tags.items():
            self.add_tag(k, (v * value))
    
    def __add_tags_list__(self, tags, value=1):
        """ adds the list of tags to the existing tag dict"""
        for t in tags:
            self.add_tag(t.lstrip().rstrip(), value)

    def __to_json__(self):
        return self.tags