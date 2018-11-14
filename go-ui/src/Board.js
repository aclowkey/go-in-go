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

export default Cell