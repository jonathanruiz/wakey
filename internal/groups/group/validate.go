package group

import (
	"fmt"
	"strings"
)

func (m *Model) groupNameValidator(value string) error {
	// Check if the value is empty
	if value == "" {
		return fmt.Errorf("group name is required")
	}

	m.err[0] = nil
	return nil
}

func (m *Model) devicesValidator(value string) error {
	if value == "" {
		m.err[1] = nil
		return nil
	}

	// Check if the value is a valid device name
	for _, deviceName := range strings.Split(value, ", ") {
		if _, ok := m.deviceNameMap[deviceName]; !ok {
			return fmt.Errorf("'%s' device does not exist", deviceName)
		}
	}

	m.err[1] = nil
	return nil
}

func (m *Model) validateInput(index int, validator func(string) error) bool {
	if err := validator(m.inputs[index].Value()); err != nil {
		m.err[index] = err
		return false
	}
	return true
}
