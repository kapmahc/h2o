import React from 'react';
import {
  Collapse,
  Navbar,
  NavbarToggler,
  NavbarBrand,
  Nav,
  NavItem,
  NavLink,
  UncontrolledDropdown,
  DropdownToggle,
  DropdownMenu,
  DropdownItem
} from 'reactstrap';
import {connect} from 'react-redux';
import Link from 'next/link';
import flush from 'styled-jsx/server';

import {backend} from '../../utils';

class Widget extends React.Component {
  state = {
    isOpen: false
  }
  // componentDidMount() {
  //   const res = await fetch(backend('/layout'));
  //   const siteInfo = await res.json();
  //    refresh(siteInfo);
  //   console.log(siteInfo);
  //   return {siteInfo};
  //    const {signIn, refresh, info, user} = this.props
  //    if (!user.uid) {
  //      var tkn = token()
  //      if (tkn) {
  //        signIn(tkn)
  //      }
  //    }
  //    if (info.languages.length === 0) {
  //      get('/layout').then((rst) => refresh(rst)).catch(message.error)
  //    }
  // }
  onToggle = () => {
    this.setState({
      isOpen: !this.state.isOpen
    });
  }
  render() {
    return (<header>
      <Navbar color="dark" fixed="top" dark={true} expand="md">
        <Link href="/">
          <NavbarBrand>reactstrap</NavbarBrand>
        </Link>
        <NavbarToggler onClick={this.onToggle}/>
        <Collapse isOpen={this.state.isOpen} navbar={true}>
          <Nav className="ml-auto" navbar={true}>
            <NavItem>
              <NavLink href="/components/">Components</NavLink>
            </NavItem>
            <NavItem>
              <NavLink href="https://github.com/reactstrap/reactstrap">Github</NavLink>
            </NavItem>
            <UncontrolledDropdown nav={true}>
              <DropdownToggle nav={true} caret={true}>
                Options
              </DropdownToggle>
              <DropdownMenu >
                <DropdownItem>
                  Option 1
                </DropdownItem>
                <DropdownItem>
                  Option 2
                </DropdownItem>
                <DropdownItem divider={true}/>
                <DropdownItem>
                  Reset
                </DropdownItem>
              </DropdownMenu>
            </UncontrolledDropdown>
          </Nav>
        </Collapse>
      </Navbar>
    </header>);
  }
};

export default connect(state => ({user: state.currentUser, site: state.siteInfo}))(Widget);
