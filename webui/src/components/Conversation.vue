<script>
import SendMessage from '../components/SendMessage.vue';
import SetComment from '../components/SetComment.vue';
import ListUser from '../components/ListUser.vue';
import VueSimpleContextMenu from 'vue-simple-context-menu'; 
import ErrorMsg from '../components/ErrorMsg.vue';

export default {
    components: {
        'sendmessage': SendMessage,
        'setcomment': SetComment,
        'listuser': ListUser,
        'vue-simple-context-menu': VueSimpleContextMenu,
        'ErrorMsg': ErrorMsg
    },
    props: [
        "selectedconvs",
        "sv",
        "login",
        "token",
    ],
    emits: [
        'list-next',
        'list-previous',
        'change-selected-conv',
        'change-selected-conv-byid',
        'show-user-info',
        'refresh',
        'logout',
    ],
    data: function() {
        return {
            showComment: false,
            showListGroup: false,
            showForwardMsgId: null,
            showForward: false,
            showReplyMsgId: null,
            Reply: null,
            showNewMessage: false,
            showNewMessageList: false,
            errormsg: null,
            selectedMessageId: null,
            selectedConversationId: null,
            selectedusers: null,
            tmpselectedusers: null,
            usersNewGroup: null,
            signMXY: true,
            users: null,
            topMX: 0,
            topMY: 0,
            optionsmessage: [
                {name: 'Comment', slug: 'comment'},
                {name: 'Uncomment', slug: 'uncomment'},
                {name: 'Forward', slug: 'forward'},
                {name: 'Reply', slug: 'reply'},
                {type: 'divider'},
                {name: '<em>Delete message</em>', slug: 'delete'},
            ],
            convs: {
                items: []
            },
        }
    },
    mounted() {
        this.refresh();
    },
    methods: {
        convertTimestamp(timestamp){
            const localDate = new Date(`${timestamp}Z`);
            const stringDate = localDate.toLocaleString("it-IT", Intl.DateTimeFormat().resolvedOptions().timeZone);
            return stringDate.slice(0, -3);
        },
        async newConversation() {
    // 1. Se siamo già in modalità "Nuovo Messaggio", torniamo indietro
    if (this.showNewMessage) {
        this.showNewMessage = false;
        this.showNewMessageList = false;
        return;
    }

    // 2. Attiviamo la modalità "Nuovo Messaggio"
    // (Questo farà sparire l'header vecchio e il pulsante 'buttonNewGroup')
    this.showNewMessage = true;
    this.showNewMessageList = true;

    let config = {
        headers: {
            Authorization: "Bearer " + this.token,
            'content-type': 'application/json'
        }
    };

    try {
        // 3. Carichiamo gli utenti
        let response = await this.$axios.get(`/users?search=`, config);
        this.selectedusers = response.data;

        // Gestione array singolo/multiplo
        if (!Array.isArray(this.selectedusers)) {
            this.selectedusers = [this.selectedusers]; // Corretto da new Array() a []
        }

        // 4. Aspettiamo che Vue aggiorni il DOM
        await this.$nextTick();

        // 5. Disabilitiamo l'input (perché devi prima scegliere un utente)
        var input = document.getElementById('inputRef');
        if (input) input.disabled = true;
        
        if (this.$refs.sendmessage) {
            this.$refs.sendmessage.isEnabled = false;
        }

        // --- RIMOSSO TUTTO IL CODICE DEL POSIZIONAMENTO ---
        // Non serve più calcolare rect, top o left.
        // La lista utenti ora si posiziona da sola grazie al CSS nel div centrale.

    } catch (e) {
        console.error("Errore newConversation:", e);
    }
},
        async addSelectedUsers(users){
            this.showListGroup = false;
            var out = [];
            for (var key in users) {
                out.push('users=' + encodeURIComponent(users[key]));
            }
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'accept': 'application/json'
                }
            };
            var convid = 0;
            try {
                let response = await this.$axios.put(`/conversations/group?${out.join('&')}`, {}, config);
                convid = response.data;
                this.$emit('change-selected-conv-byid', convid);
            } catch (e) { console.error(e); }
        },
        async addSelectedUsersNew(users){
            if (users.length > 1) {
                alert("Select only one user");
                return;
            }

            if (sessionStorage.getItem('login') == users) {
                alert("Cannot send message to yourself");
                return;
            }
            
            await this.$nextTick();
            var input = document.getElementById('inputRef');
            if (input) input.disabled = false;
            
            if (this.$refs.sendmessage) this.$refs.sendmessage.isEnabled = true;
            this.users = users;
            this.showNewMessageList = false;
        },
        async addSelectedUsersForward(users){
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'accept': 'application/json'
                }
            };
            var convid = null;
            for (var key in users) {
                if (users[key] == this.login) {
                    continue;
                }
                try {
                    let response = await this.$axios.post(
                        `/conversations/byid/${this.selectedConversationId}/message/${this.showForwardMsgId}`, users[key], config);
                    convid = response.data;
                } catch(e) { console.error(e); }
            }
            
            var element = document.getElementById('listuser');
            if(element) element.style.display = 'none';
            
            this.showForward = false;
            if(convid) this.$emit('change-selected-conv-byid', convid.conID);
        },
        cancelUsersSearch(){
            this.showListGroup = false;
        },
        cancelUsersSearchNew(){
            this.showNewMessageList = false;
            var input = document.getElementById('inputRef');
            if (input) input.disabled = false;
            if (this.$refs.sendmessage) this.$refs.sendmessage.isEnabled = true;
        },
        cancelUsersSearchForward(){
            var element = document.getElementById('listuser');
            if (element) element.style.display = 'none';
            this.showForward = false;
        },
        async newGroup(){
            this.showListGroup = true;
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'application/json'
                }
            };
            try {
                let response = await this.$axios.get(`/users?search=`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)){
                    this.selectedusers = [this.selectedusers];
                }

                await this.$nextTick();
                var element = document.getElementById('listuser');
                var button = document.getElementById('buttonNewGroup');
                if (element && button) {
                    var rect = button.getBoundingClientRect();
                    element.style.position = 'fixed';
                    element.style.top = (rect.top + 50).toString() + 'px';
                    element.style.left = (rect.left - 125).toString() + 'px';
                    element.style.display = 'block';
                }
            } catch(e) { console.error(e); }
        },
        async newForward(){
            this.showForward = true;
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'application/json'
                }
            };
            try {
                let response = await this.$axios.get(`/users?showgroups=true&search=`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)){
                    this.selectedusers = [this.selectedusers];
                }

                await this.$nextTick();
                var element = document.getElementById('listuser');
                if (element) {
                    var html = document.documentElement;
                    var X = 0, Y = 0;

                    element.style.position = 'fixed';
                    element.style.display = 'block';

                    var H = element.offsetHeight, W = element.offsetWidth;

                    Y = Math.min(this.topMY - 80 + H, html.clientHeight) - H;
                    X = Math.min(this.topMX - 80 + W, html.clientWidth) - W;

                    element.style.top = (Y).toString() + 'px';
                    element.style.left = (X).toString() + 'px';
                }
            } catch(e) { console.error(e); }
        },
        async searchUsersHelper(search, endpoint, callback) {
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'application/json'
                }
            };
            try {
                let response = await this.$axios.get(`${endpoint}search=${search}`, config);
                this.tmpselectedusers = response.data;
                if (this.tmpselectedusers === null) {
                    this.tmpselectedusers = JSON.parse('[{"login": "0", "name": "", "photo": ""}]');
                }
                if (!Array.isArray(this.tmpselectedusers)) {
                    this.tmpselectedusers = [this.tmpselectedusers];
                }
                this.selectedusers = this.tmpselectedusers;
                if (callback) await callback();
            } catch(e) { console.error(e); }
        },
        async usersSearchAgain(search){
            await this.searchUsersHelper(search, '/users?', async () => {
                await this.$nextTick();
                var element = document.getElementById('listuser');
                var button = document.getElementById('buttonNewGroup');
                if (element && button) {
                    var rect = button.getBoundingClientRect();
                    element.style.position = 'fixed';
                    element.style.top = (rect.top + 50).toString() + 'px';
                    element.style.left = (rect.left - 125).toString() + 'px';
                    element.style.display = 'block';
                }
            });
        },
        async usersSearchAgainForward(search){
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'application/json'
                }
            };
            try {
                let response = await this.$axios.get(`/users?showgroups=true&search=${search}`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)) {
                    this.selectedusers = [this.selectedusers];
                }
            } catch(e) { console.error(e); }
        },
        async usersSearchAgainNew(search){
            await this.searchUsersHelper(search, '/users?', async () => {
                await this.$nextTick();
                var element = document.getElementById('listuser');
                var button = document.getElementById('buttonNewGroup');
                if (element && button) {
                    var rect = button.getBoundingClientRect();
                    element.style.position = 'fixed';
                    element.style.top = (rect.top + 125).toString() + 'px';
                    element.style.left = (rect.left - 125).toString() + 'px';
                    element.style.display = 'block';
                }
            });
        },
        setCoordsMessage(e){
            if (this.signMXY) {
                this.topMX = e.clientX;
                this.topMY = e.clientY;
            }
        },
        openRW(conid, isOther){
            var parent = document.getElementById('sp' + conid);
            var element = document.getElementById('rw' + conid);
            if (element === null || parent === null) return;

            if (element.style.visibility != 'visible') {
                var html = document.documentElement;
                var X = 0, Y = 0;

                element.style.position = 'fixed';
                element.style.display = 'block';
                element.style.visibility = 'visible';

                var H = element.offsetHeight, W = element.offsetWidth;
                var rect = parent.getBoundingClientRect();
                
                Y = Math.min(rect.top - 0 + H, html.clientHeight) - H;
                X = Math.min(rect.left - 30 + W, html.clientWidth) - W;

                element.style.left = (X).toString() + 'px';
                element.style.top = (Y).toString() + 'px';

            } else {
                element.style.visibility = 'hidden';
            }
        },
        openEmoji(conid, commid){
            var parent = document.getElementById('em' + conid);
            var element = document.getElementById('co' + commid);

            if (element && parent) {
                if (element.style.visibility != 'visible') {
                    var html = document.documentElement;
                    var X = 0, Y = 0;

                    element.style.position = 'fixed';
                    element.style.display = 'block';
                    element.style.visibility = 'visible';

                    var H = element.offsetHeight, W = element.offsetWidth;
                    var rect = parent.getBoundingClientRect();

                    Y = Math.min(rect.top - 0 + H, html.clientHeight) - H;
                    X = Math.min(rect.left - 0 + W, html.clientWidth) - W;

                    element.style.left = (X).toString() + 'px';
                    element.style.top = (Y).toString() + 'px';
                } else {
                    element.style.visibility = 'hidden';
                }
            }
        },
        closeEmoji(){
            // Do nothing
        },
        messageAlert(status){
            if (status == 403){
                this.errormsg = "non puoi eliminare messaggi che non hai inviato te";
            } else {
                this.errormsg = "Error " + status;
            }
        },
        async deleteMessage(mesid, conid){
            this.errormsg = null;
            let config = {
                headers: { Authorization: "Bearer " + this.token }
            };
            try {
                await this.$axios.delete(`/conversations/byid/${conid}/message/${mesid}`, config)
                .catch((response) => {
                    if(response.response) this.messageAlert(response.response.status);
                    else this.messageAlert(response.status);
                });
            } catch (e) {
                if (e.response && e.response.status === 403) {
                     this.errormsg = "non puoi eliminare messaggi che non hai inviato te";
                } else {
                     this.errormsg = e.toString();
                }
            }
            this.refresh();
        },
        async deleteComment(mesid, conid){
            this.errormsg = null;
            let config = { headers: { Authorization: "Bearer " + this.token } };
            try {
                await this.$axios.delete(`/conversations/byid/${conid}/message/${mesid}/comment`, config);
            } catch (e) { this.errormsg = e.toString(); }
            this.refresh();
        },
        async optionMessageClicked(event){
            if (event.option.slug == 'comment') {
                this.showComment = true;
                await this.$nextTick();
                var element = document.getElementById('comment');
                
                if (element) {
                    var html = document.documentElement;
                    var X = 0, Y = 0;

                    element.style.position = 'fixed';
                    element.style.display = 'block';

                    var H = 366, W = 364;
                    Y = Math.min(this.topMY - 80 + H, html.clientHeight) - H;
                    X = Math.min(this.topMX - 80 + W, html.clientWidth) - W;

                    element.style.top = (Y).toString() + 'px';
                    element.style.left = (X).toString() + 'px';
                }

            } else if (event.option.slug == 'uncomment') {
                this.deleteComment(event.item.mesid, event.item.conid);
            } else if (event.option.slug == 'delete') {
                this.deleteMessage(event.item.mesid, event.item.conid);
            } else if (event.option.slug == 'forward') {
                this.showForwardMsgId = event.item.mesid;
                this.showForward = true;
                this.newForward();
            } else if (event.option.slug == 'reply') {
                globalThis.scrollTo(0, 0);
                await this.$nextTick();
                var element = document.getElementById("inputRef");
                if (element) {
                    element.value = "";
                    element.setSelectionRange(0, 0);
                }
                this.showReplyMsgId = event.item.mesid;
                this.Reply = event.item;
            }
        },
        handleMessageClick(evt, mesid){
            this.selectedMessageId = mesid.mesid;
            this.selectedConversationId = mesid.conid;
            if (this.$refs.vueMessageMenu) {
                this.$refs.vueMessageMenu.showMenu(evt, mesid);
            }
        },
        async sendMessage(message, image){
            this.errormsg = null;
            if (((message === null) || message == '') && (image === null)) {
                this.errormsg = "No message composed";
                return;
            }
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'multipart/form-data'
                }
            };
            let formdata = new FormData();
            if (message && message != '') formdata.append('text', message);
            if (image && image != '') formdata.append('image', image, 'image');

            var reply = "";
            if (this.showReplyMsgId) {
                reply = `&reply=${this.showReplyMsgId}`;
            }
            try {
                let response = await this.$axios.post(`/conversations?conversation=${this.selectedconvs.conid}${reply}`, formdata, config);
                this.$emit('change-selected-conv', response.data);
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.showReplyMsgId = null;
            this.Reply = null;
            this.refresh();
        },
        async sendMessageWithUser(message, image){
            this.errormsg = null;
            if (!this.users) {
                this.errormsg = "No recipient selected!";
                return;
            }
            if (((message === null) || message == '') && (image === null)) {
                this.errormsg = "No message composed";
                return;
            }
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'multipart/form-data'
                }
            };
            let formdata = new FormData();
            if (message && message != '') formdata.append('text', message);
            if (image && image != '') formdata.append('image', image, 'image');

            try {
                let response = await this.$axios.post(`/conversations?conversation=${this.users}`, formdata, config);
                this.showNewMessage = false;
                this.$emit('change-selected-conv', response.data);
            } catch (e) {
                this.errormsg = e.toString();
            }
            this.refresh();
        },
        async sendComment(comment, mesid, conid){
            this.showComment = false;
            if (mesid === null) return;
            this.errormsg = null;
            let config = {
                headers: {
                    Authorization: "Bearer " + this.token,
                    'content-type': 'text/plain'
                }
            };
            try {
                await this.$axios.post(`/conversations/byid/${conid}/message/${mesid}/comment`, comment.i, config);
            } catch (e) { this.errormsg = e.toString(); }
            this.refresh();
        },
        info() {
            this.$emit('show-user-info');
        },
        classReceived(rcv){
            if (rcv) {
                var i = 0;
                while (i < rcv.length) {
                    if ((rcv[i].received == 'false') && (rcv[i].name != this.login)) {
                        return "noreceived";
                    }
                    i++;
                }
                return "received";
            }
            return "noreceived";
        },
        singleClassReceived(rcv, user, login){
            if (user == login) return "received";
            if (rcv == 'true') return "received";
            return "noreceived";
        },
        classRead(rd){
            if (rd) {
                var i = 0;
                while (i < rd.length) {
                    if ((rd[i].read == 'false') && (rd[i].name != this.login)) {
                        return "noread";
                    }
                    i++;
                }
                return "read";
            }
            return "noread";
        },
        singleClassRead(rd, user, login){
            if (user == login) return "read";
            if (rd == 'true') return "read";
            return "noread";
        },
        classSenderIsMe(login){
            if (login == this.login) {
                return "message-box my-message";
            }
            return "message-box friend-message";
        },
        logout() {
            this.$emit('logout');
        },
        refresh() {
            this.$emit('refresh');
        },
        previousMessages() {
            this.$emit('list-previous');
        },
        nextMessages() {
            this.$emit('list-next');
        },
        clearError() {
            this.errormsg = null;
        },
        exitReply() {
            this.Reply = null;
            this.showReplyMsgId = null;
        }
    }
}
</script>

