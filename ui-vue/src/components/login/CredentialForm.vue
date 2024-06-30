<template>
  <form class="ui large form" id="credential-form" autocomplete="off" v-on:submit.prevent="onSubmit">
    <div class="field">
      <div class="ui left icon input">
        <i class="user icon"></i>
        <input type="text"
               v-model.trim="username"
               name="username"
               placeholder="Username"
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
               v-model.trim="otp"
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

    <div class="ui error message"></div>
  </form>
</template>

<script lang="ts">
  import { Component, Vue, Prop } from 'vue-property-decorator';

  @Component({})
  export default class CredentialForm extends Vue {
    private username: string = '';
    private password: string = '';
    private otp: string = '';

    @Prop(Function) private onComplete!: (Object) => void;
    @Prop(Boolean) private isLoading!: boolean;

    protected onSubmit(e) {
      e.preventDefault();
      this.onComplete({
        username: this.username,
        otp: this.otp,
        password: this.password,
      });
    }
  }

</script>

<style scoped lang="scss">

</style>
