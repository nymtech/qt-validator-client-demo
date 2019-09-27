// main.qml - qml definition for the gui application
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


// TODO: move stuff to separate components, it's too messy right now
Flickable {
    id: walletPage
    anchors.fill: parent

    ScrollIndicator.vertical: ScrollIndicator { }

	ColumnLayout {
		anchors.topMargin: 20
        anchors.top: parent.top
        anchors.right: parent.right
        anchors.left: parent.left

		Label {
            id: label3
            text: qsTr("Nym Wallet Demo")
            Layout.alignment: Qt.AlignLeft | Qt.AlignTop
            Layout.fillWidth: true
            // anchors.horizontalCenter: parent.horizontalCenter
            horizontalAlignment: Text.AlignHCenter
            font.weight: Font.DemiBold
            font.pointSize: 16
        }

        // I guess rather than 'visible' trick, a proper loader should have been used?
		ColumnLayout {
			id: accountFull
			Layout.alignment: Qt.AlignHCenter | Qt.AlignVCenter

			visible: false
			spacing: 5
			// anchors.fill: parent
            Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.rightMargin: 30
            Layout.leftMargin: 30
            Layout.bottomMargin: 30
            Layout.topMargin: 10

			ClientAccount {

			}
		}
		
		ColumnLayout {		
			id: configFull
			spacing: 5
			// anchors.fill: parent
			Layout.fillHeight: true
            Layout.fillWidth: true
            Layout.rightMargin: 30
            Layout.leftMargin: 30
            Layout.bottomMargin: 30
            Layout.topMargin: 10

			RowLayout {
				Layout.alignment: Qt.AlignTop
				Layout.preferredHeight: 40
				Layout.fillWidth: true
				TextField {
					id: path
					enabled: false
					text: "Please load Nym Client configuration file"
					Layout.fillWidth: true
				}
				Button {
					text: "open config"
					onClicked: fileDialog.open();
				}
			}

			ConfigSummary{
				id: configView
				opacity: 0 // temporary until I can figure out proper spacing with visible: false
			}

			RowLayout {
				Layout.fillHeight: false
				Layout.alignment: Qt.AlignBottom | Qt.AlignRight

				Layout.preferredHeight: 40
				Layout.fillWidth: true

				Button {
					Layout.alignment: Qt.AlignRight | Qt.AlignBottom
					id: configConfirmBtn
					enabled: false
					text: "confirm"
					Layout.fillHeight: false
					Layout.fillWidth: false
					onClicked: {
						QmlBridge.confirmConfig()
						configFull.visible = false
						accountFull.visible = true
					}
				}
			}
		}
	}

    QtLabs.FileDialog {
        id: fileDialog
        folder: QtLabs.StandardPaths.standardLocations(QtLabs.StandardPaths.HomeLocation)[0]
        nameFilters: [ "Config files (*.toml)", "All files (*)" ]
        // onFolderChanged: {
        //     folderModel.folder = folder;
        // }
        onAccepted: {
            QmlBridge.loadConfig(fileDialog.file)
			configConfirmBtn.enabled = true
			configView.opacity = 1
            path.text = fileDialog.file
        }
    }
    //   ProgressBar {
    //     id: loading
    //     anchors.horizontalCenter: parent.horizontalCenter
    //     visible: false
    //     indeterminate: true
    //   }
    // }



    Dialog {
        id: notificationDialog
        parent: ApplicationWindow.contentItem
        anchors.centerIn: ApplicationWindow.contentItem

        height: 200
        width: Math.min(ApplicationWindow.contentItem.width * 2/3, 800)

        modal: true

        closePolicy: Popup.CloseOnPressOutside | Popup.CloseOnEscape
        standardButtons: Dialog.Ok
        title: qsTr("Notification box")

        Label {
            id: notificationText
			wrapMode: Label.WordWrap
			// wrapMode: Label.WrapAnywhere

            width: notificationDialog.availableWidth
			height: notificationDialog.availableHeight
        }

        onAccepted: console.log("Ok clicked")
        onRejected: console.log("Cancel clicked")
    }

    Connections {
        target: QmlBridge
        onDisplayNotification: {
            notificationText.text = message
            notificationDialog.title = title

            notificationDialog.open()
        }
    }

}

























/*##^## Designer {
    D{i:0;autoSize:true;height:768;width:1024}D{i:28;anchors_height:200;anchors_width:200}
}
 ##^##*/
