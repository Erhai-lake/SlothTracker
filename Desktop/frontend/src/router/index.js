import { createRouter, createWebHashHistory } from "vue-router"

const routes = [
    {
        path: "/",
	    name: "Home",
        component: () => import("../views/Home.vue")
    },
	{
		path: "/device/:id",
		name: "Device",
		component: () => import("../views/Device.vue")
	},
	{
		path: "/share",
		name: "Share",
		component: () => import("../views/Share.vue")
	},
    {
        path: "/config",
        name: "Config",
	    component: () => import("../views/Config.vue")
    },
    {
        path: "/init",
        name: "Initialization",
        component: () => import("../views/Initialization.vue")
    }
]

const router = createRouter({
    history: createWebHashHistory(),
    routes
})

export default router
