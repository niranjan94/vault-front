import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    session: {
      email: null,
      token: null
    }
  },
  mutations: {
    setSession(state, session) {
      state.session = session;
    }
  },
  actions: {
    loginUser({ commit }, { email, token }) {
      commit('setSession', { email, token });
    },
    logoutUser({ commit }) {
      commit('setSession', { email: null, token: null });
    }
  },
});
