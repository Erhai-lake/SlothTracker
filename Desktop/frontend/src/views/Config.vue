<script>
import axios from "axios"

export default {
	name: "Config",
	data() {
		return {
			form: {
				deviceId: "",
				deviceName: "",
				description: "",
				serverUrl: "",
				refreshInterval: 5
			}
		}
	},
	created() {
		const saved = localStorage.getItem("config")
		if (saved) {
			this.form = JSON.parse(saved)
		}
	},
	methods: {
		saveConfig() {
			localStorage.setItem("config", JSON.stringify(this.form))
			alert("配置已保存!")
		},
		async writeOff() {
			if (!confirm("确认注销吗?")) {
				return
			}
			const RES = await axios.delete(`${this.form.serverUrl}/api/delete/${this.form.deviceId}`)
			localStorage.removeItem("config")
			alert(RES.data.message)
			this.$router.push("/init")
		}
	}
}
</script>

<template>
	<div class="config">
		<div class="container">
			<div class="form-item">
				<label>设备ID：</label>
				<input v-model="form.deviceId" disabled style="color: white;"/>
			</div>
			<div class="form-item">
				<label>设备名称：</label>
				<input v-model="form.deviceName" placeholder="请输入设备名"/>
			</div>
			<div class="form-item">
				<label>描述：</label>
				<input v-model="form.description" placeholder="请输入设备描述"/>
			</div>
			<div class="form-item">
				<label>服务器地址：</label>
				<input v-model="form.serverUrl" placeholder="例如：http://localhost:8080"/>
			</div>
			<div class="form-item">
				<label>自动更新时间：</label>
				<input v-model="form.refreshInterval" placeholder="请输入自动更新时间(-1 禁用)"/>
			</div>
			<button @click="saveConfig">保存</button>
			<button @click="writeOff" class="writeOff">注销</button>
		</div>
	</div>
</template>

<style scoped>
.config {
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

.form-item {
	margin-bottom: 12px;
	display: flex;
	flex-direction: column;
}

input {
	padding: 8px;
	width: 300px;
	border: 1px solid #ccc;
	border-radius: 4px;
}

button {
	margin-bottom: 12px;
	width: 300px;
	padding: 10px;
	background-color: #80ceff;
	color: white;
	border: none;
	border-radius: 4px;
	cursor: pointer;
}

button:hover {
	background-color: #66b1ff;
}

.writeOff {
	background-color: #ff8080;
}

.writeOff:hover {
	background-color: #ff6666;
}
</style>