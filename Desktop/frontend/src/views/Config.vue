<script>
import axios from "axios"
import EventBus from "../services/EventBus.js"

export default {
	name: "Config",
	data() {
		return {
			process: "regular",
			config: {},
			deviceForm: {
				deviceId: "",
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
			writeOffAccountConfirm: 1,
			// form: {
			// 	deviceId: "",
			// 	deviceName: "",
			// 	description: "",
			// 	serverUrl: "",
			// 	refreshInterval: 5
			// }
		}
	},
	created() {
		this.config = JSON.parse(localStorage.getItem("config"))
		this.getDeviceInfo()
		this.getAccountInfo()
	},
	methods: {
		// 获取设备信息
		async getDeviceInfo() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/device/${this.config.deviceId}`)
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.deviceForm.deviceId = RES.data.device.id
				this.deviceForm.name = RES.data.device.name
				this.deviceForm.description = RES.data.device.description
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
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
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
						}
					})
					if (RES.data.code !== 0) {
						this.$toast.error(RES.data.message)
						return
					}
					this.$toast.success(RES.data.message)
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
				const RES = await axios.get(`${this.config.serverUrl}/api/user/${this.config.userId}`)
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.accountForm.oldName = RES.data.user.name
				this.accountForm.newName = RES.data.user.name
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
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
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
				})
				if (RES.data.code !== 0) {
					this.$toast.error(RES.data.message)
					return
				}
				this.$toast.success(RES.data.message)
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
						}
					})
					if (RES.data.code !== 0) {
						this.$toast.error(RES.data.message)
						return
					}
					this.$toast.success(RES.data.message)
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
		<div class="type">
			<p @click="process = 'regular'">常规</p>
			<p @click="process = 'device'">设备</p>
			<p @click="process = 'account'">账户</p>
		</div>
		<div class="container" v-if="process === 'regular'">regular</div>
		<div class="container" v-if="process === 'device'">
			<div class="form-item">
				<label>设备ID：</label>
				<input v-model="deviceForm.deviceId" disabled style="color: white;"/>
			</div>
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
				<button @click="writeOffDevice" style="--primary-color: #ff8080">注销设备</button>
			</div>
		</div>
		<div class="container" v-if="process === 'account'">
			<div class="form-item">
				<label>账户名：</label>
				<input v-model="accountForm.newName" placeholder="请输入修改后的用户名"/>
			</div>
			<div class="form-item">
				<label>重置密码：</label>
				<input type="password" v-model="accountForm.oldPassword" placeholder="输入旧密码"/>
				<input type="password" v-model="accountForm.newPassword" placeholder="输入新密码"/>
			</div>
			<div class="form-item-but">
				<button @click="saveAccount">保存</button>
				<button @click="writeOffAccount" style="--primary-color: #ff8080">注销账户</button>
			</div>
		</div>

		<!--		<div class="container">-->
		<!--			<div class="form-item">-->
		<!--				<label>设备ID：</label>-->
		<!--				<input v-model="form.deviceId" disabled style="color: white;"/>-->
		<!--			</div>-->
		<!--			<div class="form-item">-->
		<!--				<label>设备名称：</label>-->
		<!--				<input v-model="form.deviceName" placeholder="请输入设备名"/>-->
		<!--			</div>-->
		<!--			<div class="form-item">-->
		<!--				<label>描述：</label>-->
		<!--				<input v-model="form.description" placeholder="请输入设备描述"/>-->
		<!--			</div>-->
		<!--			<div class="form-item">-->
		<!--				<label>服务器地址：</label>-->
		<!--				<input v-model="form.serverUrl" placeholder="例如：http://localhost:8080"/>-->
		<!--			</div>-->
		<!--			<div class="form-item">-->
		<!--				<label>自动更新时间：</label>-->
		<!--				<input v-model="form.refreshInterval" placeholder="请输入自动更新时间(-1 禁用)"/>-->
		<!--			</div>-->
		<!--			<div class="form-item-but">-->
		<!--				<button @click="saveConfig">保存</button>-->
		<!--				<button @click="writeOff" style="&#45;&#45;primary-color: #ff8080">注销</button>-->
		<!--			</div>-->
		<!--		</div>-->
	</div>
</template>

<style scoped lang="less">
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
	margin-bottom: 12px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	border-radius: 10px;
	display: flex;
	gap: 15px;

	p {
		cursor: pointer;
	}
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