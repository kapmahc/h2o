import React from 'react';
import {Container, Row, Col} from 'reactstrap';

import Layout from '../../layouts/application';
import SharedLinks from '../../components/UsersSharedLinks';

export default class Widget extends React.Component {
  render() {
    return (<Layout>
      <br/>
      <Container>
        <Row>
          <Col>
            <h2>sign in</h2>
            <hr/>
            <form></form>
            <SharedLinks/>
          </Col>
        </Row>
      </Container>
    </Layout>)
  }
}
