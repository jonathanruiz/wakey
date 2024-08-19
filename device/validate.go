package device

import "fmt"

func (m *Model) deviceNameValidator(value string) error {
	// Check if the value is empty
	if value == "" {
		return fmt.Errorf("device name is required")
	}

	m.err[0] = nil
	return nil
}

func (m *Model) descriptionValidator(value string) error {
	// Check if the value is empty
	if value == "" {
		return fmt.Errorf("description is required")
	}

	// Check max length
	if len(value) > 64 {
		return fmt.Errorf("description must be less than 64 characters")
	}

	m.err[1] = nil
	return nil
}

func (m *Model) macAddressValidator(value string) error {
	var a, b, c, d, e, f int

	// Check if the value is empty
	if value == "" {
		return fmt.Errorf("mac address is required")
	}

	// Check if the value is a valid MAC address
	if _, err := fmt.Sscanf(value, "%2x:%2x:%2x:%2x:%2x:%2x", &a, &b, &c, &d, &e, &f); err != nil {
		return fmt.Errorf("invalid mac address")
	}

	m.err[2] = nil
	return nil
}

func (m *Model) ipAddressValidator(value string) error {
	var a, b, c, d int

	// Check if the value is empty
	if value == "" {
		return fmt.Errorf("ip address is required")
	}

	// Check if the value is a valid IP address
	if _, err := fmt.Sscanf(value, "%d.%d.%d.%d", &a, &b, &c, &d); err != nil {
		return fmt.Errorf("invalid ip address")
	}

	m.err[3] = nil
	return nil
}

func (m *Model) validateInput(index int, validator func(string) error) bool {
	if err := validator(m.inputs[index].Value()); err != nil {
		m.err[index] = err
		return false
	}
	return true
}
