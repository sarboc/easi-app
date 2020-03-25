import React from 'react';
import { useOktaAuth } from '@okta/okta-react';
import Header from 'components/Header';

const Home = () => {
  const { authState } = useOktaAuth();

  return (
    <div>
      <Header />
      <div className="grid-container">
        <h1>Home</h1>
        <h3>{`A user is ${
          authState.isAuthenticated ? '' : 'NOT'
        } authenticated`}</h3>
      </div>
    </div>
  );
};

export default Home;
