package status

import (
	"github.com/distatus/battery"
	"strings"
)

type BatteryInfo struct {
	BatteryCharging    int     `json:"batteryCharging"`    // 1: 否, 2: 是
	BatteryLevel       int     `json:"batteryLevel"`       // 电量百分比
	BatteryTemperature float64 `json:"batteryTemperature"` // 暂无可靠通用 API, 返回 -1
	BatteryCapacity    int     `json:"batteryCapacity"`    // 最大容量
}

func GetBatteryInfo() (*BatteryInfo, error) {
	batteries, err := battery.GetAll()
	if err != nil || len(batteries) == 0 {
		return nil, err
	}

	b := batteries[0]

	stateStr := strings.ToLower(b.State.String())
	charging := 1
	switch stateStr {
	case "charging":
		charging = 2
	case "full":
		charging = 3
	}

	level := int(b.Current / b.Full * 100)
	capacity := int(b.Design)
	info := &BatteryInfo{
		BatteryCharging:    charging,
		BatteryLevel:       level,
		BatteryTemperature: -1,
		BatteryCapacity:    capacity,
	}

	return info, nil
}
