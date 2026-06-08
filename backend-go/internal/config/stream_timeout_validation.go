package config

import "fmt"

const (
	minStreamFirstContentTimeoutMs = 5000
	maxStreamFirstContentTimeoutMs = 300000
	minStreamInactivityTimeoutMs   = 1000
	maxStreamInactivityTimeoutMs   = 180000
	minStreamToolCallIdleTimeoutMs = 30000
	maxStreamToolCallIdleTimeoutMs = 300000
)

func validateStreamFirstContentTimeoutMs(value int) error {
	if value == 0 {
		return nil
	}
	if value < minStreamFirstContentTimeoutMs || value > maxStreamFirstContentTimeoutMs {
		return fmt.Errorf("流式首字等待超时必须为 0（继承全局）或 %d-%d 之间", minStreamFirstContentTimeoutMs, maxStreamFirstContentTimeoutMs)
	}
	return nil
}

func validateStreamInactivityTimeoutMs(value int) error {
	if value == 0 {
		return nil
	}
	if value < minStreamInactivityTimeoutMs || value > maxStreamInactivityTimeoutMs {
		return fmt.Errorf("流式断流超时必须为 0（继承全局）或 %d-%d 之间", minStreamInactivityTimeoutMs, maxStreamInactivityTimeoutMs)
	}
	return nil
}

func validateStreamToolCallIdleTimeoutMs(value int) error {
	if value == 0 {
		return nil
	}
	if value < minStreamToolCallIdleTimeoutMs || value > maxStreamToolCallIdleTimeoutMs {
		return fmt.Errorf("工具调用空闲超时必须为 0（继承全局）或 %d-%d 之间", minStreamToolCallIdleTimeoutMs, maxStreamToolCallIdleTimeoutMs)
	}
	return nil
}

func validateStreamTimeouts(firstContentTimeoutMs int, inactivityTimeoutMs int, toolCallIdleTimeoutMs int) error {
	if err := validateStreamFirstContentTimeoutMs(firstContentTimeoutMs); err != nil {
		return err
	}
	if err := validateStreamInactivityTimeoutMs(inactivityTimeoutMs); err != nil {
		return err
	}
	if err := validateStreamToolCallIdleTimeoutMs(toolCallIdleTimeoutMs); err != nil {
		return err
	}
	return nil
}
