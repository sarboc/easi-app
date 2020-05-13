import React from 'react';
import Header from 'components/Header';
import UpcomingActions from 'components/shared/UpcomingActions';
import { SecureRoute, useOktaAuth } from '@okta/okta-react';
import ActionBanner from '../../components/shared/ActionBanner/index';

// This view can be deleted whenever we're ready
// This is just a sandbox page for us to test things out

const onButtonClick = (authState: any) => {
  alert(`I WAS CLICKED! ${authState}`)
}

const Sandbox = () => {
  console.log(useOktaAuth);
  console.log(SecureRoute);

  // const foo: any = useOktaAuth()
  const authState: string = "Hello"
  // const { authState }: { authState: object } = useOktaAuth();

  return (
    <div>
      <Header />
      <div className="grid-container">
        <h1>Sandbox</h1>
        <UpcomingActions timestamp="FAKE TIME">
          <ActionBanner
            title="thing"
            helpfulText="lots of helpful text"
            label="I am a button"
            onClick={() => onButtonClick(authState)}
          />
        </UpcomingActions>
      </div>
    </div>
  );
};

export default Sandbox;
