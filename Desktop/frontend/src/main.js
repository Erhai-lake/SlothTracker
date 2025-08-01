import {createApp} from "vue"
import App from "./App.vue"
import router from "./router"
import "./assets/style/theme.css"
import ToastPlugin from "vue-toast-notification"
import "vue-toast-notification/dist/theme-bootstrap.css"

createApp(App)
	.use(router)
	.use(ToastPlugin, {
		position: "top",
		dismissible: false
	})
	.mount("#app")
