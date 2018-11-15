<template>
  <div>
    <h4 class="ui top attached header">
      Database Credentials
    </h4>
    <div class="ui attached segment" :class="{ loading: isLoading }">
      <form class="ui small form" v-on:submit.prevent="getCredentials">
        <div class="field">
          <sui-dropdown
            fluid
            :options="databases"
            placeholder="Select Role"
            search
            selection
            v-model="selectedRole"/>
        </div>
        <div class="field">
          <button type="submit" class="ui fluid green button" :disabled="isLoading || !selectedRole">
            Get credential
          </button>
        </div>
      </form>
    </div>
    <div class="ui attached segment" v-if="credentials">
      <table class="ui selectable very basic very compact small table">
        <tbody>
        <tr>
          <td class="four wide">Username</td>
          <td class="twelve wide">
            <code>{{credentials.username}}</code>
          </td>
        </tr>
        <tr>
          <td>Password</td>
          <td>
            <code>{{credentials.password}}</code>
          </td>
        </tr>
        <tr>
          <td>Validity</td>
          <td>
            <code>{{credentials.validity}}</code>
          </td>
        </tr>
        <tr>
          <td>Hostname</td>
          <td>
            <code>{{credentials.expanded.hostname}}</code>
          </td>
        </tr>
        <tr>
          <td>Port</td>
          <td>
            <code>{{credentials.expanded.port}}</code>
          </td>
        </tr>
        <tr>
          <td>URI</td>
          <td>
            <code>{{credentials.connectionUrl}}</code>
          </td>
        </tr>
        </tbody>
      </table>

    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import axios from 'axios';
import humanizeDuration from 'humanize-duration';
import Url from 'url-parse';

@Component({})
export default class DatabaseCredentials extends Vue {

  protected databases: any[] = [];

  private isLoading: boolean = true;
  private selectedRole: any = null;
  private credentials: any = null;

  protected async created() {
    this.isLoading = true;
    try {
      const response = await axios.get('databases');
      this.databases = response.data.map((role) => ({ key: role, value: role, text: role }));
    } catch (e) {
      this.$notify({ type: 'error', text: e.message });
    }
    this.isLoading = false;
  }

  protected async getCredentials() {
    this.isLoading = true;
    try {
      const response = await axios.post('databases/credentials', {
        role: this.selectedRole
      });
      this.credentials = response.data;
      this.credentials.validity = humanizeDuration(this.credentials.validity * 1000);
      this.credentials.expanded = new Url(this.credentials.connectionUrl);
    } catch (e) {
      this.$notify({ type: 'error', text: e.message });
    }
    this.isLoading = false;
  }

}
</script>

<style scoped lang="scss">
  .muted {
    color: gray;
  }
  table {
    tr {
      td:nth-child(1) {
        font-weight: bold;
      }
    }
  }
</style>
