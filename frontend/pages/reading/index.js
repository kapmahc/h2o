import React from 'react';
import {Container, Row, Col} from 'reactstrap';
import withRedux from 'next-redux-wrapper';

import makeStore from '../../store';
import Layout from '../../layouts/application';

class Widget extends React.Component {
  render() {
    return (<Layout>
      <br/>
      <Container>
        <Row>
          <Col>
            <h2>reading
            </h2>
            <hr/>
          </Col>
        </Row>
      </Container>
    </Layout>)
  }
};

export default withRedux(makeStore, (state) => ({}))(Widget);
