import React, { useState } from 'react';
import { EuiButton, EuiFieldPassword, EuiFieldText, EuiForm, EuiFormRow } from '@elastic/eui';

type PasswordResetFormProps = {
  onSubmit: (username: string, confirmPassword: string, otp: string) => Promise<void>
}


function PasswordResetForm({ onSubmit }: PasswordResetFormProps) {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [otp, setOtp] = useState('');

  const onFormSubmit = () => {
    setIsLoading(true);
    onSubmit(password, confirmPassword, otp)
      .then(() => {})
      .catch(() => console.error)
      .finally(() => setIsLoading(false));
  };

  return <EuiForm onSubmit={onFormSubmit}>
    <EuiFormRow
      label="Password"
      display="rowCompressed">
      <EuiFieldPassword name="password"
                        disabled={isLoading}
                        compressed
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}/>
    </EuiFormRow>
    <EuiFormRow
      label="Confirm Password"
      display="rowCompressed">
      <EuiFieldPassword name="password"
                        disabled={isLoading}
                        compressed
                        value={password}
                        onChange={(e) => setConfirmPassword(e.target.value)}/>
    </EuiFormRow>
    <EuiFormRow
      label="2FA Code"
      display="rowCompressed">
      <EuiFieldText name="code"
                    disabled={isLoading}
                    compressed
                    icon="grid"
                    value={otp}
                    onChange={(e) => setOtp(e.target.value)}/>
    </EuiFormRow>
    <EuiButton fullWidth type="submit" size="s" color="secondary" fill isLoading={isLoading}>
      Change Password
    </EuiButton>
  </EuiForm>;
}

export default PasswordResetForm;
