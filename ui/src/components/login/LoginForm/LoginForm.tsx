import React, { useState } from 'react';
import { EuiButton, EuiFieldPassword, EuiFieldText, EuiForm, EuiFormRow } from '@elastic/eui';

type LoginFormProps = {
  onSubmit: (username: string, password: string, otp: string) => Promise<void>
}


function LoginForm({ onSubmit }: LoginFormProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [otp, setOtp] = useState('');
  const [isLoading, setIsLoading] = useState(false);

  const onFormSubmit = () => {
    setIsLoading(true);
    onSubmit(username, password, otp)
      .then(() => {})
      .catch(() => console.error)
      .finally(() => setIsLoading(false));
  };

  return <EuiForm onSubmit={onFormSubmit}>
    <EuiFormRow
      label="Username"
      display="rowCompressed">
      <EuiFieldText name="username"
                    disabled={isLoading}
                    compressed
                    icon="user"
                    value={username}
                    onChange={(e) => setUsername(e.target.value)}/>
    </EuiFormRow>
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
      label="2FA Code"
      helpText="This can be left empty, if logging in for the first time."
      display="rowCompressed">
      <EuiFieldText name="code"
                    disabled={isLoading}
                    compressed
                    icon="grid"
                    value={otp}
                    onChange={(e) => setOtp(e.target.value)}/>
    </EuiFormRow>
    <EuiButton fullWidth type="submit" size="s" color="secondary" fill isLoading={isLoading}>
      Login
    </EuiButton>
  </EuiForm>;
}

export default LoginForm;
