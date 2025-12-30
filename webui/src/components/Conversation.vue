<script setup>
  import { ref, reactive, onMounted, nextTick, getCurrentInstance } from 'vue';
  import SendMessage from '../components/SendMessage.vue';
  import SetComment from '../components/SetComment.vue';
  import ListUser from '../components/ListUser.vue';
  
  // --- Props ---
  const props = defineProps([
      "selectedconvs",
      "sv",
      "login",
      "token",
  ]);
  
  // --- Emits ---
  const emit = defineEmits([
      'list-next',
      'list-previous',
      'change-selected-conv',
      'change-selected-conv-byid',
      'show-user-info',
      'refresh',
      'logout',
  ]);
  
  // --- Accesso all'istanza globale per $axios ---
  const { proxy } = getCurrentInstance();
  const $axios = proxy.$axios;
  
  // --- Refs per i componenti figli ---
  const sendmessageRef = ref(null);
  const vueMessageMenu = ref(null);
  
  // --- Data State (Reattivi) ---
  const showComment = ref(false);
  const showListGroup = ref(false);
  const showForwardMsgId = ref(null);
  const showForward = ref(false);
  const showReplyMsgId = ref(null);
  const Reply = ref(null);
  const showNewMessage = ref(false);
  const showNewMessageList = ref(false);
  const errormsg = ref(null);
  const selectedMessageId = ref(null);
  const selectedConversationId = ref(null);
  const selectedusers = ref(null);
  const tmpselectedusers = ref(null);
  const usersNewGroup = ref(null);
  const signMXY = ref(true);
  const users = ref(null);
  const topMX = ref(0);
  const topMY = ref(0);
  
  const optionsmessage = ref([
      { name: 'Comment', slug: 'comment' },
      { name: 'Uncomment', slug: 'uncomment' },
      { name: 'Forward', slug: 'forward' },
      { name: 'Reply', slug: 'reply' },
      { type: 'divider' },
      { name: '<em>Delete message</em>', slug: 'delete' },
  ]);
  
  const convs = reactive({
      items: []
  });
  
  
  onMounted(() => {
      refresh();
  });
  
  // metodi
  
  const convertTimestamp = (timestamp) => {
      const localDate = new Date(`${timestamp}Z`);
      const stringDate = localDate.toLocaleString("it-IT", Intl.DateTimeFormat().resolvedOptions().timeZone);
      return stringDate.slice(0, -3);
  };
  
  const newConversation = async () => {
      if (showNewMessage.value) {
          showNewMessage.value = false;
          showNewMessageList.value = false;
          return;
      }
      showNewMessage.value = true;
      showNewMessageList.value = true;
  
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'application/json'
          }
      };
  
      try {
          let response = await $axios.get(`/users?search=`, config);
          selectedusers.value = response.data;
          if (!Array.isArray(selectedusers.value)) {
              selectedusers.value = [selectedusers.value];
          }
  
          // Aspettiamo che il DOM si aggiorni per trovare l'elemento
          await nextTick();
  
          
          var element = document.getElementById('listuser');
          var input = document.getElementById('inputRef');
          if (input) input.disabled = true;
          if (sendmessageRef.value) sendmessageRef.value.isEnabled = false;
          
          var button = document.getElementById('buttonNewChat'); 
          var buttonRef = document.getElementById('buttonNewGroup'); 
          
          if (element && buttonRef) {
              var rect = buttonRef.getBoundingClientRect();
              element.style.position = 'fixed';
              element.style.top = (rect.top + 125).toString() + 'px';
              element.style.left = (rect.left - 125).toString() + 'px';
              element.style.display = 'block';
          }
      } catch (e) {
          console.error(e);
      }
  };
  
  const addSelectedUsers = async (usersList) => {
      showListGroup.value = false;
      var out = [];
      for (var key in usersList) {
          out.push('users=' + encodeURIComponent(usersList[key]));
      }
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'accept': 'application/json'
          }
      };
      var convid = 0;
      try {
          let response = await $axios.put(`/conversations/group?${out.join('&')}`, {}, config);
          convid = response.data;
          emit('change-selected-conv-byid', convid);
      } catch (e) { console.error(e); }
  };
  
  const addSelectedUsersNew = async (usersList) => {
      if (usersList.length > 1) {
          alert("Select only one user");
          return;
      }
  
      if (sessionStorage.getItem('login') == usersList) {
          alert("Cannot send message to yourself");
          return;
      }
      
      await nextTick();
      var input = document.getElementById('inputRef');
      if (input) input.disabled = false;
      
      if (sendmessageRef.value) sendmessageRef.value.isEnabled = true;
      users.value = usersList;
      showNewMessageList.value = false;
  };
  
  const addSelectedUsersForward = async (usersList) => {
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'accept': 'application/json'
          }
      };
      var convid = null;
      for (var key in usersList) {
          if (usersList[key] == props.login) {
              continue;
          }
          try {
              let response = await $axios.post(
                  `/conversations/byid/${selectedConversationId.value}/message/${showForwardMsgId.value}`, usersList[key], config);
              convid = response.data;
          } catch(e) { console.error(e); }
      }
      
      var element = document.getElementById('listuser');
      if(element) element.style.display = 'none';
      
      showForward.value = false;
      if(convid) emit('change-selected-conv-byid', convid.conID);
  };
  
  const cancelUsersSearch = () => {
      showListGroup.value = false;
  };
  
  const cancelUsersSearchNew = () => {
      showNewMessageList.value = false;
      var input = document.getElementById('inputRef');
      if (input) input.disabled = false;
      if (sendmessageRef.value) sendmessageRef.value.isEnabled = true;
  };
  
  const cancelUsersSearchForward = () => {
      var element = document.getElementById('listuser');
      if (element) element.style.display = 'none';
      showForward.value = false;
  };
  
  const newGroup = async () => {
      showListGroup.value = true;
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'application/json'
          }
      };
      try {
          let response = await $axios.get(`/users?search=`, config);
          selectedusers.value = response.data;
          if (!Array.isArray(selectedusers.value)) {
              selectedusers.value = [selectedusers.value];
          }
  
          await nextTick();
          var element = document.getElementById('listuser');
          var button = document.getElementById('buttonNewGroup');
          if (element && button) {
              var rect = button.getBoundingClientRect();
              element.style.position = 'fixed';
              element.style.top = (rect.top + 50).toString() + 'px';
              element.style.left = (rect.left - 125).toString() + 'px';
              element.style.display = 'block';
          }
      } catch(e) { console.error(e); }
  };
  
  const newForward = async () => {
      showForward.value = true;
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'application/json'
          }
      };
      try {
          let response = await $axios.get(`/users?showgroups=true&search=`, config);
          selectedusers.value = response.data;
          if (!Array.isArray(selectedusers.value)) {
              selectedusers.value = [selectedusers.value];
          }
  
          await nextTick();
          var element = document.getElementById('listuser');
          if (element) {
              var html = document.documentElement;
              var X = 0, Y = 0;
  
              element.style.position = 'fixed';
              element.style.display = 'block';
  
              var H = element.offsetHeight, W = element.offsetWidth;
  
              Y = Math.min(topMY.value - 80 + H, html.clientHeight) - H;
              X = Math.min(topMX.value - 80 + W, html.clientWidth) - W;
  
              element.style.top = (Y).toString() + 'px';
              element.style.left = (X).toString() + 'px';
          }
      } catch(e) { console.error(e); }
  };
  
  // Funzione generica per la ricerca utenti
  const searchUsersHelper = async (search, endpoint, callback) => {
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'application/json'
          }
      };
      try {
          let response = await $axios.get(`${endpoint}search=${search}`, config);
          tmpselectedusers.value = response.data;
          if (tmpselectedusers.value === null) {
              tmpselectedusers.value = JSON.parse('[{"login": "0", "name": "", "photo": ""}]');
          }
          if (!Array.isArray(tmpselectedusers.value)) {
              tmpselectedusers.value = [tmpselectedusers.value];
          }
          selectedusers.value = tmpselectedusers.value;
          if (callback) await callback();
      } catch(e) { console.error(e); }
  };
  
  const usersSearchAgain = async (search) => {
      await searchUsersHelper(search, '/users?', async () => {
          await nextTick();
          var element = document.getElementById('listuser');
          var button = document.getElementById('buttonNewGroup');
          if (element && button) {
              var rect = button.getBoundingClientRect();
              element.style.position = 'fixed';
              element.style.top = (rect.top + 50).toString() + 'px';
              element.style.left = (rect.left - 125).toString() + 'px';
              element.style.display = 'block';
          }
      });
  };
  
  const usersSearchAgainForward = async (search) => {
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'application/json'
          }
      };
      let response = await $axios.get(`/users?showgroups=true&search=${search}`, config);
      selectedusers.value = response.data;
      if (!Array.isArray(selectedusers.value)) {
          selectedusers.value = [selectedusers.value];
      }
  };
  
  const usersSearchAgainNew = async (search) => {
      await searchUsersHelper(search, '/users?', async () => {
          await nextTick();
          var element = document.getElementById('listuser');
          var button = document.getElementById('buttonNewGroup');
          if (element && button) {
              var rect = button.getBoundingClientRect();
              element.style.position = 'fixed';
              element.style.top = (rect.top + 125).toString() + 'px';
              element.style.left = (rect.left - 125).toString() + 'px';
              element.style.display = 'block';
          }
      });
  };
  
  const setCoordsMessage = (e) => {
      if (signMXY.value) {
          topMX.value = e.clientX;
          topMY.value = e.clientY;
      }
  };
  
  const openRW = (conid, isOther) => {
      var parent = document.getElementById('sp' + conid);
      var element = document.getElementById('rw' + conid);
      if (element === null || parent === null) return;
  
      if (element.style.visibility != 'visible') {
          var html = document.documentElement;
          var X = 0, Y = 0;
  
          element.style.position = 'fixed';
          element.style.display = 'block';
          element.style.visibility = 'visible';
  
          var H = element.offsetHeight, W = element.offsetWidth;
          var rect = parent.getBoundingClientRect();
          
          Y = Math.min(rect.top - 0 + H, html.clientHeight) - H;
          X = Math.min(rect.left - 30 + W, html.clientWidth) - W;
  
          element.style.left = (X).toString() + 'px';
          element.style.top = (Y).toString() + 'px';
  
      } else {
          element.style.visibility = 'hidden';
      }
  };
  
  const openEmoji = (conid, commid) => {
      var parent = document.getElementById('em' + conid);
      var element = document.getElementById('co' + commid);
  
      if (element && parent) {
          if (element.style.visibility != 'visible') {
              var html = document.documentElement;
              var X = 0, Y = 0;
  
              element.style.position = 'fixed';
              element.style.display = 'block';
              element.style.visibility = 'visible';
  
              var H = element.offsetHeight, W = element.offsetWidth;
              var rect = parent.getBoundingClientRect();
  
              Y = Math.min(rect.top - 0 + H, html.clientHeight) - H;
              X = Math.min(rect.left - 0 + W, html.clientWidth) - W;
  
              element.style.left = (X).toString() + 'px';
              element.style.top = (Y).toString() + 'px';
          } else {
              element.style.visibility = 'hidden';
          }
      }
  };
  
  const closeEmoji = () => {
      
  };
  
  const deleteMessage = async (mesid, conid) => {
      errormsg.value = null;
      let config = {
          headers: { Authorization: "Bearer " + props.token }
      };
      try {
          await $axios.delete(`/conversations/byid/${conid}/message/${mesid}`, config);
      } catch (e) {
          if (e.response && e.response.status === 403) {
               errormsg.value = "non puoi eliminare messaggi che non hai inviato te";
          } else {
               errormsg.value = e.toString();
          }
      }
      refresh();
  };
  
  const deleteComment = async (mesid, conid) => {
      errormsg.value = null;
      let config = { headers: { Authorization: "Bearer " + props.token } };
      try {
          await $axios.delete(`/conversations/byid/${conid}/message/${mesid}/comment`, config);
      } catch (e) { errormsg.value = e.toString(); }
      refresh();
  };
  
  const optionMessageClicked = async (event) => {
      if (event.option.slug == 'comment') {
          showComment.value = true;
          await nextTick();
          var element = document.getElementById('comment');
          
          if (element) {
              var html = document.documentElement;
              var X = 0, Y = 0;
  
              element.style.position = 'fixed';
              element.style.display = 'block';
  
              var H = 366, W = 364;
              Y = Math.min(topMY.value - 80 + H, html.clientHeight) - H;
              X = Math.min(topMX.value - 80 + W, html.clientWidth) - W;
  
              element.style.top = (Y).toString() + 'px';
              element.style.left = (X).toString() + 'px';
          }
  
      } else if (event.option.slug == 'uncomment') {
          deleteComment(event.item.mesid, event.item.conid);
      } else if (event.option.slug == 'delete') {
          deleteMessage(event.item.mesid, event.item.conid);
      } else if (event.option.slug == 'forward') {
          showForwardMsgId.value = event.item.mesid;
          showForward.value = true;
          newForward();
      } else if (event.option.slug == 'reply') {
          globalThis.scrollTo(0, 0);
          await nextTick();
          var element = document.getElementById("inputRef");
          if (element) {
              element.value = "";
              element.setSelectionRange(0, 0);
          }
          showReplyMsgId.value = event.item.mesid;
          Reply.value = event.item;
      }
  };
  
  const handleMessageClick = (evt, mesid) => {
      selectedMessageId.value = mesid.mesid;
      selectedConversationId.value = mesid.conid;
      if (vueMessageMenu.value) {
          vueMessageMenu.value.showMenu(evt, mesid);
      }
  };
  
  const sendMessage = async (message, image) => {
      errormsg.value = null;
      if (((message === null) || message == '') && (image === null)) {
          errormsg.value = "No message composed";
          return;
      }
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'multipart/form-data'
          }
      };
      let formdata = new FormData();
      if (message && message != '') formdata.append('text', message);
      if (image && image != '') formdata.append('image', image, 'image');
  
      var reply = "";
      if (showReplyMsgId.value) {
          reply = `&reply=${showReplyMsgId.value}`;
      }
      try {
          let response = await $axios.post(`/conversations?conversation=${props.selectedconvs.conid}${reply}`, formdata, config);
          emit('change-selected-conv', response.data);
      } catch (e) {
          errormsg.value = e.toString();
      }
      showReplyMsgId.value = null;
      Reply.value = null;
      refresh();
  };
  
  const sendMessageWithUser = async (message, image) => {
      errormsg.value = null;
      if (!users.value) {
          errormsg.value = "No recipient selected!";
          return;
      }
      if (((message === null) || message == '') && (image === null)) {
          errormsg.value = "No message composed";
          return;
      }
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'multipart/form-data'
          }
      };
      let formdata = new FormData();
      if (message && message != '') formdata.append('text', message);
      if (image && image != '') formdata.append('image', image, 'image');
  
      try {
          let response = await $axios.post(`/conversations?conversation=${users.value}`, formdata, config);
          showNewMessage.value = false;
          emit('change-selected-conv', response.data);
      } catch (e) {
          errormsg.value = e.toString();
      }
      refresh();
  };
  
  const sendComment = async (comment, mesid, conid) => {
      showComment.value = false;
      if (mesid === null) return;
      errormsg.value = null;
      let config = {
          headers: {
              Authorization: "Bearer " + props.token,
              'content-type': 'text/plain'
          }
      };
      try {
          await $axios.post(`/conversations/byid/${conid}/message/${mesid}/comment`, comment.i, config);
      } catch (e) { errormsg.value = e.toString(); }
      refresh();
  };
  
  const info = () => emit('show-user-info');
  
  const classReceived = (rcv) => {
      if (rcv) {
          var i = 0;
          while (i < rcv.length) {
              if ((rcv[i].received == 'false') && (rcv[i].name != props.login)) {
                  return "noreceived";
              }
              i++;
          }
          return "received";
      }
      return "noreceived";
  };
  
  const singleClassReceived = (rcv, user, login) => {
      if (user == login) return "received";
      if (rcv == 'true') return "received";
      return "noreceived";
  };
  
  const classRead = (rd) => {
      if (rd) {
          var i = 0;
          while (i < rd.length) {
              if ((rd[i].read == 'false') && (rd[i].name != props.login)) {
                  return "noread";
              }
              i++;
          }
          return "read";
      }
      return "noread";
  };
  
  const singleClassRead = (rd, user, login) => {
      if (user == login) return "read";
      if (rd == 'true') return "read";
      return "noread";
  };
  
  const classSenderIsMe = (login) => {
      if (login == props.login) {
          return "message-box my-message";
      }
      return "message-box friend-message";
  };
  
  const logout = () => emit('logout');
  const refresh = () => emit('refresh');
  const previousMessages = () => emit('list-previous');
  
  const nextMessages = () => {
      globalThis.scrollTo(0, 1);
      setTimeout(() => {
          emit('list-next');
      }, 1000);
  };
  
  const clearError = () => errormsg.value = null;
  
  const exitReply = () => {
      Reply.value = null;
      showReplyMsgId.value = null;
  };
  </script>
  
  <template>
      <div v-if="!showNewMessage">
          <div id="main">
              <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                  <h1 class="h4" v-if="selectedconvs && (!showNewMessage)">{{ selectedconvs.name }}</h1>
                  <h1 class="h4" v-else>&nbsp;</h1>
                  <div class="btn-toolbar mb-2 mb-md-0">
                      <div class="btn-group me-2">
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">Refresh</button>
                          <button type="button" class="btn btn-sm btn-outline-secondary" @click="info">Info</button>
                      </div>
                      <div class="btn-group me-2">
                          <button id="buttonNewGroup" type="button" class="btn btn-sm btn-outline-secondary" @click="newGroup">
                              New Group
                          </button>
                          <button id="buttonNewChat" type="button" class="btn btn-sm btn-outline-secondary" :class="{'btn-warning':showNewMessage}" @click="newConversation">
                              New Chat
                          </button>
                      </div>
                      <div class="btn-group me-2">
                          <button type="button" class="btn btn-sm btn-outline-primary" @click="logout">Logout</button>
                      </div>                    
                  </div>
              </div>
              <ErrorMsg v-if="errormsg" :msg="errormsg" @clear-error="clearError"></ErrorMsg>
          </div>
          <div v-if="showReplyMsgId">
              <div v-if="(Reply.text != '') || (Reply.photo != '')" class="replyed inforeplyed" @click="exitReply">
                  <div v-if="Reply.username != ''" class="replyedusername">{{ Reply.username }}</div>
                  <div v-if="Reply.text != ''" class="replyedtext">{{ Reply.text }}</div>
                  <div><img v-if="Reply.photo != ''" v-bind:src="'data:image/gif;base64,' + Reply.photo" class="internal"/></div>
              </div>
          </div>
          <sendmessage 
              v-if="!showNewMessage"
              v-bind:conversation="selectedconvs"
              ref="sendmessageRef"
              @send-message="sendMessage"
          />
          <listuser
              id="listuser"
              v-if="showListGroup || showForward"
              @search-again="showForward ? usersSearchAgainForward($event) : usersSearchAgain($event)"
              @send-users-list="showForward ? addSelectedUsersForward($event) : addSelectedUsers($event)"
              @search-cancel="showForward ? cancelUsersSearchForward() : cancelUsersSearch()"
              v-bind:users="selectedusers"
              v-bind:searchBar="true"
              v-bind:showOk="true"
          />
          <div id="selectedconv" v-if="!showNewMessage">
              <ul class="no-bullets" v-if="selectedconvs">
                  <li v-if="sv.start != 0" @click="previousMessages"><div class="nextmessages"> get newer messages </div></li>
                  <li v-for="cvs in selectedconvs.message" :key="cvs.mesid">
                      <div v-bind:class="classSenderIsMe(cvs.login)">
                          <div 
                              @click.prevent.stop="handleMessageClick($event, cvs)"
                              @mouseenter="setCoordsMessage"
                          >
                              <span v-if="cvs.isforwarded > 0" class="forwarded">Forwarded</span>
                              <h6 v-if="cvs.login != login">{{ cvs.username }}:</h6>
                              <div v-if="(cvs.replyedtext != '') || (cvs.replyedphoto != '')" class="replyed">
                                  <div v-if="cvs.replyedusername != ''" class="replyedusername">{{ cvs.replyedusername }}</div>
                                  <div v-if="cvs.replyedtext != ''" class="replyedtext">{{ cvs.replyedtext }}</div>
                                  <div><img v-if="cvs.replyedphoto != ''" v-bind:src="'data:image/gif;base64,' + cvs.replyedphoto" class="internal"/></div>
                              </div>
                              <img v-if="cvs.category != 'text'" v-bind:src="'data:image/gif;base64,' + cvs.photo" class="internal"/>
                              <div v-if="cvs.category != 'image'">{{ cvs.text }}<br></div>
                              <span :id="'sp' + cvs.mesid"
                                  @click.prevent.stop="openRW(cvs.mesid, cvs.login != login)"
                              >
                                  {{ convertTimestamp(cvs.timestamp) }}
                                  <i v-if="cvs.login == login" v-bind:class="classReceived(cvs.receivedarray)"></i>
                                  <i v-if="cvs.login == login" v-bind:class="classRead(cvs.readarray)"></i>
                              </span>
                              <div class="listrr" :id="'rw' + cvs.mesid"
                                  @click.prevent.stop="openRW(cvs.mesid, cvs.login != login)"
                                  v-if="cvs.readarray && cvs.login == login">
                                  <ul class="ullistrr">
                                      <li v-for="(item, i) in cvs.readarray" :key="item.name">{{ item.username }} 
                                          <i v-bind:class="singleClassReceived(cvs.receivedarray[i].received, cvs.receivedarray[i].name, login)"></i>
                                          <i v-bind:class="singleClassRead(item.read, cvs.receivedarray[i].name, login)"></i>
                                      </li>
                                  </ul>
                              </div>
                          </div>
                          <div class="emoji" :id="'em' + cvs.mesid" 
                              v-if="cvs.comments"
                              @click="openEmoji(cvs.mesid, cvs.comments[0].commId)"
                              >
                              {{ cvs.comments[0].emoji }}
                              <div class="listemoji" :id="'co' + cvs.comments[0].commId"
                                  @click="closeEmoji(cvs.comments[0].commId)"
                                  v-if="cvs.comments">
                                  <ul class="ullistemoji">
                                      <li v-for="c in cvs.comments" :key="c.commId">{{ c.emoji }} - {{ c.username }}</li>
                                  </ul>
                              </div>
                          </div> 
                      </div>
                  </li>
                  <li v-if="sv.limit == selectedconvs.message.length" @click="nextMessages"><div class="nextmessages"> get older messages </div></li>
              </ul>
          </div>
          <div class="overlayComment" :class="{ hiddenComment: !showComment }">
              <setcomment
                  id="comment"
                  v-bind:msgid="selectedMessageId"
                  v-bind:conid="selectedConversationId"
                  v-bind:topX="topMX"
                  v-bind:topY="topMY"
                  @send-comment="sendComment"
              />
          </div>
          
          <vue-simple-context-menu
              element-id="menu-message"
              :options="optionsmessage"
              ref="vueMessageMenu"
              @option-clicked="optionMessageClicked"
          />
      </div>
      <div v-if="showNewMessage">
          <div id="main2">
          </div>            
          <sendmessage 
              v-bind:conversation="showNewMessage"
              ref="sendmessageRef"
              @send-message="sendMessageWithUser"
          />
          <listuser
              id="listuser" 
              v-if="showNewMessageList"
              @search-again="usersSearchAgainNew"
              @send-users-list="addSelectedUsersNew"
              @search-cancel="cancelUsersSearchNew"
              v-bind:users="selectedusers"
              v-bind:searchBar="true"
              v-bind:showOk="true"
          />
      </div>
  </template>
  
  <style scoped>
    /* --- CONFIGURAZIONE GLOBALE & SFONDI --- */
    body, main {
        background-color: #0f172a !important; /* Blu Notte profondo */
        color: #f1f5f9; /* Bianco fumo */
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
    }
    
    /* Rimuove i bordi default e aggiunge spaziatura */
    ul.no-bullets {
        list-style-type: none;
        padding: 0;
        margin: 0;
    }   
    
    /* --- IMMAGINI INTERNE --- */
    .internal {
        max-width: 100%; /* Responsive */
        border-radius: 8px; /* Bordi arrotondati per le immagini */
        margin-top: 5px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.3);
    }
    
    /* --- BOX MESSAGGI --- */
    .message-box {
      position: relative;
      display: grid;
      grid-template-columns: auto;
      width: 100%;
      margin: 10px 0; /* Spaziatura aumentata */
      margin-top: 25px;
    }
    
    .message-box > div {
      position: relative;
      right: 0;
      text-align: left;
      max-width: 70%; /* Meno largo per estetica migliore */   
      padding: 12px 16px;
      border-radius: 12px;
      font-size: 0.95rem;
      line-height: 1.4;
      cursor: pointer;
      box-shadow: 0 2px 5px rgba(0,0,0,0.2); /* Ombra leggera */
      color: #ffffff;
    }
    
    /* --- MESSAGGIO AMICO (Ricevuto) --- */
    .friend-message {
      justify-content: flex-start;
    }
    
    .friend-message > div {
      background: #1e293b; /* Grigio bluastro scuro */
      color: #e2e8f0; /* Bianco sporco */
      border-top-left-radius: 0; /* Angolo "speech bubble" */
    }
    
    /* Coda del fumetto amico */
    .message-box.friend-message::before {
      content : '';
      position: absolute;
      top: 0;
      left: -10px;
      width: 20px;
      height: 20px;
      /* Gradiente adattato al nuovo colore di sfondo */
      background: linear-gradient(225deg, #1e293b 0%, #1e293b 50%, transparent 50%, transparent);
    }
    
    /* --- MIO MESSAGGIO (Inviato) --- */
    .my-message {
      justify-content: flex-end;
    }
    
    .message-box.my-message > div {
      /* Gradiente Verde Acqua elegante */
      background: linear-gradient(135deg, #0f766e 0%, #0d9488 100%); 
      color: #ffffff;
      border-top-right-radius: 0; /* Angolo "speech bubble" */
    }
    
    /* Coda del fumetto mio */
    .message-box.my-message > div::before {
      content : '';
      position: absolute;
      top: 0;
      right: -10px;
      width: 20px;
      height: 20px;
      /* Gradiente adattato al verde acqua */
      background: linear-gradient(135deg, #0d9488 0%, #0d9488 50%, transparent 50%, transparent);
    }
    
    /* --- NOME UTENTE --- */
    .message-box > div > h6 {
      color: #5eead4; /* Verde acqua chiaro fluo per risaltare sul blu notte */
      font-size: 0.85em;
      font-weight: 700;
      margin-bottom: 4px;
      margin-top: 0;
    }
    
    /* --- TIMESTAMPS & INFO --- */
    .message-box > div > span {
      display: block;
      margin-top: 5px;
      font-size: 0.75em;
      opacity: 0.8;
      font-style: normal;
      text-align: right;
      color: #cbd5e1;
    }
    
    /* --- INDICATORI DI LETTURA (Spunte) --- */
    /* Spunte non ricevute/non lette: Grigio chiaro */
    .noreceived::after, .noread::after {
        content: '✓';
        color: #94a3b8; 
        margin-right: -0.4em;
        font-weight: bold;
    }
    
    /* Spunte ricevute/lette: Verde Acqua Acceso */
    .received::after, .read::after {
        content: '✓';
        color: #22d3ee; /* Ciano brillante */
        margin-right: -0.4em;
        font-weight: bold;
        text-shadow: 0 0 5px rgba(34, 211, 238, 0.4); /* Effetto neon */
    }
    
    /* --- STILE RISPOSTA (REPLY) --- */
    .replyed {
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-left: 4px solid #14b8a6; /* Barra laterale Verde Acqua */
        height: auto;
        min-height: 50px;
        padding: 8px;
        margin-bottom: 8px;
        opacity: 1;
        background-color: rgba(15, 23, 42, 0.5); /* Sfondo scuro semi-trasparente */
        border-radius: 6px;
    }
    
    .replyedusername {
        font-size: 11px;
        color: #2dd4bf; /* Teal chiaro */
        font-weight: bold;
        margin-bottom: 2px;
    }
    
    .replyedtext {
        height: 20px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        max-width: 200px;
        font-size: 12px;
        color: #e2e8f0; /* Testo chiaro */
    }
    
    .inforeplyed::before {
        /* Filtro per rendere l'icona svg bianca/chiara */
        filter: invert(1) brightness(2); 
        content: url('/close.svg');
        position: absolute;
        margin-left: -22px;
        margin-top: 2px;
        cursor: pointer;
        width: 16px;
        height: 16px;
        opacity: 0.7;
    }
    
    /* --- EMOJI REACTION --- */
    .emoji {
        position: relative;
        top: -20px !important;
        border-radius: 20px;
        width: fit-content;
        height: fit-content;
        background-color: #334155 !important; /* Sfondo scuro per l'emoji */
        border: 2px solid #0f172a; /* Bordo dello stesso colore dello sfondo body per stacco */
        padding: 2px 6px !important;
        cursor: pointer;
        margin-top: 15px !important;
        font-size: 0.8em;
        box-shadow: 0 2px 4px rgba(0,0,0,0.2);
        color: #fff;
    }
    
    /* --- MENU POPUP (Emoji list, Read receipts) --- */
    .listemoji, .listrr {
        position: fixed;
        padding: 8px 12px;
        visibility: hidden;
        border-radius: 8px;
        width: fit-content;
        height: fit-content;
        /* Sfondo scuro per i popup */
        background-color: #1e293b; 
        border: 1px solid #334155;
        color: #f1f5f9;
        z-index: 100000;
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.5);
    }
    
    ul.ullistemoji li, ul.ullistrr li {
        padding: 2px 0;
        border-bottom: 1px solid #334155;
    }
    ul.ullistemoji li:last-child, ul.ullistrr li:last-child {
        border-bottom: none;
    }
    
    /* --- PULSANTE "CARICA ALTRI" --- */
    .nextmessages {
        border: 1px solid #14b8a6; /* Bordo Verde Acqua */
        width: 40%;
        margin: 20px auto;
        padding: 8px;
        text-align: center;
        border-radius: 20px;
        background-color: rgba(20, 184, 166, 0.1); /* Sfondo teal trasparente */
        color: #2dd4bf; /* Testo teal chiaro */
        cursor: pointer;
        font-size: 0.85em;
        font-weight: 600;
        transition: all 0.2s ease;
    }
    
    .nextmessages:hover {
        background-color: rgba(20, 184, 166, 0.2);
        transform: translateY(-1px);
    }
    
    /* --- ETICHETTA INOLTRATO --- */
    .forwarded {
      font-size: 0.75em;
      font-style: italic;
      margin-top: -5px !important;
      margin-bottom: 5px;
      display: block;
      opacity: 0.8 !important;
      color: #94a3b8;
    }
    .forwarded::before {
        content: "↪ ";
    }
    
    /* --- OVERLAY COMMENTO --- */
    .overlayComment {
        display: block;
        z-index: 10000;
        position: fixed;
        left: 0;
        top: 0;
        width: 100%;
        height: 4px; /* Barra sottile */
        background: #2dd4bf; /* Verde acqua */
        box-shadow: 0 0 10px #2dd4bf;
        border: none;
    }
    .hiddenComment {
      display: none;
    }
    
    .hiddenForward {
      transform: translateX(-100%);
    }
    </style>