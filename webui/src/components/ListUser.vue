<script>
    export default {
      name: "ListUser",
    
      props: {
        list: Boolean,
        users: {
          type: Array,
          required: true,
        },
        searchBar: Boolean,
        showOk: Boolean,
      },
    
      emits: ["send-users-list", "search-again", "search-cancel"],
    
      data() {
        return {
          selected: [],
          search: "",
        };
      },
    
      methods: {
        toggleUser(login) {
          if (!this.showOk) return;
    
          const index = this.selected.indexOf(login);
          if (index === -1) {
            this.selected.push(login);
          } else {
            this.selected.splice(index, 1);
          }
        },
    
        confirm() {
          if (this.selected.length > 0) {
            this.$emit("send-users-list", this.selected);
          } else {
            this.cancel();
          }
        },
    
        cancel() {
          this.$emit("search-cancel");
        },
    
        onSearch() {
          this.$emit("search-again", this.search);
        },
    
        isSelected(login) {
          return this.selected.includes(login);
        },
      },
    };
    </script>
    
    <template>
      <div id="listuser" class="list-user">
        <!-- Toolbar -->
        <div class="toolbar">
          <button class="btn ghost" @click="cancel">Cancel</button>
    
          <button
            v-if="showOk"
            class="btn primary"
            @click="confirm"
          >
            Ok
          </button>
    
          <input
            v-if="searchBar"
            v-model="search"
            @input="onSearch"
            class="search"
            type="text"
            placeholder="Search user..."
          />
        </div>
    
        <!-- Users -->
        <div class="listcontext">
          <div
            v-for="u in users"
            :key="u.login"
          >
            <!-- User -->
            <div
              v-if="u && u.login !== '0'"
              class="user-card"
              :class="{ selected: isSelected(u.login) }"
              @click="toggleUser(u.login)"
            >
              <img
                class="avatar"
                :src="u.photo
                  ? 'data:image/gif;base64,' + u.photo
                  : '/standard2.gif'"
              />
    
              <div class="details">
                <span class="name">
                  {{ u.name }}
                  <span
                    v-if="u.categoryinfo === 'user'"
                    class="login"
                  >
                    ({{ u.login }})
                  </span>
                </span>
              </div>
            </div>
    
            <!-- Empty search -->
            <div
              v-if="u && u.login === '0'"
              class="user-card empty"
            >
              <img class="avatar" src="/profile.svg" />
              <div class="details">
                <span class="name muted">
                  Search produced no results
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
    
    <style scoped>
    .list-user {
      background: white;
      border-radius: 12px;
      width: 340px;
      padding: 10px;
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    }
    
    /* Toolbar */
    .toolbar {
      display: flex;
      gap: 8px;
      align-items: center;
      margin-bottom: 10px;
    }
    
    .search {
      flex: 1;
      padding: 6px 10px;
      border-radius: 20px;
      border: 1px solid #7fdad3;
      outline: none;
    }
    
    .search:focus {
      border-color: #2bb0a6;
    }
    
    /* Buttons */
    .btn {
      padding: 6px 14px;
      border-radius: 20px;
      font-size: 0.8rem;
      border: none;
      cursor: pointer;
    }
    
    .btn.ghost {
      background: #e6f7fb;
      color: #2bb0a6;
    }
    
    .btn.primary {
      background: #2bb0a6;
      color: white;
    }
    
    /* List */
    .listcontext {
      max-height: 200px;
      overflow-y: auto;
      border-top: 1px solid #eee;
      border-bottom: 1px solid #eee;
      padding: 5px 0;
    }
    
    /* User card */
    .user-card {
      display: flex;
      align-items: center;
      height: 50px;
      padding: 6px;
      margin: 4px 0;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.2s;
    }
    
    .user-card:hover {
      background: #e6f7fb;
    }
    
    .user-card.selected {
      background: #7fdad3;
    }
    
    .user-card.empty {
      cursor: default;
    }
    
    .avatar {
      width: 38px;
      height: 38px;
      border-radius: 50%;
      margin-right: 10px;
    }
    
    /* Details */
    .details {
      display: flex;
      align-items: center;
    }
    
    .name {
      font-size: 0.9rem;
      font-weight: 500;
      color: #222;
    }
    
    .login {
      font-size: 0.75rem;
      color: #555;
      margin-left: 6px;
    }
    
    .muted {
      color: #999;
    }
    </style>
    