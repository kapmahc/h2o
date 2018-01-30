import React from 'react';
import {Container, Row, Col} from 'reactstrap';
import withRedux from 'next-redux-wrapper';

import makeStore from '../../store';
import Layout from '../../layouts/application';
import SharedLinks from '../../components/UsersSharedLinks';

class Widget extends React.Component {
  render() {
    return (<Layout>
      <br/>
      <Container>
        <Row>
          <Col>
            <h2>unlock
            </h2>
            <hr/>
            <form></form>
            <SharedLinks/>
          </Col>
        </Row>
      </Container>
    </Layout>)
  }
};

export default withRedux(makeStore, (state) => ({}))(Widget);
