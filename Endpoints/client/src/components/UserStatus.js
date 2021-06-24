import React from 'react';
import "../css/UserLabel.css";

const UserLabel = props => {
    const {color, user_name, status } = props.user_data;
    return (
        <div className="user-status">
            <div style={{color: color}} className="user-name-label">{user_name}</div>
            <div className={`us-status ${status ? "online" : "offline"}`}>*</div>
        </div>
    )
}

export default UserLabel;