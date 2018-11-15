<template>
  <div class="ui middle aligned center aligned grid">
    <div class="six wide column">
      <div class="ui segment">
        <h2>Not found</h2>
        <p>
          The page your looking for does not exist.
        </p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';

  @Component({})
  export default class IndexView extends Vue {
  }
</script>

<style scoped lang="scss">

</style>
