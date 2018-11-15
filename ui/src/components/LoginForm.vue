<template>
  <div class="ui middle aligned center aligned grid">
    <div class="six wide column">
      <h2 class="ui image header">
        <div class="content">
          Log-in
        </div>
      </h2>
      <form class="ui large form" id="login-form" autocomplete="off" v-on:submit.prevent="onSubmit">
        <div class="ui segments">
          <div class="ui segment">
            <div class="field">
              <div class="ui left icon input">
                <i class="at icon"></i>
                <input type="email"
                       v-model.trim="email"
                       name="email"
                       placeholder="E-mail address"
                       autocomplete="off"
                       required>
              </div>
            </div>
            <div class="field">
              <div class="ui left icon input">
                <i class="lock icon"></i>
                <input type="password"
                       v-model.trim="password"
                       name="password"
                       placeholder="Password"
                       autocomplete="off"
                       required>
              </div>
            </div>
            <div class="field">
              <div class="ui left icon input">
                <i class="clock icon"></i>
                <input type="password"
                       v-model.trim="code"
                       name="code"
                       placeholder="2FA Code"
                       maxlength="6"
                       minlength="6"
                       autocomplete="off">
              </div>
            </div>
            <sui-button type="submit" :loading="isLoading" :fluid="true" size="large" :disabled="isLoading" color="black">
              Login
            </sui-button>
          </div>
          <div class="ui link secondary segment">
            <router-link class="item" to="/register" active-class="active" data-test="register-button">Register</router-link>
          </div>
        </div>

        <div class="ui error message"></div>
      </form>

    </div>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator';
  import axios from 'axios';

  @Component({
    components: {}
  })
  export default class LoginForm extends Vue {

    private email: string = '';
    private password: string = '';
    private code: string = '';

    private isLoading: boolean = false;

    protected async onSubmit() {
      if (this.isLoading) {
        return;
      }
      this.isLoading = true;
      try {
        const response = await axios.post('auth/session', {
          email: this.email,
          password: this.password,
          otp: this.code
        });

        await this.$store.dispatch('loginUser', {
          email: this.email,
          token: response.data.token
        });

        this.$notify({
          title: 'You have been logged in.',
          type: 'success',
          text: 'As a security precaution, this session will not be stored by your browser and ' +
            'you will need to re-authenticate after the window is closed or refreshed.'
        });

        this.$router.push('/');

      } catch (e) {
        if (e.response.status === 401 || e.response.status === 409) {
          this.$notify({ type: 'error',  text: 'The credentials entered are incorrect.' });
        } else {
          this.$notify({ type: 'error',  text: e.message });
        }
      }

      this.password = '';
      this.code = '';

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
