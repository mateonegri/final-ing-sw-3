import logo from './logo.svg';
import React from 'react';
import './App.css';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import HomePage from './pages/home_page';
import ModifyUser from './pages/modify_user';
import AddUser from './pages/add_user';

function App() {
  return (
    <div className="App">
      <Router>
        <Routes>
          <Route path='/' element={<HomePage />}></Route>
          <Route path='/modify-user/:id' element={<ModifyUser />}></Route>
          <Route path='/add-user' element={<AddUser />}></Route>
        </Routes>
      </Router>
    </div>
  );
}

export default App;
