<script>
import EventBus from "../services/EventBus.js"
import axios from "axios"

export default {
	name: "Initialization",
	data() {
		return {
			process: "serverUrl",
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
	created() {
		const CONFIG = JSON.parse(localStorage.getItem("config"))
		if (CONFIG) {
			this.$router.push("/")
		}
	},
	methods: {
		// 测试服务器地址
		async ping() {
			if (this.serverUrl === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			try {
				const RES = await axios.get(`${this.serverUrl}/api/ping`)
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
			this.process = 'init'
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
				refreshInterval: this.initForm.refreshInterval
			}
			localStorage.setItem("config", JSON.stringify(config))
			this.process = 'loginRegistration'
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
			this.process = 'selectDevice'
			await this.getDevices()
		},
		// 获取设备列表
		async getDevices() {
			try {
				const RES = await axios.get(`${this.serverUrl}/api/devices/${JSON.parse(localStorage.getItem("config")).userId}`)
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
		async saveDeviceId() {
			if (this.deviceForm.deviceId === "") {
				this.$toast.warning("请选择设备")
				return
			}
			let config = JSON.parse(localStorage.getItem("config"))
			config = {
				...config,
				deviceId: this.deviceForm.deviceId
			}
			localStorage.setItem("config", JSON.stringify(config))
			this.$router.push("/")
			EventBus.emit("initConfig")
			EventBus.emit("sidebarOpen", true)
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
					})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
				this.deviceForm.deviceId = RES.data.device_id
				await this.saveDeviceId()
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
		<div class="container" v-if="process === 'serverUrl'">
			<h2>测试服务器地址</h2>
			<div class="form-item">
				<label>服务器地址：</label>
				<input v-model="serverUrl" placeholder="例如：http://localhost:8080"/>
			</div>
			<div class="form-item-but">
				<button @click="ping">测试并保存</button>
			</div>
		</div>
		<div class="container" v-if="process === 'init'">
			<h2>初始化设置</h2>
			<div class="form-item">
				<label>自动刷新时间(秒)：</label>
				<input v-model="initForm.refreshInterval" placeholder="-1为禁用自动刷新"/>
			</div>
			<div class="form-item-but">
				<button @click="saveInit">保存</button>
			</div>
		</div>
		<div class="container" v-if="process === 'loginRegistration'">
			<h2>登录注册</h2>
			<div class="form-item">
				<label>用户名：</label>
				<input v-model="loginRegistrationForm.name" placeholder="请输入用户名"/>
			</div>
			<div class="form-item">
				<label>密码：</label>
				<input v-model="loginRegistrationForm.password" placeholder="请输入登录密码"/>
			</div>
			<div class="form-item-but">
				<button @click="loginRegistration(1)" style="--primary-color: #3ecd39">登录</button>
				<button @click="loginRegistration(2)">注册</button>
			</div>
		</div>
		<div class="container" v-if="process === 'selectDevice'">
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
				<button @click="process = 'selectDevice'">注册新的</button>
			</div>
		</div>
		<div class="container" v-if="process === 'registrationDevice'">
			<h2>注册新设备</h2>
			<div class="form-item">
				<label>设备名称：</label>
				<input v-model="deviceForm.name" placeholder="请输入设备名"/>
			</div>
			<div class="form-item">
				<label>设备描述：</label>
				<input v-model="deviceForm.description" placeholder="请输入设备描述"/>
			</div>
			<div class="form-item-but">
				<button @click="registrationDevice" style="--primary-color: #3ecd39">注册</button>
				<button @click="process = 'selectDevice'">选择已有设备</button>
			</div>
		</div>
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
		width: 300px;
		border: 1px solid #ccc;
		border-radius: 4px;
	}
}

.form-item-but {
	display: flex;
	gap: 10px;

	button {
		width: 100%;
		padding: 10px;
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
