import React, { useState } from 'react';
import QRCode from 'qrcode.react';
import { EuiButton, EuiFieldText, EuiForm, EuiFormRow, EuiText } from '@elastic/eui';

type TwoFactorActivationFormProps = {
  activationUrl: string
  secret: string
  onSubmit: (otp: string) => Promise<void>
}

function TwoFactorActivationForm({onSubmit, activationUrl, secret}: TwoFactorActivationFormProps) {
  const [otp, setOtp] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const onFormSubmit = () => {
    setIsLoading(true);
    onSubmit(otp)
      .then(() => {})
      .catch(() => console.error)
      .finally(() => setIsLoading(false));
  };

  return <div>
    <EuiText size="s" grow={false}>
      <h4>Set up Authenticator for 2FA</h4>
      <ol>
        <li>
          Add the following code to <a
          href="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en">Google
          Authenticator (Android)</a> or <a
          href="https://itunes.apple.com/in/app/google-authenticator/id388497605?mt=8">Google Authenticator
          (iPhone)</a> or <a
          href="https://www.microsoft.com/en-in/store/p/authenticator/9nblggh08h54">Authenticator (Windows
          Phone).</a>
        </li>
        <li>
          In the app select <b>Set up account</b>.
        </li>
        <li>
          Choose <b>Scan a barcode</b> or enter code manually.
        </li>
      </ol>
      <div style={{textAlign: 'center'}}>
        <QRCode value={activationUrl} includeMargin size={256}/>

        <code>{secret}</code>
      </div>

    </EuiText>
    <br/>
    <EuiForm onSubmit={onFormSubmit}>
      <EuiFormRow
        label=""
        display="rowCompressed">
        <EuiFieldText name="code"
                      disabled={isLoading}
                      compressed
                      icon="grid"
                      placeholder="Enter the OTP from the app"
                      value={otp}
                      onChange={(e) => setOtp(e.target.value)}/>
      </EuiFormRow>
      <EuiButton fullWidth type="submit" size="s" color="secondary" fill isLoading={isLoading}>
        Verify OTP
      </EuiButton>
    </EuiForm>
  </div>;
}

export default TwoFactorActivationForm;
