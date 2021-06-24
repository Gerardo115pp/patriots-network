import React from 'react';
import browser_history from '../historial'
import { server_name } from '../server_data'
import loginLogo from '../icons/satellite.svg';
import '../css/LoginPage.css';

const LoginPage = props => {

    const loginUser = () => {
        const user_name = document.getElementById("login-input").value;
        if (user_name !== "") {
            const forma = new FormData();
            forma.append('user_name', user_name);
            forma.append('color', getRandomColor());

            const request = new Request(`${server_name}/register-user`, {method: "POST", body: forma});
            fetch(request)
                .then(promise => promise.json())
                .then(response => {
                    if (response.response !== "bad") {
                        props.userIdCallback(response.response);
                        browser_history.push("/chat", {"user_name": user_name});
                    }
                })
            }
    }

    const getRandomColor = () => {
        var letters = '0123456789ABCDEF';
        var color = '#';
        for (let i = 0; i < 6; i++) {
          color += letters[Math.floor(Math.random() * letters.length)];
        }
        return color;
    }

    const IsEnterPressed = e => {
        if (e.key.toLowerCase() === "enter") {
            loginUser()
        }
    }

    return(
        <div id="main-login-container">
            <div id="login-container">
                <h3 id="login-title">Ingresa un usuario</h3>
                <div id="login-input-container">
                    <input onKeyDown={IsEnterPressed} type="text" id="login-input" placeholder="CharlesDexterWard21..."/>
                    <object id="connect-btn" type="image/svg+xml" data={loginLogo}>error</object>
                </div>
            </div>
        </div>
    )
}

export default LoginPage;