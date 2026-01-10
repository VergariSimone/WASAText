<script>
    import SendMessage from '../components/SendMessage.vue';
    import SetComment from '../components/SetComment.vue';
    import ListUser from '../components/ListUser.vue';
    import VueSimpleContextMenu from 'vue-simple-context-menu';
    import ErrorMsg from '../components/ErrorMsg.vue';
    
    export default {
        name: 'ConversationComponent',
        components: {
            'sendmessage': SendMessage,
            'setcomment': SetComment,
            'listuser': ListUser,
            'vue-simple-context-menu': VueSimpleContextMenu,
            'ErrorMsg': ErrorMsg
        },
        props: [
            "selectedconvs", "sv", "login", "token",
        ],
        emits: [
            'list-next', 'list-previous', 'change-selected-conv',
            'change-selected-conv-byid', 'show-user-info', 'refresh', 'logout',
        ],
        data() {
            return {
                availableEmojis: ['👍', '❤️', '😂', '😮', '😢', '😡', '👏', '🔥', '🎉'],
                showEmojiModal: false,
                showNewChatModal: false,
                showGroupModal: false,
                
                // inoltro
                showForwardModal: false,
                messageToForward: null,

                targetUserNewChat: null,
                selectedUsersList: [],
                selectedMessage: null,
                replyToMessage: null,
                errormsg: null,
                optionsmessage: [
                    {name: 'Reagisci', slug: 'comment'},
                    {name: 'Rispondi', slug: 'reply'},
                    {name: 'Inoltra', slug: 'forward'},
                    {type: 'divider'},
                    {name: 'Elimina', slug: 'delete'},
                ],
            }
        },
        computed: {
            sortedMessages() {
                if (!this.selectedconvs || !this.selectedconvs.message) return [];
                return [...this.selectedconvs.message].reverse();
            }
        },
        mounted() {
            this.refresh();
            this.scrollToBottom();
        },
        watch: {
            'selectedconvs.message': {
                handler() { this.scrollToBottom(); },
                deep: true
            }
        },
        methods: {
            refresh() { this.$emit('refresh'); }, // ricarica
            logout() { this.$emit('logout'); }, // esci
            info() { this.$emit('show-user-info'); }, // info utente
            clearError() { this.errormsg = null; },
            
            scrollToBottom() {
                this.$nextTick(() => {
                    const el = this.$refs.msgContainer;
                    if (el) el.scrollTop = el.scrollHeight;
                });
            },
            
            // formato orario
            formatTime(ts) {
                if (!ts) return '';
                return new Date(ts + 'Z').toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' });
            },
    
            // nuova chat
            async startNewChat() {
                this.showNewChatModal = true;
                this.targetUserNewChat = null;
                this.selectedUsersList = [];
                try {
                    let res = await this.$axios.get(`/users?search=`, { headers: { Authorization: "Bearer " + this.token } });
                    this.selectedUsersList = Array.isArray(res.data) ? res.data : [res.data];
                } catch (e) { this.errormsg = "Errore caricamento utenti"; }
            },
    
            // utenti selezionati per la nuova chat
            onUserSelectedForChat(users) {
                let user = Array.isArray(users) ? users[0] : users;
                if (!user) return;
                if (user === this.login) { this.errormsg = "Non puoi scriverti da solo"; return; }
                this.targetUserNewChat = user;
                this.showNewChatModal = false;
            },
    
            // annulla new chat
            cancelNewChat() {
                this.showNewChatModal = false;
                this.targetUserNewChat = null;
            },
    
            // gruppi
            async startNewGroup() {
                this.showGroupModal = true;
                try {
                    let res = await this.$axios.get(`/users?search=`, { headers: { Authorization: "Bearer " + this.token } });
                    this.selectedUsersList = Array.isArray(res.data) ? res.data : [res.data];
                } catch (e) {}
            },
    
            // crea gruppo
            async createGroup(users) {
                this.showGroupModal = false;
                var out = [];
                for (var key in users) out.push('users=' + encodeURIComponent(users[key]));
                try {
                    let res = await this.$axios.put(`/conversations/group?${out.join('&')}`, {}, { headers: { Authorization: "Bearer " + this.token } });
                    this.$emit('change-selected-conv-byid', res.data);
                } catch (e) { this.errormsg = e.toString(); }
            },

            // inoltro messaggi
            async startForwardingSearch() {
                this.selectedUsersList = [];
                try {
                    //  scarica lista Utenti
                    let resUsers = await this.$axios.get(`/users?search=`, { headers: { Authorization: "Bearer " + this.token } });
                    let users = Array.isArray(resUsers.data) ? resUsers.data : [resUsers.data];

                    // scarica lista Conversazioni (Gruppi) esistenti
                    let resConvs = await this.$axios.get(`/conversations`, { headers: { Authorization: "Bearer " + this.token } });
                    let convs = resConvs.data || [];

                    let groups = convs.filter(c => c.name).map(c => {
                        return { name: c.name, isGroup: true, id: c.conid }; 
                    });

                    // Uniamo (Utenti + Gruppi)
                    this.selectedUsersList = [...users, ...groups];

                } catch (e) { 
                    this.errormsg = "Errore caricamento destinatari"; 
                }
            },

            // In Conversation.vue -> methods

            async executeForward(selection) {
                this.showForwardModal = false;
                let item = Array.isArray(selection) ? selection[0] : selection;
                
                if (!item || !this.messageToForward) return;
                let dest = item.id ? item.id : (item.name ? item.name : item);

                let formData = new FormData();
                if (this.messageToForward.text) {
                    formData.append('text', this.messageToForward.text);
                }
                formData.append('isforwarded', '1');

                try {
                    await this.$axios.post(
                        `/conversations?conversation=${dest}`, 
                        formData, 
                        { headers: { Authorization: "Bearer " + this.token } }
                    );

                    alert("Messaggio inoltrato!");
                    this.messageToForward = null;
                    this.refresh(); // aggiorna la chat corrente

                } catch (e) {
                    this.errormsg = "Errore inoltro: " + e.toString();
                }
            },
    
            // invio mess
            async handleSend(text, image) {
                this.errormsg = null;
                if (this.targetUserNewChat) await this.sendToNewUser(text, image);
                else if (this.selectedconvs) await this.sendToExistingChat(text, image);
                else this.errormsg = "Errore: Nessuna chat attiva.";
            },
    
            // invio new utente
            async sendToNewUser(text, image) {
                let formData = new FormData();
                if (text) formData.append('text', text);
                if (image) formData.append('image', image, 'image');
                try {
                    let res = await this.$axios.post(`/conversations?conversation=${this.targetUserNewChat}`, formData, { headers: { Authorization: "Bearer " + this.token } });
                    this.targetUserNewChat = null;
                    this.$emit('change-selected-conv', res.data);
                    this.refresh();
                } catch (e) { this.errormsg = "Errore invio: " + e.toString(); }
            },
    
            // invia mess a chat esistente
            async sendToExistingChat(text, image) {
                let formData = new FormData();
                if (text) formData.append('text', text);
                if (image) formData.append('image', image, 'image');
                let url = `/conversations?conversation=${this.selectedconvs.conid}`;
                if (this.replyToMessage) url += `&reply=${this.replyToMessage.mesid}`;
                try {
                    let res = await this.$axios.post(url, formData, { headers: { Authorization: "Bearer " + this.token } });
                    this.$emit('change-selected-conv', res.data);
                    this.replyToMessage = null;
                    this.refresh();
                } catch (e) { this.errormsg = "Errore invio: " + e.toString(); }
            },
    
            // menu e opzioni
            showContextMenu(e, msg) {
                this.selectedMessage = msg;
                const isMyMessage = (msg.login === this.login);

                // mio mess 
                if (isMyMessage) {
                    this.optionsmessage = [
                        { name: 'Reagisci', slug: 'comment' },
                        { name: 'Rispondi', slug: 'reply' },
                        { name: 'Inoltra', slug: 'forward' },
                        { type: 'divider' },
                        { name: 'Elimina', slug: 'delete' }
                    ];
                } else { // mess altro utente
                    this.optionsmessage = [
                        { name: 'Reagisci', slug: 'comment' },
                        { name: 'Rispondi', slug: 'reply' },
                        { name: 'Inoltra', slug: 'forward' }
                    ];
                }
                this.$refs.vueMessageMenu.showMenu(e, msg);
            },
    
            onMenuOption(event) {
                const { slug } = event.option;
                const msg = event.item;

                if (slug === 'comment') {
                    this.selectedMessage = msg;
                    this.showEmojiModal = true;
                } else if (slug === 'reply') {
                    this.replyToMessage = msg;
                } else if (slug === 'delete') {
                    this.deleteMsg(msg);
                } else if (slug === 'forward') {
                    this.messageToForward = msg;
                    this.showForwardModal = true;
                    this.startForwardingSearch();
                }
            },
    
            // invio reazione 
            async sendReaction(emojiChar) {
                if (!this.selectedMessage) return;
                
                // Chiudi il modale
                this.showEmojiModal = false;

                const existingReaction = this.selectedMessage.comments 
                    ? this.selectedMessage.comments.find(c => c.username === this.login && c.emoji === emojiChar)
                    : null;

                const conid = this.selectedMessage.conid || this.selectedconvs.conid;
                const mesid = this.selectedMessage.mesid;

                try {
                    if (existingReaction) {
                        // DELETE: rimuovi reazione
                        await this.$axios.delete(
                            `/conversations/byid/${conid}/message/${mesid}/comment`, 
                            { headers: { Authorization: "Bearer " + this.token } }
                        );

                    } else {
                        // POST: aggiungi reazione
                        await this.$axios.post(
                            `/conversations/byid/${conid}/message/${mesid}/comment`, 
                            { text: emojiChar }, 
                            { headers: { Authorization: "Bearer " + this.token } }
                        );
                    }
                    
                    // Ricarica per aggiornare la vista
                    this.refresh(); 

                } catch (e) { 
                    console.error("Errore reazione:", e);
                    this.errormsg = "Errore gestione reazione";
                }
            },
            
    
            // eliminazione messaggio 
            async deleteMsg(msg) {
                try {
                    const conversationID = msg.conid || this.selectedconvs.conid;
                    await this.$axios.delete(`/conversations/byid/${conversationID}/message/${msg.mesid}`, { 
                        headers: { Authorization: "Bearer " + this.token } 
                    });
                    
                    this.refresh();

                } catch(e) { 
                    this.errormsg = (e.response && e.response.status === 403) 
                        ? "Non puoi eliminare messaggi altrui" 
                        : "Errore durante l'eliminazione"; 
                }
            },
    
            getClass(msg) { return msg.login === this.login ? 'msg-me' : 'msg-friend'; },
            
            getCheckStatus(msg) {
                if (msg.login !== this.login) return null;
                let read = msg.readarray && msg.readarray.some(u => u.name !== this.login && u.read === 'true');
                if (read) return 'read';
                let rcv = msg.receivedarray && msg.receivedarray.some(u => u.name !== this.login && u.received === 'true');
                if (rcv) return 'received';
                return 'sent';
            },
    

            toggleEmojiPopup(commId) {
                let el = document.getElementById('co' + commId);
                if(el) el.style.visibility = (el.style.visibility == 'visible') ? 'hidden' : 'visible';
            },

            // modifica nome gruppo
            // MODIFICA QUESTA FUNZIONE IN Conversation.vue

            async editGroupName() {
                if (!this.selectedconvs) return;
                
                let newName = prompt("Inserisci nuovo nome gruppo:", this.selectedconvs.name);
                
                if (!newName || newName.trim() === "" || newName === this.selectedconvs.name) return;

                try {

                    await this.$axios.put(
                        `/conversations/byid/${this.selectedconvs.conid}/name`, 
                        {}, 
                        { 
                            params: { name: newName }, 
                            headers: { Authorization: "Bearer " + this.token }
                        }
                    );
                    
                    this.selectedconvs.name = newName;
                    this.refresh();
                    
                } catch (e) {
                    console.error("Errore modifica:", e);

                    let serverMessage = e.response && e.response.data ? JSON.stringify(e.response.data) : "";
                    alert(`Errore Server (${e.response ? e.response.status : '?'}) ${serverMessage}`);
                }
            }



        }
    }