<template>
    <div v-if="!showNewMessage" class="chat-container">
        <div id="main">
            <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom border-secondary">
                <h1 class="h4 text-light" v-if="selectedconvs && (!showNewMessage)">{{ selectedconvs.name }}</h1>
                <h1 class="h4" v-else>&nbsp;</h1>
                <div class="btn-toolbar mb-2 mb-md-0">
                    <div class="btn-group me-2">
                        <button type="button" class="btn btn-sm btn-outline-light" @click="refresh">Refresh</button>
                        <button type="button" class="btn btn-sm btn-outline-light" @click="info">Info</button>
                    </div>
                    <div class="btn-group me-2">
                        <button id="buttonNewGroup" type="button" class="btn btn-sm btn-outline-light" @click="newGroup">
                            Nuovo Gruppo
                        </button>
                        <button id="buttonNewChat" type="button" class="btn btn-sm btn-outline-light" :class="{'btn-warning text-dark':showNewMessage}" @click="newConversation">
                            Nuova Chat
                        </button>
                    </div>
                    <div class="btn-group me-2">
                        <button type="button" class="btn btn-sm btn-outline-info" @click="logout">Logout</button>
                    </div>                    
                </div>
            </div>
            <ErrorMsg v-if="errormsg" :msg="errormsg" @clear-error="clearError"></ErrorMsg>
        </div>

        <div id="selectedconv" class="message-list-area" v-if="!showNewMessage">
            <ul class="no-bullets" v-if="selectedconvs">
                <li v-if="sv.start != 0" @click="previousMessages"><div class="nextmessages"> Carica messaggi più recenti </div></li>
                <li v-for="cvs in selectedconvs.message" :key="cvs.mesid">
                    <div v-bind:class="classSenderIsMe(cvs.login)">
                        <div 
                            @click.prevent.stop="handleMessageClick($event, cvs)"
                            @mouseenter="setCoordsMessage"
                        >
                            <span v-if="cvs.isforwarded > 0" class="forwarded">Inoltrato</span>
                            <h6 v-if="cvs.login != login">{{ cvs.username }}:</h6>
                            <div v-if="(cvs.replyedtext != '') || (cvs.replyedphoto != '')" class="replyed">
                                <div v-if="cvs.replyedusername != ''" class="replyedusername">{{ cvs.replyedusername }}</div>
                                <div v-if="cvs.replyedtext != ''" class="replyedtext">{{ cvs.replyedtext }}</div>
                                <div><img v-if="cvs.replyedphoto != ''" v-bind:src="'data:image/gif;base64,' + cvs.replyedphoto" class="internal"/></div>
                            </div>
                            <img v-if="cvs.category != 'text'" v-bind:src="'data:image/gif;base64,' + cvs.photo" class="internal"/>
                            <div v-if="cvs.category != 'image'">{{ cvs.text }}<br></div>
                            <span :id="'sp' + cvs.mesid"
                                @click.prevent.stop="openRW(cvs.mesid, cvs.login != login)"
                            >
                                {{ convertTimestamp(cvs.timestamp) }}
                                <i v-if="cvs.login == login" v-bind:class="classReceived(cvs.receivedarray)"></i>
                                <i v-if="cvs.login == login" v-bind:class="classRead(cvs.readarray)"></i>
                            </span>
                            <div class="listrr" :id="'rw' + cvs.mesid"
                                @click.prevent.stop="openRW(cvs.mesid, cvs.login != login)"
                                v-if="cvs.readarray && cvs.login == login">
                                <ul class="ullistrr">
                                    <li v-for="(item, i) in cvs.readarray" :key="item.name">{{ item.username }} 
                                        <i v-bind:class="singleClassReceived(cvs.receivedarray[i].received, cvs.receivedarray[i].name, login)"></i>
                                        <i v-bind:class="singleClassRead(item.read, cvs.receivedarray[i].name, login)"></i>
                                    </li>
                                </ul>
                            </div>
                        </div>
                        <div class="emoji" :id="'em' + cvs.mesid" 
                            v-if="cvs.comments"
                            @click="openEmoji(cvs.mesid, cvs.comments[0].commId)"
                            >
                            {{ cvs.comments[0].emoji }}
                            <div class="listemoji" :id="'co' + cvs.comments[0].commId"
                                @click="closeEmoji(cvs.comments[0].commId)"
                                v-if="cvs.comments">
                                <ul class="ullistemoji">
                                    <li v-for="c in cvs.comments" :key="c.commId">{{ c.emoji }} - {{ c.username }}</li>
                                </ul>
                            </div>
                        </div> 
                    </div>
                </li>
                <li v-if="sv.limit == selectedconvs.message.length" @click="nextMessages"><div class="nextmessages"> Carica messaggi più vecchi </div></li>
            </ul>
        </div>

        <div class="chat-footer">
            <div v-if="showReplyMsgId">
                <div v-if="(Reply.text != '') || (Reply.photo != '')" class="replyed inforeplyed" @click="exitReply">
                    <div v-if="Reply.username != ''" class="replyedusername">{{ Reply.username }}</div>
                    <div v-if="Reply.text != ''" class="replyedtext">{{ Reply.text }}</div>
                    <div><img v-if="Reply.photo != ''" v-bind:src="'data:image/gif;base64,' + Reply.photo" class="internal"/></div>
                </div>
            </div>
            
            <sendmessage 
                v-if="!showNewMessage"
                v-bind:conversation="selectedconvs"
                ref="sendmessage"
                @send-message="sendMessage"
            />
        </div>

        <div class="overlayComment" :class="{ hiddenComment: !showComment }">
            <setcomment
                id="comment"
                v-bind:msgid="selectedMessageId"
                v-bind:conid="selectedConversationId"
                v-bind:topX="topMX"
                v-bind:topY="topMY"
                @send-comment="sendComment"
            />
        </div>
        
        <listuser
            id="listuser"
            v-if="showListGroup"
            @search-again="usersSearchAgain"
            @send-users-list="addSelectedUsers"
            @search-cancel="cancelUsersSearch"
            v-bind:users="selectedusers"
            v-bind:searchBar="true"
            v-bind:showOk="true"
        />

        <listuser
            id="listuser"
            v-if="showForward"
            @search-again="usersSearchAgainForward"
            @send-users-list="addSelectedUsersForward"
            @search-cancel="cancelUsersSearchForward"
            v-bind:users="selectedusers"
            v-bind:searchBar="true"
            v-bind:showOk="true"
        />

        <vue-simple-context-menu
            element-id="menu-message"
            :options="optionsmessage"
            ref="vueMessageMenu"
            @option-clicked="optionMessageClicked"
        />
    </div>

    <div v-if="showNewMessage" class="chat-container">
        <div id="main2"></div>            
        
         <div class="message-list-area">
            <listuser
                id="listuser" 
                v-if="showNewMessageList"
                @search-again="usersSearchAgainNew"
                @send-users-list="addSelectedUsersNew"
                @search-cancel="cancelUsersSearchNew"
                v-bind:users="selectedusers"
                v-bind:searchBar="true"
                v-bind:showOk="true"
            />
         </div>

        <div class="chat-footer">
            <sendmessage 
                v-bind:conversation="showNewMessage"
                ref="sendmessage"
                @send-message="sendMessageWithUser"
            />
        </div>
    </div>
