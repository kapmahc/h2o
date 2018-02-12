import {createStore, combineReducers} from "redux";
import withRedux from "next-redux-wrapper";

import reducers from './reducers';

const reducer = combineReducers({
  ...reducers
});

export default(initialState) => {
  return createStore(reducer, initialState);
};
