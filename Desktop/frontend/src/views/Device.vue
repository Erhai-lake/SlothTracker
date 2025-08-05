<script>
import {useRoute} from "vue-router"
import axios from "axios";
import EventBus from "../services/EventBus.js";

export default {
	name: "Device",
	data() {
		return {
			route: useRoute(),
			status: {
				id: null,
				device_id: null,
				timestamp: [],
				battery: {
					charging: [],
					level: null,
					temperature: null,
					capacity: null
				},
				network: {
					wifiConnected: [],
					wifiSSID: null,
					mobileDataActive: [],
					mobileSignalDbm: null,
					networkType: null,
					trafficUsedMb: null,
					uploadSpeedKbps: null,
					downloadSpeedKbps: null
				},
				foreground: {
					appName: null,
					appTitle: null,
					speakerPlaying: [],
				},
				other: {
					screenOn: [],
					isChargingViaUSB: [],
					isChargingViaAC: [],
					isLowPowerMode: []
				}
			}
		}
	},
	mounted() {
		EventBus.on("refresh", this.getStatus)
		this.getStatus()
	},
	beforeUnmount() {
		EventBus.off("refresh", this.getStatus)
	},
	methods: {
		async getStatus() {
			const CONFIG = JSON.parse(localStorage.getItem("config"))
			try {
				const ORIGINAL = this.$refs.original
				const RES = await axios.get(`${CONFIG.serverUrl}/api/status/${CONFIG.userId}/${this.route.params.id}`, {
					validateStatus: () => {
						return true
					}
				})
				// 原始数据
				ORIGINAL.innerHTML = JSON.stringify(RES.data, null, 2)
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.status = {
					id: RES.data.status.id,
					device_id: RES.data.status.device_id,
					timestamp: [
						this.formatTimestamp(RES.data.status.timestamp),
						RES.data.status.timestamp + 5 * 60 * 1000 < Date.now() ? "no" : "yes",
						this.formatTime(RES.data.status.timestamp)
					],
					battery: {
						charging: [
							RES.data.status.battery.charging === 1 ? "充电中" : RES.data.status.battery.charging === 2 ? "未充电" : "已充满",
							RES.data.status.battery.charging === 1 ? "yes" : RES.data.status.battery.charging === 2 ? "no" : ""
						],
						level: RES.data.status.battery.level + " %",
						temperature: RES.data.status.battery.temperature + " ℃",
						capacity: this.formatCapacity(RES.data.status.battery.capacity)
					},
					network: {
						wifiConnected: [
							RES.data.status.network.wifi_connected === 1 ? "已连接 ✅" : RES.data.status.network.wifi_connected === 2 ? "未连接 ❌" : "未知",
							RES.data.status.network.wifi_connected === 1 ? "yes" : RES.data.status.network.wifi_connected === 2 ? "no" : ""
						],
						wifiSSID: RES.data.status.network.wifi_ssid,
						mobileDataActive: [
							RES.data.status.network.mobile_data_active === 1 ? "已激活 ✅" : RES.data.status.network.mobile_data_active === 2 ? "未激活 ❌" : "未知",
							RES.data.status.network.mobile_data_active === 1 ? "yes" : RES.data.status.network.mobile_data_active === 2 ? "no" : ""
						],
						mobileSignalDbm: RES.data.status.network.mobile_signal_dbm + " dBm",
						networkType: RES.data.status.network.network_type,
						trafficUsedMb: this.formatTraffic(RES.data.status.network.traffic_used_mb),
						uploadSpeedKbps: this.formatSpeed(RES.data.status.network.upload_speed_kbps),
						downloadSpeedKbps: this.formatSpeed(RES.data.status.network.download_speed_kbps)
					},
					foreground: {
						appName: RES.data.status.foreground.app_name,
						appTitle: RES.data.status.foreground.app_title,
						speakerPlaying: [
							RES.data.status.foreground.speaker_playing === 1 ? "正在播放 ✅" : RES.data.status.foreground.speaker_playing === 2 ? "未播放 ❌" : "未知",
							RES.data.status.foreground.speaker_playing === 1 ? "yes" : RES.data.status.foreground.speaker_playing === 2 ? "no" : ""
						]
					},
					other: {
						screenOn: [
							RES.data.status.other.screen_on === 1 ? "屏幕已打开 ✅" : RES.data.status.other.screen_on === 2 ? "屏幕已关闭 ❌" : "未知",
							RES.data.status.other.screen_on === 1 ? "yes" : RES.data.status.other.screen_on === 2 ? "no" : ""
						],
						isChargingViaUSB: [
							RES.data.status.other.is_charging_via_usb === 1 ? "通过 USB 充电 ✅" : RES.data.status.other.is_charging_via_usb === 2 ? "未通过 USB 充电 ❌" : "未知",
							RES.data.status.other.is_charging_via_usb === 1 ? "yes" : RES.data.status.other.is_charging_via_usb === 2 ? "no" : ""
						],
						isChargingViaAC: [
							RES.data.status.other.is_charging_via_ac === 1 ? "通过 AC 充电 ✅" : RES.data.status.other.is_charging_via_ac === 2 ? "未通过 AC 充电 ❌" : "未知",
							RES.data.status.other.is_charging_via_ac === 1 ? "yes" : RES.data.status.other.is_charging_via_ac === 2 ? "no" : ""
						],
						isLowPowerMode: [
							RES.data.status.other.is_low_power_mode === 1 ? "低功耗模式 ✅" : RES.data.status.other.is_low_power_mode === 2 ? "非低功耗模式 ❌" : "未知",
							RES.data.status.other.is_low_power_mode === 1 ? "yes" : RES.data.status.other.is_low_power_mode === 2 ? "no" : ""
						]
					}
				}
			} catch (error) {
				console.error(error)
				this.$toast.error("获取状态失败")
			}
		},
		// 格式化时间戳
		formatTimestamp(timestamp) {
			const DATE = new Date(timestamp)
			const YEAR = DATE.getFullYear()
			const MONTH = String(DATE.getMonth() + 1).padStart(2, "0")
			const DAY = String(DATE.getDate()).padStart(2, "0")
			const HOURS = String(DATE.getHours()).padStart(2, "0")
			const MINUTES = String(DATE.getMinutes()).padStart(2, "0")
			const SECONDS = String(DATE.getSeconds()).padStart(2, "0")
			const MILLISECONDS = String(DATE.getMilliseconds()).padStart(3, "0")
			return `${YEAR}-${MONTH}-${DAY} ${HOURS}:${MINUTES}:${SECONDS}:${MILLISECONDS}`
		},
		// 格式化时间(超过60秒输出1分钟, 超过60分钟输出一小时以此类推, 另外输出1分钟40秒)
		formatTime(time) {
			const DATE = new Date(time)
			const NOW = new Date()
			const DIFF = NOW - DATE
			// 计算各个时间单位
			const MILLISECONDS = DIFF % 1000
			const TOTAL_SECONDS = Math.floor(DIFF / 1000)
			const SECONDS = TOTAL_SECONDS % 60
			const TOTAL_MINUTES = Math.floor(TOTAL_SECONDS / 60)
			const MINUTES = TOTAL_MINUTES % 60
			const TOTAL_HOURS = Math.floor(TOTAL_MINUTES / 60)
			const HOURS = TOTAL_HOURS % 24
			const TOTAL_DAYS = Math.floor(TOTAL_HOURS / 24)
			// 计算月和年
			const MONTHS = Math.floor(TOTAL_DAYS / 30) % 12
			const YEARS = Math.floor(TOTAL_DAYS / 365)
			const DAYS = TOTAL_DAYS % 30
			// 构建时间字符串
			const PARTS = []
			if (YEARS > 0) PARTS.push(`${YEARS} 年`)
			if (MONTHS > 0) PARTS.push(`${MONTHS} 个月`)
			if (DAYS > 0) PARTS.push(`${DAYS} 天`)
			if (HOURS > 0) PARTS.push(`${HOURS} 小时`)
			if (MINUTES > 0) PARTS.push(`${MINUTES} 分钟`)
			if (SECONDS > 0) PARTS.push(`${SECONDS} 秒`)
			if (MILLISECONDS > 0) PARTS.push(`${MILLISECONDS} 毫秒`)
			// 如果没有时间差
			if (PARTS.length === 0) return "刚刚"
			return PARTS.join(" ") + "前"
		},
		/**
		 * 格式化电池容量
		 * @param capacity 电池容量值, 单位: mAh
		 * @returns {string} 格式化后的电池容量值, 单位: mAh 或 Ah
		 */
		formatCapacity(capacity) {
			if (capacity >= 1000) {
				return (capacity / 1000).toFixed(2) + " Ah"
			}
			return capacity.toFixed(0) + " mAh"
		},
		/**
		 * 格式化流量
		 * @param valueMB 流量值, 单位: MB
		 * @returns {string} 格式化后的流量值, 单位: MB 或 GB
		 */
		formatTraffic(valueMB) {
			if (valueMB >= 1024) {
				return (valueMB / 1024).toFixed(2) + " GB"
			}
			return valueMB + " MB"
		},
		/**
		 * 格式化速度
		 * @param valueKbps 速度值, 单位: Kbps
		 * @returns {string} 格式化后的速度值, 单位: Kbps 或 Mbps 或 Gbps
		 */
		formatSpeed(valueKbps) {
			if (valueKbps >= 1000000) {
				return (valueKbps / 1000000).toFixed(2) + " Gbps"
			} else if (valueKbps >= 1000) {
				return (valueKbps / 1000).toFixed(2) + " Mbps"
			}
			return valueKbps + " Kbps"
		}
	}
}
</script>

