import '@babel/polyfill';
import 'url-polyfill';
import FastClick from 'fastclick';

import registerServiceWorker from './registerServiceWorker'
import main from './main'

main('root')
FastClick.attach(document.body);
registerServiceWorker()
