import React  from 'react';
import { EuiHeaderLink } from '@elastic/eui';
import { useHistory, useLocation } from 'react-router';
import { EuiHeaderLinkProps } from '@elastic/eui/src/components/header/header_links/header_link';

const isModifiedEvent = (event: KeyboardEvent | MouseEvent) =>
  (event.metaKey || event.altKey || event.ctrlKey || event.shiftKey);

const isLeftClickEvent = (event: MouseEvent) => event.button === 0;

const isTargetBlank = (event: any) => {
  const target = event.target.getAttribute('target');
  return target && target !== '_self';
};

type EuiRouterHeaderLink = {
  to: string
} & EuiHeaderLinkProps

function EuiRouterHeaderLink({ to, ...rest }: EuiRouterHeaderLink) {
  const history = useHistory();
  const location = useLocation();

  const onClick = (event: any) => {
    if (event.defaultPrevented) {
      return;
    }

    // Let the browser handle links that open new tabs/windows
    if (isModifiedEvent(event) || !isLeftClickEvent(event) || isTargetBlank(event)) {
      return;
    }

    // Prevent regular link behavior, which causes a browser refresh.
    event.preventDefault();

    // Push the route to the history.
    history.push(to);
  }

  // Generate the correct link href (with basename accounted for)
  const href = history.createHref({ pathname: to });
  const props = { ...rest, href, onClick };
  return <EuiHeaderLink {...props} isActive={location.pathname === to}/>;
}

export default EuiRouterHeaderLink;