<template>
	<div class="device">
		<h1>{{ this.status.timestamp[0] }}</h1>
		<div class="status-grid">
			<div class="status-card">
				<div class="title">📱 设备状态</div>
				<p :title="status.device_id">
					<span class="item-title">设备ID: </span>
					<span>{{ status.device_id }}</span>
				</p>
				<p :title="status.timestamp[2]">
					<span class="item-title">时效性: </span>
					<span :class="status.timestamp[1]">{{ status.timestamp[2] }}</span>
				</p>
			</div>
			<div class="status-card">
				<div class="title">🔋 电池状态</div>
				<p :title="status.battery.charging[0]">
					<span class="item-title">充电状态: </span>
					<span :class="status.battery.charging[1]">{{ status.battery.charging[0] }}</span>
				</p>
				<p :title="status.battery.level">
					<span class="item-title">电量: </span>
					<span>{{ status.battery.level }}</span>
				</p>
				<p :title="status.battery.temperature">
					<span class="item-title">温度: </span>
					<span>{{ status.battery.temperature }}</span>
				</p>
				<p :title="status.battery.capacity">
					<span class="item-title">电池总容量: </span>
					<span>{{ status.battery.capacity }}</span>
				</p>
			</div>
			<div class="status-card">
				<div class="title">📶 网络状态</div>
				<p :title="status.network.wifiConnected[0]">
					<span class="item-title">WiFi: </span>
					<span :class="status.network.wifiConnected[1]">{{ status.network.wifiConnected[0] }}</span>
				</p>
				<p :title="status.network.wifiSSID">
					<span class="item-title">WiFi名称: </span>
					<span>{{ status.network.wifiSSID }}</span>
				</p>
			</div>
			<div class="status-card">
				<div class="title">📶 流量状态</div>
				<p :title="status.network.mobileDataActive[0]">
					<span class="item-title">是否启用流量: </span>
					<span :class="status.network.mobileDataActive[1]">{{ status.network.mobileDataActive[0] }}</span>
				</p>
				<p :title="status.network.mobileSignalDbm">
					<span class="item-title">移动网络信号强度: </span>
					<span>{{ status.network.mobileSignalDbm }}</span>
				</p>
				<p :title="status.networkType">
					<span class="item-title">网络类型: </span>
					<span>{{ status.network.networkType }}</span>
				</p>
				<p :title="status.network.trafficUsedMb">
					<span class="item-title">今日流量使用: </span>
					<span>{{ status.network.trafficUsedMb }}</span>
				</p>
			</div>
			<div class="status-card">
				<div class="title">💨 网速状态</div>
				<p :title="status.network.uploadSpeedKbps">
					<span class="item-title">上传速度: </span>
					<span>{{ status.network.uploadSpeedKbps }}</span>
				</p>
				<p :title="status.network.downloadSpeedKbps">
					<span class="item-title">下载速度: </span>
					<span>{{ status.network.downloadSpeedKbps }}</span>
				</p>
			</div>
			<div class="status-card">
				<div class="title">🖥️ 前台应用状态</div>
				<p :title="status.foreground.appName">
					<span class="item-title">前台包名: </span>
					<span>{{ status.foreground.appName }}</span>
				</p>
				<p :title="status.foreground.appTitle">
					<span class="item-title">窗口标题: </span>
					<span>{{ status.foreground.appTitle }}</span>
				</p>
				<p :title="status.foreground.speakerPlaying[0]">
					<span class="item-title">扬声器播放: </span>
					<span :class="status.foreground.speakerPlaying[1]">{{ status.foreground.speakerPlaying[0] }}</span>
				</p>
			</div>
			<div class="status-card">
				<p class="title">⚙️ 其他状态</p>
				<p :title="status.other.screenOn[0]">
					<span class="item-title">屏幕点亮: </span>
					<span :class="status.other.screenOn[1]">{{ status.other.screenOn[0] }}</span>
				</p>
				<p :title="status.other.isChargingViaUSB[0]">
					<span class="item-title">USB充电: </span>
					<span :class="status.other.isChargingViaUSB[1]">{{ status.other.isChargingViaUSB[0] }}</span>
				</p>
				<p :title="status.other.isChargingViaAC[0]">
					<span class="item-title">AC充电: </span>
					<span :class="status.other.isChargingViaAC[1]">{{ status.other.isChargingViaAC[0] }}</span>
				</p>
				<p :title="status.other.isLowPowerMode[0]">
					<span class="item-title">省电模式: </span>
					<span :class="status.other.isLowPowerMode[1]">{{ status.other.isLowPowerMode[0] }}</span>
				</p>
			</div>
		</div>
		<div class="status-card code">
			<p class="title">🈚 原始响应</p>
			<pre><code ref="original"></code></pre>
		</div>
	</div>
</template>

<style scoped>
.device {
	padding: 16px;
}

h1 {
	margin: 0 16px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 12px;
	font-family: "JetBrains Mono", monospace;
	letter-spacing: 1px;
	color: #80ceff;
	text-align: center;
	user-select: none;
}

.status-grid {
	padding: 16px;
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
	gap: 16px;
}

.status-card {
	padding: 16px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 12px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.status-card p {
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}

.status-card .item-title {
	user-select: none;
}

.status-card .title {
	font-weight: bold;
	font-size: 18px;
	margin-bottom: 8px;
	user-select: none;
}

.status-card span.yes {
	color: #3dd702;
	font-weight: bold;
}

.status-card span.no {
	color: #ec410e;
	font-weight: bold;
}

.code {
	margin: 0 16px;
}
</style>