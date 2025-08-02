<script>
import EventBus from "./services/EventBus.js"

export default {
	name: "App",
	data() {
		return {
			sidebar: true,
			config: {},
			refreshInterval: 0
		}
	},
	mounted() {
		document.addEventListener("contextmenu", event => event.preventDefault())
		EventBus.on("initConfig", this.initConfig)
		EventBus.on("sidebarOpen", this.sidebarOpen)
	},
	beforeUnmount() {
		EventBus.off("initConfig", this.initConfig)
		EventBus.off("sidebarOpen", this.sidebarOpen)
	},
	created() {
		this.refresh()
		if (this.refreshInterval) {
			setInterval(() => {
				if (this.refreshInterval === -1) return
				if (this.refreshInterval > 0) {
					this.refreshInterval--
				} else {
					this.refresh()
				}
			}, 1000)
		}
	},
	methods: {
		initConfig() {
			const CONFIG = JSON.parse(localStorage.getItem("config"))
			const REQUIRED_FIELDS = ["serverUrl", "refreshInterval", "userId", "deviceId"]
			if (!CONFIG) {
				this.$router.push("/init")
				return
			}
			for (const FIELD of REQUIRED_FIELDS) {
				if (!CONFIG[FIELD]) {
					this.$router.push("/init")
					return
				}
			}
			this.config = CONFIG
		},
		async refresh() {
			this.initConfig()
			this.refreshInterval = Number(this.config.refreshInterval) || -1
			EventBus.emit("refresh")
			await window.go.main.App.UpdateStatus(this.config.serverUrl, this.config.deviceId)
		},
		sidebarOpen(open) {
			this.sidebar = open
		}
	}
}
</script>

<template>
	<div class="app">
		<div class="sidebar-container" v-if="sidebar">
			<router-link to="/">首页</router-link>
			<router-link to="/config">设置</router-link>
			<div></div>
			<div @click="refresh">刷新</div>
			<div>{{ refreshInterval === -1 ? "禁用" : refreshInterval }}</div>
		</div>
		<div class="view">
			<router-view/>
		</div>
		<div class="images" v-if="config.background"></div>
	</div>
</template>

<style scoped lang="less">
.app {
	width: 100%;
	height: 100vh;
	display: flex;
	overflow: hidden;
}

.sidebar-container {
	padding: 16px 0;
	width: 64px;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border-right: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	display: grid;
	grid-template-columns: 1fr;
	grid-template-rows: auto auto 1fr auto auto;
	gap: 16px;
	text-align: center;
	overflow: hidden;
	user-select: none;
	z-index: 2;
}

.view {
	width: 100%;
	height: 100%;
	overflow: auto;
	z-index: 2;
}


.images {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100vh;
	overflow: hidden;
	background-image: url("https://www.loliapi.com/acg");
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
	filter: opacity(0.8);
	z-index: 1;
}
</style>
