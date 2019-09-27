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

import QtQuick 2.13
import QtQuick.Window 2.2
import QtQuick.Controls 2.5
import QtQuick.Layouts 1.13
// import QtQuick.Dialogs 1.3
import QtQuick.Controls.Material 2.12
import Qt.labs.folderlistmodel 2.13
import Qt.labs.platform 1.1

ApplicationWindow {
	id: mainWindow

	Material.primary: "Indigo"

	visible: true
	title: qsTr("Nym Demo Application")
	minimumWidth: 900
	minimumHeight: 900

	// footer: Text {
	// 	text: "TODO: format and insert link here"
	// }

	ColumnLayout {
		id: mainView
		visible: true

		anchors.fill: parent
		anchors.margins: 20
		RowLayout {
			Layout.fillWidth: true
			Layout.alignment: Qt.AlignHCenter
			spacing: 30
			TooltipButton{
				btnText: qsTr("Demo 'Connect with Nym'")
				tooltipText: qsTr("The feature hasn't been implemented yet")
			}
			TooltipButton{
				btnText: qsTr("Demo Identity Manager")
				tooltipText: qsTr("The feature hasn't been implemented yet")
			}
			Button {
				id: walletBtn
				text: qsTr("Demo Nym Wallet")
					onClicked: {
						mainView.visible = !mainView.visible
						nymWallet.visible = !nymWallet.visible
					}
			}
		}
	}
	
	// I guess rather than 'visible' trick, a proper loader should have been used?
	NymWallet {
		id: nymWallet
		visible: false
	}
}


