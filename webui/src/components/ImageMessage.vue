<script>
    export default {
      name: "ImageMessage",
      emits: ["add-image", "add-image-cancel"],
    
      data() {
        return {
          imageFile: null,
          previewUrl: null,
          caption: "",
        };
      },
    
      computed: {
        hasContent() {
          return !!this.imageFile || this.caption.trim() !== "";
        },
      },
    
      methods: {
        openFilePicker() {
          this.$refs.fileInput.click();
        },
    
        onFileSelected(event) {
          const file = event.target.files?.[0];
          if (!file) return;
    
          this.imageFile = file;
          this.previewUrl = URL.createObjectURL(file);
        },
    
        send() {
          if (this.hasContent) {
            this.$emit("add-image", this.imageFile, this.caption);
          } else {
            this.cancel();
          }
        },
    
        cancel() {
          this.reset();
          this.$emit("add-image-cancel");
        },
    
        reset() {
          this.imageFile = null;
          this.previewUrl = null;
          this.caption = "";
        },
      },
    };
    </script>
    
    <template>
      <div class="image-message">
        <!-- Toolbar -->
        <div class="toolbar">
          <button class="btn ghost" @click="cancel">Cancel</button>
          <button class="btn ghost" @click="openFilePicker">Upload</button>
          <button class="btn primary" @click="send">Send</button>
        </div>
    
        <!-- Preview -->
        <div class="preview">
          <div class="image-box" :class="{ filled: previewUrl }">
            <img v-if="previewUrl" :src="previewUrl" alt="preview" />
            <span v-else class="placeholder">No image selected</span>
          </div>
    
          <input
            ref="fileInput"
            type="file"
            accept="image/gif,image/png,image/jpeg"
            hidden
            @change="onFileSelected"
          />
        </div>
    
        <!-- Caption -->
        <input
          v-model="caption"
          class="caption"
          type="text"
          placeholder="Add a caption…"
        />
      </div>
    </template>
    
    <style scoped>
    .image-message {
      background: #ffffff;
      border-radius: 12px;
      padding: 12px;
      box-shadow: 0 4px 10px rgba(0, 0, 0, 0.05);
    }
    
    /* Toolbar */
    .toolbar {
      display: flex;
      justify-content: space-between;
      margin-bottom: 10px;
    }
    
    .btn {
      padding: 6px 14px;
      border-radius: 20px;
      font-size: 0.85rem;
      cursor: pointer;
      border: none;
    }
    
    .btn.ghost {
      background: #e6f7fb;
      color: #2bb0a6;
    }
    
    .btn.primary {
      background: #2bb0a6;
      color: white;
    }
    
    .btn:hover {
      opacity: 0.85;
    }
    
    /* Preview */
    .preview {
      display: flex;
      justify-content: center;
      margin-bottom: 10px;
    }
    
    .image-box {
      width: 220px;
      height: 160px;
      border-radius: 10px;
      background: #e6f7fb;
      display: flex;
      align-items: center;
      justify-content: center;
      overflow: hidden;
    }
    
    .image-box.filled {
      border: 2px solid #7fdad3;
      background: white;
    }
    
    .image-box img {
      max-width: 100%;
      max-height: 100%;
      object-fit: contain;
    }
    
    .placeholder {
      color: #7aaeb4;
      font-size: 0.8rem;
    }
    
    /* Caption */
    .caption {
      width: 100%;
      border-radius: 20px;
      border: 1px solid #7fdad3;
      padding: 8px 14px;
      outline: none;
    }
    
    .caption:focus {
      border-color: #2bb0a6;
    }
    </style>
    