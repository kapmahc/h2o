import React, {Component} from 'react';
import {connect} from 'dva';
import {Link} from 'dva/router';
import {Checkbox, Alert, Icon} from 'antd';

import Login from '../../../components/Login';
import styles from '../../../routes/User/Login.less';
import {injectIntl, intlShape, FormattedMessage} from 'react-intl';

const {
  Tab,
  UserName,
  Password,
  Mobile,
  Captcha,
  Submit
} = Login;

@connect(({login, loading}) => ({login, submitting: loading.effects['login/login']}))
export default class LoginPage extends Component {
  state = {
    type: 'account',
    autoLogin: true
  }

  onTabChange = (type) => {
    this.setState({type});
  }

  handleSubmit = (err, values) => {
    const {type} = this.state;
    if (!err) {
      this.props.dispatch({
        type: 'login/login',
        payload: {
          ...values,
          type
        }
      });
    }
  }

  changeAutoLogin = (e) => {
    this.setState({autoLogin: e.target.checked});
  }

  renderMessage = (content) => {
    return (<Alert style={{
        marginBottom: 24
      }} message={content} type="error" showIcon="showIcon"/>);
  }

  render() {
    const {login, submitting} = this.props;
    const {type} = this.state;
    return (<div className={styles.main}>
      <Login defaultActiveKey={type} onTabChange={this.onTabChange} onSubmit={this.handleSubmit}>
        <Tab key="account" tab={<FormattedMessage id = "users.sign-in.by-account" />}>
          {login.status === 'error' && login.type === 'account' && !login.submitting && this.renderMessage('账户或密码错误（admin/888888）')}
          <UserName name="userName" placeholder="admin/user"/>
          <Password name="password" placeholder="888888/123456"/>
        </Tab>
        <Tab key="mobile" tab={<FormattedMessage id = "users.sign-in.by-phone" />}>
          {login.status === 'error' && login.type === 'mobile' && !login.submitting && this.renderMessage('验证码错误')}
          <Mobile name="mobile"/>
          <Captcha name="captcha"/>
        </Tab>
        <div>
          <Checkbox checked={this.state.autoLogin} onChange={this.changeAutoLogin}>
            <FormattedMessage id="users.sign-in.remember-me"/>
          </Checkbox>
          <Link style={{
              float: 'right'
            }} to="/users/forgot-password"><FormattedMessage id="users.forgot-password.title"/></Link>
        </div>
        <Submit loading={submitting}><FormattedMessage id="users.sign-in.submit"/></Submit>
        <div className={styles.other}>
          其他登录方式
          <Icon className={styles.icon} type="alipay-circle"/>
          <Icon className={styles.icon} type="taobao-circle"/>
          <Icon className={styles.icon} type="weibo-circle"/>
          <Link className={styles.register} to="/users/sign-up">注册账户</Link>
        </div>
      </Login>
    </div>);
  }
}
