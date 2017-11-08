// ++++++++++++++++++++++++++++++++++++
// TODO: 
// Add random id for RPC requests.
// Catch responses by id.
// ++++++++++++++++++++++++++++++++++++

window.onload = function () {
    // for select tag to work properly    
    $('select').material_select();

    // for modals
    $('.modal').modal();

    // for tooltips
    // $('.tooltipped').tooltip({delay: 50});
};


let ws = new WebSocket("ws://" + window.location.hostname + ":7001" + "/ws/v1");
let active_contact = {
    jid: "",
    account: ""
};

// let pwd_input = document.getElementById("password");

// FIXIT
// Does not work
// pwd_input.okeyup = function (ev) {
//     ev.preventDefault();
//     if (ev.keyCode === 13) {
//         XMPPConnectToServer();
//     }
// }

let output = document.getElementById("chat-window");
let status_select = document.getElementById("status-select");

status_select.onchange = function (ev) {
    let status = status_select.value;
    let request = {
        method: "XMPP.SetStatus",
        params: new Array({ "Show": status, "Status": "This is Cadmium IM test" }),
        id: "somerandomshit"
    }
    sendRequest(request);
}

let message_input = document.getElementById("message-input");

message_input.onkeyup = function (ev) {
    ev.preventDefault();
    if (ev.keyCode === 13) {
        sendChatMessage();
    }
}

// let connect_btn = document.getElementById("connect-btn");
let send_message_btn = document.getElementById("send-message-btn");

let nick = "";

ws.onopen = function (ev) {
    // connect_btn.classList.remove("disabled");
    console.log("WS connection opened!");
    output.innerHTML += "WS connection opened! <br/>";
    getRosters();
}
ws.onclose = function (ev) {
    // connect_btn.classList.add("disabled");
    output.innerHTML += "WS connection closed. Try to reload page..";
}
ws.onmessage = function (msg) {
    console.log(msg);
    processMessage(msg);
}


function sendRequest(request) {
    console.log("Sending to core: " + JSON.stringify(request));
    ws.send(JSON.stringify(request));
}

// function XMPPConnectToServer() {
//     let username = document.getElementById("jid").value;
//     let password = document.getElementById("password").value;
//     let domain = username.split("@")[1];
//     nick = username.split("@")[0];
//     if (username.length === 0 || password.length === 0 || domain === undefined) {
//         Materialize.toast("Username or password can not be empty", 3000);
//         return
//     }
//     let status = status_select.value;
//     let request = {
//         method: "XMPP.ConnectToServer",
//         params: new Array(
//             { "Username": username, "Domain": domain, "Password": password, "Show": status }
//         ),
//         id: "somerandomshit"
//     }
//     sendRequest(request);
// }

function sendChatMessage() {
    if (active_contact.account === "" || active_contact.jid === "") {
        return
    }
    let message = message_input.value;
    if (message.trim().length === 0) {
        return
    }
    output.innerHTML += "<b>" + nick + ":</b>" + message + "<br/>";
    let request = {
        method: "XMPP.SendMessage",
        params: new Array(
            { "From": active_contact.account, "Remote": active_contact.jid, "Text": message }
        ),
        id: "somerandomshit"
    }
    sendRequest(request);
    message_input.value = "";
}

function printChatMessage(msg_obj) {
    output.innerHTML += "<b>" + msg_obj.result.from + ":</b>" + msg_obj.result.text + "<br/>";
}

function printPresenceMessage(msg_obj) {
    output.innerHTML += "<b>" + msg_obj.result.from + " status :</b>" + msg_obj.result.show + "; " + msg_obj.result.status + "<br/>";
}

function processMessage(msg) {
    let msg_obj = JSON.parse(msg.data);
    switch (msg_obj.id) {
        case "core_broadcast":
            break;
        case "message":
            printChatMessage(msg_obj);
            break;
        case "presence":
            printPresenceMessage(msg_obj);
            break;
        case "showAccountsControlModal":
            showAccountsControlModal(msg_obj);
            break;
        case "roster":
            // TODO
            // update roster
            updateContactsWidget(msg_obj);
            break;
    }
}

function addAccount() {
    let protocol = document.getElementById("add-account-protocol-select").value;
    let jid = document.getElementById("add-account-jid-input").value;
    let username = jid.split("@")[0];
    let domain = jid.split("@")[1];
    let password = document.getElementById("add-account-password-input").value;
    let port = document.getElementById("add-account-port-input").value;    
    let create_new = document.getElementById("add-account-create-new-checkbox").checked;

    if (jid.length === 0 || password.length === 0 || domain === undefined) {
        Materialize.toast("Wrong username or password", 3000);
        return
    }

    let request = {
        method: "XMPP.AddAccount",
        params: new Array({
                "Protocol": protocol,
                "UserName": username,
                "PasswordString": password,
                "Domain": domain,
                "Port": parseInt(port, 10),
                "IsActive": true
        })
    }
    sendRequest(request);
    // TODO
    // Check if everything ok before closing window
    $("#add-account-modal").modal("close");
}

function showAccountsControlModal(msg_obj) {
    let accounts = msg_obj.result;
    console.log(accounts);

    let accounts_table = document.getElementById("accounts-table");
    let tbody = document.createElement("tbody");
    // Do stuff here
    for (let acc of accounts) {
        let tr = document.createElement("tr");
 
        // may not work in IE
        tr.innerHTML = `
        <td>
            <h5>` + acc.UserName + `@` + acc.Domain + `</h5>
        </td>
        <td>
            <div class="switch">
            <label>
                Off
                <input type="checkbox" id="account_switch_` + acc.ID +`">
                <span class="lever"></span>
                On
            </label>
            </div>
        </td>
        <td>
            <a href="#" class="waves-effect waves-light btn" id="account_edit_btn_` + acc.ID + `">Edit</a>
        </td>
        <td>
            <a href="#" class="waves-effect waves-light btn red" id="account_edit_btn_` + acc.ID + `">Delete</a>
        </td>
        `;

        tbody.appendChild(tr);
    }
    accounts_table.innerHTML = "";
    accounts_table.appendChild(tbody);


    $("#accounts-control-modal").modal('open');
}

function getAccounts(callback_func) {
    let request = {
        method: "XMPP.GetAccounts",
        params: [],
        id: callback_func,
    }
    sendRequest(request);
}

function getRosters() {
    let request = {
        method: "XMPP.GetRosters",
        params: [],
        id: "get_rosters"
    }
    sendRequest(request);
}

function updateContactsWidget(msg_obj) {
    let contacts_widget = document.getElementById("contacts-row");
    let account = msg_obj.result.from;
    let contacts = msg_obj.result.contacts;
    for (let contact of contacts) {
        let row = document.createElement("a");
        row.href = "#!";
        row.dataset.jid = contact.JID;
        row.dataset.account = account;
        row.className = "collection-item";
        row.innerHTML = "<b>" + contact.JID +"</b>";
        row.onclick = function() {
            setActiveContact(contact.JID, account);
        };
        contacts_widget.appendChild(row);

    }
}

function setActiveContact(jid ,acc) {
    let chat_title = document.getElementById("chat-title")
    chat_title.innerText = jid.toUpperCase();
    active_contact.jid = jid;
    active_contact.account = acc;
}