// credentiallistmodel.go
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
	"github.com/therecipe/qt/core"
)

func init() {
	CredentialListModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "CredentialListModel")
}

// TODO: can they be identical (value-wise) to server display roles?
const (
	CredentialRole = int(core.Qt__UserRole) + 1<<iota
	SequenceRole
	ValueRole
)

type CredentialListItem struct {
	sequence   string
	credential string
	value      uint64
}

type CredentialListModel struct {
	core.QAbstractListModel

	_         func()                        `constructor:"init"`
	_         func()                        `signal:"remove,auto"`
	_         func(item CredentialListItem) `signal:"addItem,auto"`
	modelData []CredentialListItem
}

func (m *CredentialListModel) init() {
	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectData(m.data)
}

func (m *CredentialListModel) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		CredentialRole: core.NewQByteArray2("Credential", -1),
		SequenceRole:   core.NewQByteArray2("Sequence", -1),
		ValueRole:      core.NewQByteArray2("Value", -1),
	}
}

func (m *CredentialListModel) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *CredentialListModel) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
	switch role {
	case CredentialRole:
		return core.NewQVariant1(item.credential)
	case SequenceRole:
		return core.NewQVariant1(item.sequence)
	case ValueRole:
		return core.NewQVariant1(item.value)
	}
	return core.NewQVariant()
}

func (m *CredentialListModel) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *CredentialListModel) addItem(item CredentialListItem) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(m.modelData, item)
	m.EndInsertRows()
}
