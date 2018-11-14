import React, { Component } from 'react';
import './Board.css';

class Cell extends Component {
    constructor(props){
        super(props)
        this.state = {
            piece: this.props.piece
        }
    }
    render(){
        return (
            <div class="square">{this.state.piece}</div>
        )
    }
}

class Board extends Component {
    constructor(props){
        super(props);
        this.state = {
            size: this.props.size,
            data: emptySquare(this.props.size)
        }
    }
    render(){
        return (
           <ul>
               <li><Cell piece="White"></Cell></li>
               <li><Cell piece="Black"></Cell></li>
               <li><Cell piece="Green"></Cell></li>
           </ul> 
        )
    }
}


function emptySquare(size){
    return new Array(size).fill(0).map(() => new Array(size).fill(0));
}

export default Board