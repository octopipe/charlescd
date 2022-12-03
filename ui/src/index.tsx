import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Route, Routes, useNavigate } from 'react-router-dom';
import { Provider as ReduxProvider } from 'react-redux'
import reportWebVitals from './reportWebVitals';
import Main from './modules/Main';
import { Provider as FetchProvider, IncomingOptions } from 'use-http'
import Login from './modules/Login';
import './core/components/icons/library'
import 'bootstrap/dist/css/bootstrap.min.css';
import './index.scss'
import Modules from './modules/Modules';
import Home from './modules/Home';
import store from './store'
import Diagram from './modules/CirclesMain/Diagram';
import CreateCircle from './modules/CreateCircle';
import Error from './modules/Error';
import CirclesMain from './modules/CirclesMain';


const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);

const App = () => {
  const options: IncomingOptions = {
    interceptors: {
      request: async ({ options, url, path, route }) => {
        return options
      },
      response: async ({ response }) => {
        const res = response

        if (response.status >= 500) {
          window.location.href = '/error'
        }
        return res
      }
    }
  }

  return (
    <React.StrictMode>
      <ReduxProvider store={store}>
        <FetchProvider url='http://localhost:8080' options={options}>
          <BrowserRouter>
            <Routes>
              <Route path='/login' element={<Login />} />
              <Route path='' element={<Main />}>
                <Route path='workspaces/:workspaceId' element={<Home />} />
                <Route path='workspaces/:workspaceId/circles' element={<CirclesMain />}>
                </Route>
                <Route path='workspaces/:workspaceId/circles/create' element={<CreateCircle />} />
                <Route path='workspaces/:workspaceId/modules' element={<Modules />} />
              </Route>
              <Route path='/error' element={<Error />} />
            </Routes>
            
          </BrowserRouter>
        </FetchProvider>
      </ReduxProvider>
    </React.StrictMode>
  )
}

root.render(<App />);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
