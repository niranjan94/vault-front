import VueRouter from 'vue-router';
import LoginForm from './components/LoginForm.vue';
import DashboardView from './components/DashboardView.vue';
import NotFound from './components/NotFound.vue';
import SFTPView from './components/SFTPView.vue';
import store from './store';

let router;

const authGuard = (to, from, next) => {
  if (store.state.session.token) {
    next();
  } else {
    next('/login');
  }
};

const routes = [
  { path: '/dashboard', component: DashboardView, beforeEnter: authGuard },
  { path: '/sftp', component: SFTPView, beforeEnter: authGuard },
  { path: '/login', component: LoginForm },
  { path: '*', component: NotFound, beforeEnter: authGuard }
];

router = new VueRouter({
  routes,
});

export { router };
