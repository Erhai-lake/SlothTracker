package status

import (
	"time"
	"github.com/shirou/gopsutil/v3/net"
)

type SpeedInfo struct {
	UploadSpeedKbps   int `json:"uploadSpeedKbps"`   // 上传速度 (Kbps)
	DownloadSpeedKbps int `json:"downloadSpeedKbps"` // 下载速度 (Kbps)
}

func GetSpeedInfo() (*SpeedInfo, error) {
	// 获取初始网络IO计数
	ioCountersStart, err := net.IOCounters(false)
	if err != nil || len(ioCountersStart) == 0 {
		return nil, err
	}
	time.Sleep(1 * time.Second) // 采样时间间隔
	// 获取结束时的网络IO计数
	ioCountersEnd, err := net.IOCounters(false)
	if err != nil || len(ioCountersEnd) == 0 {
		return nil, err
	}
	start := ioCountersStart[0]
	end := ioCountersEnd[0]
	// 计算字节差值
	bytesSent := end.BytesSent - start.BytesSent
	bytesRecv := end.BytesRecv - start.BytesRecv
	// 转换为 Kbps (Bytes/sec → bits/sec → kilobits/sec)
	uploadKbps := int((bytesSent * 8) / 1024)
	downloadKbps := int((bytesRecv * 8) / 1024)
	return &SpeedInfo{
		UploadSpeedKbps:   uploadKbps,
		DownloadSpeedKbps: downloadKbps,
	}, nil
}
