<script>
import axios from "axios"
import EventBus from "../services/EventBus.js"
import Tabs from "../components/Tabs.vue"
import TabsTab from "../components/TabsTab.vue"

export default {
	name: "Home",
	components: {TabsTab, Tabs},
	data() {
		return {
			activeTab: "all",
			config: {},
			allDevices: [],
			device: {
				name: "",
				platform: "",
				description: ""
			},
			accountDevices: [],
			shareDevices: []
		}
	},
	mounted() {
		EventBus.on("refresh", this.refresh)
		this.refresh()
	},
	beforeUnmount() {
		EventBus.off("refresh", this.refresh)
	},
	computed: {
		mergedDevices() {
			const ALL_DEVICES = [...this.accountDevices, ...this.shareDevices]
			// 根据设备ID去重
			return ALL_DEVICES.reduce((acc, current) => {
				const FOUND = acc.find(item => item.id === current.id)
				if (!FOUND) {
					acc.push(current)
				}
				return acc
			}, [])
		}
	},
	methods: {
		// 刷新
		refresh() {
			this.config = JSON.parse(localStorage.getItem("config"))
			this.getDevice()
			this.getAccountDevices()
			this.getShareDevices()
		},
		// 获取设备信息
		async getDevice() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/device/info`, {
					params: {
						device_id: this.config.deviceId
					},
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.device.name = RES.data.data.device.name
				this.device.platform = RES.data.data.device.platform
				this.device.description = RES.data.data.device.description
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		},
		// 获取账户设备
		async getAccountDevices() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/devices/list`, {
					params: {
						user_id: this.config.userId
					},
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.accountDevices = RES.data.data.devices
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		},
		// 获取共享设备
		async getShareDevices() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/devices/shared`, {
					params: {
						user_id: this.config.userId
					},
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.shareDevices = RES.data.data.devices
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		},
		// 获取设备类型标签
		getDeviceType(deviceId) {
			const IS_ACCOUNT_DEVICE = this.accountDevices.some(device => device.id === deviceId)
			const IS_SHARE_DEVICE = this.shareDevices.some(device => device.id === deviceId)
			if (IS_ACCOUNT_DEVICE && IS_SHARE_DEVICE) {
				return "账户 & 共享"
			} else if (IS_ACCOUNT_DEVICE) {
				return "账户设备"
			} else if (IS_SHARE_DEVICE) {
				return "共享设备"
			}
			return "未知"
		}
	}
}
</script>

<template>
	<div class="home">
		<tabs v-model="activeTab">
			<tabs-tab name="all">
				<template #label>全部</template>
				<div v-if="(mergedDevices || []).length === 0" class="default"><span>暂无设备</span></div>
				<div v-else class="item">
					<router-link
						v-for="item in mergedDevices"
						:key="item.id"
						:to="'/device/' + item.id"
						class="container">
						<p :title="device.name">{{ device.name }}</p>
						<p :title="device.platform">{{ device.platform }}</p>
						<p :title="device.description">{{ device.description }}</p>
					</router-link>
				</div>
			</tabs-tab>
			<tabs-tab name="current">
				<template #label>当前</template>
				<div class="item">
					<router-link :to="'/device/' + this.config.deviceId" class="container">
						<p :title="device.name">{{ device.name }}</p>
						<p :title="device.platform">{{ device.platform }}</p>
						<p :title="device.description">{{ device.description }}</p>
					</router-link>
				</div>
			</tabs-tab>
			<tabs-tab name="account">
				<template #label>账户</template>
				<div v-if="(accountDevices || []).length === 0" class="default">没有设备</div>
				<div v-else class="item">
					<router-link
						v-for="item in accountDevices"
						:key="item.ID"
						:to="'/device/' + item.id"
						class="container">
						<p :title="item.name">{{ item.name }}</p>
						<p :title="item.platform">{{ item.platform }}</p>
						<p :title="item.description">{{ item.description }}</p>
					</router-link>
				</div>
			</tabs-tab>
			<tabs-tab name="share">
				<template #label>共享</template>
				<div v-if="(shareDevices || []).length === 0" class="default">没有设备</div>
				<div v-else class="item">
					<router-link
						v-for="item in shareDevices"
						:key="item.ID"
						:to="'/device/' + item.id"
						class="container">
						<p :title="item.name">{{ item.name }}</p>
						<p :title="item.platform">{{ item.platform }}</p>
						<p :title="item.description">{{ item.description }}</p>
					</router-link>
				</div>
			</tabs-tab>
		</tabs>
	</div>
</template>

<style lang="less" scoped>
.home {
	padding: 16px;
	margin: 16px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 10px;
}

.default {
	padding: 16px;
	width: 100%;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 10px;
	display: flex;
	justify-content: center;
}

.item {
	user-select: none;
	display: grid;
	grid-template-columns: repeat(auto-fill, minmax(256px, 1fr));
	gap: 16px;
}

.container {
	padding: 16px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 10px;

	&:hover {
		color: #000;
		background-color: rgba(255, 255, 255, 0.4);
	}

	p {
		text-align: center;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}
}
</style>
