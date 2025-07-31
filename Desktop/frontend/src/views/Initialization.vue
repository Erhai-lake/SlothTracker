<script>
import EventBus from "../services/EventBus.js"
import axios from "axios"

export default {
	name: "Initialization",
	data() {
		return {
			form: {
				deviceName: "",
				description: "",
				serverUrl: ""
			}
		}
	},
	created() {
		const CONFIG = JSON.parse(localStorage.getItem("config"))
		if (CONFIG) {
			this.$router.push("/")
		}
	},
	methods: {
		async confirm() {
			if (this.form.deviceName === "" || this.form.description === "" || this.form.serverUrl === "") {
				alert("请填写完整信息")
				return
			}
			let deviceId = null
			try {
				const RES = await axios.post(`${this.form.serverUrl}/api/register`,
					{
						deviceName: this.form.deviceName,
						platform: "windows",
						description: this.form.description
					})
				deviceId = RES.data.deviceId
				if (deviceId === null) {
					alert("注册失败")
					return
				}
			} catch (error) {
				console.error(error)
			}
			let config = JSON.parse(localStorage.getItem("config"))
			config = {
				deviceId: deviceId,
				...this.form,
				refreshInterval: 5,
			}
			localStorage.setItem("config", JSON.stringify(config))
			EventBus.emit("initConfig")
			this.$router.push("/")
		}
	}
};
</script>

<template>
	<div class="initialization">
		<h2>初始化设置</h2>
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
		<button @click="confirm">确认</button>
	</div>
</template>

<style scoped>
.initialization {
	padding: 16px;
	width: 100%;
	height: 100%;
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
}

h2 {
	text-align: center;
	margin-bottom: 16px;
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
</style>
