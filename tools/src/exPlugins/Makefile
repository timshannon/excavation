#############################################################################
# Makefile for building: libexPlugins.so.1.0.0
# Generated by qmake (2.01a) (Qt 4.8.3) on: Sat Dec 1 13:12:44 2012
# Project:  exPlugins.pro
# Template: lib
# Command: /usr/bin/qmake -o Makefile exPlugins.pro
#############################################################################

####### Compiler, tools and options

CC            = gcc
CXX           = g++
DEFINES       = -DQT_WEBKIT -DQT_NO_DEBUG -DQT_GUI_LIB -DQT_CORE_LIB -DQT_SHARED
CFLAGS        = -m64 -pipe -O2 -Wall -W -D_REENTRANT -fPIC $(DEFINES)
CXXFLAGS      = -m64 -pipe -O2 -Wall -W -D_REENTRANT -fPIC $(DEFINES)
INCPATH       = -I/usr/share/qt4/mkspecs/linux-g++-64 -I. -I/usr/include/qt4/QtCore -I/usr/include/qt4/QtGui -I/usr/include/qt4 -I. -IQXmlTree -I.
LINK          = g++
LFLAGS        = -m64 -Wl,-O1 -shared -Wl,-soname,libexPlugins.so.1
LIBS          = $(SUBLIBS)  -L/usr/lib/x86_64-linux-gnu -lQtGui -lQtCore -lpthread 
AR            = ar cqs
RANLIB        = 
QMAKE         = /usr/bin/qmake
TAR           = tar -cf
COMPRESS      = gzip -9f
COPY          = cp -f
SED           = sed
COPY_FILE     = $(COPY)
COPY_DIR      = $(COPY) -r
STRIP         = strip
INSTALL_FILE  = install -m 644 -p
INSTALL_DIR   = $(COPY_DIR)
INSTALL_PROGRAM = install -m 755 -p
DEL_FILE      = rm -f
SYMLINK       = ln -f -s
DEL_DIR       = rmdir
MOVE          = mv -f
CHK_DIR_EXISTS= test -d
MKDIR         = mkdir -p

####### Output directory

OBJECTS_DIR   = ./

####### Files

SOURCES       = exAttachment.cpp moc_AttachmentPlugIn.cpp \
		moc_exAttachment.cpp \
		moc_QXmlTreeModel.cpp \
		moc_QXmlTreeNode.cpp \
		moc_QXmlTreeView.cpp
OBJECTS       = exAttachment.o \
		moc_AttachmentPlugIn.o \
		moc_exAttachment.o \
		moc_QXmlTreeModel.o \
		moc_QXmlTreeNode.o \
		moc_QXmlTreeView.o
DIST          = /usr/share/qt4/mkspecs/common/unix.conf \
		/usr/share/qt4/mkspecs/common/linux.conf \
		/usr/share/qt4/mkspecs/common/gcc-base.conf \
		/usr/share/qt4/mkspecs/common/gcc-base-unix.conf \
		/usr/share/qt4/mkspecs/common/g++-base.conf \
		/usr/share/qt4/mkspecs/common/g++-unix.conf \
		/usr/share/qt4/mkspecs/qconfig.pri \
		/usr/share/qt4/mkspecs/modules/qt_webkit_version.pri \
		/usr/share/qt4/mkspecs/features/qt_functions.prf \
		/usr/share/qt4/mkspecs/features/qt_config.prf \
		/usr/share/qt4/mkspecs/features/exclusive_builds.prf \
		/usr/share/qt4/mkspecs/features/default_pre.prf \
		/usr/share/qt4/mkspecs/features/release.prf \
		/usr/share/qt4/mkspecs/features/default_post.prf \
		/usr/share/qt4/mkspecs/features/unix/gdb_dwarf_index.prf \
		/usr/share/qt4/mkspecs/features/warn_on.prf \
		/usr/share/qt4/mkspecs/features/qt.prf \
		/usr/share/qt4/mkspecs/features/unix/thread.prf \
		/usr/share/qt4/mkspecs/features/moc.prf \
		/usr/share/qt4/mkspecs/features/resources.prf \
		/usr/share/qt4/mkspecs/features/uic.prf \
		/usr/share/qt4/mkspecs/features/yacc.prf \
		/usr/share/qt4/mkspecs/features/lex.prf \
		/usr/share/qt4/mkspecs/features/include_source_dir.prf \
		exPlugins.pro
QMAKE_TARGET  = exPlugins
DESTDIR       = 
TARGET        = libexPlugins.so.1.0.0
TARGETA       = libexPlugins.a
TARGETD       = libexPlugins.so.1.0.0
TARGET0       = libexPlugins.so
TARGET1       = libexPlugins.so.1
TARGET2       = libexPlugins.so.1.0

first: all
####### Implicit rules

.SUFFIXES: .o .c .cpp .cc .cxx .C

.cpp.o:
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o "$@" "$<"

.cc.o:
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o "$@" "$<"

.cxx.o:
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o "$@" "$<"

