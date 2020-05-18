import React from 'react';
import Header from 'components/Header';
import { useOktaAuth } from '@okta/okta-react';
import { DateTime } from 'luxon';
import ActionBanner from '../../components/shared/ActionBanner/index';

// This view can be deleted whenever we're ready
// This is just a sandbox page for us to test things out

const onButtonClick = async (authService: any) => {
  console.log('Fetching current token')
  const tokenManager = await authService.getTokenManager();
  // const sessionToken = tokenManager.get('sessionToken');
  // console.log(sessionToken)
  const idToken = await tokenManager.get('idToken');
  console.log(idToken)
  const token = await tokenManager.get('accessToken');
  console.log(token)
  const { expiresAt }: { expiresAt: number } = token;
  const expirationDate = DateTime.fromSeconds(expiresAt).toLocaleString(DateTime.DATETIME_FULL);
  console.log('expirationDate ', expirationDate);

  console.log('Renewing token')
  const newToken = await tokenManager.renew('accessToken')
  const { expiresAt: newExpiresAt }: { expiresAt: number } = newToken;
  const newExpirationDate = DateTime.fromSeconds(newExpiresAt).toLocaleString(DateTime.DATETIME_FULL);
  console.log('newExpirationDate ', newExpirationDate);
}

const Sandbox = () => {
  const { authService}: { authService: any } = useOktaAuth();

  return (
    <div>
      <Header />
      <div className="grid-container">
        <h1>Sandbox</h1>
          <ActionBanner
            title="thing"
            helpfulText="lots of helpful text"
            label="I am a button"
            onClick={() => onButtonClick(authService)}
          />
      </div>
    </div>
  );
};

export default Sandbox;
