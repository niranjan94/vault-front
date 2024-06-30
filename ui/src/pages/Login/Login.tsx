import React, { useEffect } from 'react';
import { EuiButton, EuiPanel, EuiText } from '@elastic/eui';

type LoginProps = {
  isLoading?: boolean
  samlCallback?: boolean
}


function Login({ isLoading, samlCallback }: LoginProps) {
  useEffect(() => {
    // This is a SAML Callback. Validate with API and get token.
  }, [samlCallback])

  return <div style={{display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh'}}>
    <EuiPanel style={{maxWidth: 330}}>
      <EuiText size="s" grow={false} style={{ textAlign: "center"}}>
        <h2 style={{ fontWeight: 300 }}>Canopy Vault</h2>
      </EuiText>
      <br/>
      <EuiButton fullWidth type="submit" size="s" color="secondary" isLoading={isLoading} fill>
        { isLoading ? "Logging you in via SSO" : "Login via SSO"}
      </EuiButton>
    </EuiPanel>
  </div>;
}

export default Login;
