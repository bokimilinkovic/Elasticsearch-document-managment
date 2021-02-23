import logo from './logo.svg';
import './App.css';
import React from 'react'
import {
  BrowserRouter as Router,
  Switch,
  Route 
} from 'react-router-dom'
import NewBook from './components/NewBook.jsx';
import Navbar from './components/navbar/Navbar';
import Books from './components/Books';
import Users from './components/Users';
import NewUser from './components/NewUser';

function App() {
  return (
    <div className="App">
      <Router>
        <Navbar />
        <Switch>
          <Route path="/book/create" exact>
            <NewBook />
          </Route>
          <Route path="/" exact>
            <Books />
          </Route>
          <Route path="/users" exact>
            <Users />
          </Route>
          <Route path="/users/create" exact>
            <NewUser />
          </Route>
        </Switch>
      </Router>
    </div>
  );
}

export default App;
