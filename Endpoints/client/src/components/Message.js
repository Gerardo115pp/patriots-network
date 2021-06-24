import React from 'react';
import fileIcon from '../icons/folder.svg';
import '../css/Message.css';
import { server_name } from '../server_data';

const Message = props => {
    let { class_name } = props;
    const { sender,
            body,
            is_control_message,
            color,
            is_file,
    } = props.message_data;
    const { download_callback } = props;

    class_name = is_control_message ? `${class_name} control-message` : class_name;

    const getFileIcon = () => {
        if (is_file) {
            return <object data={fileIcon} type="image/svg+xml">file</object>
        } else {
            return null;
        }
    }

    const isImageFile = () => /.+\.(jpg|png|gif)$/.test(body) && is_file;

    const getMessageBody = () => {
        console.log(`${body}: ${isImageFile()}`)
        if(isImageFile()) {
            let image_url = `${server_name}/static/${body}`;
            return (
                <div className="msg-image">
                    <img src={image_url} alt={image_url}/>
                </div>
            );
        } else {
            return (
                <div className="message-body">
                    <div className={`msg-body ${is_file ? "file_message" : ""}`}>
                        {body}
                    </div>
                    <div className="msg-icon">
                        {getFileIcon()}
                    </div>
                </div>
            )
        }
    }

    return(
        <div className="msg-space">
            <div onClick={is_file ? download_callback : () => true} className={`message-container ${class_name}`}>
                <div style={{color: color}} className="msg-header">
                    <span className="sender-name">{sender}</span>
                </div>
                {getMessageBody()}
                
            </div>
        </div>
    )
}

export default Message;