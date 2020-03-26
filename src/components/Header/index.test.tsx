import React from 'react';
import { BrowserRouter } from 'react-router-dom';
import { shallow, mount } from 'enzyme';
import { act } from 'react-dom/test-utils';
import AuthenticationWrapper from 'views/AuthenticationWrapper';

import { Header } from './index';

describe('The Header component', () => {
  it('renders without crashing', () => {
    shallow(
      <BrowserRouter>
        <AuthenticationWrapper>
          <Header />
        </AuthenticationWrapper>
      </BrowserRouter>
    );
  });

  describe('When logged out', () => {
    it('displays a login button', async done => {
      let component: any;
      await act(async () => {
        component = mount(
          <BrowserRouter>
            <AuthenticationWrapper>
              <Header />
            </AuthenticationWrapper>
          </BrowserRouter>
        );
      });
      setImmediate(() => {
        component.update();
        expect(component.text().includes('Login')).toBe(true);
        expect(component.text().includes('Logout')).toBe(false);
        done();
      });
    });
  });

  describe('When logged in', () => {
    jest.mock('@okta/okta-react', () => ({
      useOktaAuth: () => {
        return {
          authState: { isAuthenticated: true },
          authService: {
            getUser: () => ({
              name: 'hellloooo'
            })
          }
        };
      }
    }));

    xit('displays a login button', async done => {
      let component: any;
      await act(async () => {
        component = mount(
          <BrowserRouter>
            <AuthenticationWrapper>
              <Header />
            </AuthenticationWrapper>
          </BrowserRouter>
        );
      });

      setImmediate(() => {
        component.update();
        expect(component.text().includes('Logout')).toBe(true);
        expect(component.text().includes('Login')).toBe(false);
        done();
      });
    });

    xit('displays the users name', async done => {
      let component;

      await act(async () => {
        component = mount(
          <BrowserRouter>
            <Header />
          </BrowserRouter>
        );
      });

      setImmediate(() => {
        component.update();
        expect(component.text().includes('John Doe')).toBe(true);
        done();
      });
    });

    xit('displays dropdown when caret is clicked', async done => {
      let component;

      await act(async () => {
        component = mount(
          <BrowserRouter>
            <Header />
          </BrowserRouter>
        );
      });

      setImmediate(() => {
        component.update();
        expect(component.find('.user-actions-dropdown').exists()).toBe(false);
        component.find('.easi-header__caret').simulate('click');
        expect(component.find('.user-actions-dropdown').exists()).toBe(true);
        done();
      });
    });
  });

  it('displays children', async done => {
    let component: any;

    await act(async () => {
      component = mount(
        <BrowserRouter>
          <AuthenticationWrapper>
            <Header>
              <div className="test-class-name" />
            </Header>
          </AuthenticationWrapper>
        </BrowserRouter>
      );
    });

    setImmediate(() => {
      component.update();
      expect(component.find('.test-class-name').exists()).toBe(true);
      done();
    });
  });
});
