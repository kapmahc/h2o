import jwtDecode from 'jwt-decode'
import moment from 'moment'

import {USERS_SIGN_IN, USERS_SIGN_OUT, SITE_REFRESH} from './actions'
import {setToken} from './auth'

const currentUser = (state = {}, action) => {
  switch (action.type) {
    case USERS_SIGN_IN:
      try {
        var it = jwtDecode(action.token);
        if (moment().isBetween(moment.unix(it.nbf), moment.unix(it.exp))) {
          setToken(action.token)
          return {admin: it.admin, uid: it.uid}
        }
        setToken()
      } catch (e) {
        console.error(e)
      }
      return {}
    case USERS_SIGN_OUT:
      setToken()
      return {}
    default:
      return state
  }
}

const siteInfo = (state = {
  languages: [],
  header: [],
  footer: []
}, action) => {
  switch (action.type) {
    case SITE_REFRESH:
      // set favicon
      var link = document.querySelector("link[rel*='icon']") || document.createElement('link');
      link.type = 'image/x-icon';
      link.rel = 'shortcut icon';
      link.href = action.info.favicon;
      document.getElementsByTagName('head')[0].appendChild(link);

      return Object.assign({}, action.info)
    default:
      return state;
  }
}

export default {currentUser, siteInfo}
