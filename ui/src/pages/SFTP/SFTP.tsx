import React from 'react';
import { EuiTabbedContent, EuiSpacer, EuiFlexGroup, EuiFlexItem, EuiPanel } from '@elastic/eui';
import UserManagement from '../../components/sftp/UserManagement/UserManagement';
import PullConfiguration from '../../components/sftp/PullConfiguration/PullConfiguration';

function SFTP() {
  const tabs = [
    {
      id: 'user-management--id',
      name: 'User Management',
      content: (
        <>
          <EuiSpacer />
          <UserManagement/>
        </>
      ),
    },
    {
      id: 'pull-configuration--id',
      name: 'Pull Configuration',
      content: (
        <>
          <EuiSpacer />
          <PullConfiguration/>
        </>
      ),
    },
  ];

  return <div>
    <EuiSpacer />
    <EuiFlexGroup justifyContent="center">
      <EuiFlexItem style={{ maxWidth: 1366 }}>
        <EuiPanel paddingSize="m">
          <EuiTabbedContent
            tabs={tabs}
            initialSelectedTab={tabs[0]}
            autoFocus="selected"/>
        </EuiPanel>
      </EuiFlexItem>
    </EuiFlexGroup>
  </div>

}

export default SFTP;
