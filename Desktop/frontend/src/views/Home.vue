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
			activeTab: "current",
			config: {},
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
	},
	beforeUnmount() {
		EventBus.off("refresh", this.refresh)
	},
	created() {
		this.refresh()
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
				const RES = await axios.get(`${this.config.serverUrl}/api/device/${this.config.deviceId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.device.name = RES.data.device.name
				this.device.platform = RES.data.device.platform
				this.device.description = RES.data.device.description
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		},
		// 获取账户设备
		async getAccountDevices() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/devices/${this.config.userId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.accountDevices = RES.data.devices
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		},
		// 获取共享设备
		async getShareDevices() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/devices/shared/${this.config.userId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.shareDevices = RES.data.devices
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		}
	}
}
</script>

<template>
	<div class="home">
		<tabs v-model="activeTab">
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
				<div class="default" v-if="(accountDevices || []).length === 0">没有设备</div>
				<div class="item" v-else>
					<router-link
						class="container"
						:to="'/device/' + item.id"
						v-for="item in accountDevices"
						:key="item.ID">
						<p :title="item.name">{{ item.name }}</p>
						<p :title="item.platform">{{ item.platform }}</p>
						<p :title="item.description">{{ item.description }}</p>
					</router-link>
				</div>
			</tabs-tab>
			<tabs-tab name="share">
				<template #label>共享</template>
				<div class="default" v-if="(shareDevices || []).length === 0">没有设备</div>
				<div class="item" v-else>
					<router-link
						class="container"
						:to="'/device/' + item.id"
						v-for="item in shareDevices"
						:key="item.ID">
						<p :title="item.name">{{ item.name }}</p>
						<p :title="item.platform">{{ item.platform }}</p>
						<p :title="item.description">{{ item.description }}</p>
					</router-link>
				</div>
			</tabs-tab>
		</tabs>
	</div>
</template>

<style scoped lang="less">
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
