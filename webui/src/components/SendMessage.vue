<script>
    import EmojiPicker from "vue3-emoji-picker";
    import "vue3-emoji-picker/css";
    import ImageMessage from "./ImageMessage.vue";
    
    export default {
      name: "SendMessage",
    
      components: {
        EmojiPicker,
        ImageMessage,
      },
    
      props: {
        conversation: {
          type: Object,
          required: true,
        },
      },
    
      emits: ["send-message"],
    
      data() {
        return {
          message: "",
          image: null,
          isEnabled: true,
          showImageMessage: false,
          showEmojiPicker: false,
        };
      },
    
      methods: {
        toggleEmoji() {
          if (!this.isEnabled && !this.showImageMessage) return;
          this.showEmojiPicker = !this.showEmojiPicker;
        },
    
        toggleImageMessage() {
          if (!this.isEnabled && !this.showImageMessage) return;
    
          this.showImageMessage = !this.showImageMessage;
          this.showEmojiPicker = false;
    
          if (this.showImageMessage) {
            this.message = "";
            this.isEnabled = false;
          } else {
            this.isEnabled = true;
          }
        },
    
        addEmoji(emoji) {
          const input = this.$refs.input;
          const cursor = input.selectionStart ?? this.message.length;
    
          this.message =
            this.message.slice(0, cursor) +
            emoji.i +
            this.message.slice(cursor);
    
          this.$nextTick(() => {
            const pos = cursor + emoji.i.length;
            input.setSelectionRange(pos, pos);
            input.focus();
          });
        },
    
        send() {
          if (!this.isEnabled || this.showImageMessage) return;
          if (!this.message && !this.image) return;
    
          this.$emit("send-message", this.message, this.image);
    
          this.message = "";
          this.image = null;
        },
    
        onImageCreated(image, text) {
          this.image = image;
          this.message = text;
          this.showImageMessage = false;
          this.isEnabled = true;
          this.send();
        },
      },
    };
    </script>
    
    <template>
      <div class="chatbox-input" v-if="conversation">
        <!-- Emoji -->
        <div class="icon" @click="toggleEmoji">
          😊
          <div v-if="showEmojiPicker" class="emoji-overlay">
            <EmojiPicker
              :native="true"
              :hide-search="false"
              :disable-skin-tones="true"
              :hide-group-names="true"
              :disabled-groups="[
                'food_drink',
                'activities',
                'travel_places',
                'objects',
                'symbols',
                'flags',
              ]"
              @select="addEmoji"
            />
          </div>
        </div>
    
        <!-- Attachment -->
        <div class="icon" @click="toggleImageMessage">📎</div>
    
        <!-- Image message -->
        <ImageMessage
          v-if="showImageMessage"
          @add-image="onImageCreated"
          @add-image-cancel="toggleImageMessage"
        />
    
        <!-- Input -->
        <input
          ref="input"
          v-model="message"
          :disabled="!isEnabled"
          placeholder="Type a message"
          @keydown.enter="send"
        />
    
        <!-- Send -->
        <div class="icon send" @click="send">➤</div>
      </div>
    </template>
    
    <style scoped>
    .chatbox-input {
      display: grid;
      grid-template-columns: 32px 32px 1fr 32px;
      align-items: center;
      padding: 10px;
      background: #e6f7fb;
      border-top: 1px solid #7fdad3;
      position: relative;
    }
    
    /* Icons */
    .icon {
      font-size: 1.4rem;
      cursor: pointer;
      text-align: center;
      color: #2bb0a6;
    }
    
    .icon.send {
      color: white;
      background: #2bb0a6;
      border-radius: 50%;
      width: 32px;
      height: 32px;
      line-height: 32px;
    }
    
    /* Input */
    input {
      border-radius: 20px;
      border: 1px solid #7fdad3;
      padding: 10px 14px;
      outline: none;
      font-size: 0.95rem;
    }
    
    input:disabled {
      background: #f1f1f1;
    }
    
    /* Emoji picker */
    .emoji-overlay {
      position: absolute;
      bottom: 60px;
      left: 0;
      z-index: 1000;
    }
    </style>
    