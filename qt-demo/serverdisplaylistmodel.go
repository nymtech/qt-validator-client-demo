// serverdisplaylistmodel.go
// Copyright (C) 2019  Jedrzej Stuczynski.
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
package main

import (
	"fmt"

	"github.com/therecipe/qt/core"
)

func init() { ServerDisplayListModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "ServerDisplayListModel") }

const (
	IdentifierRole = int(core.Qt__UserRole) + 1<<iota
	AddressRole
)

type ServerListItem struct {
	identifier string
	address    string
}

type ServerDisplayListModel struct {
	core.QAbstractListModel

	_ func()                                  `constructor:"init"`
	_ func()                                  `signal:"remove,auto"`
	_ func(obj []*core.QVariant)              `signal:"add,auto"`
	_ func(identifier string, address string) `signal:"edit,auto"`

	modelData []ServerListItem
}

func (m *ServerDisplayListModel) init() {
	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectData(m.data)
}

func (m *ServerDisplayListModel) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		IdentifierRole: core.NewQByteArray2("Identifier", -1),
		AddressRole:    core.NewQByteArray2("Address", -1),
	}
}

func (m *ServerDisplayListModel) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *ServerDisplayListModel) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
	switch role {
	case IdentifierRole:
		return core.NewQVariant1(item.identifier)
	case AddressRole:
		return core.NewQVariant1(item.address)
	}
	return core.NewQVariant()

	// item := m.modelData[index.Row()]
	// return core.NewQVariant1(fmt.Sprintf("%v %v", item.identifier, item.address))
}

func (m *ServerDisplayListModel) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *ServerDisplayListModel) add(item []*core.QVariant) {
	if len(item) != 2 {
		fmt.Println("trying to invalid element")
		return
	}
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(m.modelData, ServerListItem{item[0].ToString(), item[1].ToString()})
	m.EndInsertRows()
}

func (m *ServerDisplayListModel) edit(identifier string, address string) {
	if len(m.modelData) == 0 {
		return
	}
	m.modelData[len(m.modelData)-1] = ServerListItem{identifier, address}
	m.DataChanged(m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()), m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()), []int{int(core.Qt__DisplayRole)})
}
