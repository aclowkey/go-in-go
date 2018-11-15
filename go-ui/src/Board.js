import React, { Component } from 'react';
import './Board.css';

class Cell extends Component {
    constructor(props){
        super(props)
        this.state = {
            piece: this.props.piece,
            onClick: this.props.onClick,
        }
    }
    render(){
        return (
            <div 
                className="square"
                onClick={this.state.onClick}>
                {this.state.piece}
            </div>
        )
    }
}

class Board extends Component {
    constructor(props){
        super(props);
        let size = parseInt(this.props.size)
        this.state = {
            size: size,
            data: emptySquare(size),
            whitesTurn: false
        }
    }

    cellClicked(x, y){
        return () => {
           let data = this.state.data.slice()
           if(data[y][x] !== 0){
                alert("You can't do that!")
                return;
           }
           data[y][x] =  this.state.whitesTurn ? "White" : "Black";
           this.setState({
                data: data,
                whitesTurn: !this.state.whitesTurn
           })
        }
    }

    render(){
        let board=[]
        for(let i = 0; i < this.state.size; i++){
            let row = [] 
            for(let j = 0; j < this.state.size; j++ ){
                let piece =this.state.data[i][j];
                row.push(<Cell key={i+","+j} piece={piece} onClick={this.cellClicked(j, i)}></Cell>);
            }

            board.push(
                <div key={i} style={{display: 'inline-block'}}>
                    {row}
                </div> 
            )
        }
        return (
            <div>
                {board}
            </div>
        )
    }

}


function emptySquare(size){
    return new Array(size).fill(0).map(() => new Array(size).fill(0));
}

export default Board