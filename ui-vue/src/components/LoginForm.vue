<template>
  <div class="ui middle aligned center aligned grid">
    <div class="six wide column">
      <h2 class="ui image header">
        <div class="content">
          Log-in
        </div>
      </h2>
      <div class="ui segment">
        <credential-form :is-loading="isLoading" :on-complete="onSubmit" v-if="currentView === 'credential'"/>
        <password-reset-form :is-loading="isLoading" :on-complete="onSubmit" v-if="currentView === 'password'"/>
        <verification-form :url="otpUrl" :secret="otpSecret" :is-loading="isLoading" :on-complete="onSubmit" v-if="currentView === 'two-factor'"/>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import axios from 'axios';
import CredentialForm from './login/CredentialForm.vue';
import PasswordResetForm from './login/PasswordResetForm.vue';
import VerificationForm from './login/VerificationForm.vue';

@Component({
  components: {VerificationForm, PasswordResetForm, CredentialForm}
})
export default class LoginForm extends Vue {

  private isLoading: boolean = false;
  private currentView: string = 'credential';
  private lastPayload: any;

  private otpUrl: string = '';
  private otpSecret: string = '';

  protected async onSubmit(payload: any) {
    if (this.isLoading) {
      return;
    }
    this.isLoading = true;
    try {


      if (!payload.password && !payload.username) {

        if (payload.newPassword) {
          payload.username = this.lastPayload.username;
          payload.password = this.lastPayload.password;
        }

        if (payload.code) {
          payload = {
            username: this.lastPayload.username,
            password: this.lastPayload.password,
            otp: payload.code
          };
        }
      }

      this.lastPayload = payload;
      const response = await axios.post('auth/session', payload);
      this.isLoading = false;

      if (response.status === 201) {
        this.currentView = 'credential';
        this.$notify({
          title: 'Your password has been changed.',
          type: 'success',
          text: 'Login with your new password to continue.'
        });
        return;
      }

      if (response.data.token) {
        await this.$store.dispatch('loginUser', {
          username: payload.username,
          token: response.data.token
        });
        this.currentView = '';
        this.$notify({
          title: 'You have been logged in.',
          type: 'success',
          text: 'As a security precaution, this session will not be stored by your browser and ' +
            'you will need to re-authenticate after the window is closed or refreshed.'
        });
        this.$router.push('/dashboard');
        return;
      }

      if (response.data.secret) {
        this.currentView = 'two-factor';
        this.otpSecret = response.data.secret;
        this.otpUrl = response.data.url;
        return;
      }

      if (response.data.status === 'password_expired') {
        this.$notify({
          title: 'Your password has expired.',
          type: 'warning',
          text: 'Please set a new password to continue.'
        });
        this.currentView = 'password';
        return;
      }
    } catch (e) {
      const error = (() => {
        if (this.currentView === 'two-factor' && e.response.status === 401) {
          return 'The entered OTP is incorrect.';
        } else {
          if (e.response.status === 401) {
            return 'The credentials entered are incorrect.';
          }
          if (e.response.status === 409) {
            return 'Code already used; wait until the next time period';
          }
          return e.message;
        }
      })();

      this.$notify({ type: 'error',  text: error });
    }

    this.isLoading = false;
  }
}
</script>

<style scoped lang="scss">

</style>


<style lang="scss">
  .link.secondary.segment {
    padding: 7px;
    a {
      color: gray;

      &:hover {
        text-decoration: underline;
      }
    }
  }
</style>
