import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    session: {
      username: null,
      token: null
    }
  },
  mutations: {
    setSession(state, session) {
      state.session = session;
    }
  },
  actions: {
    loginUser({ commit }, { username, token }) {
      commit('setSession', { username, token });
    },
    logoutUser({ commit }) {
      commit('setSession', { username: null, token: null });
    }
  },
});
