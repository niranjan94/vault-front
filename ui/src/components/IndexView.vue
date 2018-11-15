<template>
  <div class="ui stackable centered grid">
    <div class="eight wide column">
      <database-credentials/>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';
  import DatabaseCredentials from './DatabaseCredentials.vue';

  @Component({
    components: {DatabaseCredentials}
  })
  export default class IndexView extends Vue {
  }
</script>

<style scoped lang="scss">

</style>
