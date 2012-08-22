// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

// ****************************************************************************************
// excavation attachment plugin for storing excavation game entity data
//   in the scene tree
// 
// ****************************************************************************************
#include "exAttachment.h"

#include <QFile>
#include <QTextStream>
#include <QString>
#include <QStringList>
#include <QLineEdit>
#include <QComboBox>
#include <QXmlTree/QXmlTreeNode.h>
//#include <Qt/qinputdialog.h>
//#include <Qt/qmessagebox.h>
#include <Qt/qtextstream.h>
//#include <Qt/qdir.h>
//#include <QtGui/QWizard>
#include <QtCore/qplugin.h>
#include <QTableWidget>
#include <horde3d/Horde3D.h>

exAttachment::exAttachment(QObject* parent /*= 0*/) : AttachmentPlugIn(parent)
{
	m_widget = new QTableWidget(1,2);
	m_widget->setVisible(false);
	QStringList headers;
	headers<<"Name"<<"Value";
	m_widget->setHorizontalHeaderLabels(headers);
	//connect(m_widget, SIGNAL(itemSelectionChanged()), this, SLOT(updateValue()));
	connect(m_widget, SIGNAL(currentitemChanged(QTableWidgetItem*)), this, SLOT(updateValue()));

	m_typeCombo = new QComboBox;
	connect(m_typeCombo, SIGNAL(currentIndexChanged(int)), this, SLOT(changeType(int)));
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

	//Add Type combobox and label
	m_widget->setItem(0, 0, new QTableWidgetItem("Type", 0));
		
	//Load Entity config from file
	QFile entFile("exEntities.def");
	entFile.open(QIODevice::ReadOnly | QIODevice::Text);
	QTextStream in(&entFile);
	QString line = in.readLine();
	while (!line.isNull()) {
		QStringList entList = line.split(",");
		m_typeCombo->addItem(entList.at(0), entList);
	}
	
	m_widget->setCellWidget(0, 1,m_typeCombo);

}

void exAttachment::setCurrentNode(QXmlTreeNode* parentNode)
{	
	m_currentNode = parentNode;
	QDomElement attNode = m_currentNode->xmlNode().firstChildElement("Attachment");

	changeType(m_typeCombo->findText(attNode.attribute("type"), Qt::MatchExactly));
	//update table widget values

	for (int r = 0; r < m_widget->rowCount(); ++r)
	{
		m_widget->item(r, 1)->setText(attNode.attribute(m_widget->item(r, 0)->text(), ""));
	}
}

QXmlTreeModel* exAttachment::initExtras( const QDomElement &extraNode, QObject* parent)
{
	//There is nothing I want in the upper extras panel
	// it throws a warning on the command line because it's null
	// but it doesn't seem to break anything
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
	h3dRender(activeCameraID);
	h3dFinalizeFrame();
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
	//TODO: Fix
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

void exAttachment::changeType(int index)
{
	//Parse string into tablewidgets
	QStringList properties = m_typeCombo->itemData(index).toStringList();

	m_widget->setRowCount(properties.size());
	for (int r = 0; r < properties.size(); ++r)
	{
		m_widget->setItem(r, 0, new QTableWidgetItem(properties[r], 0));
		m_widget->setCellWidget(r, 1, new QLineEdit);
	}
	
}
void exAttachment::updateValue(QTableWidgetItem* current, QTableWidgetItem* previous)
{
	//if (m_currentNode == 0) return;

	//QDomDocument doc = m_currentNode->xmlNode().ownerDocument();
	//QDomElement node = m_currentNode->xmlNode().firstChildElement("Attachment");

	//QDomElement newNode = doc.createElement(QString("Attachment"));
	//QDomText newNodeText = doc.createTextNode(m_widget->toPlainText()); 
	//newNode.appendChild(newNodeText); 

	//m_currentNode->xmlNode().replaceChild(newNode, node); 
	
	emit modified(true);
}

Q_EXPORT_STATIC_PLUGIN(exAttachment)
Q_EXPORT_PLUGIN2(exattachment, exAttachment)
