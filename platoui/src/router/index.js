import Vue from 'vue'
import VueRouter from 'vue-router'
import Plato from '../views/Plato.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    name: 'plato',
    component: Plato
  },
  {
    path: '/plato',
    name: 'plato',
    component: Plato
  }
]

const router = new VueRouter({
  routes
})

export default router