</template>

<style scoped>
    /* --- STILI RISPOSTA (REPLY) --- */
    .replyed {
        border: 1px solid rgba(255, 255, 255, 0.1); /* Bordo sottile */
        border-left: 5px solid #10b981; /* Barra laterale Verde Smeraldo */
        height: 80px;
        opacity: 0.9;
        background-color: rgba(30, 41, 59, 0.7); /* Sfondo scuro semitrasparente */
        border-radius: 4px;
        color: #e2e8f0;
    }
    
    .inforeplyed::before {
        /* Filtro per rendere l'icona bianca */
        filter: invert(1);
        content: url('/close.svg');
        position: absolute;
        /* margin-top: -20px; */
        margin-left: -27px;
        cursor: pointer;
        width: 20px;
        height: 20px;
    }
    
    .replyed img {
        height: 40px;
    }
    
    .replyedtext {
        height: 20px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 200px;
        font-size: 12px;
        color: #cbd5e1; /* Testo grigio chiaro */
    }
    
    .replyedusername {
        font-size: 11px;
        color: #5eead4; /* Nome utente verde acqua */
        font-weight: bold;
    }
    
    /* --- SPUNTE DI LETTURA --- */
    .noreceived::after {
        content: '✓';
        color: #64748b; /* Grigio scuro (non ricevuto) */
        margin-right: -0.4em;
        font-weight: bold;
    }
    
    .received::after {
        content: '✓';
        color: #22d3ee; /* Ciano brillante (ricevuto) */
        margin-right: -0.4em;
        font-weight: bold;
    }
    
    .noread::after {
        content: '✓';
        color: #64748b; /* Grigio scuro */
        font-weight: bold;
    }
    
    .read::after {
        content: '✓';
        color: #22d3ee; /* Ciano brillante */
        font-weight: bold;
    }
    
    /* --- SFONDO GENERALE --- */
    body {
        background-color: #0f172a; /* Blu Notte Profondo */
        z-index: 0;
        color: #f1f5f9;
    }
    
    main {
        background-color: #0f172a; /* Blu Notte Profondo */
    }
    
    ul.no-bullets {
        list-style-type: none; /* Remove bullets */
        padding: 0; /* Remove padding */
        margin: 0; /* Remove margins */
    }    
    
    .internal {
        max-width: 600px;
        width: auto;
        /* max-height: 200px; */
        height: 200px;
        border-radius: 8px;
    }
    
    .message {
        min-width: 270px;
        margin-right: 5px;
        clear: both;
        font-size:  x-small;
        font-weight: 200;
        vertical-align: baseline;
        color: #94a3b8; /* Testo grigio medio */
        overflow: hidden; 
        white-space: nowrap; 
        text-overflow: ellipsis;
    }
    
    .message-box {
      position: relative;
      display: grid;
      grid-template-columns: auto;
      width: 100%;
      margin: 5px 0;
      margin-top: 20px;
    }
    
    .message-box > div {
      position: relative;
      right: 0;
      text-align: left;
      max-width: 45vmax;    
      padding: 15px;
      /* Colore di default sovrascritto sotto */
      background: #1e293b; 
      border-radius: 10px;
      font-size: 1.1em;
      cursor: pointer;
      color: #fff; /* Testo bianco */
      box-shadow: 0 2px 4px rgba(0,0,0,0.3);
    }
    
    .message-box > div > h6 {
      color: #2dd4bf; /* Verde Acqua per i nomi */
      font-size: 0.9em;
      font-weight: bold;
    }
    
    /* --- MIO MESSAGGIO (VERDE SMERALDO SCURO) --- */
    .message-box.my-message > div {
        background: linear-gradient(135deg, #059669 0%, #047857 100%);
    }
    
    .message-box.my-message > div::before {
      content : '';
      position: absolute;
      top: 0;
      right: -12px;
      width: 20px;
      height: 20px;
      /* Aggiornato il gradiente per matchare il verde smeraldo */
      background: linear-gradient(135deg, #059669 0%, #059669 50%, transparent 50%, transparent);
    }
    
    .message-box > div > span {
      display: block;
      margin-top: 5px;
      font-size: 0.8em;
      /* opacity: 0.5; */
      font-style: italic;
      color: #cbd5e1; /* Colore timestamp */
    }
    
    .forwarded {
      /* display: block; */
      /* margin-top: 5px; */
      font-size: 0.8em;
      /* opacity: 0.5; */
      font-style: italic;
      margin-top: -15px!important;
      opacity: 0.5!important;
      color: #94a3b8;
    }
    
    /* .forwarded::before{
        content: "⮫";
        font-size: 1.5em;
    } */
    .my-message {
      justify-content: flex-end;
    }
    
    /* --- MESSAGGIO AMICO (GRIGIO BLUASTRO) --- */
    .friend-message > div {
      background: #1e293b;
    }
    
    .friend-message {
      justify-content: flex-start;
    }
    
    .message-box.friend-message::before {
      content : '';
      position: absolute;
      top: 0;
      left: -12px;
      width: 20px;
      height: 20px;
      /* Aggiornato gradiente per matchare il grigio scuro */
      background: linear-gradient(225deg, #1e293b 0%, #1e293b 50%, transparent 50%, transparent);
    }
    
    .overlayComment {
        display: block;
        z-index: 10000;
        position: fixed;
      left: 0;
      top: 0;
      width: 100%;
      height: 10px;
      background: #10b981; /* Verde barra superiore */
      border-style: solid;
      border-color: #0f172a;
      border-width: 5px;
    }
    .hiddenComment {
      display: none;
    }
    
    .emoji {
        position: relative;
        top: -25px !important;
        border-radius: 50%;
        width: fit-content;
        height: fit-content;
        overflow: hidden;
        background-color: #334155 !important; /* Sfondo scuro per emoji */
        border: 1px solid #0f172a;
        opacity: 1;
        z-index: 2;
        padding: 2px !important;
        cursor: pointer;
        margin-top: 15px !important;
        color: #fff;
    }
    .emoji span {
        visibility: hidden;
    }
    .emoji::before{
        content: none!important;
        background: #fff!important;
    }
    .emoji::after{
        content: none!important;
        background: #fff!important;
    }
    
    ul.ullistemoji {
        list-style-type: none; /* Remove bullets */
        padding: 0; /* Remove padding */
        margin: 0; /* Remove margins */
    } 
    
    .listemoji {
        position: fixed;
        /* display: none; */
        /* background-color: #077A07; */
        padding: 5px;
        visibility: hidden;
        border-radius: 5px;
        width: fit-content;
        height: fit-content;
        background-color: #1e293b; /* Sfondo popup scuro */
        color: #fff;
        border: 1px solid #334155;
        z-index: 100000;
    }
    
    ul.ullistrr {
        list-style-type: none; /* Remove bullets */
        padding: 0; /* Remove padding */
        margin: 0; /* Remove margins */
    } 
    
    .listrr {
        position: fixed;
        /* display: none; */
        /* background-color: #077A07; */
        padding: 5px;
        visibility: hidden;
        border-radius: 5px;
        width: fit-content;
        height: fit-content;
        background-color: #1e293b; /* Sfondo popup scuro */
        color: #fff;
        border: 1px solid #334155;
        z-index: 100000;
    }
    
    .hiddenForward {
      transform: translateX(-100%);
    }
    
    .nextmessages {
        border: 1px solid #10b981;
        width: 30%;
        margin-left: auto;
        margin-right: auto;
        margin-bottom: 10px;
        margin-top: 10px;
        padding: 5px;
        text-align: center;
        border-radius: 5px;
        background-color: rgba(16, 185, 129, 0.2);
        color: #10b981;
        cursor: pointer;
    }
    </style>