<script>
    import ImageMessage from '../components/ImageMessage.vue';
    
    export default {
        name: 'SendMessage',
        components: {
            ImageMessage
        },
        props: ['conversation', 'isEnabled'],
        emits: ['send-message'],
        data() {
            return {
                text: "",
                showEmojiPicker: false,
                showImageMessage: false, 
                
                // lista di emoji
                emojis: [
                    '😀', '😂', '😍', '🥺', '😎', '😭', '😡', '👍', '👎', '❤️', 
                    '🔥', '🎉', '👋', '🙏', '👀', '💩', '👻', '🤖', '💀', '🤡',
                    '🌟', '🍕', '🍺', '🇮🇹', '🚀', '🙈', '🙉', '🙊'
                ]
            }
        },
        methods: {
            // invio normale 
            sendMessage() {
                // se non c'è testo e non c'è modale foto aperto, non fare nulla
                if (!this.text.trim()) return;
    
                // emette solo testo 
                this.$emit('send-message', this.text, null);
                
                this.text = "";
                this.showEmojiPicker = false;
                
                // focus
                this.$nextTick(() => {
                    if(this.$refs.msgInput) this.$refs.msgInput.focus();
                });
            },
    
            // gestise emoji
            toggleEmoji() {
                this.showImageMessage = false; 
                this.showEmojiPicker = !this.showEmojiPicker;
            },
    
            addEmoji(emoji) {
                this.text += emoji;
                this.$nextTick(() => {
                    if(this.$refs.msgInput) this.$refs.msgInput.focus();
                });
            },
    
            // gestione modale foto
            toggleImageMessage() {
                this.showEmojiPicker = false; // chiude emoji
                this.showImageMessage = !this.showImageMessage;
                if (this.showImageMessage) {
                    this.text = ""; // pulizia testo
                }
            },
    
            // callback che riceve i dati dal componente ImageMessage
            onImageReady(img, caption) {
                // emette l'evento con immagine e testo
                this.$emit('send-message', caption, img);
                
                // chiude e resetta
                this.showImageMessage = false;
                this.text = "";
            },
    
            // gestione tastiera
            onKeydown(e) {
                if (e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    this.sendMessage();
                }
            }
        }
    }
    </script>
    
    <template>
        <div class="chatbox-input">
            
            <div v-if="showEmojiPicker" class="emoji-picker-box">
                <div class="click-overlay" @click="showEmojiPicker = false"></div>
                
                <div class="emoji-content">
                    <div class="emoji-header">
                        <span>Scegli Emoji</span>
                        <button class="btn-close-small" @click="showEmojiPicker = false">✕</button>
                    </div>
                    <div class="emoji-grid">
                        <span v-for="e in emojis" :key="e" class="emoji-item" @click="addEmoji(e)">
                            {{ e }}
                        </span>
                    </div>
                </div>
            </div>
    
            <button class="icon-btn" @click="toggleEmoji" :disabled="!isEnabled" title="Emoji">
                <img src="/emoji.svg" class="icon-img" alt="emoji">
            </button>
    
            <button class="icon-btn" @click="toggleImageMessage" :disabled="!isEnabled" title="Foto">
                <img src="/attach.svg" class="icon-img" alt="foto">
            </button>
    
            <ImageMessage 
                v-if="showImageMessage"
                @add-image-cancel="toggleImageMessage"
                @add-image="onImageReady"
            />
    
            <textarea 
                ref="msgInput"
                v-model="text" 
                class="msg-input" 
                placeholder="Scrivi un messaggio..." 
                :disabled="!isEnabled"
                @keydown="onKeydown"
                rows="1"
            ></textarea>
            
            <button class="icon-btn send-btn" @click="sendMessage" :disabled="!isEnabled || !text.trim()">
                <img src="/send_message.svg" class="icon-img-send" alt="invia">
            </button>
    
        </div>
    </template>
    
    <style scoped>
    .chatbox-input {
        position: relative;
        width: 100%;
        min-height: 60px;
        background-color: #1e293b;
        border-top: 1px solid #334155;
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px 15px;
        z-index: 50; 
    }
    
    /* input*/
    .msg-input {
        flex-grow: 1;
        background: #0f172a;
        border: 1px solid #334155;
        border-radius: 20px;
        color: #f1f5f9;
        padding: 10px 15px;
        outline: none;
        resize: none;
        max-height: 100px;
        font-family: inherit;
        font-size: 0.95rem;
    }
    .msg-input:focus { border-color: #2dd4bf; }
    .msg-input:disabled { opacity: 0.6; cursor: not-allowed; }
    
    /* bottone icone */
    .icon-btn {
        background: none;
        border: none;
        cursor: pointer;
        padding: 5px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        transition: background 0.2s;
    }
    .icon-btn:hover:not(:disabled) { background-color: rgba(255, 255, 255, 0.1); }
    .icon-btn:disabled { opacity: 0.5; cursor: default; }
    
    .icon-img {
        width: 24px;
        height: 24px;
        display: block;
        filter: invert(0.8); /* inverte colori se icone scure */
        transition: filter 0.2s;
    }
    .icon-btn:hover .icon-img { filter: invert(1); }
    
    /* bottone invia */
    .send-btn {
        background-color: #0ea5e9;
        width: 40px;
        height: 40px;
    }
    .send-btn:hover:not(:disabled) { background-color: #0284c7; }
    .icon-img-send {
        width: 20px;
        height: 20px;
        filter: invert(1);
        margin-left: -2px;
    }
    
    /* popup emoji*/
    .emoji-picker-box {
        position: absolute;
        bottom: 70px;
        left: 10px;
        z-index: 1000;
    }
    
    .click-overlay {
        position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
        z-index: 999;
    }
    
    .emoji-content {
        position: relative;
        z-index: 1001;
        width: 280px;
        background-color: #1e293b;
        border: 1px solid #475569;
        border-radius: 12px;
        padding: 10px;
        box-shadow: 0 10px 25px rgba(0,0,0,0.5);
    }
    
    .emoji-header {
        display: flex; justify-content: space-between; align-items: center;
        margin-bottom: 10px; color: #94a3b8; font-size: 0.85rem;
        border-bottom: 1px solid #334155; padding-bottom: 5px;
    }
    .btn-close-small { background: none; border: none; color: #cbd5e1; cursor: pointer; font-size: 1rem; }
    
    .emoji-grid {
        display: grid;
        grid-template-columns: repeat(6, 1fr);
        gap: 8px;
        max-height: 200px;
        overflow-y: auto;
    }
    
    .emoji-item {
        font-size: 1.4rem;
        cursor: pointer;
        text-align: center;
        padding: 4px;
        border-radius: 4px;
        transition: background 0.2s;
    }
    .emoji-item:hover { background-color: rgba(255, 255, 255, 0.1); }
    </style>