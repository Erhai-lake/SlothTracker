package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Desktop/status"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) UpdateStatus(serverUrl string, deviceId string) string {
	url := fmt.Sprintf("%s/api/update/%s", serverUrl, deviceId)
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
	// 获取移动数据信息
	mobileDataInfo := status.GetMobileDataInfo()
	// 获取速度信息
	speedInfo, err := status.GetSpeedInfo()
	if err != nil {
		log.Printf("获取速度信息失败: %v", err)
		return "获取速度信息失败"
	}
	// 获取应用信息
	appInfo, err := status.GetAppInfo()
	if err != nil {
		log.Printf("获取应用信息失败: %v", err)
		return "获取应用信息失败"
	}
	// 获取其他信息
	otherInfo, err := status.GetOtherInfo()
	if err != nil {
		log.Printf("获取其他信息失败: %v", err)
		return "获取其他信息失败"
	}
	// 构造 JSON Payload
	data := map[string]any{
		"timestamp":          time.Now().Unix(),
		"batteryCharging":    batteryInfo.BatteryCharging,
		"batteryLevel":       batteryInfo.BatteryLevel,
		"batteryTemperature": batteryInfo.BatteryTemperature,
		"batteryCapacity":    batteryInfo.BatteryCapacity,
		"wifiConnected":      wifiInfo.WifiConnected,
		"wifiSSID":           wifiInfo.WifiSSID,
		"mobileDataActive":   mobileDataInfo.MobileDataActive,
		"mobileSignalDbm":    mobileDataInfo.MobileSignalDbm,
		"networkType":        mobileDataInfo.NetworkType,
		"trafficUsedMB":      mobileDataInfo.TrafficUsedMB,
		"uploadSpeedKbps":    speedInfo.UploadSpeedKbps,
		"downloadSpeedKbps":  speedInfo.DownloadSpeedKbps,
		"appName":            appInfo.AppName,
		"appTitle":           appInfo.AppTitle,
		"speakerPlaying":     appInfo.SpeakerPlaying,
		"screenOn":           otherInfo.ScreenOn,
		"isChargingViaUsb":   otherInfo.IsChargingViaUsb,
		"isChargingViaAC":    otherInfo.IsChargingViaAC,
		"isLowPowerMode":     otherInfo.IsLowPowerMode,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("JSON 编码失败: %v", err)
		return "JSON 编码失败"
	}
	// 构造请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
	message := respData["message"].(string)
	return message
}
