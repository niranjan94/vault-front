<template>
  <div class="ui inverted menu">
    <div class="header item">
      PII Vault
    </div>
    <router-link class="item" to="/" active-class="active">Dashboard</router-link>
    <div class="right menu">
      <sui-dropdown :item="true" :text="authenticatedUser.email">
        <sui-dropdown-menu>
          <sui-dropdown-item v-on:click.prevent="logout">Logout</sui-dropdown-item>
        </sui-dropdown-menu>
      </sui-dropdown>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import axios from 'axios';

@Component({})
export default class NavigationBar extends Vue {
  get authenticatedUser() {
    return this.$store.state.session;
  }

  public async logout() {
    await axios.delete('auth/session');
    await this.$store.dispatch('logoutUser');
    this.$notify({ clean: true });
    this.$router.push('/login');
  }
}
</script>

<style scoped lang="scss">

</style>
