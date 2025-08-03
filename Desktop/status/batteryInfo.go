package status

import (
	"github.com/distatus/battery"
	"strings"
)

type BatteryStatus struct {
	Charging    int     `json:"charging"`    // 是否充电中(1: 充电中, 2: 未充电, 3: 已充满)
	Level       int     `json:"level"`       // 电池电量百分比(0~100)
	Temperature float64 `json:"temperature"` // 电池温度(摄氏度)
	Capacity    int     `json:"capacity"`    // 电池设计容量或总容量(单位mAh, 可选)
}

func GetBatteryInfo() (*BatteryStatus, error) {
	batteries, err := battery.GetAll()
	if err != nil || len(batteries) == 0 {
		return nil, err
	}

	b := batteries[0]

	stateStr := strings.ToLower(b.State.String())
	charging := 2
	switch stateStr {
	case "charging":
		charging = 1
	case "full":
		charging = 3
	}

	level := int(b.Current / b.Full * 100)
	capacity := int(b.Design)
	info := &BatteryStatus{
		Charging:    charging,
		Level:       level,
		Temperature: -1,
		Capacity:    capacity,
	}

	return info, nil
}
