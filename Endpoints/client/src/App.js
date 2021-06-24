import React, { useRef } from 'react';
import { Router, Route, Switch } from 'react-router-dom';
import LoginPage from './pages/Login';
import ChatRoom from './pages/Chat'
import browserHistory from './historial';

class ChatError extends Error {
    constructor(message) {
        super(message);
        this.name = "ChatError";
        this.message = message;
    }
}

function App() {
    const user_id = useRef("");

    const updateUserId = new_id => {
        if (user_id.current === "") {
            user_id.current = new_id;
        } else if(new_id !== user_id.current){
            throw new ChatError(`attempt to change a user_id which was already defined: '${user_id.current}' to '${new_id}'`);
        }
    }

    const getUserId = () => user_id.current;

    return(
        <div className="App">
            <Router history={browserHistory}>
                <Switch>
                    <Route exact path="/" render={props => <LoginPage {...props} userIdCallback={updateUserId} />}/>
                    <Route exact path="/chat"  render={props => <ChatRoom {...props} userIdCallback={getUserId} />}/>
                </Switch>
            </Router>
        </div>
    )
}

export default App;
