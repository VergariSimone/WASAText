<script>
	import { ref, reactive, onMounted, onUnmounted, getCurrentInstance } from 'vue';
	import { useRouter } from 'vue-router';
	import SetInfo from '../components/SetInfo.vue';
	import ListUser from '../components/ListUser.vue';
	import Conversation from '../components/Conversation.vue';
	import VueSimpleContextMenu from 'vue-simple-context-menu';
	
	export default {
		components: {
			'setinfo': SetInfo,
			'listuser': ListUser,
			'conversation': Conversation,
			'vue-simple-context-menu': VueSimpleContextMenu
		},
		setup() {
			const { proxy } = getCurrentInstance();
			const $axios = proxy.$axios;
			const router = useRouter();
	
			// --- Refs ai componenti figli (per accesso metodi) ---
			const conversationRef = ref(null);
			const vueConversationMenu = ref(null);
			const vueGroupMenu = ref(null);
	
			// --- State Reattivo ---
			const sv = reactive({ start: 0, limit: 50, last: null });
			const showInfo = ref(false);
			const showUserGroup = ref(false);
			const usersGroups = ref(null);
			const ievent = ref(null);
			const itype = ref(null);
			const errormsg = ref(null);
			const token = ref(null);
			const login = ref(null);
			const loading = ref(false);
			const showList = ref(false);
			const selectedconvs = ref(null);
			const searchSelectedConv = ref(null);
			const selectedusers = ref(null);
			const tmpselectedusers = ref(null);
			const refreshInterval = ref(null);
			const signXY = ref(true);
			const topX = ref(0);
			const topY = ref(0);
	
			const convs = reactive({ items: [] });
	
			const optionsconversation = [
				{ name: 'Open', slug: 'open' },
				{ type: 'divider' },
				{ name: 'Send message', slug: 'send' },
			];
	
			const optionsgroup = [
				{ name: 'Open', slug: 'open' },
				{ name: 'Set info', slug: 'info' },
				{ name: 'Add user', slug: 'add' },
				{ name: 'Show user', slug: 'show' },
				{ type: 'divider' },
				{ name: 'Send message', slug: 'send' },
				{ type: 'divider' },
				{ name: '<em>Leave group</em>', slug: 'leave' },
			];
	
			// --- Methods ---
	
			const calcPosition = (element) => {
				var html = document.documentElement;
				var X = 0, Y = 0;
				var H = element.offsetHeight, W = element.offsetWidth;
				
				Y = Math.min(topY.value + H, html.clientHeight) - H;
				X = Math.min(topX.value - 30 + W, html.clientWidth) - W;
				return [X, Y];
			};
	
			const isPresent = (ncid, cvs_obj) => {
				if (cvs_obj && ncid && cvs_obj.items) {
					for (const element of cvs_obj.items) {
						if (element.conid == ncid) return true;
					}
				}
				return false;
			};
	
			const convertTimestamp = (timestamp) => {
				const localDate = new Date(`${timestamp}Z`);
				const stringDate = localDate.toLocaleString("it-IT", Intl.DateTimeFormat().resolvedOptions().timeZone);
				return stringDate.slice(0, -3);
			};
	
			const resetToread = (conid) => {
				for (const cvs of convs.items) {
					if (cvs.conid == conid) {
						cvs.toread = 0;
						return;
					}
				}
			};
	
			const choose = async (conid) => {
				resetToread(conid);
				errormsg.value = null;
				let config = {
					headers: { Authorization: "Bearer " + token.value, 'content-type': 'application/json' }
				};
				try {
					if (sv.last != conid) {
						sv.start = 0;
					}
					sv.last = conid;
					let response = await $axios.get(`/conversations/byid/${conid}?start=${sv.start}&limit=${sv.limit}`, config);
					selectedconvs.value = response.data;
				} catch (e) {
					errormsg.value = e.toString();
				}
			};
	
			const refresh = async (chooseFlag = true) => {
				loading.value = true;
				token.value = sessionStorage.getItem('token');
				login.value = sessionStorage.getItem('login');
				errormsg.value = null;
				let config = {
					headers: { Authorization: "Bearer " + token.value, 'content-type': 'application/json' }
				};
				try {
					let response = await $axios.get("/conversations", config);
					convs.items = response.data;
				} catch (e) {
					errormsg.value = e.toString();
				}
	
				if (selectedconvs.value) {
					if (isPresent(selectedconvs.value.conid, convs)) {
						choose(selectedconvs.value.conid);
					} else {
						if (convs.items.length > 0) {
							selectedconvs.value = convs.items[0];
							choose(convs.items[0].conid);
						}
					}
				}
				loading.value = false;
			};
	
			const listNext = () => {
				sv.start = sv.start + sv.limit;
				choose(sv.last);
			};
	
			const listPrevious = () => {
				sv.start = Math.max(sv.start - sv.limit, 0);
				choose(sv.last);
			};
	
			const addSelectedUsers = async (users) => {
				showList.value = false;
				var out = [];
				for (var key in users) {
					out.push('users=' + encodeURIComponent(users[key]));
				}
				let config = {
					headers: { Authorization: "Bearer " + token.value, 'accept': 'application/json' }
				};
				try {
					await $axios.put(`/conversations/group?groupID=${searchSelectedConv.value}&${out.join('&')}`, {}, config);
				} catch (e) { console.error(e) }
			};
	
			const leaveGroup = async (convid) => {
				let config = {
					headers: { Authorization: "Bearer " + token.value, 'accept': '*/*' }
				};
				try {
					await $axios.put(`/conversations/byid/${convid}`, {}, config);
				} catch (e) { console.error(e) }
			};
	
			const cancelUsersSearch = () => showList.value = false;
			const cancelUsersGroupSearch = () => showUserGroup.value = false;
	
			const setCoordsGroups = (e) => {
				if (signXY.value) {
					topX.value = e.clientX;
					topY.value = e.clientY;
				}
			};
	
			// Helper function to handle popup positioning logic
			const updatePopupPosition = () => {
				 var element = document.getElementById('listuser');
				 if(element) {
					element.style.position = 'fixed';
					element.style.display = 'block';
					var [X, Y] = calcPosition(element);
					element.style.top = (Y - 30).toString() + 'px';
					element.style.left = (X + 30).toString() + 'px';
				 }
			};
	
			const usersSearchAgain = async (search) => {
				let config = { headers: { Authorization: "Bearer " + token.value, 'content-type': 'application/json' } };
				let response = await $axios.get(`/users?search=${search}`, config);
				tmpselectedusers.value = response.data;
				if (tmpselectedusers.value === null) {
					tmpselectedusers.value = JSON.parse('[{"login": "0", "name": "", "photo": ""}]');
				}
				if (!Array.isArray(tmpselectedusers.value)) {
					tmpselectedusers.value = [tmpselectedusers.value];
				}
				selectedusers.value = tmpselectedusers.value;
				updatePopupPosition();
			};
	
			const usersSearch = async () => {
				signXY.value = false;
				showList.value = true;
				let config = { headers: { Authorization: "Bearer " + token.value, 'content-type': 'application/json' } };
				let response = await $axios.get(`/users?search=`, config);
				selectedusers.value = response.data;
				if (!Array.isArray(selectedusers.value)) {
					selectedusers.value = [selectedusers.value];
				}
				updatePopupPosition();
				signXY.value = true;
			};
	
			const usersGroupSearch = async (conid) => {
				signXY.value = false;
				showUserGroup.value = true;
				let config = { headers: { Authorization: "Bearer " + token.value, 'content-type': 'application/json' } };
				let response = await $axios.get(`/conversations/byid/${conid}/info`, config);
				usersGroups.value = response.data.user;
				if (!Array.isArray(usersGroups.value)) {
					usersGroups.value = [usersGroups.value];
				}
				updatePopupPosition();
				signXY.value = true;
			};
	
			const optionClicked = (event) => {
				const slug = event.option.slug;
				const item = event.item;
	
				if (slug == 'add') {
					searchSelectedConv.value = item.conid;
					usersSearch();
				} else if (slug == 'leave') {
					searchSelectedConv.value = item.conid;
					leaveGroup(item.conid);
					if (selectedconvs.value && selectedconvs.value.conid == item.conid) {
						selectedconvs.value = null;
					}
					refresh();
				} else if (slug == 'open') {
					choose(item.conid);
				} else if (slug == 'show') {
					usersGroupSearch(item.conid);
				} else if (slug == 'info') {
					itype.value = 'group';
					ievent.value = item;
					showInfo.value = true;
				} else if (slug == 'send') {
					choose(item.conid);
					// Accesso al componente figlio tramite ref
					if (conversationRef.value && conversationRef.value.$refs.sendmessage) {
						conversationRef.value.$refs.sendmessage.selectInput();
					}
				}
			};
	
			const handleClick = (evt, conv) => {
				if (conv.categoryinfo == 'conversation') {
					if(vueGroupMenu.value) vueGroupMenu.value.hideContextMenu();
					if(vueConversationMenu.value) vueConversationMenu.value.showMenu(evt, conv);
				} else {
					if(vueConversationMenu.value) vueConversationMenu.value.hideContextMenu();
					if(vueGroupMenu.value) vueGroupMenu.value.showMenu(evt, conv);
				}
			};
	
			const setInfoProps = (resp) => {
				itype.value = resp.data.categoryinfo;
				ievent.value = resp.data;
			};
	
			const reloadInfo = async (infoType, conid) => {
				let config = { headers: { Authorization: "Bearer " + token.value, 'content-type': 'text/plain' } };
				if (itype.value == 'group') {
					let response = await $axios.get(`conversations/byid/${conid}/info`, config);
					setInfoProps(response);
				} else {
					let response = await $axios.get(`/info`, config);
					setInfoProps(response);
				}
			};
	
			const info = async () => {
				if (showInfo.value) return;
				itype.value = 'conversation';
				showInfo.value = true;
				reloadInfo(null, null);
			};
	
			const classNumberToRead = (rd) => rd == 0 ? "number no" : "number yes";
	
			const logout = () => {
				token.value = null;
				clearInterval(refreshInterval.value);
				sessionStorage.removeItem('token');
				sessionStorage.removeItem('login');
				router.push({ name: 'Login' });
			};
	
			const changeSelectedConv = (newconv) => selectedconvs.value = newconv;
	
			const changeSelectedConvByid = async (newconvid) => {
				await refresh(false);
				if (isPresent(newconvid, convs)) {
					choose(newconvid);
				} else {
					if(convs.items.length > 0) choose(convs.items[0].conid);
				}
			};
	
			// --- Lifecycle ---
			onMounted(() => {
				refresh();
				if (refreshInterval.value === null) {
					refreshInterval.value = setInterval(refresh, 20 * 1000);
				}
			});
	
			onUnmounted(() => {
				if (refreshInterval.value) clearInterval(refreshInterval.value);
			});
	
			return {
				// Refs
				conversationRef,
				vueConversationMenu,
				vueGroupMenu,
				// State
				sv, showInfo, showUserGroup, usersGroups, ievent, itype, errormsg, token, login,
				loading, showList, selectedconvs, searchSelectedConv, selectedusers, tmpselectedusers,
				convs, optionsconversation, optionsgroup,
				// Methods
				listNext, listPrevious, addSelectedUsers, leaveGroup, cancelUsersSearch,
				setCoordsGroups, usersSearchAgain, usersSearch, cancelUsersGroupSearch, usersGroupSearch,
				optionClicked, handleClick, info, reloadInfo, classNumberToRead, resetToread,
				choose, logout, refresh, changeSelectedConv, changeSelectedConvByid, convertTimestamp
			};
		}
	}
	</script>
	
	<template>
		<div class="container-fluid main-wrapper">
			<div class="row h-100">
				<nav id="sidebarMenu" class="col-md-4 col-lg-4 d-md-block sidebar collapse">
					<div id="mycont">
						<ul class="no-bullets">
							<li v-for="cvs in convs.items" :key="cvs.conid">
								<div class="conversation" 
									 :class="{ active: selectedconvs && selectedconvs.conid === cvs.conid }"
									 @click="choose(cvs.conid)">
									<div class="chat">
										<img class="chatimg" v-if="cvs.photo == ''" src="/standard2.gif"/>
										<img class="chatimg" v-if="cvs.photo != ''" v-bind:src="'data:image/gif;base64,' + cvs.photo"/>
										
										<div class="details">
											<div class="firstline">
												<span class="titlename">{{ cvs.name }}</span>
												<span class="titledate">{{ convertTimestamp(cvs.mtimestamp) }}</span>
											</div>
											<div class="message" v-if="cvs.mcategory == 'text'">
												<b>{{cvs.sender}}:</b> {{ cvs.mtext }}
											</div>
											<div class="message" v-if="cvs.mcategory != 'text'">
												<b>{{cvs.sender}}:</b> <img src="/image.svg" class="imagemessage"/>
											</div>
										</div>
									</div>
									<div class="numandmenu">
										<div v-bind:class="classNumberToRead(cvs.toread)">&nbsp;</div>
										<div class="menu"
											 @click.prevent.stop="handleClick($event, cvs)"
											 @mouseenter="setCoordsGroups"
										></div>
									</div>
								</div>
							</li>
						</ul>
	
						<listuser
							id="listuser"
							v-if="showList"
							@search-again="usersSearchAgain"
							@send-users-list="addSelectedUsers"    
							@search-cancel="cancelUsersSearch"
							v-bind:users="selectedusers" 
							v-bind:searchBar="true"
							v-bind:showOk="true"
						/>
						<listuser
							id="listuser"
							v-if="showUserGroup"
							@send-users-list="cancelUsersGroupSearch"
							@search-cancel="cancelUsersGroupSearch"
							v-bind:users="usersGroups"
							v-bind:searchBar="false"
							v-bind:showOk="false"
						/>
					</div>
				</nav>
	
				<vue-simple-context-menu
					element-id="menu-conversation"
					:options="optionsconversation"
					ref="vueConversationMenu"
					@option-clicked="optionClicked"
					class="context-menu-dark"
				/>
				<vue-simple-context-menu
					element-id="menu-group"
					:options="optionsgroup"
					ref="vueGroupMenu"
					@option-clicked="optionClicked"
					class="context-menu-dark"
				/>
				
				<main class="col-md-8 ms-sm-auto col-lg-8 px-md-4 chat-area">
					<conversation
						ref="conversationRef"
						v-bind:sv="sv"
						v-bind:selectedconvs="selectedconvs"
						v-bind:login="login"
						v-bind:token="token"
						@list-next="listNext"
						@list-previous="listPrevious"
						@change-selected-conv="changeSelectedConv"
						@change-selected-conv-byid="changeSelectedConvByid"
						@show-user-info="info"
						@refresh="refresh"
						@logout="logout"
					/>
				</main>
	
				<div v-if="showInfo" class="overlayInfo" :class="{ hiddenInfo: !showInfo }">
					<setinfo 
						:key="ievent"
						v-bind:itype="itype"
						v-bind:ievent="ievent"
						@close-info="showInfo = false"
						@set-info="reloadInfo"
					/>
				</div>
			</div>            
		</div>
	</template>
	
	<style>
	/* --- GLOBALI & LAYOUT --- */
	html, body {
	  margin: 0;
	  padding: 0;
	  height: 100vh;
	  background-color: #0f172a; /* Blu Notte Sfondo */
	  color: #f1f5f9;
	  overflow: hidden; /* Evita scroll pagina intera */
	}
	
	.main-wrapper {
		height: 100vh;
		padding: 0;
		overflow: hidden;
	}
	
	.row.h-100 {
		height: 100%;
		margin: 0;
	}
	
	/* --- SIDEBAR LISTA CHAT --- */
	#sidebarMenu {
		background-color: #1e293b; /* Grigio Bluastro */
		border-right: 1px solid #334155;
		padding: 0;
		height: 100%;
		overflow-y: hidden;
	}
	
	#mycont {
		width: 100%;
		height: 100%;
		overflow-y: auto;
		overflow-x: hidden;
		padding-top: 10px;
	}
	
	/* Scrollbar Sidebar */
	#mycont::-webkit-scrollbar {
		width: 6px;
	}
	#mycont::-webkit-scrollbar-thumb {
		background-color: #334155;
		border-radius: 3px;
	}
	
	ul.no-bullets {
		list-style-type: none;
		padding: 0;
		margin: 0;
	}   
	
	/* --- CHAT ITEM (Card Conversazione) --- */
	div.conversation {
		width: 95%;
		height: 70px; /* Più alto per respiro */
		display: flex;
		background-color: #0f172a; /* Card scura */
		border: 1px solid #334155;
		border-radius: 12px;
		margin: 6px auto;
		transition: all 0.2s ease;
		align-items: center;
		padding-left: 5px;
	}
	
	div.conversation:hover, div.conversation.active {
		background-color: #334155; /* Highlight hover */
		border-color: #2dd4bf; /* Verde Acqua */
		cursor: pointer;
	}
	
	div.chat {
		flex-grow: 1;
		height: 100%;
		display: flex;
		align-items: center;
		background-color: transparent;
		overflow: hidden;
	}
	
	/* Avatar Immagine */
	img.chatimg {
		height: 48px;
		width: 48px;
		min-width: 48px;
		border-radius: 50%;
		object-fit: cover;
		margin-right: 12px;
		border: 2px solid #475569;
	}
	
	div.conversation:hover img.chatimg {
		border-color: #2dd4bf;
	}
	
	/* Dettagli Testo Chat */
	div.details {
		display: flex;
		flex-direction: column;
		justify-content: center;
		width: 100%;
		overflow: hidden;
		padding-right: 10px;
	}
	
	.firstline {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin-bottom: 4px;
	}
	
	.titlename { 
		font-family: 'Segoe UI', sans-serif;
		font-size: 0.95rem;
		font-weight: 600;
		color: #e2e8f0; /* Bianco sporco */
		white-space: nowrap; 
		overflow: hidden; 
		text-overflow: ellipsis;
	}
	
	.titledate { 
		font-size: 0.75rem;
		font-weight: 400;
		color: #94a3b8;
		white-space: nowrap;
	}
	
	.message {
		font-size: 0.85rem;
		font-weight: 400;
		color: #cbd5e1; /* Grigio chiaro */
		white-space: nowrap; 
		overflow: hidden; 
		text-overflow: ellipsis;
		display: block;
	}
	
	.message b {
		color: #2dd4bf; /* Mittente in verde acqua */
		font-weight: 600;
		margin-right: 4px;
	}
	
	.imagemessage {
		height: 16px;
		vertical-align: text-bottom;
		filter: invert(1); /* Icona bianca */
	}
	
	/* --- MENU E COUNTER --- */
	.numandmenu {
		display: flex;
		flex-flow: column;
		align-items: center;
		justify-content: center;
		width: 30px;
		margin-right: 5px;
	}
	
	/* Menu Dots */
	.menu:after {
		content: '\2807'; /* Tre puntini verticali */
		font-size: 1.2rem;
		color: #64748b;
		transition: color 0.2s;
	}
	
	div.conversation:hover .menu:after {
		color: #f1f5f9; /* Bianco su hover */
	}
	
	/* Badge Non letti */
	.number {
		width: 18px;
		height: 18px;
		line-height: 18px;
		border-radius: 50%;
		text-align: center;
		font-size: 0.7rem;
		margin-bottom: 5px;
	}
	
	.yes {
		background-color: #2dd4bf; /* Verde Acqua */
		color: #0f172a; /* Testo scuro */
		font-weight: bold;
		box-shadow: 0 0 5px rgba(45, 212, 191, 0.5);
	}
	
	.no {
		background: transparent;
	}
	
	/* --- AREA CHAT CENTRALE --- */
	.chat-area {
		background-color: #0f172a;
		padding: 0 !important;
		display: flex;
		flex-direction: column;
		height: 100vh;
		border-left: 1px solid #334155;
	}
	
	/* --- OVERLAY INFO --- */
	.overlayInfo {
		z-index: 200;
		transition: 0.3s opacity;
		background: rgba(15, 23, 42, 0.95); /* Backdrop scuro */
		width: 100%;
		height: 100%;
		position: absolute;
		top: 0;
		left: 0;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	
	.hiddenInfo {
		opacity: 0;
		pointer-events: none;
		z-index: -1;
	}
	
	/* --- CONTEXT MENU DARK OVERRIDE --- */
	/* Stili specifici per vue-simple-context-menu in dark mode */
	.vue-simple-context-menu {
		background-color: #1e293b !important;
		border: 1px solid #334155 !important;
		border-radius: 6px !important;
		padding: 5px 0 !important;
		box-shadow: 0 10px 15px rgba(0,0,0,0.5) !important;
	}
	
	.vue-simple-context-menu__item {
		color: #e2e8f0 !important;
		padding: 8px 15px !important;
	}
	
	.vue-simple-context-menu__item:hover {
		background-color: #334155 !important;
		color: #2dd4bf !important;
	}
	
	.vue-simple-context-menu__divider {
		background-color: #475569 !important;
	}
	</style>