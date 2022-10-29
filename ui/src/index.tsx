import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import { Provider as ReduxProvider } from 'react-redux'
import reportWebVitals from './reportWebVitals';
import Main from './modules/Main';
import { Provider as FetchProvider } from 'use-http'
import Login from './modules/Login';
import './core/components/icons/library'
import 'bootstrap/dist/css/bootstrap.min.css';
import './index.scss'
import Circles from './modules/Circles';
import Modules from './modules/Modules';
import Home from './modules/Home';
import store from './store'
import Circle from './modules/Circle';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <ReduxProvider store={store}>
      <FetchProvider url='http://localhost:8080'>
        <BrowserRouter>
          <Routes>
            <Route path='/login' element={<Login />} />
            <Route path='' element={<Main />}>
              <Route path='/' element={<Home />} />
              <Route path='circles' element={<Circles />} />
              <Route path='modules' element={<Modules />} />
            </Route>
            <Route path='/circles/:name' element={<Circle />}>

            </Route>
          </Routes>
        </BrowserRouter>
      </FetchProvider>
    </ReduxProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
