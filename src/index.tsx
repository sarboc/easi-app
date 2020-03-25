import React from 'react';
import ReactDOM from 'react-dom';
import { Provider } from 'react-redux';
import { TextEncoder } from 'text-encoding';
import 'uswds';
import App from './views/App';
import store from './store';
import './index.scss';
import * as serviceWorker from './serviceWorker';

if (process.env.NODE_ENV === 'development') {
  import('react-axe').then(axe => {
    axe.default(React, ReactDOM, 1000);
  });
}

if (typeof (window as any).TextEncoder === 'undefined') {
  (window as any).TextEncoder = TextEncoder;
}

ReactDOM.render(
  <Provider store={store}>
    <App />
  </Provider>,
  document.getElementById('root')
);

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
