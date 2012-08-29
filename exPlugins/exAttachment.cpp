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
#include <QCoreApplication>
#include <QXmlTree/QXmlTreeNode.h>
//#include <Qt/qinputdialog.h>
//#include <Qt/qmessagebox.h>
#include <Qt/qtextstream.h>
//#include <Qt/qdir.h>
//#include <QtGui/QWizard>
#include <QtCore/qplugin.h>
#include <QTableWidget>
#include <QtDebug>
#include <horde3d/Horde3D.h>

exAttachment::exAttachment(QObject* parent /*= 0*/) : AttachmentPlugIn(parent)
{
	m_widget = new QTableWidget(1,2);
	m_widget->setVisible(false);
	QStringList headers;
	headers<<"Name"<<"Value";
	m_widget->setHorizontalHeaderLabels(headers);
	//connect(m_widget, SIGNAL(currentCellChanged(int, int, int, int)), this, SLOT(updateValue(int, int, int, int)));

	m_typeCombo = new QComboBox();
	connect(m_typeCombo, SIGNAL(currentIndexChanged(int)), this, SLOT(changeType(int)));
	m_currentNode = 0;
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
	qDebug() << "init";
	m_sceneFile = file;

	m_typeCombo->addItem("NONE");
	//Add Type combobox and label
	m_widget->setItem(0, 0, new QTableWidgetItem("Type", 0));
	
	//Load Entity config from file
	QString path = QCoreApplication::applicationDirPath();
	path.append("/exEntities.def");
	QFile entFile(path);
	entFile.open(QIODevice::ReadOnly | QIODevice::Text);
	QTextStream in(&entFile);
	QString line = in.readLine();
	while (!line.isNull()) {
		QStringList entList = line.split(",");
		m_typeCombo->addItem(entList[0], entList);
		line = in.readLine();
	}
	
	m_widget->setCellWidget(0, 1,m_typeCombo);

}

void exAttachment::setCurrentNode(QXmlTreeNode* parentNode)
{	
	qDebug() << "setCurrentNode";
	m_currentNode = parentNode;

	if (m_currentNode == 0) {
		changeType(0);
		return;
	}
	QDomElement attNode = m_currentNode->xmlNode().firstChildElement("Attachment");

	if (attNode.isNull()) {
		qDebug() << "setCurrentNode: No attachment";
		changeType(0);
		return;
	}
	
	qDebug() << "setCurrentNode: Attachment Exists";
	changeType(m_typeCombo->findText(attNode.attribute("type"), Qt::MatchExactly));
	//update table widget values

	for (int r = 1; r < m_widget->rowCount(); ++r)
	{
		m_widget->item(r, 1)->setText(attNode.attribute(m_widget->item(r, 0)->text(), ""));
	}
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
	qDebug() << "initNodeAttachment";
	Q_ASSERT(!sceneNode->xmlNode().firstChildElement("Attachment").isNull());
	//Nothing?
}

void exAttachment::destroyNodeAttachment(QXmlTreeNode* sceneNode)
{
	//Nothing?
}

void exAttachment::createNodeAttachment()
{	
	qDebug() << "createNodeAttachment";
	Q_ASSERT(m_currentNode != 0);	
	QDomElement node = m_currentNode->xmlNode().insertBefore(QDomDocument().createElement("Attachment"), QDomNode()).toElement();
	node.setAttribute("type", "NONE");
	initNodeAttachment(m_currentNode);
	setCurrentNode(m_currentNode);
}

void exAttachment::removeNodeAttachment()
{
	qDebug() << "removeNodeAttachment";
	QDomElement node = m_currentNode->xmlNode().firstChildElement("Attachment");
	m_currentNode->xmlNode().removeChild(node);	
	
	setCurrentNode(m_currentNode);
}

QXmlTreeModel* exAttachment::initExtras( const QDomElement &extraNode, QObject* parent)
{
	//There is nothing I want in the upper extras panel
	// it throws a warning on the command line because it's null
	// but it doesn't seem to break anything
	
	qDebug() << "initExtras";
	return NULL;
}

void exAttachment::sceneFileConfig()
{
	qDebug() << "sceneFileConfig";
}
void exAttachment::registerLuaFunctions(lua_State* lua) {}
QFileInfoList exAttachment::findReferences(const QDomElement &node) const {}

void exAttachment::changeType(int index)
{
	qDebug() << "Entering changeType";
	if (m_currentNode == 0) return;

	//set type
	QDomElement node = m_currentNode->xmlNode().firstChildElement("Attachment");
	
	qDebug() << "changeType";
	if (node.isNull()) return;
	node.setAttribute("type", m_typeCombo->currentText());

	if (m_typeCombo->currentText() == "NONE") {
		m_widget->setRowCount(1);
		return;
	}

	//Parse string into tablewidgets
	QStringList properties = m_typeCombo->itemData(index).toStringList();

	m_widget->setRowCount(properties.size());
	for (int r = 0; r < properties.size(); ++r)
	{
		m_widget->setItem(r, 0, new QTableWidgetItem(properties[r], 0));
		QLineEdit* lineWidget = new QLineEdit;
		connect(lineWidget, SIGNAL(editingFinished()), this, SLOT(updateValue()));
		m_widget->setCellWidget(r, 1, lineWidget);
		
	}
	
}
void exAttachment::updateValue()
{
	qDebug() << "updateValue";
	if (m_currentNode == 0) return;

	QDomElement node = m_currentNode->xmlNode().firstChildElement("Attachment");
	if (node.isNull()) return;

	for (int r = 1; r < m_widget->rowCount(); ++r)
	{
		qDebug() << "updating attribute: "; 
		qDebug() << m_widget->item(r, 1)->text();
		node.setAttribute(m_widget->item(r, 0)->text(), m_widget->item(r, 1)->text());
	}


	emit modified(true);
}

Q_EXPORT_STATIC_PLUGIN(exAttachment)
Q_EXPORT_PLUGIN2(exattachment, exAttachment)
