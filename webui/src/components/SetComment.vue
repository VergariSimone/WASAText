<script setup>
    import { defineEmits } from 'vue';
    import EmojiPicker from 'vue3-emoji-picker';
    import 'vue3-emoji-picker/css';
    
    // Definiamo gli eventi
    const emit = defineEmits(['send-comment']);
    
    // Funzione chiamata al click sull'emoji
    const onSelectEmoji = (emoji) => {
        // Emettiamo l'oggetto emoji
        emit('send-comment', emoji);
    };
    </script>
    
    <template>
        <div class="set-comment-container">
            <EmojiPicker 
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
    .set-comment-container {
        display: flex;
        justify-content: center;
        align-items: center;
        width: 100%;
        background: transparent;
    }
    
    /* Stili per adattare il picker al tema scuro */
    :deep(.v3-emoji-picker) {
        width: 100% !important;
        max-width: 320px;
        height: 320px;
        background: transparent !important;
        border: none !important;
        box-shadow: none !important;
        
        --ep-color-bg: #1e293b;
        --ep-color-text: #f1f5f9;
        --ep-color-border: #334155;
        --ep-color-hover: #0f172a;
        --ep-color-active: #2dd4bf;
        --ep-search-bg: #0f172a;
        --ep-search-text: #f1f5f9;
        --ep-search-border: #334155;
    }
    </style>