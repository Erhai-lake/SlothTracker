<script>
import EventBus from "../services/EventBus.js"
import axios from "axios"
import Tabs from "../components/Tabs.vue"
import TabsTab from "../components/TabsTab.vue"

export default {
	name: "Initialization",
	components: {TabsTab, Tabs},
	data() {
		return {
			activeTab: "serverUrl",
			serverUrl: "",
			initForm: {
				refreshInterval: 5
			},
			loginRegistrationForm: {
				name: "",
				password: ""
			},
			deviceForm: {
				deviceId: "",
				name: "",
				platform: "windows",
				description: ""
			},
			devices: []
		}
	},
	mounted() {
		this.init()
	},
	methods: {
		// 初始化
		init() {
			EventBus.emit("sidebarOpen", false)
			const CONFIG = JSON.parse(localStorage.getItem("config"))
			const REQUIRED_FIELDS = [["serverUrl", "serverUrl"], ["refreshInterval", "init"], ["userId", "loginRegistration"], ["deviceId", "selectDevice"]]
			if (!CONFIG) {
				this.$router.push("/init")
				return
			}
			for (const FIELD of REQUIRED_FIELDS) {
				if (!CONFIG[FIELD[0]]) {
					this.$router.push("/init")
					return
				}
			}
			this.$router.push("/")
			EventBus.emit("initConfig")
			EventBus.emit("sidebarOpen", true)
		},
		// 完成
		async complete() {
			this.init()
		},
		// 测试服务器地址
		async ping() {
			if (this.serverUrl === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			try {
				const RES = await axios.get(`${this.serverUrl}/api/ping`, {
					validateStatus: () => {
						return true
					}
				})
				if (RES.data.code !== 0) {
					this.$toast.error("服务器地址错误")
					return
				}
				this.$toast.success(`服务器延迟: ${RES.data.message}`)
			} catch (error) {
				console.error(error)
				this.$toast.error("服务器地址错误")
				return
			}
			let config = JSON.parse(localStorage.getItem("config"))
			config = {
				...config,
				serverUrl: this.serverUrl
			}
			localStorage.setItem("config", JSON.stringify(config))
			this.activeTab = "init"
		},
		// 初始化设置
		async saveInit() {
			if (this.initForm.refreshInterval === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			let config = JSON.parse(localStorage.getItem("config"))
			config = {
				...config,
				refreshInterval: this.initForm.refreshInterval,
				background: true
			}
			localStorage.setItem("config", JSON.stringify(config))
			this.activeTab = "loginRegistration"
		},
		// 登录注册
		async loginRegistration(type) {
			if (this.loginRegistrationForm.name === "" || this.loginRegistrationForm.password === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			let url = ""
			if (type === 1) {
				url = `${this.serverUrl}/api/user/login`
			} else {
				url = `${this.serverUrl}/api/user/register`
			}
			let userId = null
			try {
				const RES = await axios.post(url,
					{
						name: this.loginRegistrationForm.name,
						password: this.loginRegistrationForm.password
					}, {
						validateStatus: () => {
							return true
						}
					})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
				userId = RES.data.user_id
			} catch (error) {
				console.error(error)
				this.$toast.error("登录注册错误")
				return
			}
			let config = JSON.parse(localStorage.getItem("config"))
			config = {
				...config,
				userId: userId
			}
			localStorage.setItem("config", JSON.stringify(config))
			this.activeTab = "selectDevice"
			await this.getDevices()
		},
		// 获取设备列表
		async getDevices() {
			try {
				const RES = await axios.get(`${this.serverUrl}/api/devices/${JSON.parse(localStorage.getItem("config")).userId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
				this.devices = RES.data.devices
			} catch (error) {
				console.error(error)
			}
		},
		// 保存设备ID
		async saveDeviceId(type = 1) {
			if (this.deviceForm.deviceId === "") {
				this.$toast.warning("请选择设备")
				return
			}
			if (type === 2) {
				try {
					const RES = await axios.put(`${this.serverUrl}/api/device/update`, {
						deviceId: this.deviceForm.deviceId,
						name: this.deviceForm.name,
						platform: this.deviceForm.platform,
						description: this.deviceForm.description
					}, {
						validateStatus: () => {
							return true
						}
					})
					if (RES.data.code !== 0) {
						this.$toast.error(RES.data.message)
						return
					}
				} catch (error) {
					console.error(error)
					this.$toast.error("保存设备平台错误")
					return
				}
			}
			let config = JSON.parse(localStorage.getItem("config"))
			config = {
				...config,
				deviceId: this.deviceForm.deviceId
			}
			localStorage.setItem("config", JSON.stringify(config))
			await this.complete()
		},
		// 注册设备
		async registrationDevice() {
			if (this.deviceForm.name === "" || this.deviceForm.description === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			try {
				const RES = await axios.post(`${this.serverUrl}/api/device/register`,
					{
						ownerID: JSON.parse(localStorage.getItem("config")).userId,
						deviceName: this.deviceForm.name,
						platform: this.deviceForm.platform,
						description: this.deviceForm.description
					}, {
						validateStatus: () => {
							return true
						}
					})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
				this.deviceForm.deviceId = RES.data.device_id
				await this.saveDeviceId(2)
			} catch (error) {
				console.error(error)
				this.$toast.error("注册设备错误")
			}
		}
	}
}
</script>

<template>
	<div class="initialization">
		<tabs v-model="activeTab" class="type">
			<tabs-tab name="serverUrl">
				<template #label>服务器地址</template>
				<div class="container">
					<h2>测试服务器地址</h2>
					<div class="form-item">
						<label>
							服务器地址：
							<input v-model="serverUrl" placeholder="例如：http://localhost:8080"/>
						</label>
					</div>
					<div class="form-item-but">
						<button @click="ping">测试并保存</button>
					</div>
				</div>
			</tabs-tab>
			<tabs-tab name="init">
				<template #label>初始化设置</template>
				<div class="container">
					<h2>初始化设置</h2>
					<div class="form-item">
						<label>
							自动刷新时间(秒)：
							<input v-model="initForm.refreshInterval" placeholder="-1为禁用自动刷新"/>
						</label>
					</div>
					<div class="form-item-but">
						<button @click="saveInit">保存</button>
					</div>
				</div>
			</tabs-tab>
			<tabs-tab name="loginRegistration">
				<template #label>登录注册</template>
				<div class="container">
					<h2>登录注册</h2>
					<div class="form-item">
						<label>
							用户名：
							<input v-model="loginRegistrationForm.name" placeholder="请输入用户名"/>
						</label>
					</div>
					<div class="form-item">
						<label>
							密码：
							<input v-model="loginRegistrationForm.password" placeholder="请输入登录密码"/>
						</label>
					</div>
					<div class="form-item-but">
						<button @click="loginRegistration(1)" style="--primary-color: #3ecd39">登录</button>
						<button @click="loginRegistration(2)">注册</button>
					</div>
				</div>
			</tabs-tab>
			<tabs-tab name="selectDevice">
				<template #label>选择已有设备</template>
				<div class="container">
					<h2>选择已有设备</h2>
					<div class="form-item">
						<select v-model="deviceForm.deviceId">
							<option value="" disabled>请选择设备</option>
							<option v-for="device in devices" :key="device.id" :value="device.id">
								{{ device.name }}
							</option>
						</select>
					</div>
					<div class="form-item-but">
						<button @click="saveDeviceId" style="--primary-color: #3ecd39">确定</button>
						<button @click="activeTab = 'registrationDevice'">注册新的</button>
					</div>
				</div>
			</tabs-tab>
			<tabs-tab name="registrationDevice">
				<template #label>注册新设备</template>
				<div class="container">
					<h2>注册新设备</h2>
					<div class="form-item">
						<label>
							设备名称：
							<input v-model="deviceForm.name" placeholder="请输入设备名"/>
						</label>
					</div>
					<div class="form-item">
						<label>
							设备描述：
							<input v-model="deviceForm.description" placeholder="请输入设备描述"/>
						</label>
					</div>
					<div class="form-item-but">
						<button @click="registrationDevice" style="--primary-color: #3ecd39">注册</button>
						<button @click="activeTab = 'selectDevice'">选择已有设备</button>
					</div>
				</div>
			</tabs-tab>
		</tabs>
	</div>
</template>

<style scoped lang="less">
.initialization {
	padding: 16px;
	width: 100%;
	height: 100%;
	display: flex;
	justify-content: center;
	align-items: center;
}

.type {
	padding: 10px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 10px;
}

.container {
	padding: 26px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 10px;
	display: flex;
	flex-direction: column;
}

h2 {
	text-align: center;
	margin-bottom: 16px;
}

.form-item {
	margin-bottom: 12px;
	display: flex;
	flex-direction: column;

	input {
		padding: 8px;
		margin-top: 5px;
		width: 300px;
		border: 1px solid #ccc;
		border-radius: 4px;

		&[type="checkbox"] {
			width: auto;
		}
	}
}

.form-item-but {
	margin-bottom: 10px;
	display: flex;
	gap: 10px;

	button {
		padding: 10px;
		width: 100%;
		background-color: var(--primary-color, #80ceff);
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		white-space: nowrap;

		&:hover {
			background-color: #66b1ff;
		}
	}
}
</style>
