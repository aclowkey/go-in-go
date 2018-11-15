import React, { Component } from "react";
import Board from "./Board.js";

class Game extends Component {
  constructor(props) {
    super(props);
    let size = parseInt(this.props.size);
    this.state = {
      size: size,
    };
  }

  render() {
    return (
      <div>
        <h1>Let's Go!</h1>
        <Board size={this.state.size}  />
      </div>
    );
  }
}

export default Game;
