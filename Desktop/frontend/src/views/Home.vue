<script>
import axios from "axios"
import EventBus from "../services/EventBus.js"

export default {
	name: "Home",
	data() {
		return {
			devices: []
		}
	},
	mounted() {
		EventBus.on("refresh", this.getDevices)
	},
	beforeUnmount() {
		EventBus.off("refresh", this.getDevices)
	},
	created() {
		this.getDevices()
	},
	methods: {
		async getDevices() {
			const CONFIG = JSON.parse(localStorage.getItem("config"))
			try {
				const RES = await axios.get(`${CONFIG.serverUrl}/api/devices`)
				this.devices = RES.data.devices
			} catch (error) {
				console.error(error)
			}
		}
	}
}
</script>

<template>
	<div class="home">
		<div class="default" v-if="(devices || []).length === 0">没有任何设备</div>
		<div class="item">
			<router-link
				:to="'/' + item.ID"
				class="container"
				v-for="item in devices"
				:key="item.ID">
				<p :title="item.name">{{ item.name }}</p>
				<p :title="item.platform">{{ item.platform }}</p>
				<p :title="item.description">{{ item.description }}</p>
			</router-link>
		</div>
	</div>
</template>

<style scoped>
.home {
	padding: 16px;
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
}

.container:hover {
	color: #000;
	background-color: rgba(255, 255, 255, 0.4);
}

.container p {
	text-align: center;
	white-space: nowrap;
	overflow: hidden;
	text-overflow: ellipsis;
}
</style>
