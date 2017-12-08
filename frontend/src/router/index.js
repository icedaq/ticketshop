import Vue from 'vue'
import Router from 'vue-router'
import Welcome from '@/components/Welcome'
import Artists from '@/components/Artists'
import Shows from '@/components/Shows'
import Tickets from '@/components/Tickets'
import Buy from '@/components/Buy'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Welcome',
      component: Welcome
    },
    {
      path: '/artists',
      name: 'Artists',
      component: Artists
    },
    {
      path: '/shows/:id',
      name: 'Shows',
      component: Shows
    },
    {
      path: '/tickets/:id',
      name: 'Tickets',
      component: Tickets
    },
    {
      path: '/buy/:id',
      name: 'Buy',
      component: Buy
    }
  ]
})