.C.o:
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o "$@" "$<"

.c.o:
	$(CC) -c $(CFLAGS) $(INCPATH) -o "$@" "$<"

####### Build rules

all: Makefile  $(TARGET)

$(TARGET):  $(OBJECTS) $(SUBLIBS) $(OBJCOMP)  
	-$(DEL_FILE) $(TARGET) $(TARGET0) $(TARGET1) $(TARGET2)
	$(LINK) $(LFLAGS) -o $(TARGET) $(OBJECTS) $(LIBS) $(OBJCOMP)
	-ln -s $(TARGET) $(TARGET0)
	-ln -s $(TARGET) $(TARGET1)
	-ln -s $(TARGET) $(TARGET2)



staticlib: $(TARGETA)

$(TARGETA):  $(OBJECTS) $(OBJCOMP) 
	-$(DEL_FILE) $(TARGETA) 
	$(AR) $(TARGETA) $(OBJECTS)

Makefile: exPlugins.pro  /usr/share/qt4/mkspecs/linux-g++-64/qmake.conf /usr/share/qt4/mkspecs/common/unix.conf \
		/usr/share/qt4/mkspecs/common/linux.conf \
		/usr/share/qt4/mkspecs/common/gcc-base.conf \
		/usr/share/qt4/mkspecs/common/gcc-base-unix.conf \
		/usr/share/qt4/mkspecs/common/g++-base.conf \
		/usr/share/qt4/mkspecs/common/g++-unix.conf \
		/usr/share/qt4/mkspecs/qconfig.pri \
		/usr/share/qt4/mkspecs/modules/qt_webkit_version.pri \
		/usr/share/qt4/mkspecs/features/qt_functions.prf \
		/usr/share/qt4/mkspecs/features/qt_config.prf \
		/usr/share/qt4/mkspecs/features/exclusive_builds.prf \
		/usr/share/qt4/mkspecs/features/default_pre.prf \
		/usr/share/qt4/mkspecs/features/release.prf \
		/usr/share/qt4/mkspecs/features/default_post.prf \
		/usr/share/qt4/mkspecs/features/unix/gdb_dwarf_index.prf \
		/usr/share/qt4/mkspecs/features/warn_on.prf \
		/usr/share/qt4/mkspecs/features/qt.prf \
		/usr/share/qt4/mkspecs/features/unix/thread.prf \
		/usr/share/qt4/mkspecs/features/moc.prf \
		/usr/share/qt4/mkspecs/features/resources.prf \
		/usr/share/qt4/mkspecs/features/uic.prf \
		/usr/share/qt4/mkspecs/features/yacc.prf \
		/usr/share/qt4/mkspecs/features/lex.prf \
		/usr/share/qt4/mkspecs/features/include_source_dir.prf \
		/usr/lib/x86_64-linux-gnu/libQtGui.prl \
		/usr/lib/x86_64-linux-gnu/libQtCore.prl
	$(QMAKE) -o Makefile exPlugins.pro
/usr/share/qt4/mkspecs/common/unix.conf:
/usr/share/qt4/mkspecs/common/linux.conf:
/usr/share/qt4/mkspecs/common/gcc-base.conf:
/usr/share/qt4/mkspecs/common/gcc-base-unix.conf:
/usr/share/qt4/mkspecs/common/g++-base.conf:
/usr/share/qt4/mkspecs/common/g++-unix.conf:
/usr/share/qt4/mkspecs/qconfig.pri:
/usr/share/qt4/mkspecs/modules/qt_webkit_version.pri:
/usr/share/qt4/mkspecs/features/qt_functions.prf:
/usr/share/qt4/mkspecs/features/qt_config.prf:
/usr/share/qt4/mkspecs/features/exclusive_builds.prf:
/usr/share/qt4/mkspecs/features/default_pre.prf:
/usr/share/qt4/mkspecs/features/release.prf:
/usr/share/qt4/mkspecs/features/default_post.prf:
/usr/share/qt4/mkspecs/features/unix/gdb_dwarf_index.prf:
/usr/share/qt4/mkspecs/features/warn_on.prf:
/usr/share/qt4/mkspecs/features/qt.prf:
/usr/share/qt4/mkspecs/features/unix/thread.prf:
/usr/share/qt4/mkspecs/features/moc.prf:
/usr/share/qt4/mkspecs/features/resources.prf:
/usr/share/qt4/mkspecs/features/uic.prf:
/usr/share/qt4/mkspecs/features/yacc.prf:
/usr/share/qt4/mkspecs/features/lex.prf:
/usr/share/qt4/mkspecs/features/include_source_dir.prf:
/usr/lib/x86_64-linux-gnu/libQtGui.prl:
/usr/lib/x86_64-linux-gnu/libQtCore.prl:
qmake:  FORCE
	@$(QMAKE) -o Makefile exPlugins.pro

