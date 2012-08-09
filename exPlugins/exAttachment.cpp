// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

// ****************************************************************************************
// excavation attachment plugin for storing excavation game entity data
//   in the scene tree
// 
// ****************************************************************************************
#include "exAttachment.h"

//#include "GameControllerWidget.h"
#include "AttachmentTreeModel.h"

#include "ExtraTreeModel.h"
#include "SceneFile.h"
#include "PlugInManager.h"
//#include "CustomAttachmentTypes.h"

//#include "QExtraNode.h"
//#include "QGameEntityNode.h"

//#include <QPropertyEditor/QPropertyEditorWidget.h>

//#include <GameEngine/GameEngine.h>
//#include <GameEngine/GameEngine_BulletPhysics.h>
//#include <GameEngine/GameEngine_Sound.h>
//#include <GameEngine/GameEngine_SceneGraph.h>

#include <Qt/qinputdialog.h>
#include <Qt/qmessagebox.h>
#include <Qt/qtextstream.h>
#include <Qt/qdir.h>
#include <QtGui/QWizard>
#include <QtCore/qplugin.h>
#include <QTextEdit>
//extern "C"
//{
//	#include "Lua/lua.h"
//}

//#include <horde3d/horde3dutils.h>

exAttachment::exAttachment(QObject* parent /*= 0*/) : AttachmentPlugIn(parent)
{
	m_widget = new QTextEdit();
	//m_widget->setVisible(false);
	connect(m_widget, SIGNAL(modified(bool)), this, SIGNAL(modified(bool)));
	//CustomAttachmentTypes::registerTypes();
}

exAttachment::~exAttachment() 
{
	delete m_widget;
}

QWidget* exAttachment::configurationWidget()
{
	return m_widget;
}

void exAttachment::init(SceneFile* file, QPropertyEditorWidget* widget) 
{

	if (file)
	{
	}
	else
	{ 
	}
	m_sceneFile = file;
}

void exAttachment::setCurrentNode(QXmlTreeNode* parentNode)
{	

}

void exAttachment::update()
{
}

void exAttachment::render(int activeCameraID)
{
}

void exAttachment::initNodeAttachment(QXmlTreeNode* sceneNode)
{	
	Q_ASSERT(!sceneNode->xmlNode().firstChildElement("Attachment").isNull());

}

void exAttachment::destroyNodeAttachment(QXmlTreeNode* sceneNode)
{
}

void exAttachment::createNodeAttachment()
{	
	Q_ASSERT(m_currentNode != 0);	
}

void exAttachment::removeNodeAttachment()
{
}

QXmlTreeModel* exAttachment::initExtras( const QDomElement &extraNode, QObject* parent)
{
	ExtraTreeModel* model = new ExtraTreeModel(m_sceneFile->pluginManager(), extraNode, parent);	
	return model;
}

void exAttachment::sceneFileConfig()
{
}

void exAttachment::registerLuaFunctions(lua_State* lua)
{
}

QFileInfoList exAttachment::findReferences(const QDomElement &node) const
{
	QFileInfoList references;	
	return references;
}



Q_EXPORT_STATIC_PLUGIN(exAttachment)
Q_EXPORT_PLUGIN2(exAttachment, exAttachment)
