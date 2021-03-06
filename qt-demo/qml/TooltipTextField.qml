// TooltipTextField.qml - qml definition for textfield with tooltip
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

Item {
    property string textFieldText: ""
    property string textFieldPlaceholderText: ""
    property string tooltipText: ""
    implicitHeight: textField.implicitHeight
    implicitWidth: textField.implicitWidth
    TextField {
        anchors.fill: parent
        id: textField
        enabled: false
        text: textFieldText
        placeholderText: textFieldPlaceholderText
    }
    MouseArea {
        anchors.fill: textField
        acceptedButtons: Qt.NoButton
        cursorShape: Qt.PointingHandCursor
        hoverEnabled: true
        ToolTip.text: tooltipText
        ToolTip.visible: containsMouse
    }
}
