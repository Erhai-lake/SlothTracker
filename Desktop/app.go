package main

import (
	"Desktop/status"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) UpdateStatus(serverUrl string, userId string, deviceId string) any {
	url := fmt.Sprintf("%s/api/status/update/%s/%s", serverUrl, userId, deviceId)
	// 获取电池信息
	batteryInfo, err := status.GetBatteryInfo()
	if err != nil {
		log.Printf("获取电池信息失败: %v", err)
		return "获取电池信息失败"
	}
	// 获取 WiFi 信息
	wifiInfo, err := status.GetWifiInfo()
	if err != nil {
		log.Printf("获取 WiFi 信息失败: %v", err)
		return "获取 WiFi 信息失败"
	}
	// 获取应用信息
	foregroundStatus, err := status.GetForegroundStatus()
	if err != nil {
		log.Printf("获取前台状态失败: %v", err)
		return "获取前台状态失败"
	}
	// 获取其他信息
	otherInfo, err := status.GetOtherInfo()
	if err != nil {
		log.Printf("获取其他信息失败: %v", err)
		return "获取其他信息失败"
	}
	// 构造 JSON Payload
	data := map[string]any{
		"battery": map[string]any{
			"charging":    batteryInfo.Charging,
			"level":       batteryInfo.Level,
			"temperature": batteryInfo.Temperature,
			"capacity":    batteryInfo.Capacity,
		},
		"network": map[string]any{
			"wifi_connected":      wifiInfo.WifiConnected,
			"wifi_ssid":           wifiInfo.WifiSSId,
			"mobile_data_active":  wifiInfo.MobileDataActive,
			"mobile_signal_dbm":   wifiInfo.MobileSignalDbm,
			"network_type":        wifiInfo.NetworkType,
			"traffic_used_mb":     wifiInfo.TrafficUsedMB,
			"upload_speed_kbps":   wifiInfo.UploadSpeedKbps,
			"download_speed_kbps": wifiInfo.DownloadSpeedKbps,
		},
		"foreground": map[string]any{
			"app_name":        foregroundStatus.AppName,
			"app_title":       foregroundStatus.AppTitle,
			"speaker_playing": foregroundStatus.SpeakerPlaying,
		},
		"other": map[string]any{
			"screen_on":           otherInfo.ScreenOn,
			"is_charging_via_usb": otherInfo.IsChargingViaUSB,
			"is_charging_via_ac":  otherInfo.IsChargingViaAC,
			"is_low_power_mode":   otherInfo.IsLowPowerMode,
		},
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("JSON 编码失败: %v", err)
		return "JSON 编码失败"
	}
	// 构造请求
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("请求创建失败: %v", err)
		return "请求创建失败"
	}
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("请求发送失败: %v", err)
		return "请求失败"
	}
	defer resp.Body.Close()
	// 读取响应
	respBody := new(bytes.Buffer)
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		log.Printf("读取响应失败: %v", err)
		return "读取响应失败"
	}
	var respData map[string]any
	err = json.Unmarshal(respBody.Bytes(), &respData)
	if err != nil {
		log.Printf("JSON 解码失败: %v", err)
		return "JSON 解码失败"
	}
	return respData
}
