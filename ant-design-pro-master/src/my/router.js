import React from 'react';
import {routerRedux, Route, Switch} from 'dva/router';
import {LocaleProvider, Spin} from 'antd';
import dynamic from 'dva/dynamic';
import {addLocaleData, IntlProvider} from 'react-intl';

import {getRouterData} from './common/router';
import Authorized from './utils/Authorized';
import styles from '../index.less';
import {get as getLocale} from './locales';

const {ConnectedRouter} = routerRedux;
const {AuthorizedRoute} = Authorized;
dynamic.setDefaultLoadingComponent(() => {
  return <Spin size="large" className={styles.globalSpin}/>;
});

function RouterConfig({history, app}) {
  const routerData = getRouterData(app);
  const BasicLayout = routerData['/'].component;

  const user = getLocale();
  addLocaleData(user.data);

  return (<LocaleProvider locale={user.antd}>
    <IntlProvider locale={user.locale} messages={user.messages}>
      <ConnectedRouter history={history}>
        <Switch>
          <Route path="/users/sign-in" component={routerData['/users/sign-in'].component}/>
          <AuthorizedRoute path="/" render={props => <BasicLayout {...props}/>} redirectPath="/users/sign-in"/>
        </Switch>
      </ConnectedRouter>
    </IntlProvider>
  </LocaleProvider>);
}

export default RouterConfig;
