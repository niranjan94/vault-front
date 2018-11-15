<template>
  <form id="password-reset-form" class="ui form" autocomplete="off" v-on:submit.prevent="onSubmit">
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

    <div class="field">
      <div class="ui left icon input">
        <i class="clock icon"></i>
        <input type="password"
               v-model.trim="otp"
               name="otp"
               placeholder="2FA Code"
               maxlength="6"
               minlength="6"
               required
               autocomplete="off">
      </div>
    </div>

    <div style="color: #9f3a38; text-align: left; margin-bottom: 12px;" v-if="passwordErrorMessage">
      {{passwordErrorMessage}}
    </div>

    <sui-button type="submit" :loading="isLoading" :fluid="true" size="large" :disabled="isLoading" color="black">
      Change Password
    </sui-button>

  </form>
</template>

<script lang="ts">
  import { Component, Prop, Vue } from 'vue-property-decorator';

  const passwordRegex = /^(?=.*[0-9])(?=.*[A-Z])(?=.*[!@#$%^&*])[a-zA-Z0-9!@#$%^&*]{8,}$/;

  @Component({})
  export default class DetailsForm extends Vue {

    @Prop(Function) private onComplete!: (Object) => void;
    @Prop(Boolean) private isLoading!: boolean;

    private password: string = '';
    private passwordConfirm: string = '';
    private otp: string = '';

    get passwordErrorMessage() {
      if (this.password === '') {
        return false;
      }
      if (this.password.length < 8) {
        return 'Password should be atleast 8 characters long';
      }
      if (!passwordRegex.test(this.password)) {
        return 'Password should be a mix of lower/upper case letters, numbers & special characters.';
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
        newPassword: this.password,
        otp: this.otp,
      });
    }

  }

</script>

<style scoped lang="scss">

</style>
