// Copyright 2012 Tim Shannon. All rights reserved. 
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file. 

// ****************************************************************************************
// excavation attachment plugin for storing excavation game entity data
//   in the scene tree
// 
// ****************************************************************************************

#ifndef EXATTACHMENT_H_
#define EXATTACHMENT_H_

#include "AttachmentPlugIn.h"
#include <QTableWidget>
#include <QComboBox>
/*class exWidget;*/

class exAttachment : public AttachmentPlugIn
{
	Q_OBJECT
	Q_INTERFACES(AttachmentPlugIn)

public:
	exAttachment(QObject* parent = 0);

	virtual ~exAttachment();

	QString plugInName() const {return QString("exAttachment");}

	QWidget* configurationWidget();

	void init(SceneFile* file, QPropertyEditorWidget* widget);

	/**
	 * \brief Indicates the currently selected node. 
	 * 
	 * This method has to be called whenever a new node has been selected whose attachment should be managed by this plugin
	 * @param parentNode the node whose attachment should be configured
	 */
	void setCurrentNode(QXmlTreeNode* parentNode);


	/**
	 * Updates the component managers of the game engine
	 */
	void update();

	/**
	 * Updates the active camera
	 */
	void render( int activeCameraID );

	/**
	 * \brief let the plugin initialize the node's attachment data when loading the node
	 * 
	 */
	void initNodeAttachment(QXmlTreeNode* sceneNode);

	/**
	 * \brief let the plugin destroy the node's attachment data when removing the whole scene node from the graph
	 */
	void destroyNodeAttachment(QXmlTreeNode* sceneNode);

	/**
	 * \brief create a new attachment for the current scene node
	 */
	void createNodeAttachment();

	/**
	 * \brief remove an existing attachment from the current scene node
	 */
	void removeNodeAttachment();


	/**
	 * Loads the scene file's extra node and returns the newly created QXmlTreeNode instance of it
	 */
	virtual QXmlTreeModel* initExtras(  const QDomElement& extraNode, QObject* parent );

	/**
	 * Calls the configuration wizard for this plugin when a new scene file has been created
	 */
	void sceneFileConfig();

	
	void registerLuaFunctions(lua_State* lua);
	/**
	 * \brief Returns all files referenced in the given node
	 *
	 * If a save as or clean scene command is triggered, the Editor will parse the complete scene file and
	 * scene graph files. If it encounters an Attachment node, it will call this methods that should return all files
	 * referenced in the given node
	 * @param node the attachment node that will be parsed for file references
	 * @return QFileInfoList all files referenced in the given node
	 */
	QFileInfoList findReferences(const QDomElement& node) const;

signals:
	void modified(bool);
	
private slots:
	void updateValue();
	void changeType(int);

protected:
	SceneFile*	m_sceneFile;
	
	QComboBox* 	m_typeCombo;

	QTableWidget*	m_widget;

	QXmlTreeNode*	m_currentNode;

};
#endif
