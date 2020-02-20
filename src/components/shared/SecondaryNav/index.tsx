import React from 'react';
import {Link} from 'react-router-dom';
import { ShortSystem } from 'types/systems'
import './index.scss';

type SecondaryNavProps = {
  secondaryNavList?: any[];
  activeNavItem?: string | undefined;
  searchSystem: ShortSystem;
};

const SecondaryNav = ({
  secondaryNavList = [],
  activeNavItem = '',
  searchSystem = {id: '', name: '', acronym: '', slug: '', link: '' }
}: SecondaryNavProps) => {
  let secondaryNavLength: number = secondaryNavList.length
  if (searchSystem.id && secondaryNavList.length === 10) {
    secondaryNavLength = 9
  }
  return (
    <nav aria-label="Primary navigation" className="secondary-nav">
      <div className="usa-nav__inner">
        <ul className="usa-nav__primary usa-accordion">
          {secondaryNavList.slice(0, secondaryNavLength).map(item => (
            <li
              key={item.id}
              className={`usa-nav__primary-item ${
                activeNavItem === item.slug ? 'usa-current' : ''
              }`.trim()}
              data-testid="header-nav-item"
            >
              <Link className="secondary-nav__link" to={item.link}>
                <span>{item.name}</span>
              </Link>
            </li>
          ))}
          {searchSystem.id && (
            <li
              key={searchSystem.id}
              className={`easi_search-item ${
                activeNavItem === searchSystem.slug ? 'usa-current' : ''
              }`.trim()}
              data-testid="header-nav-item"
            >
              <Link className="secondary-nav__link easi_search-link" to={searchSystem.link}>
                <span>{searchSystem.name}</span>
              </Link>
            </li>
          )}
        </ul>
      </div>
    </nav>
  );
};

export default SecondaryNav;
