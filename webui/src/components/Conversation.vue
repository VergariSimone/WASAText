<script>
    import SendMessage from './SendMessage.vue'
    import SetComment from './SetComment.vue'
    import ListUser from './ListUser.vue'
    
    export default {
      name: 'Conversation',
    
      components: {
        SendMessage,
        SetComment,
        ListUser
      },
    
      props: {
        selectedconvs: Object,
        sv: Object,
        login: String,
        token: String
      },
    
      emits: [
        'list-next',
        'list-previous',
        'change-selected-conv',
        'change-selected-conv-byid',
        'show-user-info',
        'refresh',
        'logout'
      ],
    
      data() {
        return {
          errormsg: null,
    
          // UI state
          showComment: false,
          showForward: false,
          showNewMessage: false,
          showNewMessageList: false,
    
          // selections
          selectedMessageId: null,
          selectedConversationId: null,
          selectedusers: null,
    
          // reply / forward
          replyMessage: null,
          replyMsgId: null,
          forwardMsgId: null,
    
          // mouse coords
          topMX: 0,
          topMY: 0,
    
          optionsmessage: [
            { name: 'Comment', slug: 'comment' },
            { name: 'Forward', slug: 'forward' },
            { name: 'Reply', slug: 'reply' },
            { type: 'divider' },
            { name: '<em>Delete</em>', slug: 'delete' }
          ]
        }
      },
    
      mounted() {
        this.refresh()
      },
    
      methods: {
        /* =======================
           UTIL
        ======================= */
        convertTimestamp(ts) {
          return new Date(`${ts}Z`).toLocaleString('it-IT').slice(0, -3)
        },
    
        classSender(login) {
          return login === this.login
            ? 'message mine'
            : 'message other'
        },
    
        /* =======================
           MESSAGE ACTIONS
        ======================= */
        handleMessageClick(evt, msg) {
          this.selectedMessageId = msg.mesid
          this.selectedConversationId = msg.conid
          this.$refs.messageMenu.showMenu(evt, msg)
        },
    
        async sendMessage(text, image) {
          if (!text && !image) return
    
          const form = new FormData()
          if (text) form.append('text', text)
          if (image) form.append('image', image)
    
          const reply = this.replyMsgId ? `&reply=${this.replyMsgId}` : ''
    
          await this.$axios.post(
            `/conversations?conversation=${this.selectedconvs.conid}${reply}`,
            form,
            { headers: { Authorization: `Bearer ${this.token}` } }
          )
    
          this.replyMsgId = null
          this.replyMessage = null
          this.refresh()
        },
    
        async deleteMessage(mesid, conid) {
          await this.$axios.delete(
            `/conversations/byid/${conid}/message/${mesid}`,
            { headers: { Authorization: `Bearer ${this.token}` } }
          )
          this.refresh()
        },
    
        /* =======================
           CONTEXT MENU
        ======================= */
        optionMessageClicked({ option, item }) {
          switch (option.slug) {
            case 'reply':
              this.replyMsgId = item.mesid
              this.replyMessage = item
              break
    
            case 'forward':
              this.forwardMsgId = item.mesid
              this.showForward = true
              break
    
            case 'comment':
              this.showComment = true
              break
    
            case 'delete':
              this.deleteMessage(item.mesid, item.conid)
              break
          }
        },
    
        /* =======================
           NAVIGATION
        ======================= */
        nextMessages() {
          this.$emit('list-next')
        },
        previousMessages() {
          this.$emit('list-previous')
        },
    
        refresh() {
          this.$emit('refresh')
        },
        logout() {
          this.$emit('logout')
        }
      }
    }
 </script>


<template>
    <div class="conversation-wrapper" v-if="selectedconvs">
  
      <!-- HEADER -->
      <header class="conversation-header">
        <h3>{{ selectedconvs.name }}</h3>
  
        <div class="actions">
          <button @click="refresh">⟳</button>
          <button @click="$emit('show-user-info')">ℹ</button>
          <button @click="logout">Logout</button>
        </div>
      </header>
  
      <!-- REPLY PREVIEW -->
      <div v-if="replyMessage" class="reply-preview" @click="replyMessage=null">
        <b>{{ replyMessage.username }}</b>
        <p>{{ replyMessage.text }}</p>
      </div>
  
      <!-- MESSAGES -->
      <div class="messages">
        <div
          v-for="msg in selectedconvs.message"
          :key="msg.mesid"
          :class="classSender(msg.login)"
          @contextmenu.prevent="handleMessageClick($event, msg)"
        >
          <h6 v-if="msg.login !== login">{{ msg.username }}</h6>
  
          <img
            v-if="msg.category !== 'text'"
            :src="'data:image/gif;base64,' + msg.photo"
          />
  
          <p>{{ msg.text }}</p>
  
          <small>{{ convertTimestamp(msg.timestamp) }}</small>
        </div>
  
        <div
          v-if="sv.limit === selectedconvs.message.length"
          class="load-more"
          @click="nextMessages"
        >
          Carica altri messaggi
        </div>
      </div>
  
      <!-- SEND MESSAGE -->
      <SendMessage
        ref="sendmessage"
        @send-message="sendMessage"
      />
  
      <!-- COMMENT -->
      <SetComment
        v-if="showComment"
        @send-comment="showComment=false"
      />
  
      <!-- CONTEXT MENU -->
      <vue-simple-context-menu
        ref="messageMenu"
        element-id="menu-message"
        :options="optionsmessage"
        @option-clicked="optionMessageClicked"
      />
  
    </div>
  
    <div v-else class="empty-chat">
      Seleziona una conversazione
    </div>
</template>

<style scoped>
    .conversation-wrapper {
      display: flex;
      flex-direction: column;
      height: 100%;
    }
    
    .conversation-header {
      padding: 12px;
      background: #ededed;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .actions button {
      margin-left: 8px;
    }
    
    .messages {
      flex: 1;
      padding: 15px;
      overflow-y: auto;
      background: #f0f2f5;
    }
    
    .message {
      max-width: 60%;
      padding: 10px;
      border-radius: 8px;
      margin-bottom: 10px;
    }
    
    .message.mine {
      margin-left: auto;
      background: #dcf8c6;
    }
    
    .message.other {
      background: white;
    }
    
    .reply-preview {
      background: #e6c9fd;
      padding: 10px;
      cursor: pointer;
    }
    
    .load-more {
      text-align: center;
      padding: 10px;
      cursor: pointer;
      color: #555;
    }
    
    .empty-chat {
      text-align: center;
      margin-top: 50px;
      color: #888;
    }
</style>
    
  
    