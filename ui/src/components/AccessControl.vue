<template>
  <div>
    <h4 class="ui top attached header">
      Access Control
    </h4>
    <div class="ui attached segment" v-if="selectableRoles.length > 0">
      <form class="ui small form" v-on:submit.prevent="addRole">
        <div class="fields">
          <div class="twelve wide required field">
            <label>Role</label>
            <sui-dropdown
              multiple
              fluid
              :options="selectableRoles"
              placeholder="Roles"
              search
              selection
              v-model="selectedRoles"/>
          </div>
          <div class="four wide field">
            <label>&nbsp;</label>
            <button type="submit" class="ui fluid button" :disabled="isSaving">Add</button>
          </div>
        </div>
      </form>
    </div>
    <div class="ui bottom attached segment" v-if="roles.length > 0">
      <table class="ui very compact very basic small fluid table">
        <thead>
        <tr>
          <th>Role</th>
          <th class="two wide"> </th>
        </tr>
        </thead>
        <tbody>
          <tr v-for="(role) in roles">
            <td>{{role.replace('app-', '')}}</td>
            <td>
              <button class="ui mini red icon button"
                      data-tooltip="Delete"
                      :disabled="isSaving"
                      data-position="right center"
                      v-on:click.prevent="deleteRole(role)"
                      data-inverted="">
                <i class="trash icon"></i>
              </button>
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

  @Component({})
  export default class AccessControl extends Vue {
    private roles: string[] = [];

    private isSaving: boolean = false;

    private selectedRoles: string[] = [];

    private existingRoles: any[] = [];

    get selectableRoles() {
      return this.existingRoles.filter((role) => !this.roles.includes(role.value));
    }

    protected async created() {
      const apps = await axios.get('acl/apps');
      this.existingRoles = [];
      if (apps.data) {
        this.existingRoles = (apps.data as string[]).map((policy) => ({
          key: policy,
          text: policy.replace('app-', ''),
          value: policy
        }));
      }
      const roles = await axios.get('acl');
      this.roles = roles.data || [];
    }

    protected async addRole() {
      const expanded = this.roles.concat(this.selectedRoles);
      await this.save(expanded);
    }

    protected async deleteRole(role) {
      const filtered = this.roles.filter((r) => r !== role);
      await this.save(filtered);
    }

    private async save(roles) {
      this.isSaving = true;
      try {
        await axios.post('acl', roles);
        this.roles = roles;
        this.selectedRoles = [];
      } catch (e) {
        this.$notify({ type: 'error', text: e.message });
      }
      this.isSaving = false;
    }

  }
</script>

<style scoped lang="scss">

</style>
