// ConfigSummary.qml - client config summary view
// Copyright (C) 2018-2019  Jedrzej Stuczynski.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import QtQuick 2.12
import QtQuick.Controls 2.5
import QtQuick.Layouts 1.12
import Qt.labs.platform 1.1 as QtLabs
import CustomQmlTypes 1.0


ColumnLayout {
    property var config: ConfigBridge
    Layout.topMargin: 10
    Layout.alignment: Qt.AlignLeft | Qt.AlignTop
    spacing: 20

    id: configView
    Layout.fillWidth: true
    Layout.fillHeight: true

    Label {
        id: label2
        text: qsTr("Configuration Details")
        Layout.rowSpan: 2
        Layout.columnSpan: 4
        font.pointSize: 14
        font.weight: Font.DemiBold
        Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter
    }

    GridLayout {
        id: gridLayout
        width: 100
        height: 100
        rowSpacing: 5
        columnSpacing: 5
        rows: 2
        columns: 5
        Layout.fillHeight: true
        Layout.fillWidth: true

        Label {
            text: "Client id: "
            font.weight: Font.DemiBold
        }

        TextField {
            id: clientId
            enabled: false
            text: configView.config.identifier
            Layout.fillWidth: false
        }

        ToolSeparator {
            id: toolSeparator
            opacity: 0
        }

        Label {
            text: "Client address: "
            font.weight: Font.DemiBold
        }

        TextField {
            id: clientAddress
            enabled: false
            text: configView.config.address
            Layout.fillWidth: true
        }

        Label {
            text: "Keyfile: "
            font.weight: Font.DemiBold
        }

        TextField {
            id: clientKeyfile
            enabled: false
            text: configView.config.keyfile
            Layout.columnSpan: 4
            Layout.fillWidth: true
        }

    }

    GroupBox {
        id: groupBox
        width: 200
        height: 200
        enabled: true
        Layout.fillHeight: false
        Layout.fillWidth: true
        title: qsTr("Ethereum Configuration")

        GridLayout {
            id: gridLayout1
            x: -12
            width: 100
            anchors.horizontalCenter: parent.horizontalCenter
            anchors.left: parent.left
            anchors.bottom: parent.bottom
            anchors.top: parent.top
            anchors.bottomMargin: 10
            anchors.topMargin: 10
            rows: 3
            columns: 2
            Layout.fillHeight: true
            Layout.fillWidth: true

            Label {
                text: "Node address: "
                font.weight: Font.DemiBold
            }

            TextField {
                id: ethereumNodeName
                text: configView.config.ethereumNode
                enabled: false
                Layout.fillWidth: true
            }

            Label {
                text: "Nym ERC20 Contract:"
                font.weight: Font.DemiBold
            }

            TextField {
                id: nymERC20Address
                text: configView.config.nymERC20
                enabled: false
                Layout.fillWidth: true
            }

            Label {
                text: "Pipe Account Contract:"
                font.weight: Font.DemiBold
            }

            TextField {
                id: pipeAccountAddress
                text: configView.config.pipeAccount
                enabled: false
                Layout.fillWidth: true
            }
        }
    }

    ServerDisplayListModel {
        id: nymValidatorsListModel
    }

    ServerDisplayListModel {
        id: tendermintValidatorsListModel
    }

    RowLayout {
        id: rowLayout
        // Layout.rowSpan: 1
        Layout.fillWidth: true
        Layout.preferredWidth: 200

        GroupBox {
            id: groupBox1
            width: 200
            height: 200
            Layout.preferredHeight: 200
            Layout.minimumHeight: 200
            Layout.maximumHeight: 300
            // Layout.fillHeight: false
            // Layout.fillWidth: true
            Layout.preferredWidth: parent.width/2
            title: qsTr("Nym Validator Nodes")
            ScrollView {
                id: scrollView2
                anchors.bottomMargin: 5
                anchors.topMargin: 5
                anchors.fill: parent

                ListView {
                    id: nymValidatorsList
                    x: 160
                    y: -2
                    width: 110
                    height: 160
                    clip: true
                    anchors.top: parent.top
                    anchors.topMargin: 0
                    anchors.horizontalCenterOffset: 0
                    anchors.horizontalCenter: parent.horizontalCenter
                
            
                    model: nymValidatorsListModel
                    
                    delegate: Item {
                        x: 5
                        width: 80
                        height: 20
                        Row {
                            spacing: 5
                            Label {
                                text: Identifier
                                font.weight: Font.DemiBold
                            }
                            Text {
                                text: Address
                            }
                        }
                    }
                }
            }
        }

        GroupBox {
            id: groupBox2
            width: 200
            height: 200
            Layout.preferredHeight: 200
            Layout.minimumHeight: 200
            Layout.maximumHeight: 300
            // Layout.fillHeight: false
            // Layout.fillWidth: true
            Layout.preferredWidth: parent.width/2
            title: qsTr("Tendermint Validator Nodes")

            ScrollView {
                id: scrollView
                anchors.bottomMargin: 5
                anchors.topMargin: 5
                anchors.fill: parent


                ListView {
                    id: tendermintValidatorsList
                    x: 160
                    y: -2
                    width: 110
                    height: 160
                    clip: true
                    anchors.top: parent.top
                    anchors.topMargin: 0
                    anchors.horizontalCenterOffset: 0
                    anchors.horizontalCenter: parent.horizontalCenter
                
            
                    model: tendermintValidatorsListModel
                    
                    delegate: Item {
                        x: 5
                        width: 80
                        height: 20
                        Row {
                            spacing: 5
                            Label {
                                text: Identifier
                                font.weight: Font.DemiBold
                            }
                            Text {
                                text: Address
                            }
                        }
                    }
                }
            }
        }
    }

    Dialog {
        id: newKeyDialog
        parent: ApplicationWindow.contentItem
        anchors.centerIn: ApplicationWindow.contentItem

        height: 200
        width: Math.min(ApplicationWindow.contentItem.width * 2/3, 800)

        modal: true

        closePolicy: Popup.NoAutoClose	
        standardButtons: Dialog.Ok | Dialog.Close	
        title: qsTr("No compatible key detected")

        Label {
			wrapMode: Label.WordWrap
			// wrapMode: Label.WrapAnywhere
            text: qsTr("The keyfile specified by your configuration file could not be loaded - was its path specified correctly?\nDo you want to generate a fresh keypair and save it to the the specified location? If rejected the application will be terminated.")

            width: newKeyDialog.availableWidth
			height: newKeyDialog.availableHeight
        }

        onAccepted: QmlBridge.generateNewKey()
        onRejected: Qt.quit()
    }


    Connections {
        target: QmlBridge
		onNewNymValidator: {
            nymValidatorsListModel.add([identifier, address])
		}

        onNewTendermintValidator: {
            tendermintValidatorsListModel.add([identifier, address])
        }

        onShowNewKeyDialog: {
            newKeyDialog.open()
        }
    }

}