import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import CircleDiagram from './modules/CircleDiagram';
import CircleDiagramResourceSidebar from './modules/CircleDiagram/ResourceSidebar';
import Circles from './modules/Circles';
import Dashboard from './modules/Dashboard/ index';
import Modules from './modules/Modules';
import reportWebVitals from './reportWebVitals';
import Main from './modules/Main';
import './core/components/icons/library'
import './index.css'
import { createTheme, CssBaseline, ThemeProvider } from '@mui/material';

const darkTheme = createTheme({
  palette: {
    mode: 'dark',
  },
});

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <ThemeProvider theme={darkTheme}>
      <CssBaseline />
      <BrowserRouter>
        <Routes>
          <Route path='/' element={<Main />}>
            <Route index element={<Dashboard />} />
            <Route path='circles' element={<Circles />} />
            <Route path="circles/:circle" element={<CircleDiagram />}>
              <Route path="namespaces/:namespace/ref/:ref/kind/:kind/resource/:resource" element={<CircleDiagramResourceSidebar />} />
            </Route>
            <Route path='modules' element={<Modules />} />
          </Route>
          
        </Routes>
      </BrowserRouter>
    </ThemeProvider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
