import { EuiButton, EuiComboBox, EuiDescriptionList, EuiPanel } from '@elastic/eui';
import React from 'react';
import ClickToCopy from '../../utils/ClickToCopy/ClickToCopy';

const options = [
  {
    label: 'Tethys'
  },
  {
    label: 'Hyperion'
  }
];

const favoriteVideoGame = [
  {
    title: 'Username',
    description: <ClickToCopy text="v-userpass-ml-devel-nXw1vRYpxA8eBzcHuRvr-1595272732"/>,
  },
  {
    title: 'Password',
    description: <ClickToCopy text="A1a-240Zi001BAfb6iF3"/>,
  },
  {
    title: 'Validity',
    description: '1 day',
  }
];


function DatabaseCredentials() {

  return <EuiPanel paddingSize="m" betaBadgeLabel={'Windows Credentials'}>
    <br/>
    <EuiComboBox
      placeholder="Select a role"
      singleSelection={{asPlainText: true}}
      options={options}
      compressed
      fullWidth
      isClearable={false}/>
    <br/>
    <EuiButton fullWidth color="secondary" size="s">
      Get credentials
    </EuiButton>
    <br/>
    <EuiDescriptionList listItems={favoriteVideoGame} />


  </EuiPanel>;
}

export default DatabaseCredentials;
