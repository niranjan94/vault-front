import React from 'react';
import './App.scss';
import { EuiHeader, EuiHeaderLink, EuiHeaderLinks, EuiHeaderLogo } from '@elastic/eui';
import { useDispatch, useSelector } from 'react-redux';
import { isAuthenticatedSelector } from './redux/auth/selectors';
import { Route, Redirect, RouteProps, Switch } from 'react-router-dom';
import Login from './pages/Login/Login';
import { logout } from './redux/auth/slice';
import Credentials from './pages/Credentials/Credentials';
import EuiRouterHeaderLink from './components/utils/EuiRouterHeaderLink/EuiRouterHeaderLink';
import SFTP from './pages/SFTP/SFTP';

function App() {

  const isAuthenticated = useSelector(isAuthenticatedSelector);
  const dispatch = useDispatch();

  const onLogout = () => {
    dispatch(logout());
  }

  const headerSections = [
    <EuiHeaderLogo
      iconType="securityApp"
      href="/"
      aria-label="Goes to home">
      Vault
    </EuiHeaderLogo>,
    <EuiHeaderLinks aria-label="App navigation dark theme example">
      <EuiRouterHeaderLink to="/">Credentials</EuiRouterHeaderLink>
      <EuiRouterHeaderLink to="/sftp">SFTP</EuiRouterHeaderLink>
    </EuiHeaderLinks>
  ];

  const headerRightSections = [
    <EuiHeaderLink iconType="push" href="#" onClick={onLogout}>
      Logout
    </EuiHeaderLink>
  ];

  return (
    <div className="App">
      {
        isAuthenticated ?
          <EuiHeader theme="dark"
                     sections={[
                       {items: headerSections, borders: 'right'},
                       {items: headerRightSections, borders: 'none'}
                     ]}>
          </EuiHeader> : null
      }
      <div className="container">
        <Switch>
          <PrivateRoute exact path="/">
            <Credentials/>
          </PrivateRoute>
          <PrivateRoute exact path="/sftp">
            <SFTP/>
          </PrivateRoute>
          <GuestOnlyRoute exact path="/login">
            <Login/>
          </GuestOnlyRoute>
          <GuestOnlyRoute exact path="/saml/callback">
            <Login isLoading={true} samlCallback={true}/>
          </GuestOnlyRoute>
        </Switch>
      </div>
    </div>
  );
}

// @ts-ignore
function PrivateRoute({ children, ...rest }: RouteProps) {
  const isAuthenticated = useSelector(isAuthenticatedSelector);
  return (
    <Route
      {...rest}
      render={({ location }) =>
        isAuthenticated ? (
          children
        ) : (
          <Redirect
            to={{
              pathname: '/login',
              state: { from: location },
            }}
          />
        )
      }
    />
  );
}

function GuestOnlyRoute({ children, ...rest }: RouteProps) {
  const isAuthenticated = useSelector(isAuthenticatedSelector);
  return (
    <Route
      {...rest}
      render={() =>
        !isAuthenticated ? children : <Redirect to="/" />
      }
    />
  );
}


export default App;
