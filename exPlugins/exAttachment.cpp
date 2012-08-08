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

//extern "C"
//{
//	#include "Lua/lua.h"
}

#include <Horde3D/Horde3DUtils.h>

exAttachment::exAttachment(QObject* parent /*= 0*/) : AttachmentPlugIn(parent)
{
	m_widget = new exWidget();
	m_widget->setVisible(false);
	connect(m_widget, SIGNAL(modified(bool)), this, SIGNAL(modified(bool)));
	//CustomAttachmentTypes::registerTypes();
}

exAttachment::~exAttachment() 
{
	delete m_widget;
	GameEngine::release();
}

QWidget* exAttachment::configurationWidget()
{
	return m_widget;
}

void exAttachment::init(SceneFile* file, QPropertyEditorWidget* widget)
{
	if (file)
	{
		 //Horde3D specific initialization has been done by the HordeSceneEditor already
		 //so we only have to initialize directories and the attachment callback
		GameEngine::setSoundResourceDirectory( qPrintable(file->sceneFileDom().documentElement().firstChildElement("EnginePath").attribute("mediapath")) );

		file->pluginManager()->registerExtraNode("GameEntity", QGameEntityNode::loadNode, QGameEntityNode::createNode);
		//if( widget ) widget->registerCustomPropertyCB(CustomAttachmentTypes::createCustomProperty);

		GameEngine::init();
		m_widget->init();
	}
	else
	{ 
		if( m_sceneFile )
		{
			m_sceneFile->pluginManager()->unregisterExtraNode("GameEntity");
		}
		//if( widget ) widget->unregisterCustomPropertyCB(CustomAttachmentTypes::createCustomProperty);
		GameEngine::release();
		m_widget->release();
	}
	m_sceneFile = file;
}

void exAttachment::setCurrentNode(QXmlTreeNode* parentNode)
{	
	if( parentNode && parentNode->property("__AttachmentModel").isValid() && parentNode->property("__AttachmentModel").value<void*>() != 0)
		m_widget->setCurrentNode(parentNode);
	else	
		m_widget->setCurrentNode(0);			
	m_currentNode = parentNode;
}

void exAttachment::update()
{
	GameEngine::update();
}

void exAttachment::render(int activeCameraID)
{
	GameEngine::setActiveCamera( activeCameraID );
}

void exAttachment::initNodeAttachment(QXmlTreeNode* sceneNode)
{	
	Q_ASSERT(!sceneNode->xmlNode().firstChildElement("Attachment").isNull());

	QString entityName = sceneNode->xmlNode().firstChildElement("Attachment").attribute("name");
	// If there is already an entity with this name we have to rename this one
	while (GameEngine::entityWorldID(qPrintable(entityName)) != 0)
		entityName = entityName + QString::number(sceneNode->property("ID").toInt());
	if (entityName != sceneNode->xmlNode().firstChildElement("Attachment").attribute("name"))
		sceneNode->xmlNode().firstChildElement("Attachment").setAttribute("name", entityName);
	// Create the Entity within the GameEngine
	unsigned int entityID = GameEngine::createGameEntity( qPrintable(QString("<GameEntity type=\"%1\" name=\"%2\"/>").arg(plugInName()).arg(entityName)) );
	// Add Scene Graph Component
	GameEngine::setComponentData( entityID, "Horde3D", qPrintable( QString("<Horde3D id=%1 />").arg(sceneNode->property("ID").toInt()) ) ); 
	// Create a treemodel for the Attachment Widget
	AttachmentTreeModel* model = new AttachmentTreeModel(sceneNode, sceneNode->xmlNode().firstChildElement("Attachment"));
	// Store treemodel as dynamic property
	sceneNode->setProperty("__AttachmentModel", QVariant::fromValue<void*>(model));	
}

void exAttachment::destroyNodeAttachment(QXmlTreeNode* sceneNode)
{
	unsigned int entityID = GameEngine::entityWorldID( qPrintable(sceneNode->xmlNode().firstChildElement("Attachment").attribute("name")) );
	GameEngine::removeGameEntity(entityID);
}

void exAttachment::createNodeAttachment()
{	
	Q_ASSERT(m_currentNode != 0);	
	QDomElement node = m_currentNode->xmlNode().insertBefore(QDomDocument().createElement("Attachment"), QDomNode()).toElement();
	node.setAttribute("type", plugInName());
	node.setAttribute("name", m_currentNode->property("Name").toString() + "_" + m_currentNode->property("ID").toString());
	initNodeAttachment(m_currentNode);
	setCurrentNode(m_currentNode);
}

void exAttachment::removeNodeAttachment()
{
	QDomElement node = m_currentNode->xmlNode().firstChildElement("Attachment");
	unsigned int entityID = GameEngine::entityWorldID( qPrintable(node.attribute("name")) );
	// Reset scene graph component to avoid removement of Horde's scenegraph node
	GameEngine::setEntitySceneGraphID( entityID, 0 );
	m_currentNode->xmlNode().removeChild(node);	
	if( m_currentNode->property("__AttachmentModel").isValid() )
	{
		delete static_cast<AttachmentTreeModel*>(m_currentNode->property("__AttachmentModel").value<void*>());
		m_currentNode->setProperty("__AttachmentModel", QVariant::fromValue<void*>(0));
		GameEngine::removeGameEntity( entityID );
	}
	setCurrentNode(m_currentNode);
}

QXmlTreeModel* exAttachment::initExtras( const QDomElement &extraNode, QObject* parent)
{
	ExtraTreeModel* model = new ExtraTreeModel(m_sceneFile->pluginManager(), extraNode, parent);	
	return model;
}

void exAttachment::sceneFileConfig()
{
	QDomDocument sceneFile(m_sceneFile->sceneFileDom());
	QDomElement pathNode(sceneFile.documentElement().firstChildElement("EnginePath"));
	// Create a wizard for the configuration of the directories
	QWizard wizard;
	PathPage* page = new PathPage(&wizard);	
	page->setDirectories( 
		pathNode.attribute("mediapath"), 
		pathNode.attribute("trackerpath"), 
		pathNode.attribute("scriptpath")
	);
	wizard.addPage(page);
	if (wizard.exec() == QDialog::Accepted)
	{
		pathNode.setAttribute("mediapath", wizard.field("mediadir").toString());
		//pathNode.setAttribute("trackerpath", wizard.field("trackerdir").toString());
		pathNode.setAttribute("scriptpath", wizard.field("scriptdir").toString());
	}
}

void exAttachment::registerLuaFunctions(lua_State* lua)
{
	GameEngine::registerLuaStack(lua);
}

QFileInfoList exAttachment::findReferences(const QDomElement &node) const
{
	QFileInfoList references;	
	if (node.tagName() == "Sound3D" && node.hasAttribute("file"))
	{
		QFileInfo file(node.attribute("file"));
		references.append(file);
	}
	if( node.tagName() == "StaticAnimation" && node.hasAttribute("file"))
	{
		QFileInfo file(node.attribute("file"));
		references.append(file);
	}
	QDomNodeList children = node.childNodes();
	for (int i=0; i<children.size(); ++i)
		references << findReferences(children.at(i).toElement());
	return references;
}



Q_EXPORT_STATIC_PLUGIN(exAttachment)
Q_EXPORT_PLUGIN2(exAttachment, exAttachment)
