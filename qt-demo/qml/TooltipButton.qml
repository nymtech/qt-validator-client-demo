// TooltipButton.qml - qml definition for button with tooltip
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
    property string btnText: ""
    property string tooltipText: ""
    implicitHeight: btn.implicitHeight
    implicitWidth: btn.implicitWidth
    Button {
        anchors.fill: parent
        id: btn
        enabled: false
        text: btnText
    }
    MouseArea {
        anchors.fill: btn
        acceptedButtons: Qt.NoButton
        cursorShape: Qt.PointingHandCursor
        hoverEnabled: true
        ToolTip.text: tooltipText
        ToolTip.visible: containsMouse
    }
}
