// class TrainsmiterException extends Error {
//     constructor(message) {
//         super(message)
//         this.message = message
//         this.name = "TransmiterException"
//     }
// }

class Trainsmiter {

    constructor() {
        this.status = 0;
        this.host = null;
        this.is_connected = false;
        this.socket = null;
        this.disconnection_callback = undefined;
        this.connection_callback = undefined;
        this.message_callback = undefined;
        this.error_callback = undefined;
    }

    isConnected = () => this.status === 1 ? true : false;

    connect = (host, disconnection_callback) => {
        this.host = host;
        this.socket = new WebSocket(host);
        this.socket.onopen = this.onConnect;
        this.socket.onmessage = this.onMessage;
        this.socket.onclose = this.onDisconnect;
    }



    setDisconnectionCallback = callback => {
        if (!this.is_connected) {
            this.disconnection_callback = callback;
        }
    }

    setConnectionCallback = callback => {
        if (!this.is_connected) {
            this.connection_callback = callback;
        }
    }

    onConnect = e => {
        if(this.connection_callback !== undefined) {
            this.connection_callback(e);
        }
    }

    onDisconnect = e => {
        if(this.disconnection_callback !== undefined) {
            this.disconnection_callback(e);
        }
    }

    Close = () => this.socket.close(1000, "User logged out")

    onMessage = message_event => {
        if (this.message_callback !== undefined) {
            this.message_callback(message_event)
        } 
    }

    setMessageCallback = callback => this.message_callback = callback;
    
    onError = callback => this.error_callback = callback;

    emit = message => {
        if(message.length > 0) {
            this.socket.send(message);
        } else {
            console.warn(`Message '${message}' is invalid`);
        }
    }
}


export default Trainsmiter;