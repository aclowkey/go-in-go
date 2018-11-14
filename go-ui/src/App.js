import React, { Component } from 'react';
import Board from './Board.js'

class App extends Component {
  render() {
    return (
      <div>
        <h1>Let's Go! </h1>
        <Board size="5"></Board>
      </div>
    );
  }
}

export default App;
