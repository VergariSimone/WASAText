<script>
    import { ref } from 'vue';
    import EmojiPicker from 'vue3-emoji-picker';
    import 'vue3-emoji-picker/css';
    
    export default {
        name: 'SetComment',
        components: {
            'emojipicker': EmojiPicker,
        },
        // Non servono più msgid/conid/coords qui, li gestisce il padre (Chat.vue)
        emits: [
            'send-comment'
        ],
        setup(props, { emit }) {
            
            const onSelectEmoji = (emoji) => {
                // Emettiamo l'oggetto emoji direttamente
                // Il padre (Chat.vue) sa già a quale messaggio riferirsi
                emit('send-comment', emoji);
            };
    
            return {
                onSelectEmoji
            };
        }
    }
    </script>
    
    <template>
        <div class="set-comment-container">
            <emojipicker 
                :native="true"
                :hide-search="false"
                :disable-skin-tones="true"
                :hide-group-names="true"
                :disabled-groups='["food_drink", "activities", "travel_places", "objects", "symbols", "flags"]'
                @select="onSelectEmoji"
                class="custom-picker"
            />
        </div>
    </template>
    
    <style scoped>
    /* --- CONTENITORE --- */
    .set-comment-container {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 100%;
        background-color: #1e293b;
        border-radius: 8px;
    }
    
    /* --- OVERRIDE STILI LIBRERIA (DARK MODE PERFETTA) --- */
    :deep(.v3-emoji-picker) {
        width: 100% !important;
        max-width: 350px;
        height: 300px; /* Altezza fissa per evitare salti */
        border: none !important;
        border-radius: 8px;
        background: transparent !important;
        box-shadow: none !important;
    
        /* Variabili CSS della libreria forzate al tuo tema */
        --ep-color-bg: #1e293b;       /* Sfondo uguale alla modale */
        --ep-color-text: #f1f5f9;     /* Testo bianco */
        --ep-color-border: #334155;   /* Bordi */
        --ep-color-hover: #0f172a;    /* Hover più scuro */
        --ep-color-active: #2dd4bf;   /* Verde acqua */
        --ep-search-bg: #0f172a;      /* Barra ricerca scura */
        --ep-search-text: #f1f5f9;
        --ep-search-border: #334155;
    }
    
    /* Nasconde scrollbar brutte se appaiono */
    :deep(.v3-emoji-picker .v3-body-inner::-webkit-scrollbar) {
        width: 4px;
    }
    :deep(.v3-emoji-picker .v3-body-inner::-webkit-scrollbar-thumb) {
        background-color: #334155;
        border-radius: 4px;
    }
    </style>