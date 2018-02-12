import React from 'react';
import {Container} from 'reactstrap';

export default class Widget extends React.Component {
  render() {
    return (<Container>
      <hr/>
      <footer>
        <p className="float-right">
          <a href="#">Back to top</a>
        </p>
        <p>&copy; 2017-2018 Company, Inc. &middot;
          <a href="#">Privacy</a>
          &middot;
          <a href="#">Terms</a>
        </p>
      </footer>
    </Container>)
  }
}
