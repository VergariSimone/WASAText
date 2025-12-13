<script>
    export default {
      name: "SetInfo",
    
      props: {
        itype: String,
        ievent: Object,
      },
    
      emits: ["set-info", "close-info"],
    
      data() {
        return {
          name: this.ievent?.name || "",
          photo: this.ievent?.photo || "",
          selectedFile: null,
        };
      },
    
      methods: {
        chooseFile() {
          this.$refs.fileInput.click();
        },
    
        async setName() {
          if (!this.name) return;
    
          const config = {
            headers: { Authorization: "Bearer " + sessionStorage.getItem("token") },
          };
    
          try {
            if (this.itype === "group") {
              await this.$axios.put(
                `/conversations/byid/${this.ievent.conid}/name?name=${this.name}`,
                this.name,
                config
              );
            } else {
              await this.$axios.put("/info/username", this.name, config);
            }
    
            this.$emit("set-info", this.itype, this.ievent.conid);
          } catch (err) {
            if (err.response?.status === 406) {
              alert("Username not acceptable or already used");
            }
          }
        },
    
        async setPhoto(event) {
          const files = event.target.files;
          if (!files.length) return;
          const file = files[0];
    
          const config = {
            headers: {
              Authorization: "Bearer " + sessionStorage.getItem("token"),
              "content-type": "image/gif",
            },
          };
    
          try {
            if (this.itype === "group") {
              await this.$axios.put(
                `/conversations/byid/${this.ievent.conid}/photo`,
                file,
                config
              );
            } else {
              await this.$axios.put("/info/userphoto", file, config);
            }
    
            this.photo = await this.readFileAsDataURL(file);
            this.$emit("set-info", this.itype, this.ievent.conid);
          } catch (err) {
            console.error(err);
          }
        },
    
        readFileAsDataURL(file) {
          return new Promise((resolve, reject) => {
            const reader = new FileReader();
            reader.onload = (e) => resolve(e.target.result);
            reader.onerror = (e) => reject(e);
            reader.readAsDataURL(file);
          });
        },
    
        close() {
          this.$emit("close-info");
        },
      },
    };
    </script>
    
    <template>
      <div class="info-container">
        <h5>View / Set Information</h5>
    
        <div class="info-row">
          <label>Name:</label>
          <input type="text" v-model="name" maxlength="32" />
          <button class="btn" @click="setName">Save Name</button>
        </div>
    
        <div class="info-row">
          <label>Photo:</label>
          <div class="camera photo-preview">
            <img
              v-if="photo"
              :src="photo.startsWith('data') ? photo : 'data:image/gif;base64,' + photo"
            />
            <img v-else src="/standard2.gif" />
          </div>
          <input
            ref="fileInput"
            type="file"
            accept="image/gif"
            hidden
            @change="setPhoto"
          />
          <button class="btn" @click="chooseFile">Upload Image</button>
        </div>
    
        <div class="info-row">
          <button class="btn close-btn" @click="close">Close</button>
        </div>
      </div>
    </template>
    
    <style scoped>
    .info-container {
      display: grid;
      gap: 15px;
      padding: 15px;
      background: #f0fefc; /* verde acqua chiaro */
      border-radius: 12px;
      width: 320px;
      box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12);
      position: relative;
      top: 50px;
      color: #013440; /* verde acqua scuro per testi */
    }
    
    h5 {
      color: #00bcd4; /* celeste */
      font-weight: bold;
      text-align: center;
      margin-bottom: 10px;
    }
    
    .info-row {
      display: flex;
      align-items: center;
      gap: 10px;
    }
    
    label {
      min-width: 50px;
      font-weight: 500;
    }
    
    input[type="text"] {
      flex: 1;
      padding: 6px 8px;
      border-radius: 8px;
      border: 1px solid #7fdad3;
      background-color: white;
    }
    
    .btn {
      padding: 6px 10px;
      border-radius: 8px;
      border: 1px solid #7fdad3;
      background-color: #e0f7f5;
      cursor: pointer;
      transition: 0.2s;
    }
    
    .btn:hover {
      background-color: #d0f1ee;
    }
    
    .photo-preview {
      width: 100px;
      height: 100px;
      overflow: hidden;
      border-radius: 50%;
      display: flex;
      justify-content: center;
      align-items: center;
    }
    
    .photo-preview img {
      max-width: 100%;
      max-height: 100%;
      object-fit: cover;
    }
    
    .close-btn {
      background-color: #00bcd4;
      color: white;
      border: none;
      width: 100%;
    }
    </style>
    