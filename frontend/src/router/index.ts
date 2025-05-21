import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/Admin.vue')
  },
  {
    path: '/incidents/:id',
    name: 'IncidentDetails',
    component: () => import('../views/IncidentDetails.vue')
  },
  {
    path: '/components',
    name: 'Components',
    component: () => import('../views/Components.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
