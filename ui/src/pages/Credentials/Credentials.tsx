import React from 'react';
import { EuiFlexGroup, EuiFlexItem, EuiSpacer } from '@elastic/eui';
import style from './Credentials.module.scss';
import DatabaseCredentials from '../../components/credentials/DatabaseCredentials/DatabaseCredentials';
import WindowsCredentials from '../../components/credentials/WindowsCredentials/WindowsCredentials';

const maxWidth = 550;

function Credentials() {
  return <div className={style.credentials}>
    <EuiSpacer />
    <EuiFlexGroup justifyContent="center">
      <EuiFlexItem style={{ maxWidth, display: "block" }}>
        <DatabaseCredentials/>
      </EuiFlexItem>
      <EuiFlexItem style={{ maxWidth, display: "block" }}>
        <WindowsCredentials/>
      </EuiFlexItem>
    </EuiFlexGroup>
  </div>;
}

export default Credentials;
