import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import CircleDiagram from './modules/CircleDiagram';
import CircleDiagramSidebar from './modules/CircleDiagram/Sidebar';
import Circles from './modules/Circles';
import Dashboard from './modules/Dashboard/ index';
import Modules from './modules/Modules';
import reportWebVitals from './reportWebVitals';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        <Route index element={<Dashboard />} />
        <Route path='circles' element={<Circles />} />
        <Route path="circles/:id/diagram" element={<CircleDiagram />}>
          <Route path=":object" element={<CircleDiagramSidebar />} />
        </Route>
        <Route path='modules' element={<Modules />} />
      </Routes>
    </BrowserRouter>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
