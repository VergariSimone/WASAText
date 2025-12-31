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
            "selectedconvs",
            "sv", // Scroll/Pagination info
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
        data() {
            return {
                // --- STATI MODALI ---
                showEmojiModal: false,
                showNewChatModal: false, // Per nuova chat singola
                showGroupModal: false,   // Per nuovi gruppi
    
                // --- STATI DATI ---
                targetUserNewChat: null, // Utente selezionato per nuova chat
                selectedUsersList: [],   // Lista utenti caricata dall'API per le modali
                
                // --- STATI MESSAGGIO ---
                selectedMessage: null,   // Messaggio su cui si sta agendo (es. click destro)
                replyToMessage: null,    // Messaggio a cui si sta rispondendo
                
                errormsg: null,
    
                // --- MENU CONTESTUALE ---
                optionsmessage: [
                    {name: 'Reagisci', slug: 'comment'},
                    {name: 'Rispondi', slug: 'reply'},
                    {name: 'Inoltra', slug: 'forward'}, // (Logica base)
                    {type: 'divider'},
                    {name: 'Elimina', slug: 'delete'},
                ],
            }
        },
    
        mounted() {
            this.refresh();
            this.scrollToBottom();
        },
    
        watch: {
            // Scrolla in basso quando arrivano nuovi messaggi
            'selectedconvs.message': {
                handler() { this.scrollToBottom(); },
                deep: true
            }
        },
    
        methods: {
            // ================= UTILITIES =================
            refresh() { this.$emit('refresh'); },
            logout() { this.$emit('logout'); },
            info() { this.$emit('show-user-info'); },
            clearError() { this.errormsg = null; },
            
            scrollToBottom() {
                this.$nextTick(() => {
                    const el = this.$refs.msgContainer;
                    if (el) el.scrollTop = el.scrollHeight;
                });
            },
            
            formatTime(ts) {
                if (!ts) return '';
                return new Date(ts + 'Z').toLocaleTimeString('it-IT', { hour: '2-digit', minute: '2-digit' });
            },
    
            // ================= GESTIONE NUOVA CHAT =================
            async startNewChat() {
                this.showNewChatModal = true;
                this.targetUserNewChat = null;
                this.selectedUsersList = [];
                
                try {
                    let res = await this.$axios.get(`/users?search=`, { 
                        headers: { Authorization: "Bearer " + this.token } 
                    });
                    // Normalizza in array
                    this.selectedUsersList = Array.isArray(res.data) ? res.data : [res.data];
                } catch (e) {
                    this.errormsg = "Errore caricamento utenti";
                }
            },
    
            onUserSelectedForChat(users) {
                // ListUser restituisce un array. Prendiamo il primo elemento.
                let user = Array.isArray(users) ? users[0] : users;
                
                if (!user) return;
                if (user === this.login) { 
                    this.errormsg = "Non puoi scriverti da solo"; 
                    return; 
                }
    
                this.targetUserNewChat = user; // Salviamo la stringa "mario"
                this.showNewChatModal = false; // Chiudiamo la modale
                
                // Focus automatico input (opzionale, richiede accesso al ref figlio)
            },
    
            cancelNewChat() {
                this.showNewChatModal = false;
                this.targetUserNewChat = null;
            },
    
            // ================= GESTIONE GRUPPI =================
            async startNewGroup() {
                this.showGroupModal = true;
                try {
                    let res = await this.$axios.get(`/users?search=`, { 
                        headers: { Authorization: "Bearer " + this.token } 
                    });
                    this.selectedUsersList = Array.isArray(res.data) ? res.data : [res.data];
                } catch (e) {}
            },
    
            async createGroup(users) {
                this.showGroupModal = false;
                var out = [];
                for (var key in users) out.push('users=' + encodeURIComponent(users[key]));
                
                try {
                    let res = await this.$axios.put(`/conversations/group?${out.join('&')}`, {}, { 
                        headers: { Authorization: "Bearer " + this.token } 
                    });
                    this.$emit('change-selected-conv-byid', res.data);
                } catch (e) { this.errormsg = e.toString(); }
            },
    
            // ================= GESTIONE INVIO MESSAGGI =================
            async handleSend(text, image) {
                this.errormsg = null;
                
                // CASO 1: Nuova Chat (Utente selezionato ma conversazione non ancora esistente)
                if (this.targetUserNewChat) {
                    await this.sendToNewUser(text, image);
                } 
                // CASO 2: Chat Esistente
                else if (this.selectedconvs) {
                    await this.sendToExistingChat(text, image);
                } else {
                    this.errormsg = "Errore: Nessuna chat attiva.";
                }
            },
    
            async sendToNewUser(text, image) {
                let formData = new FormData();
                if (text) formData.append('text', text);
                if (image) formData.append('image', image, 'image');
    
                try {
                    // API call verso l'utente target
                    let res = await this.$axios.post(`/conversations?conversation=${this.targetUserNewChat}`, formData, {
                        headers: { Authorization: "Bearer " + this.token }
                    });
                    
                    // Successo: Resettiamo lo stato "Nuova Chat"
                    this.targetUserNewChat = null;
                    // Emettiamo la nuova conversazione arrivata dal server
                    this.$emit('change-selected-conv', res.data);
                    this.refresh();
                } catch (e) {
                    this.errormsg = "Errore invio: " + e.toString();
                }
            },
    
            async sendToExistingChat(text, image) {
                let formData = new FormData();
                if (text) formData.append('text', text);
                if (image) formData.append('image', image, 'image');
    
                let url = `/conversations?conversation=${this.selectedconvs.conid}`;
                if (this.replyToMessage) url += `&reply=${this.replyToMessage.mesid}`;
    
                try {
                    let res = await this.$axios.post(url, formData, {
                        headers: { Authorization: "Bearer " + this.token }
                    });
                    this.$emit('change-selected-conv', res.data);
                    this.replyToMessage = null; // Chiudi reply bar
                } catch (e) {
                    this.errormsg = "Errore invio: " + e.toString();
                }
            },
    
            // ================= GESTIONE MENU & AZIONI =================
            showContextMenu(e, msg) {
                this.selectedMessage = msg;
                this.$refs.vueMessageMenu.showMenu(e, msg);
            },
    
            onMenuOption(event) {
                const { slug } = event.option;
                const msg = event.item;
    
                if (slug === 'comment') {
                    this.selectedMessage = msg; // Salviamo il msg per la modale emoji
                    this.showEmojiModal = true;
                } else if (slug === 'reply') {
                    this.replyToMessage = msg;
                    // Scroll up o focus input opzionale
                } else if (slug === 'delete') {
                    this.deleteMsg(msg);
                } else if (slug === 'forward') {
                    // Semplificazione: mostra alert o logica complessa
                    alert("Funzionalità Inoltra da implementare nella modale utenti");
                }
            },
    
            async sendReaction(emojiData) {
                // setComment emette l'emoji. selectedMessage è salvato nello stato.
                if (!this.selectedMessage) return;
                this.showEmojiModal = false;
                
                try {
                    await this.$axios.post(
                        `/conversations/byid/${this.selectedMessage.conid}/message/${this.selectedMessage.mesid}/comment`, 
                        emojiData.i, // emojiData.i contiene il carattere emoji
                        { headers: { Authorization: "Bearer " + this.token, 'content-type': 'text/plain' } }
                    );
                    this.refresh();
                } catch (e) { this.errormsg = e.toString(); }
            },
    
            async deleteMsg(msg) {
                 try {
                    await this.$axios.delete(`/conversations/byid/${msg.conid}/message/${msg.mesid}`, {
                         headers: { Authorization: "Bearer " + this.token }
                    });
                    this.refresh();
                } catch(e) { 
                    if(e.response && e.response.status === 403) this.errormsg = "Non puoi eliminare messaggi altrui";
                    else this.errormsg = "Errore eliminazione"; 
                }
            },
    
            // ================= UI HELPERS =================
            getClass(msg) {
                return msg.login === this.login ? 'msg-me' : 'msg-friend';
            },
            
            getCheckStatus(msg) {
                if (msg.login !== this.login) return null;
                // 1. Letto (Blu/Verde)
                let read = msg.readarray && msg.readarray.some(u => u.name !== this.login && u.read === 'true');
                if (read) return 'read'; 
                // 2. Ricevuto (Doppio Grigio)
                let rcv = msg.receivedarray && msg.receivedarray.some(u => u.name !== this.login && u.received === 'true');
                if (rcv) return 'received'; 
                // 3. Inviato (Singolo)
                return 'sent'; 
            },
    
            // Gestione popup emoji esistenti (chi ha reagito)
            toggleEmojiPopup(commId) {
                let el = document.getElementById('co' + commId);
                if(el) el.style.visibility = (el.style.visibility == 'visible') ? 'hidden' : 'visible';
            }
        }
    }
    </script>
    
    <template>
        <div class="chat-wrapper">
            
            <header class="chat-header">
                <div class="header-content">
                    <div class="header-title text-truncate">
                        <h5 v-if="targetUserNewChat" class="m-0 text-warning">Nuova Chat: {{ targetUserNewChat }}</h5>
                        <h5 v-else-if="selectedconvs" class="m-0">{{ selectedconvs.name }}</h5>
                        <h5 v-else class="m-0">&nbsp;</h5>
                    </div>
    
                    <div class="header-actions">
                        <button class="btn btn-icon" @click="refresh" title="Refresh">↻</button>
                        <button class="btn btn-icon" @click="info" title="Info Utente">ℹ️</button>
                        
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
                    
                    <div v-if="sv && sv.start != 0" class="load-more" @click="$emit('list-previous')">
                        ▲ Carica precedenti
                    </div>
                    
                    <div v-for="msg in selectedconvs.message" :key="msg.mesid" 
                         class="msg-row" :class="getClass(msg)">
                        
                        <div class="bubble" @contextmenu.prevent="showContextMenu($event, msg)">
                            
                            <div v-if="msg.login !== login" class="sender-name">{{ msg.username }}</div>
                            <div v-if="msg.isforwarded > 0" class="forwarded-tag">↩ Inoltrato</div>
    
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
    
                            <div v-if="msg.comments" class="reactions-badge" @click.stop="toggleEmojiPopup(msg.comments[0].commId)">
                                {{ msg.comments[0].emoji }}
                                <div :id="'co' + msg.comments[0].commId" class="reactions-popup">
                                    <div v-for="c in msg.comments" :key="c.commId">{{ c.emoji }} {{ c.username }}</div>
                                </div>
                            </div>
                        </div>
                    </div>
    
                    <div v-if="sv && sv.limit == selectedconvs.message.length" class="load-more" @click="$emit('list-next')">
                        ▼ Carica successivi
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
                    <button class="btn-close-reply" @click="replyToMessage = null">✕</button>
                </div>
    
                <sendmessage 
                    :conversation="targetUserNewChat ? {name: targetUserNewChat} : selectedconvs"
                    :isEnabled="!!(targetUserNewChat || selectedconvs)"
                    @send-message="handleSend"
                />
            </footer>
    
            <div v-if="showNewChatModal" class="modal-overlay">
                <div class="modal-box">
                    <listuser 
                        :users="selectedUsersList" 
                        :showOk="false" 
                        @send-users-list="onUserSelectedForChat" 
                        @search-cancel="showNewChatModal = false"
                    />
                </div>
            </div>
    
            <div v-if="showGroupModal" class="modal-overlay">
                <div class="modal-box">
                    <listuser 
                        :users="selectedUsersList" 
                        :showOk="true" 
                        @send-users-list="createGroup" 
                        @search-cancel="showGroupModal = false"
                    />
                </div>
            </div>
    
            <div v-if="showEmojiModal" class="modal-overlay">
                <div class="modal-box small-box">
                    <div class="d-flex justify-content-between mb-2">
                        <h6 class="m-0 text-light">Reagisci</h6>
                        <button class="btn-close-custom" @click="showEmojiModal = false">✕</button>
                    </div>
                    <setcomment @send-comment="sendReaction" />
                </div>
            </div>
    
            <vue-simple-context-menu element-id="menu-message" :options="optionsmessage" ref="vueMessageMenu" @option-clicked="onMenuOption" />
    
        </div>
    </template>
    
    <style scoped>
    /* ================= LAYOUT ================= */
    .chat-wrapper {
        display: flex;
        flex-direction: column;
        height: 100vh; /* O 100dvh se supportato */
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
        overflow-y: auto; /* Scroll principale */
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
    
    /* ================= MESSAGGI ================= */
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
    
    /* Spunte */
    .tick-read { color: #22d3ee; }
    .tick-received { color: #cbd5e1; }
    .tick-sent { color: #94a3b8; }
    
    /* Reazioni */
    .reactions-badge {
        position: absolute; bottom: -12px; right: -5px;
        background: #1e293b; border: 2px solid #0f172a; border-radius: 12px;
        padding: 2px 6px; font-size: 13px; cursor: pointer; z-index: 5;
        box-shadow: 0 2px 4px rgba(0,0,0,0.2);
    }
    .reactions-popup {
        position: absolute; bottom: 25px; right: 0; background: #334155;
        padding: 5px 10px; border-radius: 6px; visibility: hidden;
        white-space: nowrap; font-size: 0.8rem; z-index: 10;
        box-shadow: 0 5px 10px rgba(0,0,0,0.3);
    }
    
    /* Reply Preview */
    .reply-preview {
        border-left: 3px solid #cbd5e1; padding-left: 8px; margin-bottom: 6px;
        background: rgba(0,0,0,0.2); border-radius: 0 4px 4px 0; font-size: 0.85rem;
    }
    .reply-author { font-weight: bold; color: #5eead4; }
    
    /* Reply Bar (Footer) */
    .reply-bar {
        background: #0f172a; padding: 6px 10px; border-radius: 6px;
        margin-bottom: 8px; display: flex; justify-content: space-between; align-items: center;
        border: 1px solid #334155;
    }
    .btn-close-reply { background: none; border: none; color: #94a3b8; cursor: pointer; font-size: 1.2rem; }
    
    /* ================= MODALI ================= */
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
    
    /* ================= UTILS ================= */
    .load-more {
        text-align: center; color: #10b981; padding: 8px; cursor: pointer;
        background: rgba(16,185,129,0.1); border-radius: 6px; margin: 10px auto; width: 80%;
        font-size: 0.85rem;
    }
    .empty-state { margin: auto; color: #64748b; font-size: 1.1rem; }
    .btn-icon { background: none; border: 1px solid #475569; color: #fff; border-radius: 4px; cursor: pointer; margin-right: 5px; }
    .btn-close-custom { background: none; border: none; color: white; cursor: pointer; }
    </style>