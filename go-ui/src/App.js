import React, { Component } from 'react';
import Cell from './Board.js'

class App extends Component {
  render() {
    return (
      <div>
        <h1>Let's Go! </h1>
        <Cell piece="White"></Cell>
      </div>
    );
  }
}

export default App;
