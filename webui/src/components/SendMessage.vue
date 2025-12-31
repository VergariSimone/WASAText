<script>
    import { ref, nextTick } from 'vue';
    import EmojiPicker from 'vue3-emoji-picker';
    import 'vue3-emoji-picker/css';
    import ImageMessage from '../components/ImageMessage.vue';
    
    export default {
        name: 'SendMessage',
        components: {
            'emojipicker': EmojiPicker,
            'imagemessage': ImageMessage,
        },
        props: [
            'conversation', // Può essere oggetto conversazione o {name: 'utente'}
            'isEnabled'     // Opzionale, per disabilitare esternamente
        ],
        emits: [
            'send-message',
        ],
        setup(props, { emit }) {
            
            //State
            const message = ref("");
            const isEmojiPickerOpen = ref(false);
            const showImageMessage = ref(false);
            const image = ref(null);
    
            // Funzioni
            const toggleEmoji = () => {
                isEmojiPickerOpen.value = !isEmojiPickerOpen.value;
            };
    
            const toggleImageMessage = () => {
                showImageMessage.value = !showImageMessage.value;
                // Se apro immagine, pulisco testo e chiudo emoji
                if(showImageMessage.value) {
                    message.value = "";
                    isEmojiPickerOpen.value = false;
                }
            };
    
            const sendMessage = () => {
                // Validazione base
                if (!message.value && !image.value) return;
    
                emit('send-message', message.value, image.value);
                
                // Reset
                message.value = "";
                image.value = null;
                showImageMessage.value = false;
                isEmojiPickerOpen.value = false;
            };
    
            const onImageReady = (img, text) => {
                // Callback dal componente ImageMessage quando l'utente conferma l'upload
                image.value = img;
                message.value = text || "";
                sendMessage();
            };
    
            const onSelectEmoji = (emoji) => {
                // Aggiunge emoji alla posizione del cursore
                const input = document.getElementById("inputRef");
                if (!input) {
                    message.value += emoji.i;
                    return;
                }
    
                const start = input.selectionStart;
                const end = input.selectionEnd;
                const text = message.value;
                
                const before = text.substring(0, start);
                const after = text.substring(end, text.length);
                
                message.value = before + emoji.i + after;
                
                // Riposiziona cursore
                nextTick(() => {
                    input.focus();
                    input.setSelectionRange(start + emoji.i.length, start + emoji.i.length);
                });
            };
    
            return {
                message,
                isEmojiPickerOpen,
                showImageMessage,
                image,
                toggleEmoji,
                toggleImageMessage,
                sendMessage,
                onSelectEmoji,
                onImageReady
            };
        }
    }
    </script>
    
    <template>
        <div class="chatbox-input">
            
            <div v-if="isEmojiPickerOpen" class="emoji-wrapper">
                 <div class="click-outside" @click="isEmojiPickerOpen = false"></div>
                
                <emojipicker 
                    :native="true"
                    :hide-search="false"
                    :disable-skin-tones="true"
                    :hide-group-names="true"
                    :disabled-groups='["food_drink", "activities", "travel_places", "objects", "symbols", "flags"]' 
                    @select="onSelectEmoji"
                    class="picker-custom"
                />
            </div>
    
            <i class="grin" @click="toggleEmoji" title="Emoji"></i>
    
            <i class="paperclip" @click="toggleImageMessage" title="Allega Immagine"></i>
            
            <imagemessage 
                v-if="showImageMessage"
                @add-image-cancel="toggleImageMessage"
                @add-image="onImageReady"
            />
    
            <input 
                id="inputRef" 
                type="text" 
                placeholder="Scrivi un messaggio..." 
                v-model="message" 
                @keyup.enter="sendMessage"
                autocomplete="off"
            >
            
            <i class="paperplane" @click="sendMessage" title="Invia"></i>
        </div>
    </template>
    
    <style scoped>
    /* --- CONTENITORE --- */
    .chatbox-input {
        position: relative;
        width: 100%;
        min-height: 60px;
        background-color: #1e293b;
        border-top: 1px solid #334155;
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 10px;
        z-index: 50;
    }
    
    /* --- ICONE --- */
    .chatbox-input i {
        cursor: pointer;
        width: 40px;
        height: 40px;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;
        transition: background 0.2s;
        flex-shrink: 0;
    }
    
    .chatbox-input i:hover {
        background-color: rgba(255, 255, 255, 0.1);
    }
    
    /* Icone SVG tramite Pseudo-elementi */
    .grin::after { content: url('/emoji.svg'); filter: invert(0.6); width: 24px; }
    .paperclip::after { content: url('/attach.svg'); filter: invert(0.6); width: 24px; }
    .paperplane::after { content: url('/send_message.svg'); filter: invert(0.6); width: 24px; margin-left: 4px; }
    
    /* Hover effect sulle icone (Verde Acqua) */
    .chatbox-input i:hover::after {
        filter: invert(73%) sepia(35%) saturate(827%) hue-rotate(124deg) brightness(93%) contrast(88%);
    }
    
    /* --- INPUT TEXT --- */
    .chatbox-input input {
        flex-grow: 1;
        padding: 10px 15px;
        border-radius: 20px;
        font-size: 0.95em;
        border: 1px solid #334155;
        background-color: #0f172a;
        color: #f1f5f9;
        outline: none;
    }
    .chatbox-input input:focus {
        border-color: #2dd4bf;
    }
    
    /* --- EMOJI PICKER WRAPPER --- */
    .emoji-wrapper {
        position: absolute;
        bottom: 70px; /* Sopra la barra input */
        left: 10px;
        z-index: 100;
    }
    
    /* Sfondo invisibile per chiudere cliccando fuori */
    .click-outside {
        position: fixed;
        top: 0; left: 0; width: 100vw; height: 100vh;
        z-index: -1; 
    }
    
    /* Stili libreria picker (deep selector) */
    :deep(.v3-emoji-picker) {
        background: #0f172a;
        border: 1px solid #334155;
        --ep-color-bg: #0f172a;
        --ep-color-text: #f1f5f9;
        --ep-color-border: #334155;
        --ep-color-hover: #1e293b;
        box-shadow: 0 5px 15px rgba(0,0,0,0.5);
        height: 300px;
        width: 300px;
    }
    </style>