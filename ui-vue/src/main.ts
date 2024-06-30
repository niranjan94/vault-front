import 'es6-promise/auto';
import Vue from 'vue';
import SuiVue from 'semantic-ui-vue';
import VueRouter from 'vue-router';
import Notifications from 'vue-notification';
import axios from 'axios';

import 'semantic-ui-css/semantic.min.css';

import App from './App.vue';
import store from './store';
import { router } from './router';

Vue.config.productionTip = false;

Vue.use(SuiVue);
Vue.use(VueRouter);
Vue.use(Notifications);

if (process.env.VUE_APP_BASE_URL && process.env.VUE_APP_BASE_URL !== '') {
  axios.defaults.baseURL = process.env.VUE_APP_BASE_URL;
} else {
  const l = document.location as Location;
  axios.defaults.baseURL = `${l.protocol}//${l.host}/api/v1`;
}

axios.interceptors.request.use((config) => {
  if (store.state.session && store.state.session.token) {
    config.headers = config.headers || {};
    config.headers.Authorization = 'Bearer ' + store.state.session.token;
  }
  return config;
});

// noinspection JSUnusedGlobalSymbols
new Vue({
  store,
  router,
  render: (h) => h(App),
}).$mount('#app');
