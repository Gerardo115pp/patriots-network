import React, { Component } from 'react';
import MessageComponent from '../components/Message';
import UserLabelComponent from '../components/UserStatus';
import Transmisor from '../types/Transmiter';
import telefakeLogo from './../imgs/Untitled-1.png';
import sendMsgLogo from './../icons/paper-plane.svg';
import { server_name, ws_server } from '../server_data';
import './../css/Chat.css';

class ChatRoom extends Component {
    constructor(props) {
        super(props);
        this.transmisor = undefined;
        this.user_name = props.location.state.user_name;
        this.state = {
            user_id: "",
            logged: false,
            messages: [],
            users: []
        };
    }

    componentDidMount() {
        if (this.state.user_id === "") {
            this.requestUsers();
            const user_id = this.props.userIdCallback();
            this.setState({
                ...this.state,
                user_id
            }, this.requestMessagesHistory);
        }
    }
    

    // componentWillUnmount() {
    //     this.transmisor.Close();
    // }

    componentDidUpdate() {
        if(this.transmisor === undefined && this.state.user_id.length >= 40) {
            this.transmisor = new Transmisor();
            this.transmisor.setConnectionCallback(() => this.toggleLogin(true));
            this.transmisor.setDisconnectionCallback(() => this.toggleLogin(false));
            this.transmisor.setMessageCallback(this.updateMessages);
            this.transmisor.connect(`${ws_server}?user_name=${this.user_name}`, null);
        }
        
    }

    getUserFile = () => {
        const file_input = document.getElementById("file-input");
        file_input.click()
    }

    getStatusElement = () => {
        let status_class = this.state.logged ? 'connected' : 'disconnected';
        return (
            <div className={`status-entry ${status_class}`}>
                {status_class}
            </div>
        )
    }

    toggleLogin = is_logged => {
        if (this.state.logged !== is_logged) {
            this.setState({
                ...this.state,
                logged: is_logged
            })
        }
    }

    isMessageDone = e => {
        if(e.key.toLowerCase() === "enter") {
            const messager = document.getElementById("transmisor-input");
            this.transmisor.emit(messager.value);
            messager.value = "";
        }
    }

    uploadFile = () => {
        const file = document.getElementById("file-input").files[0];
        if(file) {
            const forma = new FormData();
            forma.append("file", file);
            forma.append("file_name", file.name);
            forma.append("user_uuid", this.state.user_id);

            const request = new Request(`${server_name}/send-file`, {method: "POST", body: forma});
            fetch(request)
                .then(promise => promise.json())
                .then(response => {
                    if(response.response === "bad") {
                        console.warn("Server error while sending file")
                    }
                })
        }
    }

    updateMessages = message => {
        const { messages } = this.state;
        const json_message = JSON.parse(message.data);
        if(json_message.is_control_message) {
            window.setTimeout(this.requestUsers, 400);
        }
        messages.push(json_message);
        this.setState({
            ...this.state,
            messages
        }, this.scrollToBottom);
    }

    scrollToBottom = () => {
        const chat_element = document.getElementById("chat-massages-container");
        chat_element.scrollTop = chat_element.scrollHeight;
    }

    downloadFileCallback = message_data => {
        const file_name = message_data.body;
        let file_url = `${server_name}/static/${file_name}`;
        fetch(file_url)
            .then(promise => promise.blob())
            .then(file_blob => {
                let blob_url = window.URL.createObjectURL(file_blob);
                const download_element = document.createElement("a");
                download_element.href = blob_url;
                download_element.download = file_name;
                download_element.click();
                download_element.remove();
            })

    }

    requestUsers = () => fetch(`${server_name}/get-users`)
                                .then(promise => promise.json())
                                .then(response => {
                                    this.setState({
                                        ...this.state,
                                        users: response
                                    })
                                });

    requestMessagesHistory = () => {
        let api_url = `${server_name}/retrive-messages-history`;
        fetch(api_url)
            .then(promise => promise.json())
            .then(response => {
                if(response.length > 0) {
                    this.setState({
                        ...this.state,
                        messages: response
                    });
                }
            })
    }

    render() {
        const { messages, users } = this.state;
        return(
            <div id="chat-main-container">
                <header>
                    <div id="chat-logo"><img src={telefakeLogo} alt="telefake.png" /></div>
                    <h3 id="user-name">Usuario: {this.user_name}</h3>
                    <h3 id="telefake-label">Telefake</h3>
                </header>
                <main id="central-display">
                    <section id="room-content">
                        <div id="room-header">Usuarios</div>
                        <div id="room-users-container">
                            {users.map((usr, h) => {
                                if (usr.user_name === this.user_name) {
                                    return null;
                                }
                                return <UserLabelComponent key={`user-${h}`} user_data={usr}/>;
                            })}
                        </div>
                    </section>
                    <section id="chat-messages">
                        <div id="chat-massages-container">
                            {messages.map(m => <MessageComponent download_callback={() => this.downloadFileCallback(m)} key={m.uuid} class_name={ m.sender === this.user_name ? "sent" : "recivied"} message_data={m}/>)}
                        </div>
                        <div id="chat-transmisor">
                            <input onKeyDown={this.isMessageDone} type="text" id="transmisor-input" maxLength="512" placeholder="mensaje..."/>
                            <div onClick={this.getUserFile} className="circular-input-btn">
                                <object color="white" data={sendMsgLogo} type="image/svg+xml" id="circular-input-btn-svg">error</object>
                            </div>
                        </div>
                    </section>
                </main>
                <footer id="status-bar">
                    {this.getStatusElement()}
                    <div className="status-entry"><input onChange={this.uploadFile} id="file-input" type="file"/></div>
                </footer>
            </div>  
        );
    }
}

export default ChatRoom;