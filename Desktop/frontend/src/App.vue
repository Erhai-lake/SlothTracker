<script>
import EventBus from "./services/EventBus.js"

// 定义防抖函数
function debounce(func, delay) {
	let timer = null
	return function (...args) {
		if (timer) clearTimeout(timer)
		timer = setTimeout(() => {
			func.apply(this, args)
		}, delay)
	}
}

export default {
	name: "App",
	data() {
		return {
			sidebar: true,
			interval: null,
			config: {},
			go: window.go,
			refreshInterval: 0,
			backgroundUrl: "https://www.loliapi.com/acg",
			isDragging: false,
			dragStartX: 0,
			dragStartY: 0,
			windowStartX: 0,
			windowStartY: 0
		}
	},
	mounted() {
		document.addEventListener("contextmenu", event => event.preventDefault())
		EventBus.on("initConfig", this.initConfig)
		EventBus.on("sidebarOpen", this.sidebarOpen)
		if (!this.go) {
			this.$toast.warning("非客户端环境无法同步设备状态!")
		}
		this.interval = setInterval(() => {
			if (this.refreshInterval === -1) return
			if (this.refreshInterval > 0) {
				this.refreshInterval--
			} else {
				this.refresh()
			}
		}, 1000)
	},
	beforeUnmount() {
		EventBus.off("initConfig", this.initConfig)
		EventBus.off("sidebarOpen", this.sidebarOpen)
		if (this.interval) {
			clearInterval(this.interval)
		}
		this.removeDragListeners()
	},
	methods: {
		minimizeWindow() {
			window.runtime.WindowMinimise()
		},
		maximizeWindow() {
			window.runtime.WindowToggleMaximise()
		},
		hideWindow() {
			window.runtime.WindowHide()
		},
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
		refreshBackground: debounce(function () {
			this.backgroundUrl = `https://www.loliapi.com/acg?${Date.now()}`
		}, 1000),
		refresh: debounce(async function () {
			this.initConfig()
			this.refreshInterval = Number(this.config.refreshInterval) || -1
			if (this.go) {
				if (this.config.serverUrl && this.config.userId && this.config.deviceId) {
					const RES = await this.go.main.App.UpdateStatus(this.config.serverUrl, this.config.userId, this.config.deviceId)
					if (!RES.success) {
						this.$toast.error(RES.data.message)
						return
					}
				}
				EventBus.emit("refresh")
			} else {
				EventBus.emit("refresh")
			}
		}, 500),
		sidebarOpen(open) {
			this.sidebar = open
		},
		startDrag(event) {
			this.isDragging = true
			this.dragStartX = event.clientX
			this.dragStartY = event.clientY
			// 添加事件监听器
			this.addDragListeners()
			// 阻止文本选择
			event.preventDefault()
		},
		onDrag(event) {
			if (!this.isDragging) return
			const DELTA_X = event.clientX - this.dragStartX
			const DELTA_Y = event.clientY - this.dragStartY
			if (window.runtime && window.runtime.WindowSetPosition) {
				requestAnimationFrame(() => {
					window.runtime.WindowGetPosition().then(position => {
						const NEW_X = position.x + DELTA_X
						const NEW_Y = position.y + DELTA_Y
						window.runtime.WindowSetPosition(NEW_X, NEW_Y)
					})
				})
			}
			// 更新起始位置为当前位置
			this.dragStartX = event.clientX
			this.dragStartY = event.clientY
		},
		endDrag() {
			this.isDragging = false
			this.removeDragListeners()
		},
		addDragListeners() {
			document.addEventListener("mousemove", this.onDrag)
			document.addEventListener("mouseup", this.endDrag)
		},
		removeDragListeners() {
			document.removeEventListener("mousemove", this.onDrag)
			document.removeEventListener("mouseup", this.endDrag)
		},
		// 双击最大化/还原
		handleDoubleClick() {
			this.maximizeWindow()
		}
	}
}
</script>

