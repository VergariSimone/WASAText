<script setup>
    import { ref, nextTick, defineProps, defineEmits } from 'vue';
    import EmojiPicker from 'vue3-emoji-picker';
    import 'vue3-emoji-picker/css';
    import ImageMessage from '../components/ImageMessage.vue'; // Assicurati che questo percorso sia corretto
    
    // Props e Emits
    const props = defineProps(['conversation', 'isEnabled']);
    const emit = defineEmits(['send-message']);
    
    // Stati
    const message = ref("");
    const isEmojiPickerOpen = ref(false);
    const showImageMessage = ref(false);
    const image = ref(null);
    const inputRef = ref(null); // Riferimento all'input HTML
    
    // --- METODI ---
    
    const toggleEmoji = () => {
        if (showImageMessage.value) return;
        isEmojiPickerOpen.value = !isEmojiPickerOpen.value;
    };
    
    const toggleImageMessage = () => {
        showImageMessage.value = !showImageMessage.value;
        if(showImageMessage.value) {
            isEmojiPickerOpen.value = false;
            message.value = "";
        }
    };
    
    const onImageReady = (img, text) => {
        image.value = img;
        message.value = text || "";
        sendMessage();
    };
    
    const onSelectEmoji = (emoji) => {
        // emoji.i è il carattere (es. 😀)
        if (!emoji || !emoji.i) return;
        
        const char = emoji.i;
        const input = inputRef.value;
    
        if (input) {
            // Inserimento alla posizione del cursore
            const start = input.selectionStart || 0;
            const end = input.selectionEnd || 0;
            const text = message.value;
    
            const before = text.substring(0, start);
            const after = text.substring(end, text.length);
    
            message.value = before + char + after;
    
            // Rimette il focus e sposta il cursore dopo l'emoji
            nextTick(() => {
                input.focus();
                const newPos = start + char.length;
                input.setSelectionRange(newPos, newPos);
            });
        } else {
            // Fallback se l'input non è pronto
            message.value += char;
        }
    };
    
    const sendMessage = () => {
        // Evita invio vuoto
        if ((!message.value || message.value.trim() === "") && !image.value) return;
    
        emit('send-message', message.value, image.value);
    
        // Reset
        message.value = "";
        image.value = null;
        showImageMessage.value = false;
        isEmojiPickerOpen.value = false;
    
        // Focus
        nextTick(() => {
            if(inputRef.value) inputRef.value.focus();
        });
    };
    </script>
    
    <template>
        <div class="chatbox-input">
            
            <div v-if="isEmojiPickerOpen">
                <div class="click-overlay" @click="isEmojiPickerOpen = false"></div>
                
                <div class="emoji-popup">
                    <EmojiPicker 
                        :native="true"
                        :hide-search="false"
                        :disable-skin-tones="true"
                        :hide-group-names="true"
                        :disabled-groups='["food_drink", "activities", "travel_places", "objects", "symbols", "flags"]'
                        @select="onSelectEmoji"
                    />
                </div>
            </div>
    
            <i class="icon-btn grin" @click="toggleEmoji" title="Emoji"></i>
            <i class="icon-btn clip" @click="toggleImageMessage" title="Foto"></i>
    
            <ImageMessage 
                v-if="showImageMessage"
                @add-image-cancel="toggleImageMessage"
                @add-image="onImageReady"
            />
    
            <input 
                ref="inputRef"
                type="text" 
                class="msg-input"
                placeholder="Scrivi un messaggio..." 
                v-model="message" 
                @keyup.enter="sendMessage"
                autocomplete="off"
            >
            
            <i class="icon-btn plane" @click="sendMessage" title="Invia"></i>
        </div>
    </template>
    
    <style scoped>
    /* STESSI STILI DI PRIMA, MA SICURI */
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
    
    /* OVERLAY SFONDO */
    .click-overlay {
        position: fixed;
        top: 0; left: 0; width: 100vw; height: 100vh;
        z-index: 900;
    }
    
    /* POPUP EMOJI */
    .emoji-popup {
        position: absolute;
        bottom: 75px; /* Sopra la barra input */
        left: 10px;
        z-index: 1000;
        box-shadow: 0 10px 25px rgba(0,0,0,0.5);
        border-radius: 8px;
        overflow: hidden;
        background: #1e293b;
    }
    
    .msg-input {
        flex-grow: 1;
        padding: 10px 15px;
        border-radius: 20px;
        border: 1px solid #334155;
        background-color: #0f172a;
        color: #f1f5f9;
        outline: none;
        font-size: 0.95rem;
    }
    .msg-input:focus { border-color: #2dd4bf; }
    
    .icon-btn {
        width: 36px; height: 36px; cursor: pointer;
        background-repeat: no-repeat; background-position: center; background-size: 24px;
        filter: invert(0.7); transition: 0.2s;
    }
    .icon-btn:hover { filter: invert(1); }
    
    /* Assicurati che i file SVG esistano nella cartella public */
    .grin { background-image: url('/emoji.svg'); }
    .clip { background-image: url('/attach.svg'); }
    .plane { background-image: url('/send_message.svg'); margin-left: 5px; }
    
    /* Override CSS Picker */
    :deep(.v3-emoji-picker) {
        background: #1e293b;
        border: 1px solid #475569;
        --ep-color-bg: #1e293b;
        --ep-color-text: #f1f5f9;
        --ep-color-border: #334155;
        --ep-color-hover: #334155;
        --ep-search-bg: #0f172a;
        --ep-search-text: #f1f5f9;
        height: 300px;
        width: 300px;
    }
    </style>