</script>
    
    <template>
        <div class="chat-wrapper">
            
            <header class="chat-header">
                <div class="header-content">
                    <div class="header-title text-truncate">
                        <div class="small text-secondary" v-if="login">
                            Loggato come: <strong>{{ login }}</strong>
                        </div>
                        <h5 v-if="targetUserNewChat" class="m-0 text-warning">Nuova Chat: {{ targetUserNewChat }}</h5>
                        <h5 v-else-if="selectedconvs" class="m-0">{{ selectedconvs.name }}</h5>
                        <button v-if="!targetUserNewChat" @click="editGroupName" class="btn-edit-group" title="Modifica nome">✏️</button>

                        <h5 v-else class="m-0">&nbsp;</h5>
                    </div>
                    <div class="header-actions">
                        <button class="btn btn-icon" @click="refresh" title="Refresh">↻</button>
                        <button class="btn btn-icon" @click="info" title="Info Utente">Profilo</button>
                        <button v-if="targetUserNewChat" class="btn btn-danger btn-sm ms-2" @click="cancelNewChat">Annulla</button>
                        <div v-else class="d-inline-block ms-2">
                            <button class="btn btn-outline-light btn-sm me-1" @click="startNewGroup">Gruppo</button>
                            <button class="btn btn-primary btn-sm" @click="startNewChat">Nuova Chat</button>
                        </div>
                        <button class="btn btn-outline-info btn-sm ms-2" @click="logout">Esci</button>
                    </div>
                </div>
                <ErrorMsg v-if="errormsg" :msg="errormsg" @clear-error="clearError" />
            </header>
    
            <main class="chat-body" ref="msgContainer">
                
                <div v-if="!targetUserNewChat && selectedconvs">
                    
                    <div v-if="sv && sv.limit == selectedconvs.message.length" class="load-more" @click="$emit('list-next')">
                        carica messaggi precedenti
                    </div>
                    
                    <div v-for="msg in sortedMessages" :key="msg.mesid" 
                         class="msg-row" :class="getClass(msg)">
                        
                        <div class="bubble" @contextmenu.prevent="showContextMenu($event, msg)">
                            <div v-if="msg.login !== login" class="sender-name">{{ msg.username }}</div>
                            <div v-if="msg.isforwarded > 0" class="forwarded-tag">Inoltrato</div>
    
                            <div v-if="msg.replyedtext || msg.replyedphoto" class="reply-preview">
                                <div class="reply-author">{{ msg.replyedusername }}</div>
                                <div class="reply-text">{{ msg.replyedtext }}</div>
                            </div>
    
                            <img v-if="msg.photo" :src="'data:image/jpeg;base64,'+msg.photo" class="msg-img">
                            <div v-if="msg.text" class="msg-txt">{{ msg.text }}</div>
    
                            <div class="msg-meta">
                                <span class="time">{{ formatTime(msg.timestamp) }}</span>
                                <span v-if="msg.login === login" class="ticks ms-1">
                                    <span v-if="getCheckStatus(msg) === 'read'" class="tick-read">✓✓</span>
                                    <span v-else-if="getCheckStatus(msg) === 'received'" class="tick-received">✓✓</span>
                                    <span v-else class="tick-sent">✓</span>
                                </span>
                            </div>
    
                            <div v-if="msg.comments && msg.comments.length > 0" 
                                    class="reactions-badge" 
                                    @click.stop="toggleEmojiPopup(msg.comments[0].commId)">
                                    
                                    {{ msg.comments[msg.comments.length - 1].emoji }}
                                    
                                    <div :id="'co' + msg.comments[0].commId" class="reactions-popup">
                                        <div v-for="c in msg.comments" :key="c.commId">
                                            {{ c.emoji }} {{ c.username }}
                                        </div>
                                    </div>
                            </div>
                        </div>
                    </div>
    
                    <div v-if="sv && sv.start != 0" class="load-more" @click="$emit('list-previous')">
                        torna ai messaggi recenti
                    </div>
                </div>
    
                <div v-else-if="targetUserNewChat" class="empty-state">
                    <div class="text-center text-muted">
                        <h4>Chat con {{ targetUserNewChat }}</h4>
                        <p>Scrivi il primo messaggio per iniziare.</p>
                    </div>
                </div>
                
                <div v-else class="empty-state">
                    <p>Nessuna conversazione selezionata</p>
                </div>
            </main>
    
            <footer class="chat-footer">
                <div v-if="replyToMessage" class="reply-bar">
                    <div class="text-truncate">
                        <span class="text-info fw-bold">{{ replyToMessage.username }}</span>
                        <span class="ms-2 text-muted small">{{ replyToMessage.text }}</span>
                    </div>
                    <button class="btn-close-reply" @click="replyToMessage = null">x</button>
                </div>
    
                <sendmessage 
                    :conversation="targetUserNewChat ? {name: targetUserNewChat} : selectedconvs"
                    :isEnabled="!!(targetUserNewChat || selectedconvs)"
                    @send-message="handleSend"
                />
            </footer>
    
            <div v-if="showNewChatModal" class="modal-overlay">
                <div class="modal-box">
                    <listuser :users="selectedUsersList" :showOk="false" @send-users-list="onUserSelectedForChat" @search-cancel="showNewChatModal = false"/>
                </div>
            </div>
    
            <div v-if="showGroupModal" class="modal-overlay">
                <div class="modal-box">
                    <listuser :users="selectedUsersList" :showOk="true" @send-users-list="createGroup" @search-cancel="showGroupModal = false"/>
                </div>
            </div>
    
            <div v-if="showEmojiModal" class="modal-overlay" @click.self="showEmojiModal = false">
                <div class="modal-box small-box">
                    <div class="d-flex justify-content-between mb-3">
                        <h6 class="m-0 text-light">Scegli reazione</h6>
                        <button class="btn-close-custom" @click="showEmojiModal = false">✕</button>
                    </div>
        
                    <setcomment @send-comment="sendReaction" />

                </div>
            </div>
            <div v-if="showForwardModal" class="modal-overlay">
                <div class="modal-box">
                    <div class="d-flex justify-content-between mb-2">
                        <h5 class="text-light">Inoltra a...</h5>
                        <button class="btn-close-custom" @click="showForwardModal = false">✕</button>
                    </div>
                    <listuser 
                        :users="selectedUsersList" 
                        :showOk="false" 
                        @send-users-list="executeForward" 
                        @search-cancel="showForwardModal = false"
                    />
                </div>
            </div>
    
            <vue-simple-context-menu element-id="menu-message" :options="optionsmessage" ref="vueMessageMenu" @option-clicked="onMenuOption" />
        </div>
    </template>
    
    <style>
        .vue-simple-context-menu {
                position: fixed !important; /*si fissa alla finestra non a div*/
                z-index: 10000 !important;  /* sia sopra a tutto */
                max-height: 100vh;
                overflow-y: auto;
        }
    </style>
    <style scoped>
    /*layout*/
    .chat-wrapper {
        display: flex;
        flex-direction: column;
        height: 100dvh;
        background-color: #0f172a;
        color: #f1f5f9;
        overflow: hidden;
    }
    
    .chat-header {
        flex-shrink: 0;
        padding: 10px 15px;
        background-color: #1e293b;
        border-bottom: 1px solid #334155;
    }
    .header-content { display: flex; justify-content: space-between; align-items: center; }
    
    .chat-body {
        flex-grow: 1;
        overflow-y: auto;
        padding: 15px;
        display: flex;
        flex-direction: column;
    }
    
    .chat-footer {
        flex-shrink: 0;
        padding: 10px;
        background-color: #1e293b;
        border-top: 1px solid #334155;
    }
    
    /* messaggi */
    .msg-row { display: flex; margin-bottom: 15px; width: 100%; }
    .msg-me { justify-content: flex-end; }
    .msg-friend { justify-content: flex-start; }
    
    .bubble {
        max-width: 80%;
        min-width: 120px;
        padding: 10px 14px;
        border-radius: 12px;
        position: relative;
        box-shadow: 0 1px 3px rgba(0,0,0,0.3);
        overflow: visible;  
    }
    
    .msg-me .bubble { background-color: #059669; border-bottom-right-radius: 2px; }
    .msg-friend .bubble { background-color: #334155; border-bottom-left-radius: 2px; }
    
    .sender-name { font-size: 0.8rem; font-weight: bold; color: #5eead4; margin-bottom: 4px; }
    .forwarded-tag { font-size: 0.75rem; font-style: italic; opacity: 0.7; margin-bottom: 4px; }
    
    .msg-txt { word-wrap: break-word; line-height: 1.4; }
    .msg-img { max-width: 100%; border-radius: 8px; margin-bottom: 5px; }
    
    .msg-meta {
        text-align: right; font-size: 0.7rem; opacity: 0.8; margin-top: 4px;
        display: flex; justify-content: flex-end; align-items: center;
    }
    
    /* spunte per i messaggi */
    .tick-read { color: #22d3ee; }
    .tick-received { color: #cbd5e1; }
    .tick-sent { color: #94a3b8; }
    
    /* reaction */
    .reactions-badge {
        position: absolute;
        bottom: -10px; 
        right: -5px;
        background: #1e293b;
        border: 2px solid #0f172a;
        border-radius: 12px;
        padding: 2px 6px;
        font-size: 1.1rem; 
        cursor: pointer;
        z-index: 10; 
        box-shadow: 0 2px 5px rgba(0,0,0,0.3);
    }
    .reactions-popup {
        position: absolute; bottom: 25px; right: 0; background: #334155;
        padding: 5px 10px; border-radius: 6px; visibility: hidden;
        white-space: nowrap; font-size: 0.8rem; z-index: 10;
        box-shadow: 0 5px 10px rgba(0,0,0,0.3);
    }
    
    /* anteprima "rispondi"*/
    .reply-preview {
        border-left: 3px solid #cbd5e1; padding-left: 8px; margin-bottom: 6px;
        background: rgba(0,0,0,0.2); border-radius: 0 4px 4px 0; font-size: 0.85rem;
    }
    .reply-author { font-weight: bold; color: #5eead4; }
    
    /* barra "rispondi" */
    .reply-bar {
        background: #0f172a; padding: 6px 10px; border-radius: 6px;
        margin-bottom: 8px; display: flex; justify-content: space-between; align-items: center;
        border: 1px solid #334155;
    }
    .btn-close-reply { background: none; border: none; color: #94a3b8; cursor: pointer; font-size: 1.2rem; }
    
    /*modali*/
    .modal-overlay {
        position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
        background: rgba(0,0,0,0.7); z-index: 1000;
        display: flex; justify-content: center; align-items: center;
    }
    .modal-box {
        background: #1e293b; padding: 20px; border-radius: 10px;
        width: 400px; max-width: 95%; max-height: 85vh;
        border: 1px solid #475569; display: flex; flex-direction: column;
    }
    .small-box { width: auto; min-width: 250px; }
    
    /* utilities */
    .load-more {
        text-align: center; color: #10b981; padding: 8px; cursor: pointer;
        background: rgba(16,185,129,0.1); border-radius: 6px; margin: 10px auto; width: 80%;
        font-size: 0.85rem;
    }
    .empty-state { margin: auto; color: #64748b; font-size: 1.1rem; }
    .btn-icon { background: none; border: 1px solid #475569; color: #fff; border-radius: 4px; cursor: pointer; margin-right: 5px; }
    .btn-close-custom { background: none; border: none; color: white; cursor: pointer; }

    /* emoji */
    .emoji-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr); /* 3 per riga */
    gap: 15px;
    padding: 10px;
    }

    .emoji-item {
        font-size: 2rem;
        cursor: pointer;
        text-align: center;
        transition: transform 0.2s;
        padding: 5px;
        border-radius: 50%;
    }

    .emoji-item:hover {
        transform: scale(1.3);
        background-color: rgba(255, 255, 255, 0.1);
    }
    </style>