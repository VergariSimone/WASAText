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
    props: [ "selectedconvs", "sv", "login", "token" ],
    emits: [ 'list-next', 'list-previous', 'change-selected-conv', 'change-selected-conv-byid', 'show-user-info', 'refresh', 'logout' ],
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
            convs: { items: [] },
        }
    },
    mounted() {
        this.refresh();
        this.scrollToBottom();
    },
    watch: {
        // Quando cambiano i messaggi, scendi in basso
        'selectedconvs.message': {
            handler() {
                this.scrollToBottom();
            },
            deep: true
        }
    },
    methods: {
        // --- NUOVO: Funzione per scollare in fondo alla chat ---
        scrollToBottom() {
            this.$nextTick(() => {
                const container = this.$refs.messageListContainer;
                if (container) {
                    container.scrollTop = container.scrollHeight;
                }
            });
        },
        convertTimestamp(timestamp){
            const localDate = new Date(`${timestamp}Z`);
            const stringDate = localDate.toLocaleString("it-IT", Intl.DateTimeFormat().resolvedOptions().timeZone);
            return stringDate.slice(0, -3);
        },
        async newConversation() {
            if (this.showNewMessage) {
                this.showNewMessage = false;
                this.showNewMessageList = false;
                return;
            }
            this.showNewMessage = true;
            this.showNewMessageList = true;

            let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'application/json' } };
            try {
                let response = await this.$axios.get(`/users?search=`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)) {
                    this.selectedusers = [this.selectedusers];
                }
                await this.$nextTick();
                var input = document.getElementById('inputRef');
                if (input) input.disabled = true;
                if (this.$refs.sendmessage) {
                    this.$refs.sendmessage.isEnabled = false;
                }
            } catch (e) { console.error("Errore newConversation:", e); }
        },
        async addSelectedUsers(users){
            this.showListGroup = false;
            var out = [];
            for (var key in users) { out.push('users=' + encodeURIComponent(users[key])); }
            let config = { headers: { Authorization: "Bearer " + this.token, 'accept': 'application/json' } };
            try {
                let response = await this.$axios.put(`/conversations/group?${out.join('&')}`, {}, config);
                this.$emit('change-selected-conv-byid', response.data);
            } catch (e) { console.error(e); }
        },
        async addSelectedUsersNew(users){
            if (users.length > 1) { alert("Select only one user"); return; }
            if (sessionStorage.getItem('login') == users) { alert("Cannot send message to yourself"); return; }
            await this.$nextTick();
            var input = document.getElementById('inputRef');
            if (input) input.disabled = false;
            if (this.$refs.sendmessage) this.$refs.sendmessage.isEnabled = true;
            this.users = users;
            this.showNewMessageList = false;
        },
        async addSelectedUsersForward(users){
            let config = { headers: { Authorization: "Bearer " + this.token, 'accept': 'application/json' } };
            var convid = null;
            for (var key in users) {
                if (users[key] == this.login) continue;
                try {
                    let response = await this.$axios.post(`/conversations/byid/${this.selectedConversationId}/message/${this.showForwardMsgId}`, users[key], config);
                    convid = response.data;
                } catch(e) { console.error(e); }
            }
            this.showForward = false;
            if(convid) this.$emit('change-selected-conv-byid', convid.conID);
        },
        cancelUsersSearch(){ this.showListGroup = false; },
        cancelUsersSearchNew(){
            this.showNewMessageList = false;
            var input = document.getElementById('inputRef');
            if (input) input.disabled = false;
            if (this.$refs.sendmessage) this.$refs.sendmessage.isEnabled = true;
        },
        cancelUsersSearchForward(){ this.showForward = false; },
        async newGroup(){
            this.showListGroup = true;
            let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'application/json' } };
            try {
                let response = await this.$axios.get(`/users?search=`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)) this.selectedusers = [this.selectedusers];
                await this.$nextTick();
                // Posizionamento mantenuto come da tua logica originale
                var element = document.getElementById('listuser');
                var button = document.getElementById('buttonNewGroup');
                if (element && button) {
                    var rect = button.getBoundingClientRect();
                    element.style.position = 'fixed';
                    element.style.top = (rect.top + 50) + 'px';
                    element.style.left = (rect.left - 125) + 'px';
                    element.style.display = 'block';
                }
            } catch(e) { console.error(e); }
        },
        async newForward(){
            this.showForward = true;
            let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'application/json' } };
            try {
                let response = await this.$axios.get(`/users?showgroups=true&search=`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)) this.selectedusers = [this.selectedusers];
                await this.$nextTick();
                // Logica posizionamento originale
                var element = document.getElementById('listuser');
                if (element) {
                    var html = document.documentElement;
                    var X = 0, Y = 0;
                    element.style.position = 'fixed';
                    element.style.display = 'block';
                    var H = element.offsetHeight, W = element.offsetWidth;
                    Y = Math.min(this.topMY - 80 + H, html.clientHeight) - H;
                    X = Math.min(this.topMX - 80 + W, html.clientWidth) - W;
                    element.style.top = (Y) + 'px';
                    element.style.left = (X) + 'px';
                }
            } catch(e) { console.error(e); }
        },
        async searchUsersHelper(search, endpoint, callback) {
             let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'application/json' } };
            try {
                let response = await this.$axios.get(`${endpoint}search=${search}`, config);
                this.tmpselectedusers = response.data || JSON.parse('[{"login": "0", "name": "", "photo": ""}]');
                if (!Array.isArray(this.tmpselectedusers)) this.tmpselectedusers = [this.tmpselectedusers];
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
                    element.style.top = (rect.top + 50) + 'px';
                    element.style.left = (rect.left - 125) + 'px';
                    element.style.display = 'block';
                }
            });
        },
        async usersSearchAgainForward(search){
             let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'application/json' } };
            try {
                let response = await this.$axios.get(`/users?showgroups=true&search=${search}`, config);
                this.selectedusers = response.data;
                if (!Array.isArray(this.selectedusers)) this.selectedusers = [this.selectedusers];
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
                    element.style.top = (rect.top + 125) + 'px';
                    element.style.left = (rect.left - 125) + 'px';
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
            if (!element || !parent) return;
            if (element.style.visibility != 'visible') {
                var html = document.documentElement;
                element.style.position = 'fixed';
                element.style.display = 'block';
                element.style.visibility = 'visible';
                var H = element.offsetHeight, W = element.offsetWidth;
                var rect = parent.getBoundingClientRect();
                var Y = Math.min(rect.top - 0 + H, html.clientHeight) - H;
                var X = Math.min(rect.left - 30 + W, html.clientWidth) - W;
                element.style.left = (X) + 'px';
                element.style.top = (Y) + 'px';
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
                    element.style.position = 'fixed';
                    element.style.display = 'block';
                    element.style.visibility = 'visible';
                    var H = element.offsetHeight, W = element.offsetWidth;
                    var rect = parent.getBoundingClientRect();
                    var Y = Math.min(rect.top - 0 + H, html.clientHeight) - H;
                    var X = Math.min(rect.left - 0 + W, html.clientWidth) - W;
                    element.style.left = (X) + 'px';
                    element.style.top = (Y) + 'px';
                } else {
                    element.style.visibility = 'hidden';
                }
            }
        },
        closeEmoji(){ },
        messageAlert(status){
            this.errormsg = (status == 403) ? "Non puoi eliminare messaggi che non hai inviato te" : "Error " + status;
        },
        async deleteMessage(mesid, conid){
            this.errormsg = null;
            let config = { headers: { Authorization: "Bearer " + this.token } };
            try {
                await this.$axios.delete(`/conversations/byid/${conid}/message/${mesid}`, config);
            } catch (e) {
                 this.errormsg = (e.response && e.response.status === 403) ? "Non puoi eliminare messaggi che non hai inviato te" : e.toString();
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
                // Posizionamento commento semplificato
                var element = document.getElementById('comment');
                if (element) {
                    element.style.top = '50%';
                    element.style.left = '50%';
                    element.style.transform = 'translate(-50%, -50%)'; // Centra nello schermo
                    element.style.position = 'fixed';
                    element.style.display = 'block';
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
                // Rimossa scroll globale, focus su input
                await this.$nextTick();
                var element = document.getElementById("inputRef");
                if (element) {
                    element.value = "";
                    element.focus();
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
                this.errormsg = "No message composed"; return;
            }
            let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'multipart/form-data' } };
            let formdata = new FormData();
            if (message && message != '') formdata.append('text', message);
            if (image && image != '') formdata.append('image', image, 'image');
            var reply = this.showReplyMsgId ? `&reply=${this.showReplyMsgId}` : "";
            
            try {
                let response = await this.$axios.post(`/conversations?conversation=${this.selectedconvs.conid}${reply}`, formdata, config);
                this.$emit('change-selected-conv', response.data);
            } catch (e) { this.errormsg = e.toString(); }
            this.showReplyMsgId = null;
            this.Reply = null;
            this.refresh();
        },
        async sendMessageWithUser(message, image){
            this.errormsg = null;
            if (!this.users) { this.errormsg = "No recipient selected!"; return; }
            if (((message === null) || message == '') && (image === null)) { this.errormsg = "No message composed"; return; }
            let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'multipart/form-data' } };
            let formdata = new FormData();
            if (message && message != '') formdata.append('text', message);
            if (image && image != '') formdata.append('image', image, 'image');
            try {
                let response = await this.$axios.post(`/conversations?conversation=${this.users}`, formdata, config);
                this.showNewMessage = false;
                this.$emit('change-selected-conv', response.data);
            } catch (e) { this.errormsg = e.toString(); }
            this.refresh();
        },
        async sendComment(comment, mesid, conid){
            this.showComment = false;
            if (mesid === null) return;
            this.errormsg = null;
            let config = { headers: { Authorization: "Bearer " + this.token, 'content-type': 'text/plain' } };
            try {
                await this.$axios.post(`/conversations/byid/${conid}/message/${mesid}/comment`, comment.i, config);
            } catch (e) { this.errormsg = e.toString(); }
            this.refresh();
        },
        info() { this.$emit('show-user-info'); },
        classReceived(rcv){
            if (rcv) {
                for(let r of rcv) {
                    if (r.received == 'false' && r.name != this.login) return "noreceived";
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
                for(let r of rd) {
                    if (r.read == 'false' && r.name != this.login) return "noread";
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
            return (login == this.login) ? "message-box my-message" : "message-box friend-message";
        },
        logout() { this.$emit('logout'); },
        refresh() { this.$emit('refresh'); },
        previousMessages() { this.$emit('list-previous'); },
        nextMessages() { this.$emit('list-next'); },
        clearError() { this.errormsg = null; },
        exitReply() { this.Reply = null; this.showReplyMsgId = null; }
    }
}
</script>

<template>
    <div class="chat-wrapper">
        <div v-if="!showNewMessage" class="chat-container">
            <div id="main" class="chat-header">
                <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-2 pb-2 mb-2 border-bottom border-secondary">
                    <h1 class="h5 text-light mb-0 ps-2" v-if="selectedconvs">{{ selectedconvs.name }}</h1>
                    <h1 class="h5 mb-0" v-else>&nbsp;</h1>
                    <div class="btn-toolbar mb-2 mb-md-0 pe-2">
                        <div class="btn-group me-2">
                            <button type="button" class="btn btn-sm btn-outline-light" @click="refresh">Refresh</button>
                            <button type="button" class="btn btn-sm btn-outline-light" @click="info">Info</button>
                        </div>
                        <div class="btn-group me-2">
                            <button id="buttonNewGroup" type="button" class="btn btn-sm btn-outline-light" @click="newGroup">
                                Nuovo Gruppo
                            </button>
                            <button id="buttonNewChat" type="button" class="btn btn-sm btn-outline-light" @click="newConversation">
                                Nuova Chat
                            </button>
                        </div>
                        <div class="btn-group">
                            <button type="button" class="btn btn-sm btn-outline-info" @click="logout">Logout</button>
                        </div>                    
                    </div>
                </div>
                <ErrorMsg v-if="errormsg" :msg="errormsg" @clear-error="clearError"></ErrorMsg>
            </div>

            <div id="selectedconv" class="message-list-area" ref="messageListContainer">
                <ul class="no-bullets" v-if="selectedconvs">
                    <li v-if="sv.start != 0" @click="previousMessages"><div class="nextmessages"> Carica messaggi più recenti </div></li>
                    
                    <li v-for="cvs in selectedconvs.message" :key="cvs.mesid">
                        <div v-bind:class="classSenderIsMe(cvs.login)">
                            <div class="message-bubble"
                                @click.prevent.stop="handleMessageClick($event, cvs)"
                                @mouseenter="setCoordsMessage"
                            >
                                <div v-if="cvs.isforwarded > 0" class="forwarded-label">
                                    <small>Inoltrato</small>
                                </div>
                                
                                <h6 v-if="cvs.login != login">{{ cvs.username }}:</h6>

                                <div v-if="(cvs.replyedtext != '') || (cvs.replyedphoto != '')" class="replyed-preview mb-2">
                                    <div v-if="cvs.replyedusername != ''" class="replyedusername">{{ cvs.replyedusername }}</div>
                                    <div v-if="cvs.replyedtext != ''" class="replyedtext">{{ cvs.replyedtext }}</div>
                                    <div v-if="cvs.replyedphoto != ''"><img v-bind:src="'data:image/gif;base64,' + cvs.replyedphoto" class="internal-thumb"/></div>
                                </div>

                                <img v-if="cvs.category != 'text'" v-bind:src="'data:image/gif;base64,' + cvs.photo" class="internal-img"/>
                                
                                <div v-if="cvs.category != 'image'" class="message-text">
                                    {{ cvs.text }}
                                </div>

                                <div class="message-meta" :id="'sp' + cvs.mesid" @click.prevent.stop="openRW(cvs.mesid, cvs.login != login)">
                                    {{ convertTimestamp(cvs.timestamp) }}
                                    <span v-if="cvs.login == login" class="checks">
                                        <i v-bind:class="classReceived(cvs.receivedarray)"></i>
                                        <i v-bind:class="classRead(cvs.readarray)"></i>
                                    </span>
                                </div>

                                <div class="listrr" :id="'rw' + cvs.mesid" @click.prevent.stop="openRW(cvs.mesid, cvs.login != login)" v-if="cvs.readarray && cvs.login == login">
                                    <ul class="ullistrr">
                                        <li v-for="(item, i) in cvs.readarray" :key="item.name">{{ item.username }} 
                                            <i v-bind:class="singleClassReceived(cvs.receivedarray[i].received, cvs.receivedarray[i].name, login)"></i>
                                            <i v-bind:class="singleClassRead(item.read, cvs.receivedarray[i].name, login)"></i>
                                        </li>
                                    </ul>
                                </div>

                                <div class="emoji-badge" :id="'em' + cvs.mesid" 
                                    v-if="cvs.comments"
                                    @click="openEmoji(cvs.mesid, cvs.comments[0].commId)"
                                >
                                    {{ cvs.comments[0].emoji }}
                                    <div class="listemoji" :id="'co' + cvs.comments[0].commId" @click.stop="closeEmoji(cvs.comments[0].commId)" v-if="cvs.comments">
                                        <ul class="ullistemoji">
                                            <li v-for="c in cvs.comments" :key="c.commId">{{ c.emoji }} - {{ c.username }}</li>
                                        </ul>
                                    </div>
                                </div> 
                            </div>
                        </div>
                    </li>
                    <li v-if="sv.limit == selectedconvs.message.length" @click="nextMessages"><div class="nextmessages"> Carica messaggi più vecchi </div></li>
                </ul>
            </div>

            <div class="chat-footer">
                <div v-if="showReplyMsgId" class="reply-bar">
                    <div class="replyed inforeplyed" @click="exitReply">
                        <div v-if="Reply.username != ''" class="replyedusername">{{ Reply.username }}</div>
                        <div v-if="Reply.text != ''" class="replyedtext">{{ Reply.text }}</div>
                        <div v-if="Reply.photo != ''"><img v-bind:src="'data:image/gif;base64,' + Reply.photo" class="internal-thumb"/></div>
                    </div>
                </div>
                
                <sendmessage 
                    v-if="!showNewMessage"
                    v-bind:conversation="selectedconvs"
                    ref="sendmessage"
                    @send-message="sendMessage"
                />
            </div>
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
        
        <listuser id="listuser" v-if="showListGroup" @search-again="usersSearchAgain" @send-users-list="addSelectedUsers" @search-cancel="cancelUsersSearch" v-bind:users="selectedusers" v-bind:searchBar="true" v-bind:showOk="true" />
        <listuser id="listuser" v-if="showForward" @search-again="usersSearchAgainForward" @send-users-list="addSelectedUsersForward" @search-cancel="cancelUsersSearchForward" v-bind:users="selectedusers" v-bind:searchBar="true" v-bind:showOk="true" />

        <vue-simple-context-menu element-id="menu-message" :options="optionsmessage" ref="vueMessageMenu" @option-clicked="optionMessageClicked" />

        <div v-if="showNewMessage" class="chat-container">
             <div class="chat-header p-3 border-bottom border-secondary">
                 <button class="btn btn-outline-secondary btn-sm" @click="newConversation">Indietro</button>
                 <span class="ms-3 text-light">Nuova Conversazione</span>
             </div>
             <div class="message-list-area">
                <listuser id="listuser" v-if="showNewMessageList" @search-again="usersSearchAgainNew" @send-users-list="addSelectedUsersNew" @search-cancel="cancelUsersSearchNew" v-bind:users="selectedusers" v-bind:searchBar="true" v-bind:showOk="true" />
             </div>
            <div class="chat-footer">
                <sendmessage v-bind:conversation="showNewMessage" ref="sendmessage" @send-message="sendMessageWithUser" />
            </div>
        </div>
    </div>
</template>

<style scoped>
    /* --- LAYOUT GENERALE --- */
    /* Occupa tutto lo schermo ed è flessibile */
    .chat-wrapper {
        height: 100vh;
        width: 100%;
        background-color: #0f172a;
        color: #f1f5f9;
        overflow: hidden; /* Evita scroll body */
    }

    .chat-container {
        display: flex;
        flex-direction: column;
        height: 100%;
    }

    .chat-header {
        flex-shrink: 0; /* Non rimpicciolire */
        background-color: #0f172a;
        z-index: 10;
    }

    .message-list-area {
        flex-grow: 1; /* Occupa lo spazio rimanente */
        overflow-y: auto; /* Scrolla solo qui */
        padding: 10px;
        padding-bottom: 20px;
        scroll-behavior: smooth;
    }

    .chat-footer {
        flex-shrink: 0;
        background-color: #1e293b;
        padding: 10px;
        border-top: 1px solid rgba(255,255,255,0.1);
    }

    /* --- MESSAGGI E BOLLE --- */
    ul.no-bullets {
        list-style-type: none;
        padding: 0;
        margin: 0;
    }

    .message-box {
        display: flex;
        width: 100%;
        margin-bottom: 15px; /* Spazio tra messaggi */
        position: relative;
    }

    .my-message {
        justify-content: flex-end; /* Allinea a destra */
    }

    .friend-message {
        justify-content: flex-start; /* Allinea a sinistra */
    }

    .message-bubble {
        position: relative;
        max-width: 75%; /* Larghezza massima più sensata */
        min-width: 120px;
        padding: 12px 15px;
        border-radius: 12px;
        color: #fff;
        box-shadow: 0 1px 2px rgba(0,0,0,0.2);
        cursor: pointer;
    }

    /* Colori Bolle */
    .friend-message .message-bubble {
        background-color: #1e293b; /* Grigio Blu */
        border-bottom-left-radius: 2px; /* Piccola punta */
    }

    .my-message .message-bubble {
        background-color: #059669; /* Verde Smeraldo */
        border-bottom-right-radius: 2px; /* Piccola punta */
    }

    .message-text {
        font-size: 15px;
        line-height: 1.4;
        word-wrap: break-word;
    }

    .message-meta {
        display: block;
        text-align: right;
        font-size: 10px;
        opacity: 0.7;
        margin-top: 4px;
        margin-bottom: -5px; /* Compatta in basso */
    }

    .forwarded-label {
        font-size: 11px;
        font-style: italic;
        opacity: 0.6;
        margin-bottom: 4px;
        display: block;
    }

    /* --- IMMAGINI --- */
    .internal-img {
        max-width: 100%;
        height: auto;
        border-radius: 8px;
        margin-bottom: 5px;
        display: block;
    }

    .internal-thumb {
        height: 35px;
        width: auto;
        border-radius: 4px;
        margin-top: 2px;
    }

    /* --- REPLY PREVIEW (DENTRO MESSAGGIO) --- */
    .replyed-preview {
        background: rgba(0,0,0,0.2);
        border-left: 3px solid rgba(255,255,255,0.5);
        padding: 5px 8px;
        border-radius: 4px;
        font-size: 12px;
    }
    
    /* --- REPLY BAR (FOOTER) --- */
    .reply-bar {
        background: #0f172a;
        padding: 5px;
        border-radius: 5px;
        margin-bottom: 5px;
        border: 1px solid #334155;
    }

    .replyed {
        border-left: 4px solid #10b981;
        padding-left: 8px;
        position: relative;
    }

    .inforeplyed {
        cursor: pointer;
    }

    .replyedusername {
        color: #5eead4;
        font-weight: bold;
        font-size: 11px;
    }
    .replyedtext {
        color: #cbd5e1;
        font-size: 12px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    /* --- EMOJI BADGE --- */
    .emoji-badge {
        position: absolute;
        bottom: -10px; /* Esce leggermente sotto */
        right: -5px;   /* Esce leggermente a destra */
        background-color: #334155;
        border: 2px solid #0f172a; /* Bordo color sfondo pagina per stacco */
        border-radius: 12px;
        padding: 2px 6px;
        font-size: 14px;
        z-index: 5;
        cursor: pointer;
        box-shadow: 0 2px 4px rgba(0,0,0,0.3);
    }

    /* --- LISTE POPUP (Lettura/Emoji) --- */
    .listrr, .listemoji {
        position: fixed;
        background-color: #1e293b;
        padding: 8px;
        border-radius: 6px;
        border: 1px solid #475569;
        box-shadow: 0 4px 6px rgba(0,0,0,0.3);
        z-index: 9999;
        min-width: 100px;
    }
    .ullistrr, .ullistemoji {
        list-style: none;
        padding: 0;
        margin: 0;
        font-size: 12px;
    }

    /* --- SPUNTE --- */
    .checks {
        margin-left: 4px;
    }
    .noreceived::after, .noread::after { content: '✓'; color: #64748b; font-weight: bold; margin-left: 2px; }
    .received::after { content: '✓'; color: #22d3ee; font-weight: bold; margin-left: 2px; }
    .read::after { content: '✓'; color: #22d3ee; font-weight: bold; margin-left: 2px; }
    /* Doppio check per read? Normalmente si usano due span o icone diverse, qui mantengo logica CSS */

    /* --- BOTTONE CARICA ALTRO --- */
    .nextmessages {
        border: 1px solid #10b981;
        color: #10b981;
        background: rgba(16, 185, 129, 0.1);
        text-align: center;
        padding: 8px;
        margin: 10px auto;
        width: 80%;
        border-radius: 20px;
        cursor: pointer;
        font-size: 13px;
    }

    /* --- OLD CLASSES COMPATIBILITY --- */
    .overlayComment { z-index: 10000; }
    .hiddenComment { display: none; }
</style>
