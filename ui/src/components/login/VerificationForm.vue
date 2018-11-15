<template>
  <form class="ui form" autocomplete="off" v-on:submit.prevent="onSubmit">
    <h5>Set up Authenticator for 2FA</h5>
    <div style="text-align: left">
      <ol class='mini'>
        <li>
          Add the following code to
          <a href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en">Google
            Authenticator (Android)</a> or <a
          href="https://itunes.apple.com/in/app/google-authenticator/id388497605?mt=8">Google Authenticator (iPhone)</a>
          or <a href="https://www.microsoft.com/en-in/store/p/authenticator/9nblggh08h54">Authenticator (Windows
          Phone).</a>
        </li>
        <li>
          In the app select <b>Set up account</b>.
        </li>
        <li>
          Choose <b>Scan a barcode</b> or enter code manually.
        </li>
      </ol>
    </div>

    <img :src="qr"/>

    <div>
      <code>{{secret}}</code>
    </div>
    <br>

    <div class="field">
      <div class="ui left icon input">
        <i class="clock icon"></i>
        <input type="password"
               v-model.trim="code"
               name="code"
               placeholder="Verification Code"
               maxlength="6"
               minlength="6"
               autocomplete="off"
               required>
      </div>
    </div>
    <sui-button type="submit" :loading="isLoading" :fluid="true" size="large" :disabled="isLoading" color="black">
      Verify OTP
    </sui-button>
  </form>
</template>

<script lang="ts">
  import { Component, Prop, Vue } from 'vue-property-decorator';
  import QrCode from 'qrious';

  @Component({})
  export default class VerificationForm extends Vue {

    @Prop(String) private url!: string;
    @Prop(String) private secret!: string;
    @Prop(Function) private onComplete!: (Object) => void;
    @Prop(Boolean) private isLoading!: boolean;

    private code: string = '';

    protected get qr() {
      const q = new QrCode();
      q.set({
        padding: 15,
        size: 250,
        value: this.url
      });
      return q.toDataURL();
    }

    protected onSubmit(e) {
      e.preventDefault();
      this.onComplete({
        code: this.code
      });
    }
  }

</script>

<style scoped lang="scss">

</style>