dist: 
	@$(CHK_DIR_EXISTS) .tmp/exPlugins1.0.0 || $(MKDIR) .tmp/exPlugins1.0.0 
	$(COPY_FILE) --parents $(SOURCES) $(DIST) .tmp/exPlugins1.0.0/ && $(COPY_FILE) --parents AttachmentPlugIn.h exAttachment.h QXmlTree/QXmlNodePropertyCommand.h QXmlTree/QXmlTreeModel.h QXmlTree/QXmlTreeNode.h QXmlTree/QXmlTreeUndoCommands.h QXmlTree/QXmlTreeView.h .tmp/exPlugins1.0.0/ && $(COPY_FILE) --parents exAttachment.cpp .tmp/exPlugins1.0.0/ && (cd `dirname .tmp/exPlugins1.0.0` && $(TAR) exPlugins1.0.0.tar exPlugins1.0.0 && $(COMPRESS) exPlugins1.0.0.tar) && $(MOVE) `dirname .tmp/exPlugins1.0.0`/exPlugins1.0.0.tar.gz . && $(DEL_FILE) -r .tmp/exPlugins1.0.0


clean:compiler_clean 
	-$(DEL_FILE) $(OBJECTS)
	-$(DEL_FILE) *~ core *.core


####### Sub-libraries

distclean: clean
	-$(DEL_FILE) $(TARGET) 
	-$(DEL_FILE) $(TARGET0) $(TARGET1) $(TARGET2) $(TARGETA)
	-$(DEL_FILE) Makefile


check: first

mocclean: compiler_moc_header_clean compiler_moc_source_clean

mocables: compiler_moc_header_make_all compiler_moc_source_make_all

compiler_moc_header_make_all: moc_AttachmentPlugIn.cpp moc_exAttachment.cpp moc_QXmlTreeModel.cpp moc_QXmlTreeNode.cpp moc_QXmlTreeView.cpp
compiler_moc_header_clean:
	-$(DEL_FILE) moc_AttachmentPlugIn.cpp moc_exAttachment.cpp moc_QXmlTreeModel.cpp moc_QXmlTreeNode.cpp moc_QXmlTreeView.cpp
moc_AttachmentPlugIn.cpp: AttachmentPlugIn.h
	/usr/bin/moc-qt4 $(DEFINES) $(INCPATH) AttachmentPlugIn.h -o moc_AttachmentPlugIn.cpp

moc_exAttachment.cpp: AttachmentPlugIn.h \
		exAttachment.h
	/usr/bin/moc-qt4 $(DEFINES) $(INCPATH) exAttachment.h -o moc_exAttachment.cpp

moc_QXmlTreeModel.cpp: QXmlTree/QXmlTreeModel.h
	/usr/bin/moc-qt4 $(DEFINES) $(INCPATH) QXmlTree/QXmlTreeModel.h -o moc_QXmlTreeModel.cpp

moc_QXmlTreeNode.cpp: QXmlTree/QXmlTreeNode.h
	/usr/bin/moc-qt4 $(DEFINES) $(INCPATH) QXmlTree/QXmlTreeNode.h -o moc_QXmlTreeNode.cpp

moc_QXmlTreeView.cpp: QXmlTree/QXmlTreeView.h
	/usr/bin/moc-qt4 $(DEFINES) $(INCPATH) QXmlTree/QXmlTreeView.h -o moc_QXmlTreeView.cpp

compiler_rcc_make_all:
compiler_rcc_clean:
compiler_image_collection_make_all: qmake_image_collection.cpp
compiler_image_collection_clean:
	-$(DEL_FILE) qmake_image_collection.cpp
compiler_moc_source_make_all:
compiler_moc_source_clean:
compiler_uic_make_all:
compiler_uic_clean:
compiler_yacc_decl_make_all:
compiler_yacc_decl_clean:
compiler_yacc_impl_make_all:
compiler_yacc_impl_clean:
compiler_lex_make_all:
compiler_lex_clean:
compiler_clean: compiler_moc_header_clean 

####### Compile

exAttachment.o: exAttachment.cpp exAttachment.h \
		AttachmentPlugIn.h \
		QXmlTree/QXmlTreeNode.h
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o exAttachment.o exAttachment.cpp

moc_AttachmentPlugIn.o: moc_AttachmentPlugIn.cpp 
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o moc_AttachmentPlugIn.o moc_AttachmentPlugIn.cpp

moc_exAttachment.o: moc_exAttachment.cpp 
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o moc_exAttachment.o moc_exAttachment.cpp

moc_QXmlTreeModel.o: moc_QXmlTreeModel.cpp 
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o moc_QXmlTreeModel.o moc_QXmlTreeModel.cpp

moc_QXmlTreeNode.o: moc_QXmlTreeNode.cpp 
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o moc_QXmlTreeNode.o moc_QXmlTreeNode.cpp

moc_QXmlTreeView.o: moc_QXmlTreeView.cpp 
	$(CXX) -c $(CXXFLAGS) $(INCPATH) -o moc_QXmlTreeView.o moc_QXmlTreeView.cpp

####### Install

install:   FORCE

uninstall:   FORCE

FORCE:

