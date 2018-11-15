<template>
  <div id="app">
    <NavigationBar v-if="authenticatedUser.token"/>
    <div class="ui container">
      <router-view></router-view>
    </div>
    <notifications position="bottom right" :duration="5000"/>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';
  import NavigationBar from './components/NavigationBar.vue';

  @Component({
    components: {
      NavigationBar
    }
  })
  export default class App extends Vue {
    get authenticatedUser() {
      return this.$store.state.session;
    }
  }
</script>

<style lang="scss">
  #app {
    > .ui.container {
      margin-top: 40px;
    }
  }

  *:not(.icon) {
    border-radius: 0 !important;
  }

  .ui.button {
    border-radius: 0 !important;
  }
</style>
