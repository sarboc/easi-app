import React from 'react';
import { Security } from '@okta/okta-react';
import { detect } from 'detect-browser';

// This can do anything. It doesn't have to redirect
// It can be a pop up modal, alert message, etc.
function onAuthRequired({ history }: any): void {
  history.push('/login');
}

type AuthenticationWrapperProps = {
  children: React.ReactNode;
};

const AuthenticationWrapper = ({ children }: AuthenticationWrapperProps) => {
  const browser: any = detect();
  const enablePkce =
    process.env.REACT_APP_ENVIRONMENT !== 'local' ||
    (browser && browser.name !== 'ie');
  return (
    <Security
      issuer={process.env.REACT_APP_OKTA_ISSUER}
      clientId={process.env.REACT_APP_OKTA_CLIENT_ID}
      redirectUri={process.env.REACT_APP_OKTA_REDIRECT_URI}
      onAuthRequired={onAuthRequired}
      responseType={['code']}
      pkce={enablePkce}
    >
      {children}
    </Security>
  );
};

export default AuthenticationWrapper;
