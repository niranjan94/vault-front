<template>
  <div>
    <h4 class="ui top attached header">
      SSH Key Signing
    </h4>
    <div class="ui attached segment" :class="{ loading: isLoading }" v-if="instances.length > 0">
      <form class="ui small form" v-on:submit.prevent="getCredentials">
        <div class="field">
          <label>Instance</label>
          <sui-dropdown
            fluid
            :options="instances"
            placeholder="Select Instance"
            search
            selection
            v-model="selectedRole"/>
        </div>
        <div class="field">
          <label>OpenSSH-compatible Public Key to sign</label>
          <input type="file" name="public-key" placeholder="Public Key" @change="onPublicKeySelect" required>
        </div>
        <div class="field">
          <button type="submit" class="ui fluid green button" :disabled="isLoading || !selectedRole">
            Get credential
          </button>
        </div>
      </form>
    </div>
    <div class="ui attached secondary segment" v-if="instances.length === 0">
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
          <td>Name</td>
          <td>
            <code>{{credentials.name}}</code>
          </td>
        </tr>
        <tr>
          <td>Validity</td>
          <td>
            <code>{{credentials.validity}}</code>
          </td>
        </tr>
        <tr>
          <td>Signed Key</td>
          <td>
            <a href="#" v-on:click.prevent="downloadSignedKey">Download signed key</a>
          </td>
        </tr>
        <tr>
          <td>Usage</td>
          <td>
            <code>{{credentials.usage}}</code>
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
import { saveAs } from 'file-saver';

@Component({})
export default class SshSigning extends Vue {

  protected instances: any[] = [];

  private isLoading: boolean = true;
  private selectedRole: any = null;
  private credentials: any = null;
  private publicKey: any = '';

  protected async created() {
    this.isLoading = true;
    try {
      const response = await axios.get('ssh');
      this.instances = response.data.map((role) => ({ key: role, value: role, text: role }));
    } catch (e) {
      this.$notify({ type: 'error', text: e.message });
    }
    this.isLoading = false;
  }

  protected async getCredentials() {
    this.isLoading = true;
    const publicKey = this.publicKey.trim();

    if (publicKey === '') {
      this.isLoading = false;
      return this.$notify({ type: 'error', text: 'A valid public key is required' });
    }
    if (!publicKey.startsWith('ssh-rsa')) {
      this.isLoading = false;
      return this.$notify({ type: 'error', text: 'A valid OpenSSH-compatible public key is required.' });
    }

    try {
      const response = await axios.post('ssh/sign', {
        role: this.selectedRole,
        publicKey: this.publicKey
      });
      this.credentials = response.data;
      this.credentials.validity = humanizeDuration(this.credentials.validity * 1000);
      this.credentials.fileName = '';
      let fileName = 'signed-cert';
      if (this.credentials.serial && this.credentials.serial !== '') {
        fileName += `-${this.credentials.serial}`
      } else {
        fileName += `-${Math.floor(Date.now() / 1000)}`
      }
      fileName = `${fileName}.pub`;
      this.credentials.fileName = fileName;
      this.credentials.usage = `ssh -i ${fileName} -i <path-to-private-key> ${this.credentials.username}@<hostname>`

    } catch (e) {
      this.$notify({ type: 'error', text: e.message });
    }
    this.isLoading = false;
  }

  protected async downloadSignedKey() {
    const blob = new Blob([this.credentials.signedKey], {type: 'text/plain'});
    saveAs(blob, this.credentials.fileName);
  }

  private readFile(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = function() {
        resolve({
          name: file.name,
          size: file.size,
          type: file.type,
          content: this.result
        });
      };
      reader.onerror = reject;
      reader.readAsBinaryString(file);
    })
  }

  protected async onPublicKeySelect(e) {
    const files = e.target.files || e.dataTransfer.files;
    if (!files.length) {
      return;
    }
    const file = await this.readFile(files[0]);
    this.publicKey = file.content;
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
