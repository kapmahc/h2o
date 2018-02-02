import React from 'react'
import {Route} from 'react-router'
import Loadable from 'react-loadable';

import nut from './nut'
import forum from './forum'
import survey from './survey'
import reading from './reading'
import mall from './mall'
import pos from './pos'
import erp from './erp'
import ops_vpn from './ops/vpn'
import ops_mail from './ops/mail'

const plugins = [
  nut,
  forum,
  survey,
  reading,
  mall,
  pos,
  erp,
  ops_vpn,
  ops_mail
]

const dynamicWrapper = (w) => Loadable({
  loader: () => w,
  loading: () => <div>Loading...</div>
});

export default {
  routes: plugins.reduce((ar, it) => ar.concat(it.routes), []).map((it) => {
    return (< Route key = {
      it.path
    }
    exact = {
      true
    }
    path = {
      it.path
    }
    component = {
      dynamicWrapper(it.component)
    } />)
  }).concat([<Route key="not-found" component={dynamicWrapper(import ('./NotFound'))}/>])
}
