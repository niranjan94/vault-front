import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import * as serviceWorker from './serviceWorker';
import store from './redux/store';
import { Provider } from 'react-redux'
import { BrowserRouter } from 'react-router-dom';

const render = (AppComponent: any) => {
  return ReactDOM.render(
    <React.StrictMode>
      <Provider store={store}>
        <BrowserRouter>
          <AppComponent />
        </BrowserRouter>
      </Provider>
    </React.StrictMode>,
    document.getElementById('root')
  );
};

render(App);

if ((module as any).hot) {
  (module as any).hot.accept('./App', () => {
    const NextApp = require('./App').default;
    render(NextApp);
  });
}


// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
