import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'


const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: HomeView},
		{name: 'Login', path:'/login', component: LoginView},
	]
});

router.beforeEach((to, from, next) => {
	const token = sessionStorage.getItem('token');
	if (to.name !== 'Login' && !token) next({ name: 'Login' })
	else next()
});

export default router
