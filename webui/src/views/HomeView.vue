<script>
	import { ref, reactive, onMounted } from "vue";
	import SetInfo from "../components/SetInfo.vue";
	import ListUser from "../components/ListUser.vue";
	import SendMessage from "../components/SendMessage.vue";
	import SetComment from "../components/SetComment.vue";
	import Conversation from "../components/Conversation.vue";
	
	export default {
	  components: {
		SetInfo,
		ListUser,
		SendMessage,
		SetComment,
		Conversation,
	  },
	
	  setup(_, { emit }) {
		// Stato principale
		const convs = reactive({ items: [] });
		const selectedConv = ref(null);
		const showInfo = ref(false);
		const ievent = ref(null);
		const itype = ref(null);
		const showUserList = ref(false);
		const usersList = ref([]);
		const showUserGroup = ref(false);
		const usersGroups = ref([]);
		const showSetComment = ref(false);
		const commentData = reactive({ msgid: null, conid: null, topX: 0, topY: 0 });
		const token = ref(sessionStorage.getItem("token"));
		const login = ref(sessionStorage.getItem("login"));
		const refreshInterval = ref(null);
		const errormsg = ref(null);
	
		const sv = reactive({ start: 0, limit: 50, last: null });
	
		// Funzioni
		const fetchConversations = async () => {
		  try {
			const config = { headers: { Authorization: `Bearer ${token.value}` } };
			const response = await fetch("/conversations", config);
			const data = await response.json();
			convs.items = data;
		  } catch (e) {
			errormsg.value = e.toString();
		  }
		};
	
		const chooseConv = async (conid) => {
		  sv.last = conid;
		  sv.start = sv.last !== conid ? 0 : sv.start;
		  try {
			const config = { headers: { Authorization: `Bearer ${token.value}` } };
			const response = await fetch(`/conversations/byid/${conid}?start=${sv.start}&limit=${sv.limit}`, config);
			selectedConv.value = await response.json();
		  } catch (e) {
			errormsg.value = e.toString();
		  }
		};
	
		const openUserList = async (conid = null) => {
		  showUserList.value = true;
		  const config = { headers: { Authorization: `Bearer ${token.value}` } };
		  const response = await fetch("/users?search=", config);
		  usersList.value = await response.json();
		  if (!Array.isArray(usersList.value)) usersList.value = [usersList.value];
		};
	
		const openUserGroup = async (conid) => {
		  showUserGroup.value = true;
		  const config = { headers: { Authorization: `Bearer ${token.value}` } };
		  const response = await fetch(`/conversations/byid/${conid}/info`, config);
		  usersGroups.value = (await response.json()).user;
		  if (!Array.isArray(usersGroups.value)) usersGroups.value = [usersGroups.value];
		};
	
		const optionClicked = (event, item) => {
		  switch (event.option.slug) {
			case "add":
			  openUserList(item.conid);
			  break;
			case "leave":
			  leaveGroup(item.conid);
			  break;
			case "open":
			  chooseConv(item.conid);
			  break;
			case "show":
			  openUserGroup(item.conid);
			  break;
			case "info":
			  itype.value = "group";
			  ievent.value = item;
			  showInfo.value = true;
			  break;
			case "send":
			  chooseConv(item.conid);
			  break;
		  }
		};
	
		const leaveGroup = async (conid) => {
		  const config = { headers: { Authorization: `Bearer ${token.value}` } };
		  await fetch(`/conversations/byid/${conid}`, { method: "PUT", headers: config.headers });
		  selectedConv.value = null;
		  fetchConversations();
		};
	
		const sendComment = (comment, msgid, conid) => {
		  if (!comment) return;
		  // simulazione aggiunta commento
		  if (selectedConv.value) selectedConv.value.messages.push({ id: Date.now(), user: login.value, text: comment });
		  showSetComment.value = false;
		};
	
		const logout = () => {
		  token.value = null;
		  sessionStorage.removeItem("token");
		  sessionStorage.removeItem("login");
		  window.location.href = "/login";
		};
	
		onMounted(() => {
		  fetchConversations();
		  refreshInterval.value = setInterval(fetchConversations, 20000);
		});
	
		return {
		  convs,
		  selectedConv,
		  chooseConv,
		  showInfo,
		  ievent,
		  itype,
		  showUserList,
		  usersList,
		  showUserGroup,
		  usersGroups,
		  showSetComment,
		  commentData,
		  optionClicked,
		  sendComment,
		  openUserList,
		  logout,
		  login,
		};
	  },
	};
	</script>
	
	<template>
	  <div class="home-container">
		<!-- Sidebar conversazioni -->
		<aside class="sidebar">
		  <ul>
			<li v-for="cvs in convs.items" :key="cvs.conid" @click="chooseConv(cvs.conid)">
			  <img :src="cvs.photo ? 'data:image/gif;base64,' + cvs.photo : '/standard2.gif'" class="avatar" />
			  <span class="conv-name">{{ cvs.name }}</span>
			</li>
		  </ul>
	
		  <!-- Lista utenti selezione -->
		  <list-user v-if="showUserList" :users="usersList" @send-users-list="showUserList=false" @search-cancel="showUserList=false"/>
		  <list-user v-if="showUserGroup" :users="usersGroups" @send-users-list="showUserGroup=false" @search-cancel="showUserGroup=false"/>
		</aside>
	
		<!-- Main Chat -->
		<main class="chat-area">
		  <conversation v-if="selectedConv" :selectedconvs="selectedConv" :login="login"/>
		  <send-message v-if="selectedConv" :conversation="selectedConv"/>
		</main>
	
		<!-- Modali -->
		<set-info v-if="showInfo" :itype="itype" :ievent="ievent" @close-info="showInfo=false"/>
		<set-comment v-if="showSetComment" v-bind="commentData" @send-comment="sendComment"/>
	  </div>
	</template>
	
	<style scoped>
	.home-container {
	  display: flex;
	  height: 100vh;
	  background-color: #f0f0f0;
	}
	
	.sidebar {
	  width: 300px;
	  overflow-y: auto;
	  background: #fff;
	  border-right: 1px solid #ccc;
	  padding: 10px;
	}
	
	.sidebar ul {
	  list-style: none;
	  padding: 0;
	  margin: 0;
	}
	
	.sidebar li {
	  display: flex;
	  align-items: center;
	  gap: 10px;
	  padding: 5px;
	  cursor: pointer;
	}
	
	.sidebar li:hover {
	  background-color: #f0e9f4;
	}
	
	.avatar {
	  width: 40px;
	  height: 40px;
	  border-radius: 50%;
	}
	
	.conv-name {
	  font-weight: bold;
	}
	
	.chat-area {
	  flex: 1;
	  display: flex;
	  flex-direction: column;
	  padding: 10px;
	  overflow-y: auto;
	}
</style>
	