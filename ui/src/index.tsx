import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Route, Routes, useNavigate } from 'react-router-dom';
import { Provider as ReduxProvider } from 'react-redux'
import reportWebVitals from './reportWebVitals';
import Main from './modules/Main';
import Login from './modules/Login';
import Home from './modules/Home';
import store from './store'
import Error from './modules/Error';
import Circles from './modules/Circles';
import { ToastContainer, toast } from "react-toastify";
import ModulesMain from './modules/Modules';
import Root from './modules/Root';
import Workspaces from './modules/Workspaces';
import { ROUTES } from './core/constants/routes';

import './core/components/icons/library'
import 'bootstrap/dist/css/bootstrap.min.css';
import 'react-toastify/dist/ReactToastify.css';
import './core/style/base.scss'
import GlobalError from './modules/GlobalError';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const App = () => {
  return (
    <React.StrictMode>
      <ReduxProvider store={store}>
        <BrowserRouter>
          <ToastContainer autoClose={2000} hideProgressBar theme='dark'/>
          <Routes>
            <Route path={ROUTES.LOGIN} element={<Login />} />
            <Route path={ROUTES.ROOT} element={<Root />}>
              <Route path='' element={<Workspaces />} />
            </Route>
            <Route path={ROUTES.MAIN} element={<Main />}>
              <Route path='' element={<Home />} />
              <Route path='circles' element={<Circles />} />
              <Route path='modules' element={<ModulesMain />} />
            </Route>
            <Route path='/error' element={<Error />} />
          </Routes>
        </BrowserRouter>
      </ReduxProvider>
    </React.StrictMode>
  )
}

root.render(<App />);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
