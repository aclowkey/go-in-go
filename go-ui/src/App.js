import React, { Component } from 'react';
import Game from './Game.js'

class App extends Component {
  render() {
    return (
      <div>
        <Game size="9"></Game>
      </div>
    );
  }
}

export default App;
