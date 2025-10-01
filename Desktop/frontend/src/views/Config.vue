<script>
import axios from "axios"
import EventBus from "../services/EventBus.js"
import Tabs from "../components/Tabs.vue"
import TabsTab from "../components/TabsTab.vue"

export default {
	name: "Config",
	components: {TabsTab, Tabs},
	data() {
		return {
			activeTab: "regular",
			config: {},
			regularForm: {
				serverUrl: "",
				refreshInterval: 5,
				background: true
			},
			deviceForm: {
				name: "",
				description: ""
			},
			writeOffDeviceConfirm: 1,
			accountForm: {
				oldName: "",
				newName: "",
				oldPassword: "",
				newPassword: ""
			},
			writeOffAccountConfirm: 1
		}
	},
	mounted() {
		this.config = JSON.parse(localStorage.getItem("config"))
		this.getRegularInfo()
		this.getDeviceInfo()
		this.getAccountInfo()
	},
	methods: {
		// 获取常规信息
		async getRegularInfo() {
			this.regularForm.serverUrl = this.config.serverUrl
			this.regularForm.refreshInterval = this.config.refreshInterval
			this.regularForm.background = this.config.background
		},
		// 测试服务器地址
		async ping() {
			if (this.regularForm.serverUrl === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			try {
				const RES = await axios.get(`${this.regularForm.serverUrl}/api/ping`, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error("服务器地址错误")
					return
				}
				this.$toast.success(`服务器延迟: ${RES.data.data.latency}`)
			} catch (error) {
				console.error(error)
				this.$toast.error("服务器地址错误")
			}
		},
		// 保存常规信息
		async saveRegular() {
			if (this.regularForm.serverUrl === "" || this.regularForm.refreshInterval === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			localStorage.setItem("config", JSON.stringify({
				...this.config,
				serverUrl: this.regularForm.serverUrl,
				refreshInterval: this.regularForm.refreshInterval,
				background: this.regularForm.background
			}))
			this.$toast.success("保存成功")
		},
		// 刷新页面
		refresh() {
			location.reload()
		},
		// 重置设置
		resetConfig() {
			localStorage.removeItem("config")
			location.reload()
		},
		// 获取设备信息
		async getDeviceInfo() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/device/${this.config.deviceId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.deviceForm.name = RES.data.data.device.name
				this.deviceForm.description = RES.data.data.device.description
			} catch (error) {
				console.error(error)
				this.$toast.error("获取设备信息错误")
			}
		},
		// 保存设备信息
		async saveDevice() {
			if (this.deviceForm.name === "" || this.deviceForm.description === "") {
				this.$toast.warning("请填写完整信息")
				return
			}
			try {
				const RES = await axios.put(`${this.config.serverUrl}/api/device/update`, {
					deviceId: this.config.deviceId,
					name: this.deviceForm.name,
					platform: "windows",
					description: this.deviceForm.description
				}, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.$toast.success(RES.data.data.message)
			} catch (error) {
				console.error(error)
				this.$toast.error("保存设备信息错误")
			}
		},
		// 注销设备
		async writeOffDevice() {
			if (this.writeOffDeviceConfirm === 1) {
				this.writeOffDeviceConfirm = 2
				this.$toast.warning("请再次点击确认是否继续注销, 3秒后重置注销状态")
				this.$toast.error("注销账户后将删除当前设备的设备信息且无法恢复")
				setTimeout(() => {
					this.writeOffDeviceConfirm = 1
				}, 3000)
			} else if (this.writeOffDeviceConfirm === 2) {
				this.$toast.clear()
				try {
					const RES = await axios.delete(`${this.config.serverUrl}/api/device/delete`, {
						data: {
							id: this.config.deviceId,
						},
						validateStatus: () => {
							return true
						}
					})
					if (!RES.data.success) {
						this.$toast.error(RES.data.data.message)
						return
					}
					this.$toast.success(RES.data.data.message)
				} catch (error) {
					console.error(error)
					this.$toast.error("注销账设备错误")
					return
				}
				const CONFIG = JSON.parse(localStorage.getItem("config"))
				localStorage.setItem("config", JSON.stringify({
					...CONFIG,
					deviceId: ""
				}))
				this.$router.push("/init")
				EventBus.emit("initConfig")
				this.$toast.warning("感谢使用, 再见!")
			}
		},
		// 获取账户信息
		async getAccountInfo() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/user/${this.config.userId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.accountForm.oldName = RES.data.data.user.name
				this.accountForm.newName = RES.data.data.user.name
			} catch (error) {
				console.error(error)
				this.$toast.error("获取账户信息错误")
			}
		},
		// 保存账户信息
		saveAccount() {
			if (this.accountForm.name !== "" && this.accountForm.oldName !== this.accountForm.newName) {
				this.resetAccountName()
			}
			if (this.accountForm.oldPassword !== "" && this.accountForm.newPassword !== "") {
				this.resetAccountPassword()
			}
		},
		// 重置用户名
		async resetAccountName() {
			try {
				const RES = await axios.put(`${this.config.serverUrl}/api/user/reset_name`, {
					id: this.config.userId,
					name: this.accountForm.newName
				}, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.$toast.success(RES.data.data.message)
				this.accountForm.oldName = this.accountForm.newName
			} catch (error) {
				console.error(error)
				this.$toast.error("保存账户信息错误")
			}
		},
		// 重置密码
		async resetAccountPassword() {
			try {
				const RES = await axios.put(`${this.config.serverUrl}/api/user/reset_password`, {
					id: this.config.userId,
					old_password: this.accountForm.oldPassword,
					new_password: this.accountForm.newPassword
				}, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.$toast.success(RES.data.data.message)
			} catch (error) {
				console.error(error)
				this.$toast.error("保存账户信息错误")
			}
		},
		// 注销账户
		async writeOffAccount() {
			if (this.writeOffAccountConfirm === 1) {
				this.writeOffAccountConfirm = 2
				this.$toast.warning("请再次点击确认是否继续注销, 3秒后重置注销状态")
				this.$toast.error("注销账户后将删除所有设备信息且无法恢复")
				setTimeout(() => {
					this.writeOffAccountConfirm = 1
				}, 3000)
			} else if (this.writeOffAccountConfirm === 2) {
				this.$toast.clear()
				if (this.accountForm.oldPassword === "" || this.accountForm.newPassword === "") {
					this.$toast.warning("请在密码框中输入两遍密码!")
					return
				}
				if (this.accountForm.oldPassword !== this.accountForm.newPassword) {
					this.$toast.warning("两次输入的密码不一致!")
					return
				}
				try {
					const RES = await axios.delete(`${this.config.serverUrl}/api/user/delete`, {
						data: {
							id: this.config.userId,
							password: this.accountForm.oldPassword
						},
						validateStatus: () => {
							return true
						}
					})
					if (!RES.data.success) {
						this.$toast.error(RES.data.data.message)
						return
					}
					this.$toast.success(RES.data.data.message)
				} catch (error) {
					console.error(error)
					this.$toast.error("注销账户错误")
					return
				}
				localStorage.removeItem("config")
				this.$router.push("/init")
				EventBus.emit("initConfig")
				this.$toast.warning("感谢使用, 再见!")
			}
		}
	}
}
</script>

<template>
	<div class="config">
		<tabs v-model="activeTab" class="type">
			<tabs-tab name="regular">
				<template #label>常规</template>
				<div class="container">
					<div class="form-item">
						<label>服务器地址：</label>
						<input v-model="regularForm.serverUrl"/>
					</div>
					<div class="form-item">
						<label>自动刷新时间(秒)：</label>
						<input v-model="regularForm.refreshInterval"/>
					</div>
					<div class="form-item">
						<label>
							背景图开关：
							<input v-model="regularForm.background" type="checkbox"/>
						</label>
					</div>
					<div class="form-item-but">
						<button style="--primary-color: #3ecd39" @click="ping">测试地址</button>
						<button @click="saveRegular">保存</button>
					</div>
					<div class="form-item-but">
						<button @click="refresh">刷新页面</button>
						<button style="--primary-color: #ff8080" @click="resetConfig">重置设置</button>
					</div>
				</div>
			</tabs-tab>
			<tabs-tab name="device">
				<template #label>设备</template>
				<div class="container">
					<div class="form-item">
						<label>设备名称：</label>
						<input v-model="deviceForm.name" placeholder="请输入设备名"/>
					</div>
					<div class="form-item">
						<label>描述：</label>
						<input v-model="deviceForm.description" placeholder="请输入设备描述"/>
					</div>
					<div class="form-item-but">
						<button @click="saveDevice">保存</button>
						<button style="--primary-color: #ff8080" @click="writeOffDevice">注销设备</button>
					</div>
				</div>
			</tabs-tab>
			<tabs-tab name="account">
				<template #label>账户</template>
				<div class="container">
					<div class="form-item">
						<label>账户名：</label>
						<input v-model="accountForm.newName" placeholder="请输入修改后的用户名"/>
					</div>
					<div class="form-item">
						<label>重置密码：</label>
						<input v-model="accountForm.oldPassword" placeholder="输入旧密码" type="password"/>
						<input v-model="accountForm.newPassword" placeholder="输入新密码" type="password"/>
					</div>
					<div class="form-item-but">
						<button @click="saveAccount">保存</button>
						<button style="--primary-color: #ff8080" @click="writeOffAccount">注销账户</button>
					</div>
				</div>
			</tabs-tab>
		</tabs>
	</div>
</template>

<style lang="less" scoped>
.config {
	padding: 16px;
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: column;
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