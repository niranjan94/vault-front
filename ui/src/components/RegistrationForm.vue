<template>
  <div class="ui middle aligned center aligned grid">
    <div class="six wide column">
      <h2 class="ui image header">
        <div class="content">
          Register
        </div>
      </h2>
      <div class="ui segments">
        <div class="ui segment" v-if="!isSecondStep">
          <details-form :on-complete="onDetailsSubmit" />
        </div>
        <div class="ui segment" v-if="isSecondStep">
          <verification-form :url="url" :secret="secret" :on-complete="onVerificationSubmit" />
        </div>
        <div class="ui link secondary segment">
          <router-link class="item" to="/login" active-class="active">Login</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import axios, { AxiosResponse } from 'axios';
import VerificationForm from './registration/VerificationForm.vue';
import DetailsForm from './registration/DetailsForm.vue';

@Component({
  components: {DetailsForm, VerificationForm}
})
export default class RegistrationForm extends Vue {

  private isSecondStep: boolean = false;
  private url: string = '';
  private secret: string = '';

  private userInformationPayload: any;

  protected async onDetailsSubmit(payload) {
    try {
      const response: AxiosResponse = await axios.post('auth/register', {
        email: payload.email,
        username: payload.username,
        password: payload.password,
        twoFactor: payload.twoFactor
      });

      switch (response.status) {
        case 200: {
          this.url = response.data.url;
          this.secret = response.data.secret;
          this.userInformationPayload = payload;
          this.isSecondStep = true;
          break;
        }
        case 201: {
          this.$notify({ type: 'success',  text: 'Registration successful. Please login' });
          this.$router.push('/login');
          break;
        }
        default: {
          this.$notify({ type: 'error',  text: 'Unrecognized response from server.' });
        }
      }
    } catch (e) {
      switch (e.response.status) {
        case 409: {
          this.$notify({ type: 'warning',  text: 'The email/username already exists.' });
          break;
        }
        default: {
          this.$notify({ type: 'error',  text: e.message });
        }
      }
    }
  }

  protected async onVerificationSubmit(code) {
    try {
      const response: AxiosResponse = await axios.post('auth/register', {
        email: this.userInformationPayload.email,
        username: this.userInformationPayload.username,
        password: this.userInformationPayload.password,
        twoFactor: this.userInformationPayload.twoFactor,
        otp: code
      });
      switch (response.status) {
        case 201: {
          this.userInformationPayload = null;
          this.$notify({ type: 'success',  text: 'Registration successful. Please login' });
          this.$router.push('/login');
          break;
        }
        default: {
          this.$notify({ type: 'error',  text: 'Unrecognized response from server.' });
        }
      }
    } catch (e) {
      switch (e.response.status) {
        case 422: {
          this.$notify({ type: 'error',  text: 'Invalid OTP entered. Please try again' });
          break;
        }
        case 409: {
          this.$notify({ type: 'error',  text: 'The email/username already exists.' });
          break;
        }
        default: {
          this.$notify({ type: 'error',  text: e.message });
        }
      }

    }
  }
}
</script>

<style scoped lang="scss">

</style>
