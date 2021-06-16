import React, { Component } from "react";
import "./ChatInput.scss";
import upload from '../../Images/uploadimage.svg';
import DragImage from '../DragImage'


class ChatInput extends Component {
  state = {
    files: []
  }
 handleDrop = (files) => {
    let fileList = this.state.files
    for (var i = 0; i < files.length; i++) {
      if (!files[i].name) return
      fileList.push(files[i].name)
    }
    this.setState({files: fileList})
    console.log(files)
  }
  render() {
    return (
      <div className="ChatInput"
        // onDragOver={dragOver}
        // onDragEnter={dragEnter}
        // onDragLeave={dragLeave}
        // onDrop={fileDrop}
      >
        {/* <div className="upload-image"> */}
          {/* <DragImage handleDrop = {this.handleDrop}> */}
            <input onKeyDown={this.props.send} />
            {/* <img alt="" className="upload-icon" src={upload}/> */}
          {/* </DragImage> */}
        {/* </div> */}
      </div>
    );
  }
}

export default ChatInput;