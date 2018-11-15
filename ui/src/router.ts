import VueRouter from 'vue-router';
import LoginForm from './components/LoginForm.vue';
import IndexView from './components/IndexView.vue';
import NotFound from './components/NotFound.vue';
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
  { path: '/', component: IndexView, beforeEnter: authGuard },
  { path: '/login', component: LoginForm },
  { path: '*', component: NotFound, beforeEnter: authGuard }
];

router = new VueRouter({
  routes,
});

export { router };
