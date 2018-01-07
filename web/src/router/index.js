import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/components/Login'
import Records from '@/components/Records'
import Logout from '@/components/Logout'
import Signup from '@/components/Signup'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Login',
      component: Login
    },
    {
      path: '/records',
      name: 'Records',
      component: Records
    },
    {
      path: '/logout',
      name: 'Logout',
      component: Logout
    },
    {
      path: '/signup',
      name: 'Signup',
      component: Signup
    }
  ]
})
