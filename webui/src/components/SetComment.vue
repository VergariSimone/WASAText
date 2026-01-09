<script>
    export default {
        name: 'SetComment',
        emits: ['send-comment'],
        data() {
            return {
                // La lista delle emoji 
                availableEmojis: ['👍', '❤️', '😂', '😮', '😢', '😡', '👏', '🔥', '🎉']
            }
        },
        methods: {
            selectEmoji(emoji) {
                // Emette direttamente la stringa (es. "👍")
                this.$emit('send-comment', emoji);
            }
        }
    }
    </script>
    
    <template>
        <div class="emoji-container">
            <div class="emoji-grid">
                <span 
                    v-for="emoji in availableEmojis" 
                    :key="emoji" 
                    class="emoji-item" 
                    @click="selectEmoji(emoji)">
                    {{ emoji }}
                </span>
            </div>
        </div>
    </template>
    
    <style scoped>
    .emoji-container {
        padding: 10px;
        display: flex;
        justify-content: center;
    }
    
    .emoji-grid {
        display: grid;
        grid-template-columns: repeat(3, 1fr); /* 3 per riga */
        gap: 15px;
    }
    
    .emoji-item {
        font-size: 2rem;
        cursor: pointer;
        text-align: center;
        transition: transform 0.2s;
        padding: 5px;
        border-radius: 5px;
        user-select: none;
    }
    
    .emoji-item:hover {
        transform: scale(1.2);
        background-color: rgba(255, 255, 255, 0.1);
    }
    </style>