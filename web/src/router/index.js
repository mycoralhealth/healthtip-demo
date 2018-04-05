import Vue from 'vue';
import Router from 'vue-router';
import Login from '@/components/Login';
import Callback from '@/components/Callback';
import Records from '@/components/Records';
import Logout from '@/components/Logout';
import Signup from '@/components/Signup';
import Forgot from '@/components/Forgot';
import ChangePassword from '@/components/ChangePassword';

Vue.use(Router);

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'Login',
      component: Login,
    },
    {
      path: '/callback',
      name: 'Callback',
      component: Callback,
    },
    {
      path: '/records',
      name: 'Records',
      component: Records,
    },
    {
      path: '/logout',
      name: 'Logout',
      component: Logout,
    },
    {
      path: '/signup',
      name: 'Signup',
      component: Signup,
    },
    {
      path: '/forgot',
      name: 'Forgot',
      component: Forgot,
    },
    {
      path: '/changePass',
      name: 'ChangePassword',
      component: ChangePassword,
    },
  ],
});
