<template>
  <div class="ui stackable centered grid">
    <div class="eight wide column">
      <database-credentials/>
      <br>
      <system-credentials/>
    </div>
    <div class="eight wide column">
      <ssh-signing />
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';
  import DatabaseCredentials from './DatabaseCredentials.vue';
  import SshSigning from './SshSigning.vue';
  import SystemCredentials from './SystemCredentials.vue';

  @Component({
    components: {SystemCredentials, SshSigning, DatabaseCredentials}
  })
  export default class IndexView extends Vue {
  }
</script>

<style scoped lang="scss">

</style>
