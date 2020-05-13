import React from 'react';
import { BrowserRouter, Switch, Route, Redirect } from 'react-router-dom';
import { SecureRoute, LoginCallback } from '@okta/okta-react';
import AuthenticationWrapper from 'views/AuthenticationWrapper';
import Home from 'views/Home';
import Login from 'views/Login';
import BusinessCase from 'views/BusinessCase';
import GRTSystemIntakeReview from 'views/GRTSystemIntakeReview';
import SystemIntake from 'views/SystemIntake';
import Sandbox from 'views/Sandbox';

import './index.scss';

type MainState = {};

type MainProps = {};

// eslint-disable-next-line react/prefer-stateless-function
class App extends React.Component<MainProps, MainState> {
  render() {
    return (
      <div>
        <div className="usa-overlay" />
        <BrowserRouter>
          <AuthenticationWrapper>
            <Switch>
              <Route path="/" exact component={Home} />
              <Route path="/login" exact component={Login} />
              <SecureRoute path="/sandbox" exact component={Sandbox} />
              <SecureRoute
                exact
                path="/system/:systemId/grt-review"
                render={({ component }: any) => component()}
                component={GRTSystemIntakeReview}
              />
              <Redirect
                exact
                from="/system/:systemId"
                to="/system/:systemId/contact-details"
              />
              <SecureRoute
                path="/system/:systemId/:formPage"
                render={({ component }: any) => component()}
                component={SystemIntake}
              />
              {/* <SecureRoute
                path="/system/all"
                exact
                component={SystemProfiles}
              /> */}
              {/* <SecureRoute
                path="/system/:profileId"
                component={SystemProfile}
              /> */}
              <SecureRoute
                path="/business/new"
                render={({ component }: any) => component()}
                component={BusinessCase}
              />
              <Route path="/implicit/callback" component={LoginCallback} />
            </Switch>
          </AuthenticationWrapper>
        </BrowserRouter>
      </div>
    );
  }
}

export default App;
