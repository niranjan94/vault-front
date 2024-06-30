import { EuiCode, EuiCopy } from '@elastic/eui';
import React from 'react';

type ClickToCopyProps = {
  text: string
}

function ClickToCopy({ text }: ClickToCopyProps) {
  return <EuiCopy textToCopy={text} beforeMessage="Click to copy">
    {copy => <EuiCode style={{ cursor: "pointer" }} onClick={copy}>{text}</EuiCode>}
  </EuiCopy>
}

export default ClickToCopy;
