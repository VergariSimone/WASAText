<script>
    import EmojiPicker from "vue3-emoji-picker";
    import "vue3-emoji-picker/css";
    
    export default {
      name: "SetComment",
    
      components: {
        EmojiPicker,
      },
    
      props: {
        msgid: {
          type: Number,
          required: true,
        },
        conid: {
          type: Number,
          required: true,
        },
        topX: Number,
        topY: Number,
      },
    
      emits: ["send-comment"],
    
      methods: {
        selectEmoji(emoji) {
          this.$emit("send-comment", emoji.i, this.msgid, this.conid);
        },
    
        close() {
          this.$emit("send-comment", null, null, null);
        },
      },
    };
    </script>
    
    <template>
      <div
        class="reaction-container"
        :style="{ top: `${topY}px`, left: `${topX}px` }"
      >
        <div class="reaction-header">
          <span class="close" @click="close">✕</span>
        </div>
    
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
          @select="selectEmoji"
        />
      </div>
    </template>
    
    <style scoped>
    /* Container */
    .reaction-container {
      position: fixed;
      z-index: 1000;
      background: white;
      border-radius: 14px;
      padding: 8px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.15);
      border: 1px solid #7fdad3;
    }
    
    /* Header */
    .reaction-header {
      display: flex;
      justify-content: flex-end;
      margin-bottom: 6px;
    }
    
    /* Close button */
    .close {
      cursor: pointer;
      font-size: 0.9rem;
      color: #2bb0a6;
      font-weight: bold;
    }
    
    /* Emoji picker tweaks */
    .v3-emoji-picker {
      box-shadow: none;
      background: transparent;
    }
</style>
    