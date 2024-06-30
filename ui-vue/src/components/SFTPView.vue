<template>
  <div class="ui stackable centered grid">
    <div class="sixteen wide column">
      <sui-tab :menu="{ secondary: true }">
        <sui-tab-pane title="Users" :attached="false">
          <sui-button content="Add User" primary icon="user" size="tiny" label-position="left" style="float: right" />
          <br/>
          <br/>
          <sui-table basic="very">
            <sui-table-header>
              <sui-table-row>
                <sui-table-header-cell>Username</sui-table-header-cell>
                <sui-table-header-cell>Environment</sui-table-header-cell>
                <sui-table-header-cell>Auth mode</sui-table-header-cell>
                <sui-table-header-cell>IP Restricted</sui-table-header-cell>
                <sui-table-header-cell>No. of directories</sui-table-header-cell>
                <sui-table-header-cell text-align="right">Actions</sui-table-header-cell>
              </sui-table-row>
            </sui-table-header>
            <sui-table-body>
              <sui-table-row>
                <sui-table-cell>niranjan</sui-table-cell>
                <sui-table-cell>Shared</sui-table-cell>
                <sui-table-cell>Public key</sui-table-cell>
                <sui-table-cell>Yes</sui-table-cell>
                <sui-table-cell>5</sui-table-cell>
                <sui-table-cell text-align="right">
                  <sui-button-group size="tiny">
                    <sui-button icon="code" />
                    <sui-button icon="edit" />
                    <sui-button icon="trash" />
                  </sui-button-group>
                </sui-table-cell>
              </sui-table-row>
            </sui-table-body>
          </sui-table>
        </sui-tab-pane>
        <sui-tab-pane title="Remote Pulls" :attached="false">
          Box
        </sui-tab-pane>
      </sui-tab>
    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';

  @Component({
    components: {}
  })
  export default class SFTPView extends Vue {
  }
</script>

<style scoped lang="scss">

</style>
