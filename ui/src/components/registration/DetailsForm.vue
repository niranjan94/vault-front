<template>
  <form id="registration-details-form" class="ui form" autocomplete="off" v-on:submit.prevent="onSubmit">
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
        <i class="user icon"></i>
        <input type="text"
               v-model.trim="username"
               name="username"
               placeholder="Canopy Username"
               autocomplete="off"
               required>
      </div>
    </div>
    <div class="field" v-bind:class="{ error: passwordErrorMessage }">
      <div class="ui left icon input">
        <i class="lock icon"></i>
        <input type="password"
               v-model.trim="password"
               name="password"
               placeholder="Password"
               minlength="8"
               autocomplete="new-password"
               required>
      </div>
    </div>
    <div class="field" v-bind:class="{ error: passwordErrorMessage }">
      <div class="ui left icon input">
        <i class="lock icon"></i>
        <input type="password"
               v-model.trim="passwordConfirm"
               name="passwordConfirm"
               placeholder="Confirm password"
               minlength="8"
               autocomplete="off"
               required>
      </div>
    </div>
    <div style="color: #9f3a38; text-align: left; margin-bottom: 12px;" v-if="passwordErrorMessage">
      {{passwordErrorMessage}}
    </div>
    <div class="field ">
      <sui-checkbox class="slider" label="Enable 2FA (Recommended)" v-model="twoFactor"/>
    </div>

    <sui-button type="submit" :loading="isLoading" :fluid="true" size="large" :disabled="isLoading" color="black">
      Register
    </sui-button>

  </form>
</template>

<script lang="ts">
  import { Component, Prop, Vue } from 'vue-property-decorator';

  const passwordRegex = /^(?=.*[0-9])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{8,}$/;

  @Component({})
  export default class DetailsForm extends Vue {

    @Prop(Function) private onComplete!: (Object) => void;
    @Prop(Boolean) private isLoading!: boolean;

    private email: string = '';
    private username: string = '';
    private password: string = '';
    private passwordConfirm: string = '';
    private twoFactor: boolean = true;

    get passwordErrorMessage() {
      if (this.password === '') {
        return false;
      }
      if (this.password.length < 8) {
        return 'Password should be atleast 8 characters long';
      }
      if (!passwordRegex.test(this.password)) {
        return 'Password should be a mix of letter, numbers & special characters.';
      }
      if (this.password !== this.passwordConfirm) {
        return 'Password and its confirm do not match.';
      }
      return false;
    }

    protected onSubmit(e) {
      e.preventDefault();
      if (this.passwordErrorMessage !== false) {
        return;
      }
      this.onComplete({
        email: this.email,
        username: this.username,
        password: this.password,
        twoFactor: this.twoFactor
      });
    }

  }

</script>

<style scoped lang="scss">

</style>
