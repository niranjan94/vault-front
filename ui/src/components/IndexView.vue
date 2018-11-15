<template>
  <div class="ui stackable grid">
    <div class="ten wide column">
      <WordCloud/>
    </div>
    <div class="six wide column">
      <AccessControl/>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';
  import WordCloud from './WordCloud.vue';
  import AccessControl from './AccessControl.vue';

  @Component({
    components: {AccessControl, WordCloud}
  })
  export default class IndexView extends Vue {
  }
</script>

<style scoped lang="scss">

</style>
