(this.webpackJsonptelefake=this.webpackJsonptelefake||[]).push([[0],{30:function(e,t,n){},31:function(e,t,n){},32:function(e,t,n){},33:function(e,t,n){},39:function(e,t,n){"use strict";n.r(t);var s=n(0),c=n(1),a=n.n(c),i=n(20),o=n.n(i),r=n(3),l=n(6),u=n(9),d=n(14),h=n(24),j=n(41),b=n(21),m=Object(b.a)(),g="http://0.0.0.0:5000",f=n.p+"static/media/satellite.78a72fb7.svg",v=(n(30),function(e){var t=function(){for(var e="0123456789ABCDEF",t="#",n=0;n<6;n++)t+=e[Math.floor(Math.random()*e.length)];return t};return Object(s.jsx)("div",{id:"main-login-container",children:Object(s.jsxs)("div",{id:"login-container",children:[Object(s.jsx)("h3",{id:"login-title",children:"Ingresa un usuario"}),Object(s.jsxs)("div",{id:"login-input-container",children:[Object(s.jsx)("input",{onKeyDown:function(n){"enter"===n.key.toLowerCase()&&function(){var n=document.getElementById("login-input").value;if(""!==n){var s=new FormData;s.append("user_name",n),s.append("color",t());var c=new Request("".concat(g,"/register-user"),{method:"POST",body:s});fetch(c).then((function(e){return e.json()})).then((function(t){"bad"!==t.response&&(e.userIdCallback(t.response),m.push("/chat",{user_name:n}))}))}}()},type:"text",id:"login-input",placeholder:"CharlesDexterWard21..."}),Object(s.jsx)("object",{id:"connect-btn",type:"image/svg+xml",data:f,children:"error"})]})]})})}),O=n(22),p=n.p+"static/media/folder.5ac5f545.svg",x=(n(31),function(e){var t=e.class_name,n=e.message_data,c=n.sender,a=n.body,i=n.is_control_message,o=n.color,r=n.is_file,l=e.download_callback;t=i?"".concat(t," control-message"):t;var u=function(){return/.+\.(jpg|png|gif)$/.test(a)&&r};return Object(s.jsx)("div",{className:"msg-space",children:Object(s.jsxs)("div",{onClick:r?l:function(){return!0},className:"message-container ".concat(t),children:[Object(s.jsx)("div",{style:{color:o},className:"msg-header",children:Object(s.jsx)("span",{className:"sender-name",children:c})}),function(){if(console.log("".concat(a,": ").concat(u())),u()){var e="".concat(g,"/static/").concat(a);return Object(s.jsx)("div",{className:"msg-image",children:Object(s.jsx)("img",{src:e,alt:e})})}return Object(s.jsxs)("div",{className:"message-body",children:[Object(s.jsx)("div",{className:"msg-body ".concat(r?"file_message":""),children:a}),Object(s.jsx)("div",{className:"msg-icon",children:r?Object(s.jsx)("object",{data:p,type:"image/svg+xml",children:"file"}):null})]})}()]})})}),k=(n(32),function(e){var t=e.user_data,n=t.color,c=t.user_name,a=t.status;return Object(s.jsxs)("div",{className:"user-status",children:[Object(s.jsx)("div",{style:{color:n},className:"user-name-label",children:c}),Object(s.jsx)("div",{className:"us-status ".concat(a?"online":"offline"),children:"*"})]})}),_=function e(){var t=this;Object(l.a)(this,e),this.isConnected=function(){return 1===t.status},this.connect=function(e,n){t.host=e,t.socket=new WebSocket(e),t.socket.onopen=t.onConnect,t.socket.onmessage=t.onMessage,t.socket.onclose=t.onDisconnect},this.setDisconnectionCallback=function(e){t.is_connected||(t.disconnection_callback=e)},this.setConnectionCallback=function(e){t.is_connected||(t.connection_callback=e)},this.onConnect=function(e){void 0!==t.connection_callback&&t.connection_callback(e)},this.onDisconnect=function(e){void 0!==t.disconnection_callback&&t.disconnection_callback(e)},this.Close=function(){return t.socket.close(1e3,"User logged out")},this.onMessage=function(e){void 0!==t.message_callback&&t.message_callback(e)},this.setMessageCallback=function(e){return t.message_callback=e},this.onError=function(e){return t.error_callback=e},this.emit=function(e){e.length>0?t.socket.send(e):console.warn("Message '".concat(e,"' is invalid"))},this.status=0,this.host=null,this.is_connected=!1,this.socket=null,this.disconnection_callback=void 0,this.connection_callback=void 0,this.message_callback=void 0,this.error_callback=void 0},y=n.p+"static/media/Untitled-1.6fdffd06.png",w=n.p+"static/media/paper-plane.8efd8a9b.svg",C=(n(33),function(e){Object(u.a)(n,e);var t=Object(d.a)(n);function n(e){var c;return Object(l.a)(this,n),(c=t.call(this,e)).getUserFile=function(){document.getElementById("file-input").click()},c.getStatusElement=function(){var e=c.state.logged?"connected":"disconnected";return Object(s.jsx)("div",{className:"status-entry ".concat(e),children:e})},c.toggleLogin=function(e){c.state.logged!==e&&c.setState(Object(r.a)(Object(r.a)({},c.state),{},{logged:e}))},c.isMessageDone=function(e){if("enter"===e.key.toLowerCase()){var t=document.getElementById("transmisor-input");c.transmisor.emit(t.value),t.value=""}},c.uploadFile=function(){var e=document.getElementById("file-input").files[0];if(e){var t=new FormData;t.append("file",e),t.append("file_name",e.name),t.append("user_uuid",c.state.user_id);var n=new Request("".concat(g,"/send-file"),{method:"POST",body:t});fetch(n).then((function(e){return e.json()})).then((function(e){"bad"===e.response&&console.warn("Server error while sending file")}))}},c.updateMessages=function(e){var t=c.state.messages,n=JSON.parse(e.data);n.is_control_message&&window.setTimeout(c.requestUsers,400),t.push(n),c.setState(Object(r.a)(Object(r.a)({},c.state),{},{messages:t}),c.scrollToBottom)},c.scrollToBottom=function(){var e=document.getElementById("chat-massages-container");e.scrollTop=e.scrollHeight},c.downloadFileCallback=function(e){var t=e.body,n="".concat(g,"/static/").concat(t);fetch(n).then((function(e){return e.blob()})).then((function(e){var n=window.URL.createObjectURL(e),s=document.createElement("a");s.href=n,s.download=t,s.click(),s.remove()}))},c.requestUsers=function(){return fetch("".concat(g,"/get-users")).then((function(e){return e.json()})).then((function(e){c.setState(Object(r.a)(Object(r.a)({},c.state),{},{users:e}))}))},c.requestMessagesHistory=function(){var e="".concat(g,"/retrive-messages-history");fetch(e).then((function(e){return e.json()})).then((function(e){e.length>0&&c.setState(Object(r.a)(Object(r.a)({},c.state),{},{messages:e}))}))},c.transmisor=void 0,c.user_name=e.location.state.user_name,c.state={user_id:"",logged:!1,messages:[],users:[]},c}return Object(O.a)(n,[{key:"componentDidMount",value:function(){if(""===this.state.user_id){this.requestUsers();var e=this.props.userIdCallback();this.setState(Object(r.a)(Object(r.a)({},this.state),{},{user_id:e}),this.requestMessagesHistory)}}},{key:"componentDidUpdate",value:function(){var e=this;void 0===this.transmisor&&this.state.user_id.length>=40&&(this.transmisor=new _,this.transmisor.setConnectionCallback((function(){return e.toggleLogin(!0)})),this.transmisor.setDisconnectionCallback((function(){return e.toggleLogin(!1)})),this.transmisor.setMessageCallback(this.updateMessages),this.transmisor.connect("".concat("ws://0.0.0.0:5000/login","?user_name=").concat(this.user_name),null))}},{key:"render",value:function(){var e=this,t=this.state,n=t.messages,c=t.users;return Object(s.jsxs)("div",{id:"chat-main-container",children:[Object(s.jsxs)("header",{children:[Object(s.jsx)("div",{id:"chat-logo",children:Object(s.jsx)("img",{src:y,alt:"telefake.png"})}),Object(s.jsxs)("h3",{id:"user-name",children:["Usuario: ",this.user_name]}),Object(s.jsx)("h3",{id:"telefake-label",children:"Telefake"})]}),Object(s.jsxs)("main",{id:"central-display",children:[Object(s.jsxs)("section",{id:"room-content",children:[Object(s.jsx)("div",{id:"room-header",children:"Usuarios"}),Object(s.jsx)("div",{id:"room-users-container",children:c.map((function(t,n){return t.user_name===e.user_name?null:Object(s.jsx)(k,{user_data:t},"user-".concat(n))}))})]}),Object(s.jsxs)("section",{id:"chat-messages",children:[Object(s.jsx)("div",{id:"chat-massages-container",children:n.map((function(t){return Object(s.jsx)(x,{download_callback:function(){return e.downloadFileCallback(t)},class_name:t.sender===e.user_name?"sent":"recivied",message_data:t},t.uuid)}))}),Object(s.jsxs)("div",{id:"chat-transmisor",children:[Object(s.jsx)("input",{onKeyDown:this.isMessageDone,type:"text",id:"transmisor-input",maxLength:"512",placeholder:"mensaje..."}),Object(s.jsx)("div",{onClick:this.getUserFile,className:"circular-input-btn",children:Object(s.jsx)("object",{color:"white",data:w,type:"image/svg+xml",id:"circular-input-btn-svg",children:"error"})})]})]})]}),Object(s.jsxs)("footer",{id:"status-bar",children:[this.getStatusElement(),Object(s.jsx)("div",{className:"status-entry",children:Object(s.jsx)("input",{onChange:this.uploadFile,id:"file-input",type:"file"})})]})]})}}]),n}(c.Component)),N=function(e){Object(u.a)(n,e);var t=Object(d.a)(n);function n(e){var s;return Object(l.a)(this,n),(s=t.call(this,e)).name="ChatError",s.message=e,s}return n}(Object(h.a)(Error));var D=function(){var e=Object(c.useRef)(""),t=function(t){if(""===e.current)e.current=t;else if(t!==e.current)throw new N("attempt to change a user_id which was already defined: '".concat(e.current,"' to '").concat(t,"'"))},n=function(){return e.current};return Object(s.jsx)("div",{className:"App",children:Object(s.jsx)(j.b,{history:m,children:Object(s.jsxs)(j.c,{children:[Object(s.jsx)(j.a,{exact:!0,path:"/",render:function(e){return Object(s.jsx)(v,Object(r.a)(Object(r.a)({},e),{},{userIdCallback:t}))}}),Object(s.jsx)(j.a,{exact:!0,path:"/chat",render:function(e){return Object(s.jsx)(C,Object(r.a)(Object(r.a)({},e),{},{userIdCallback:n}))}})]})})})},M=function(e){e&&e instanceof Function&&n.e(3).then(n.bind(null,42)).then((function(t){var n=t.getCLS,s=t.getFID,c=t.getFCP,a=t.getLCP,i=t.getTTFB;n(e),s(e),c(e),a(e),i(e)}))};o.a.render(Object(s.jsx)(a.a.StrictMode,{children:Object(s.jsx)(D,{})}),document.getElementById("root")),M()}},[[39,1,2]]]);
//# sourceMappingURL=main.de75485e.chunk.js.map