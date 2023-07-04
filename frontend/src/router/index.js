import { createRouter, createWebHistory } from 'vue-router'
import SimpleApp from '../views/SimpleApp.vue'

const routes = [
  {
    path: '/',
    name: 'simpleapp',
    component: SimpleApp
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
