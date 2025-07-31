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
				timestamp: null,
				batteryCharging: null,
				batteryLevel: null,
				batteryTemperature: null,
				batteryCapacity: null,
				wifiConnected: null,
				wifiSSID: null,
				mobileDataActive: null,
				mobileSignalDbm: null,
				networkType: null,
				uploadSpeedKbps: null,
				downloadSpeedKbps: null,
				trafficUsedMB: null,
				appName: null,
				appTitle: null,
				speakerPlaying: null,
				screenOn: null,
				isChargingViaUSB: null,
				isChargingViaAC: null,
				isLowPowerMode: null,
				isAirplaneMode: null
			}
		}
	},
	mounted() {
		EventBus.on("refresh", this.getStatus)
	},
	beforeUnmount() {
		EventBus.off("refresh", this.getStatus)
	},
	created() {
		this.getStatus()
	},
	methods: {
		async getStatus() {
			const CONFIG = JSON.parse(localStorage.getItem("config"))
			try {
				const RES = await axios.get(`${CONFIG.serverUrl}/api/status/${this.route.params.id}`)
				this.status = RES.data.status
			} catch (error) {
				console.error(error)
			}
		},
		formatTimestamp(timestamp) {
			const DATE = new Date(timestamp * 1000)
			const YEAR = DATE.getFullYear()
			const MONTH = String(DATE.getMonth() + 1).padStart(2, "0")
			const DAY = String(DATE.getDate()).padStart(2, "0")
			const HOURS = String(DATE.getHours()).padStart(2, "0")
			const MINUTES = String(DATE.getMinutes()).padStart(2, "0")
			const SECONDS = String(DATE.getSeconds()).padStart(2, "0")
			return `${YEAR}-${MONTH}-${DAY} ${HOURS}:${MINUTES}:${SECONDS}`
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
		<h1>{{ formatTimestamp(this.status.timestamp) }}</h1>
		<div class="status-grid">
			<div class="status-card">
				<div class="title">🔋 电池状态</div>
				<p>
					<span class="item-title">充电状态: </span>
					<span :class="{ no: status.batteryCharging === 1, yes: status.batteryCharging === 2 }"
						  :title="status.batteryCharging === 1 ? '未充电' : status.batteryCharging === 2 ? '充电中' : '已充满'">
						{{
							status.batteryCharging === 1 ? "未充电" : status.batteryCharging === 2 ? "充电中" : "已充满"
						}}
					</span>
				</p>
				<p :title="status.batteryLevel">
					<span class="item-title">电量: </span>
					{{ status.batteryLevel }}%
				</p>
				<p :title="status.batteryTemperature">
					<span class="item-title">温度: </span>
					{{ status.batteryTemperature }} ℃
				</p>
				<p :title="formatCapacity(status.batteryCapacity)">
					<span class="item-title">电池总容量: </span>
					{{ formatCapacity(status.batteryCapacity) }}
				</p>
			</div>
			<div class="status-card">
				<div class="title">📶 网络状态</div>
				<p>
					<span class="item-title">WiFi: </span>
					<span
						:class="{ no: status.wifiConnected === 1, yes: status.wifiConnected === 2 }"
						:title="status.wifiConnected === 1 ? '未连接' : '已连接'">
						{{ status.wifiConnected === 1 ? "未连接 ❌" : "已连接 ✅" }}
					</span>
				</p>
				<p :title="status.wifiSSID">
					<span class="item-title">WiFi名称: </span>
					{{ status.wifiSSID }}
				</p>
			</div>
			<div class="status-card">
				<div class="title">📶 流量状态</div>
				<p>
					<span class="item-title">是否启用流量: </span>
					<span
						:class="{ no: status.mobileDataActive === 1, yes: status.mobileDataActive === 2 }"
						:title="status.mobileDataActive === 1 ? '否' : '是'">
						{{ status.mobileDataActive === 1 ? "否" : "是" }}
					</span>
				</p>
				<p :title="status.mobileSignalDbm">
					<span class="item-title">移动网络信号强度: </span>
					{{ status.mobileSignalDbm }} dBm
				</p>
				<p :title="status.networkType">
					<span class="item-title">网络类型: </span>
					{{ status.networkType }}
				</p>
				<p :title="formatTraffic(status.trafficUsedMB)">
					<span class="item-title">今日流量使用: </span>
					{{ formatTraffic(status.trafficUsedMB) }}
				</p>
			</div>
			<div class="status-card">
				<div class="title">💨 网速状态</div>
				<p :title="formatSpeed(status.uploadSpeedKbps)">
					<span class="item-title">上传速度: </span>
					{{ formatSpeed(status.uploadSpeedKbps) }}
				</p>
				<p :title="formatSpeed(status.downloadSpeedKbps)">
					<span class="item-title">下载速度: </span>
					{{ formatSpeed(status.downloadSpeedKbps) }}
				</p>
			</div>
			<div class="status-card">
				<div class="title">🖥️ 前台应用状态</div>
				<p :title="status.appName">
					<span class="item-title">前台包名: </span>
					{{ status.appName }}
				</p>
				<p :title="status.appTitle">
					<span class="item-title">窗口标题: </span>
					{{ status.appTitle }}
				</p>
				<p :title="status.speakerPlaying === 1 ? '否' : status.speakerPlaying === 2 ? '是' : '未知'">
					<span class="item-title">扬声器播放: </span>
					<span :class="{ no: status.speakerPlaying === 1 || -1, yes: status.speakerPlaying === 2 }">
						{{ status.speakerPlaying === 1 ? "否" : status.speakerPlaying === 2 ? "是" : "未知" }}
					</span>
				</p>
			</div>
			<div class="status-card">
				<p class="title">⚙️ 其他状态</p>
				<p>
					<span class="item-title">屏幕点亮: </span>
					<span :class="{ no: status.screenOn === 1, yes: status.screenOn === 2 }">
						{{ status.screenOn === 1 ? "否" : "是" }}
					</span>
				</p>
				<p>
					<span class="item-title">USB充电: </span>
					<span :class="{ no: status.isChargingViaUSB === 1, yes: status.isChargingViaUSB === 2 }">
						{{ status.isChargingViaUSB === 1 ? "否" : "是" }}
					</span>
				</p>
				<p>
					<span class="item-title">AC充电: </span>
					<span :class="{ no: status.isChargingViaAC === 1, yes: status.isChargingViaAC === 2 }">
						{{ status.isChargingViaAC === 1 ? "否" : "是" }}
					</span>
				</p>
			</div>
		</div>
		<div class="status-card code">
			<p class="title">🈚 原始响应</p>
			<pre><code>{{ status }}</code></pre>
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
	grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
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
	color: #b02c04;
	font-weight: bold;
}

.code {
	margin: 0 16px;
}
</style>