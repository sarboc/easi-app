// src/OktaSignInWidget.js

import React, { Component } from 'react';
import OktaSignIn from '@okta/okta-signin-widget/dist/js/okta-sign-in.min';
import { detect } from 'detect-browser';

type OktaSignInWidgetProps = {
  onSuccess: (auth: any) => any;
  onError: () => void;
};

export default class OktaSignInWidget extends Component<
  OktaSignInWidgetProps,
  {}
> {
  widget: any;

  componentDidMount() {
    this.widget = new OktaSignIn({
      baseUrl: process.env.REACT_APP_OKTA_DOMAIN,
      authParams: { pkce: this.enablePkce() },
      el: '#sign-in-widget'
    });
    this.widget.showSignInToGetTokens({
      authorizationServerId: process.env.REACT_APP_OKTA_SERVER_ID,
      clientId: process.env.REACT_APP_OKTA_CLIENT_ID,
      redirectUri: process.env.REACT_APP_OKTA_REDIRECT_URI,
      scope: 'openid profile email'
    });
  }

  componentWillUnmount() {
    this.widget.remove();
  }

  enablePkce = () => {
    const browser: any = detect();

    return (
      process.env.REACT_APP_ENVIRONMENT !== 'local' || browser.name !== 'ie'
    );
  };

  render() {
    return (
      <div>
        <div id="sign-in-widget" />
      </div>
    );
  }
}
