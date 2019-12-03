/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package xsettings

import (
	"errors"
	"fmt"
)

var (
	errPropInvalid      = fmt.Errorf("Invalid xsettings property")
	errPropNotFound     = fmt.Errorf("This property not found")
	errPropTypeNotMatch = fmt.Errorf("This property's type not match")
)

func (m *XSManager) ListProps() string {
	datas, err := getSettingPropValue(m.owner, m.conn)
	if err != nil {
		return ""
	}

	infos := unmarshalSettingData(datas)
	if infos == nil || len(infos.items) == 0 {
		return ""
	}
	return infos.items.listProps()
}

func (m *XSManager) SetInteger(prop string, v int32) error {
	var setting = xsSetting{
		sType: settingTypeInteger,
		prop:  prop,
		value: v,
	}

	err := m.setSettings([]xsSetting{setting})
	if err != nil {
		logger.Debugf("Set '%s' to '%v' failed: %v", prop, v, err)
		return err
	}
	err = m.setGSettingsByXProp(prop, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *XSManager) GetInteger(prop string) (int32, error) {
	v, sType, err := m.getSettingValue(prop)
	if err != nil {
		logger.Debugf("Get '%s' value failed: %v", prop, err)
		return -1, err
	}

	if sType != settingTypeInteger {
		return -1, errPropTypeNotMatch
	}

	return v.(*integerValueInfo).value, nil
}

func (m *XSManager) SetString(prop, v string) error {
	var setting = xsSetting{
		sType: settingTypeString,
		prop:  prop,
		value: v,
	}

	err := m.setSettings([]xsSetting{setting})
	if err != nil {
		logger.Debugf("Set '%s' to '%v' failed: %v", prop, v, err)
		return err
	}
	err = m.setGSettingsByXProp(prop, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *XSManager) GetString(prop string) (string, error) {
	v, sType, err := m.getSettingValue(prop)
	if err != nil {
		logger.Debugf("Get '%s' value failed: %v", prop, err)
		return "", err
	}

	if sType != settingTypeString {
		return "", errPropTypeNotMatch
	}

	return v.(*stringValueInfo).value, nil
}

func (m *XSManager) SetColor(prop string, v []uint16) error {
	if len(v) != 4 {
		return errors.New("length of value is not 4")
	}

	var val [4]uint16
	copy(val[:], v)

	var setting = xsSetting{
		sType: settingTypeColor,
		prop:  prop,
		value: val,
	}

	err := m.setSettings([]xsSetting{setting})
	if err != nil {
		logger.Debugf("Set '%s' to '%v' failed: %v", prop, val, err)
		return err
	}
	err = m.setGSettingsByXProp(prop, val)
	if err != nil {
		return err
	}

	return nil
}

func (m *XSManager) GetColor(prop string) ([]uint16, error) {
	v, sType, err := m.getSettingValue(prop)
	if err != nil {
		logger.Debugf("Get '%s' value failed: %v", prop, err)
		return nil, err
	}

	if sType != settingTypeColor {
		return nil, errPropTypeNotMatch
	}

	tmp := v.(*colorValueInfo)

	return []uint16{tmp.red, tmp.blue, tmp.green, tmp.alpha}, nil
}

func (m *XSManager) getSettingValue(prop string) (interface{}, uint8, error) {
	m.settingsLocker.RLock()
	defer m.settingsLocker.RUnlock()
	datas, err := getSettingPropValue(m.owner, m.conn)
	if err != nil {
		return nil, 0, err
	}

	xsInfo := unmarshalSettingData(datas)
	item := xsInfo.getPropItem(prop)
	if item == nil {
		return nil, 0, errPropNotFound
	}

	return item.value, item.header.sType, nil
}

func (m *XSManager) setGSettingsByXProp(prop string, v interface{}) error {
	info := gsInfos.getByXSKey(prop)
	if info == nil {
		return errPropNotFound
	}

	return info.setValue(m.gs, v)
}

func (m *XSManager) GetScaleFactor() float64 {
	return m.getScaleFactor()
}

func (m *XSManager) SetScaleFactor(scale float64) error {
	primary, err := getPrimaryScreenName(m.conn)
	if err != nil {
		return err
	}

	err = m.setScreenScaleFactors(map[string]float64{primary: scale}, true)
	return err
}

func (m *XSManager) SetScreenScaleFactors(factors map[string]float64) error {
	err := m.setScreenScaleFactors(factors, true)
	return err
}

func (m *XSManager) GetScreenScaleFactors() (map[string]float64, error) {
	v := m.getScreenScaleFactors()
	return v, nil
}
