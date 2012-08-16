// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

// ****************************************************************************************
// excavation attachment plugin for storing excavation game entity data
//   in the scene tree
// 
// ****************************************************************************************
#include "exAttachment.h"

#include <QXmlTree/QXmlTreeNode.h>
//#include <Qt/qinputdialog.h>
//#include <Qt/qmessagebox.h>
#include <Qt/qtextstream.h>
//#include <Qt/qdir.h>
//#include <QtGui/QWizard>
#include <QtCore/qplugin.h>
#include <QTextEdit>
//#include <horde3d/horde3dutils.h>

exAttachment::exAttachment(QObject* parent /*= 0*/) : AttachmentPlugIn(parent)
{
	m_widget = new QTextEdit();
	m_widget->setVisible(false);
	//connect(m_widget, SIGNAL(modified(bool)), this, SIGNAL(modified(bool)));
	connect(m_widget, SIGNAL(textChanged()), this, SLOT(updateValue()));
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
	m_sceneFile = file;
}

void exAttachment::setCurrentNode(QXmlTreeNode* parentNode)
{	
	m_currentNode = parentNode;
	m_widget->setPlainText(parentNode->xmlNode().text());
}

QXmlTreeModel* exAttachment::initExtras( const QDomElement &extraNode, QObject* parent)
{
	return NULL;
}

void exAttachment::sceneFileConfig()
{
}
void exAttachment::update()
{
	//Nothing to update
}

void exAttachment::render(int activeCameraID)
{
	//Nothing to render
}

void exAttachment::initNodeAttachment(QXmlTreeNode* sceneNode)
{	
	Q_ASSERT(!sceneNode->xmlNode().firstChildElement("Attachment").isNull());
	//m_widget->setPlainText(sceneNode->xmlNode().text());
}

void exAttachment::destroyNodeAttachment(QXmlTreeNode* sceneNode)
{
	//Nothing?
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
	m_currentNode->xmlNode().removeChild(node);	
	
	setCurrentNode(m_currentNode);
}


void exAttachment::registerLuaFunctions(lua_State* lua) {}
QFileInfoList exAttachment::findReferences(const QDomElement &node) const {}
void exAttachment::updateValue()
{
	if (m_currentNode == 0) return;
	QDomElement oldNode = m_currentNode.firstChildElement("firstchild"); 
	
	QDomElement newNode = m_currentNode.creat createElement(QString("firstchild"));  
	//  QDomText newNodeText = doc.createTextNode(QString("New Text")); 
	//  newNodeTag.appendChild(newNodeText); 
	//   
	//   // replace existing node with new node 
	//   root.replaceChild(newNodeTag, nodeTag); 
	//
	emit modified(true);
}

Q_EXPORT_STATIC_PLUGIN(exAttachment)
Q_EXPORT_PLUGIN2(exattachment, exAttachment)