<template>
	<div class="app">
		<div v-if="sidebar" class="sidebar-container">
			<router-link class="nav-item" to="/">首页</router-link>
			<router-link class="nav-item" to="/share">共享</router-link>
			<router-link class="nav-item" to="/config">设置</router-link>
			<div class="spacer"></div>
			<div class="nav-item clickable" @click="refreshBackground">背景</div>
			<div class="nav-item clickable" @click="refresh">刷新</div>
			<div class="nav-item status">{{ refreshInterval === -1 ? "禁用" : refreshInterval }}</div>
		</div>
		<div v-else></div>
		<div class="container">
			<div class="window-controls" v-if="go" @dblclick="handleDoubleClick" @mousedown="startDrag">
				<div></div>
				<div>
					<button class="control-btn minimize" @click="minimizeWindow">−</button>
					<button class="control-btn maximize" @click="maximizeWindow">□</button>
					<button class="control-btn close" @click="hideWindow">×</button>
				</div>
			</div>
			<div class="view">
				<router-view/>
			</div>
		</div>
		<div
			v-if="config.background" :style="{ backgroundImage: `url(${backgroundUrl})` }"
			class="images images1">
		</div>
		<div v-if="config.background" class="images images2"></div>
	</div>
</template>

<style lang="less" scoped>
.app {
	width: 100%;
	height: 100vh;
	display: grid;
	grid-template-columns: auto 1fr;
	overflow: hidden;
}

.window-controls {
	padding: 10px;
	margin: 10px;
	display: grid;
	grid-template-columns: 1fr auto;
	border-radius: 10px;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border-right: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	z-index: 3;

	.control-btn {
		width: 30px;
		height: 30px;
		border: none;
		background: rgba(255, 255, 255, 0.1);
		color: white;
		cursor: pointer;
		margin-left: 5px;
		border-radius: 3px;
		font-size: 14px;
		font-weight: bold;
	}

	.control-btn:hover {
		background: rgba(255, 255, 255, 0.2);
	}

	.control-btn.close:hover {
		background: #ff5f57;
	}
}

.container {
	display: grid;
	grid-template-rows: auto 1fr;
	overflow: auto;
}

.sidebar-container {
	padding: 16px 0;
	width: 64px;
	height: 100%;
	border-radius: 0 10px 10px 0;
	background-color: rgba(0, 0, 0, 0.4);
	backdrop-filter: blur(5px);
	border-right: 1px solid var(--border-color);
	box-shadow: rgba(142, 142, 142, 0.2) 0 6px 15px 0;
	display: grid;
	grid-template-columns: 1fr;
	grid-template-rows: auto auto auto 1fr auto auto auto;
	gap: 16px;
	text-align: center;
	overflow: hidden;
	user-select: none;
	z-index: 3;

	.nav-item {
		display: flex;
		align-items: center;
		justify-content: center;
		height: 40px;
		font-size: 14px;
		color: #ddd;
		border-radius: 6px;
		margin: 0 8px;
		cursor: pointer;
		transition: all 0.2s ease-in-out;

		&:hover {
			background-color: rgba(255, 255, 255, 0.1);
			color: #fff;
		}

		&.router-link-exact-active {
			background-color: #80ceff;
			color: #fff;
			font-weight: bold;
		}
	}

	.clickable {
		&:active {
			transform: scale(0.96);
		}
	}

	.status {
		font-size: 12px;
		color: #aaa;
		pointer-events: none;
	}
}

.view {
	margin: 0 10px;
	box-sizing: border-box;
	height: 100%;
	overflow: auto;
	z-index: 3;
}

.images {
	position: absolute;
	top: 0;
	left: 0;
	width: 100%;
	height: 100vh;
	overflow: hidden;
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
}

.images1 {
	background-image: url("assets/images/background.png");
	z-index: 2;
}

.images2 {
	background-image: url("assets/images/background.png");
	z-index: 1;
}
</style>
