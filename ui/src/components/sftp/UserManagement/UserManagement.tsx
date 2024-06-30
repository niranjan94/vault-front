import React, { useState } from 'react';
import { EuiBasicTable, EuiHealth, EuiLink, EuiPopover, EuiText, EuiButton, EuiSpacer } from '@elastic/eui';
import { EuiBasicTableColumn } from "@elastic/eui/src/components/basic_table/basic_table";
import ClickToCopy from '../../utils/ClickToCopy/ClickToCopy';


function MultiValueCell({ values, label }: { values: string[], label: string }) {
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);
  const onButtonClick = () => setIsPopoverOpen(isPopoverOpen => !isPopoverOpen);
  const closePopover = () => setIsPopoverOpen(false);

  const listItems = values.map(value => <div key={value} style={{ marginTop: 6}}>
    <ClickToCopy text={value}/>
  </div>)

  return <EuiPopover
    button={<EuiLink onClick={onButtonClick}>{label}</EuiLink>}
    isOpen={isPopoverOpen}
    closePopover={closePopover}>
    <EuiText>
      {listItems}
    </EuiText>
  </EuiPopover>
}


const columns: Array<EuiBasicTableColumn<any>> = [
  {
    field: 'username',
    name: 'Username',
    render: (value: string) => <ClickToCopy text={value}/>
  },
  {
    field: 'group',
    name: 'Group',
  },
  {
    field: 'publicKeys',
    name: 'Public Keys',
    render: (value: any) => {
      return <MultiValueCell values={value} label={`${value.length} keys`}/>
    }
  },
  {
    field: 'ipAddresses',
    name: 'IP Addresses',
    render: (value: any) => {
      return <MultiValueCell values={value} label={`${value.length} addresses`}/>
    }
  },
  {
    field: 'directories',
    name: 'Directories',
    render: (value: any) => {
      return <MultiValueCell values={value} label={`${value.length} directories`}/>
    }
  },
  {
    field: 'home',
    name: 'Home',
    render: (value: string) => <ClickToCopy text={value}/>
  },
  {
    field: 'enabled',
    name: 'Enabled',
    dataType: 'boolean',
    render: (enabled: boolean) => {
      const color = enabled ? 'success' : 'danger';
      const label = enabled ? 'Enabled' : 'Disabled';
      return <EuiHealth color={color}>{label}</EuiHealth>;
    },
  },
  {
    name: 'Actions',
    actions: [
      {
        name: 'View',
        description: 'View configuration',
        icon: 'eye',
        color: 'success',
        type: 'icon',
        onClick: () => {},
        isPrimary: true,
      },
      {
        name: 'Edit',
        isPrimary: true,
        description: 'Edit this user',
        icon: 'pencil',
        type: 'icon',
        onClick: () => {},
      },
      {
        name: 'Delete',
        description: 'Delete this user',
        icon: 'trash',
        color: 'danger',
        type: 'icon',
        onClick: () => {},
      },
    ]
  },
];

function UserManagement() {
  const items = [
    {
      username: 'sandbox',
      group: 'shared',
      publicKeys: [
        'ssh-rsa abcdefghijklmnopqrstuvwxyz123345',
        'ssh-rsa abcdefghijklmnopqrstuvwxyz12334',
        'ssh-rsa abcdefghijklmnopqrstuvwxyz1233',
        'ssh-rsa abcdefghijklmnopqrstuvwxyz12',
      ],
      ipAddresses: [
        '127.0.0.1/32',
        '127.0.2.1/32',
        '127.0.3.1/32',
      ],
      directories: [
        'rw:intesa',
        'rwd:capital-company/efg'
      ],
      home: 'rwd:sandbox',
      enabled: true
    }
  ]
  return <div>
    <EuiButton
      color="secondary"
      size="s"
      iconSide="left"
      iconType="user"
      onClick={() => window.alert('Button clicked')}>
      Add new user
    </EuiButton>
    <EuiSpacer size="l"/>
    <EuiBasicTable
      items={items}
      rowHeader="firstName"
      columns={columns}
    />
  </div>
}

export default UserManagement;
