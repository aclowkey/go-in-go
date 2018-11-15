import React, { Component } from "react";
import Board from "./Board.js";

class Game extends Component {
  constructor(props) {
    super(props);
    let size = parseInt(this.props.size);
    this.state = {
      size: size,
      data: emptySquare(size),
      whitesTurn: false
    };
  }

  cellClicked(x, y) {
    return () => {
      let data = this.state.data.slice();
      if (data[y][x] !== "Empty") {
        alert("You can't do that!");
        return;
      }
      data[y][x] = this.state.whitesTurn ? "white" : "black";
      this.setState({
        data: data,
        whitesTurn: !this.state.whitesTurn
      });
    };
  }

  render() {
    return (
      <div>
        <h1>Let's Go!</h1>
        <Board size={this.state.size} data={this.state.data} />
      </div>
    );
  }
}

function emptySquare(size) {
  return new Array(size).fill("Empty").map(() => new Array(size).fill("Empty"));
}

export default Game;
