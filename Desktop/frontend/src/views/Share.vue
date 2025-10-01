<script>
import axios from "axios"
import EventBus from "../services/EventBus.js"
import Tabs from "../components/Tabs.vue"
import TabsTab from "../components/TabsTab.vue"

export default {
	name: "Share",
	components: {TabsTab, Tabs},
	data() {
		return {
			activeTab: "applyForSharing",
			config: {},
			deviceId: "",
			shareAuthorizations: [],
			userApplications: []
		}
	},
	mounted() {
		EventBus.on("refresh", this.refresh)
		this.refresh()
	},
	beforeUnmount() {
		EventBus.off("refresh", this.refresh)
	},
	methods: {
		// 刷新
		refresh() {
			this.config = JSON.parse(localStorage.getItem("config"))
			this.getShareAuthorizationList()
			this.getUserApplicationsList()
		},
		// 格式化时间
		formatTime(time) {
			const DATE = new Date(time)
			const YEAR = DATE.getFullYear()
			const MONTH = String(DATE.getMonth() + 1).padStart(2, "0")
			const DAY = String(DATE.getDate()).padStart(2, "0")
			const HOURS = String(DATE.getHours()).padStart(2, "0")
			const MINUTES = String(DATE.getMinutes()).padStart(2, "0")
			const SECONDS = String(DATE.getSeconds()).padStart(2, "0")
			return `${YEAR}-${MONTH}-${DAY} ${HOURS}:${MINUTES}:${SECONDS}`
		},
		// 申请共享
		async applyForSharing() {
			if (!this.deviceId) {
				this.$toast.error("请输入对方的设备ID")
				return
			}
			try {
				const RES = await axios.post(`${this.config.serverUrl}/api/share/apply`, {
					deviceId: this.deviceId,
					viewerId: this.config.userId
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
				this.deviceId = ""
			} catch (error) {
				console.error(error)
				this.$toast.error("申请共享错误")
			}
		},
		// 获取共享授权列表
		async getShareAuthorizationList() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/share/authorizations/${this.config.userId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.shareAuthorizations = RES.data.data.authorizations
			} catch (error) {
				console.error(error)
				this.$toast.error("获取共享授权列表错误")
			}
		},
		// 授权/删除
		async authorization(id, status) {
			if (status === 3) {
				try {
					const RES = await axios.delete(`${this.config.serverUrl}/api/share/delete`, {
						data: {
							id: id
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
					await this.getShareAuthorizationList()
				} catch (error) {
					console.error(error)
					this.$toast.error("删除错误")
				}
			} else {
				try {
					const RES = await axios.put(`${this.config.serverUrl}/api/share/authorize`, {
						id: id,
						status: status
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
					await this.getShareAuthorizationList()
				} catch (error) {
					console.error(error)
					this.$toast.error("授权错误")
				}
			}
		},
		// 获取用户申请的授权列表
		async getUserApplicationsList() {
			try {
				const RES = await axios.get(`${this.config.serverUrl}/api/share/query/${this.config.userId}`, {
					validateStatus: () => {
						return true
					}
				})
				if (!RES.data.success) {
					this.$toast.error(RES.data.data.message)
					return
				}
				this.userApplications = RES.data.data.applications
			} catch (error) {
				console.error(error)
				this.$toast.error("获取用户申请的授权列表错误")
			}
		}
	}
}
</script>

<template>
	<div class="share">
		<tabs v-model="activeTab">
			<tabs-tab name="applyForSharing">
				<template #label>申请共享</template>
				<div class="container">
					<div class="form-item">
						<label>
							我的设备ID：
							<input v-model="config.deviceId" disabled style="color: white;"/>
						</label>
					</div>
					<div class="form-item">
						<label>
							对方的设备ID：
							<input v-model="deviceId" placeholder="请输入对方的设备ID"/>
						</label>
					</div>
					<div class="form-item-but">
						<button @click="applyForSharing">申请共享</button>
					</div>
				</div>
				<div class="container">
					<div class="default" v-if="(userApplications || []).length === 0">没有设备</div>
					<table v-else>
						<thead>
						<tr>
							<th>申请的用户名</th>
							<th>申请的设备名称</th>
							<th>申请时间</th>
							<th>操作</th>
						</tr>
						</thead>
						<tbody>
						<tr
							v-for="item in userApplications"
							:key="item.id"
							:class="{'status-yes': item.status === 1, 'status-no': item.status === 2}">
							<td>{{ item.user_name }}</td>
							<td>{{ item.device_name }}</td>
							<td>{{ formatTime(item.created_at) }}</td>
							<td>
								<i></i>
								<router-link :to="'/device/' + item.device_id">
									<button>查看</button>
								</router-link>
								<button @click="authorization(item.id, 3)">删除</button>
							</td>
						</tr>
						</tbody>
					</table>
				</div>
			</tabs-tab>
			<tabs-tab name="sharedManagement">
				<template #label>共享管理</template>
				<div class="default" v-if="(shareAuthorizations || []).length === 0">没有设备</div>
				<table v-else>
					<thead>
					<tr>
						<th>申请人的用户名</th>
						<th>申请的设备名称</th>
						<th>申请时间</th>
						<th>操作</th>
					</tr>
					</thead>
					<tbody>
					<tr
						v-for="item in shareAuthorizations"
						:key="item.id"
						:class="{'status-yes': item.status === 1, 'status-no': item.status === 2}">
						<td>{{ item.user_name }}</td>
						<td>{{ item.device_name }}</td>
						<td>{{ formatTime(item.created_at) }}</td>
						<td>
							<button @click="authorization(item.id, 1)">同意</button>
							<button @click="authorization(item.id, 2)">拒绝</button>
							<button @click="authorization(item.id, 3)">删除</button>
						</td>
					</tr>
					</tbody>
				</table>
			</tabs-tab>
		</tabs>
	</div>
</template>

<style scoped lang="less">
.share {
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

.container {
	padding: 26px;
	margin-bottom: 16px;
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

table {
	width: 100%;
	border-collapse: collapse;
	table-layout: fixed;

	th, td {
		padding: 10px;
		box-sizing: border-box;
		border: 1px solid var(--border-color);
		text-align: center;
		white-space: nowrap;
		overflow: hidden;
		text-overflow: ellipsis;
	}

	.status-yes {
		background-color: rgba(0, 255, 0, 0.2);
	}

	.status-no {
		background-color: rgba(255, 0, 0, 0.2);
	}

	tbody tr {
		transition: background-color 0.2s;
		cursor: pointer;

		&:hover {
			background-color: rgba(255, 255, 255, 0.4);
		}

		button {
			padding: 5px 10px;
			margin: 0 5px;
			color: white;
			border: none;
			border-radius: 4px;
			cursor: pointer;
			white-space: nowrap;

			&:nth-child(1) {
				background-color: #66cc66;

				&:hover {
					background-color: #4d994d;
				}
			}

			&:nth-child(2), &:nth-child(3) {
				background-color: #ff6666;

				&:hover {
					background-color: #ff4d4d;
				}
			}
		}
	}
}
</style>