<template>
  <div>
    <h4 class="ui top attached header">
      System Credentials
    </h4>
    <div class="ui attached segment" :class="{ loading: isLoading }" v-if="roles.length > 0">
      <form class="ui small form" v-if="roles.length > 0" v-on:submit.prevent="getCredentials">
        <div class="field">
          <sui-dropdown
            fluid
            :options="roles"
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
    <div class="ui attached secondary segment" v-if="roles.length === 0">
      <p>You do not have permissions to access any instances.</p>
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
export default class SystemCredentials extends Vue {

  protected roles: any[] = [];

  private isLoading: boolean = true;
  private selectedRole: any = null;
  private credentials: any = null;

  protected async created() {
    this.isLoading = true;
    try {
      const response = await axios.get('system');
      this.roles = response.data.map((role) => ({
        key: role, value: role, text: role.replace('instance-creds/data/', '').split('/').join(' — ')
      }));
    } catch (e) {
      this.$notify({ type: 'error', text: e.message });
    }
    this.isLoading = false;
  }

  protected async getCredentials() {
    this.isLoading = true;
    try {
      const response = await axios.post('system/credentials', {
        role: this.selectedRole
      });
      this.credentials = response.data;
      this.credentials.validity = humanizeDuration(this.credentials.validity * 1000);
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
