import React from 'react';
import {ListGroup, ListGroupItem} from 'reactstrap';
import Link from 'next/link';

export default class Widget extends React.Component {
  render() {
    return (<ListGroup>
      {
        [
          {
            href: '/users/sign-in',
            label: 'sign in'
          }, {
            href: '/users/sign-up',
            label: 'sign up'
          }, {
            href: '/users/forgot-password',
            label: 'forgot password'
          }, {
            href: '/users/confirm',
            label: 'confirm'
          }, {
            href: '/users/unlock',
            label: 'unlock'
          }
        ].map((l, i) => (<Link href={l.href} key={i}>
          <ListGroupItem tag="a" action={true}>
            {l.label}
          </ListGroupItem>
        </Link>))
      }
    </ListGroup>)
  }
}
