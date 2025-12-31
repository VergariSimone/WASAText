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
	
			// --- Refs ai componenti figli ---
			const conversationRef = ref(null);
			const vueConversationMenu = ref(null);
			const vueGroupMenu = ref(null);
	
			// --- State Reattivo ---
			const sv = reactive({ start: 0, limit: 50, last: null });
			
			// Modali
			const showInfo = ref(false);
			const showUserGroup = ref(false); // Mostra membri gruppo
			const showList = ref(false);      // Mostra ricerca utenti generica (es. aggiungi a gruppo)
	
			// Dati
			const token = ref(null);
			const login = ref(null);
			const loading = ref(false);
			const errormsg = ref(null);
			
			// Selezioni
			const selectedconvs = ref(null);
			const searchSelectedConv = ref(null); // ID conversazione target per operazioni (add user, etc)
			
			// Liste Dati
			const usersGroups = ref([]);    // Membri di un gruppo
			const selectedusers = ref([]);  // Risultati ricerca utenti
			const convs = reactive({ items: [] }); // Lista conversazioni sidebar
	
			// Info Panel
			const ievent = ref(null);
			const itype = ref(null);
	
			// Utils
			const refreshInterval = ref(null);
			const topX = ref(0);
			const topY = ref(0);
	
			// Menu Options
			const optionsconversation = [
				{ name: 'Apri', slug: 'open' },
				{ type: 'divider' },
				{ name: 'Invia messaggio', slug: 'send' },
			];
	
			const optionsgroup = [
				{ name: 'Apri', slug: 'open' },
				{ name: 'Info Gruppo', slug: 'info' },
				{ name: 'Membri', slug: 'show' },
				{ name: 'Aggiungi utente', slug: 'add' },
				{ type: 'divider' },
				{ name: 'Abbandona gruppo', slug: 'leave' },
			];
	
			// --- METHODS ---
	
			// Posizionamento Popup (obsoleto con le nuove modali centrate, ma mantenuto per compatibilità)
			const calcPosition = (element) => {
				if(!element) return [0,0];
				var html = document.documentElement;
				var X = 0, Y = 0;
				var H = element.offsetHeight, W = element.offsetWidth;
				Y = Math.min(topY.value + H, html.clientHeight) - H;
				X = Math.min(topX.value - 30 + W, html.clientWidth) - W;
				return [X, Y];
			};
	
			const isPresent = (ncid, cvs_obj) => {
				if (cvs_obj && ncid && cvs_obj.items) {
					return cvs_obj.items.some(el => el.conid == ncid);
				}
				return false;
			};
	
			const convertTimestamp = (timestamp) => {
				if(!timestamp) return "";
				const localDate = new Date(`${timestamp}Z`);
				return localDate.toLocaleString("it-IT", { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' });
			};
	
			const resetToread = (conid) => {
				let cvs = convs.items.find(c => c.conid == conid);
				if (cvs) cvs.toread = 0;
			};
	
			// Carica messaggi di una conversazione specifica
			const choose = async (conid) => {
				resetToread(conid);
				errormsg.value = null;
				let config = { headers: { Authorization: "Bearer " + token.value } };
				
				try {
					if (sv.last != conid) sv.start = 0;
					sv.last = conid;
					
					let response = await $axios.get(`/conversations/byid/${conid}?start=${sv.start}&limit=${sv.limit}`, config);
					selectedconvs.value = response.data;
				} catch (e) { errormsg.value = e.toString(); }
			};
	
			// Ricarica la lista conversazioni (Sidebar)
			const refresh = async (chooseFlag = true) => {
				loading.value = true;
				token.value = sessionStorage.getItem('token');
				login.value = sessionStorage.getItem('login');
				errormsg.value = null;
				
				let config = { headers: { Authorization: "Bearer " + token.value } };
				
				try {
					let response = await $axios.get("/conversations", config);
					convs.items = response.data;
				} catch (e) { errormsg.value = e.toString(); }
	
				// Se c'è una conv selezionata, ricaricala per vedere nuovi messaggi
				if (selectedconvs.value) {
					if (isPresent(selectedconvs.value.conid, convs)) {
						// Mantiene la selezione
						// Opzionale: chiamare choose() qui se vuoi auto-aggiornare i messaggi aperti
					} else if (convs.items.length > 0 && !selectedconvs.value) {
						// Seleziona la prima se nulla è selezionato
						selectedconvs.value = convs.items[0];
						choose(convs.items[0].conid);
					}
				}
				loading.value = false;
			};
	
			const listNext = () => { sv.start += sv.limit; choose(sv.last); };
			const listPrevious = () => { sv.start = Math.max(sv.start - sv.limit, 0); choose(sv.last); };
	
			// Aggiungi utente a gruppo esistente (da Sidebar Context Menu)
			const addSelectedUsers = async (users) => {
				showList.value = false;
				var out = [];
				for (var key in users) out.push('users=' + encodeURIComponent(users[key]));
				
				try {
					await $axios.put(`/conversations/group?groupID=${searchSelectedConv.value}&${out.join('&')}`, {}, 
						{ headers: { Authorization: "Bearer " + token.value } }
					);
					alert("Utente aggiunto!");
				} catch (e) { console.error(e); }
			};
	
			const leaveGroup = async (convid) => {
				try {
					await $axios.put(`/conversations/byid/${convid}`, {}, { headers: { Authorization: "Bearer " + token.value } });
				} catch (e) { console.error(e); }
			};
	
			const setCoordsGroups = (e) => {
				topX.value = e.clientX;
				topY.value = e.clientY;
			};
	
			// Logica Ricerca Utenti (Generica)
			const searchUsersAPI = async (search = "") => {
				let config = { headers: { Authorization: "Bearer " + token.value } };
				let response = await $axios.get(`/users?search=${search}`, config);
				let data = response.data;
				
				if(!data) return [];
				return Array.isArray(data) ? data : [data];
			};
	
			const usersSearchAgain = async (search) => {
				selectedusers.value = await searchUsersAPI(search);
			};
	
			const usersSearch = async () => {
				showList.value = true;
				selectedusers.value = await searchUsersAPI("");
			};
	
			const usersGroupSearch = async (conid) => {
				showUserGroup.value = true;
				let config = { headers: { Authorization: "Bearer " + token.value } };
				let response = await $axios.get(`/conversations/byid/${conid}/info`, config);
				let users = response.data.user;
				usersGroups.value = Array.isArray(users) ? users : [users];
			};
	
			// Gestione Click Menu Contestuale
			const optionClicked = (event) => {
				const slug = event.option.slug;
				const item = event.item;
	
				if (slug == 'add') {
					searchSelectedConv.value = item.conid;
					usersSearch();
				} else if (slug == 'leave') {
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
				}
			};
	
			const handleClick = (evt, conv) => {
				// Mostra menu diverso se gruppo o user
				if (conv.categoryinfo == 'conversation') {
					if(vueGroupMenu.value) vueGroupMenu.value.hideContextMenu();
					if(vueConversationMenu.value) vueConversationMenu.value.showMenu(evt, conv);
				} else {
					if(vueConversationMenu.value) vueConversationMenu.value.hideContextMenu();
					if(vueGroupMenu.value) vueGroupMenu.value.showMenu(evt, conv);
				}
			};
	
			// Info Panel
			const reloadInfo = async (infoType, conid) => {
				let config = { headers: { Authorization: "Bearer " + token.value } };
				if (itype.value == 'group') {
					let response = await $axios.get(`conversations/byid/${conid}/info`, config);
					ievent.value = response.data;
				} else {
					let response = await $axios.get(`/info`, config);
					ievent.value = response.data;
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
	
			// Eventi emessi da Conversation.vue
			const changeSelectedConv = (newconv) => {
				selectedconvs.value = newconv;
				refresh(false); // Aggiorna sidebar per muovere la chat in alto
			};
	
			const changeSelectedConvByid = async (newconvid) => {
				await refresh(false);
				choose(newconvid);
			};
	
			// Lifecycle
			onMounted(() => {
				refresh();
				if (refreshInterval.value === null) {
					refreshInterval.value = setInterval(() => refresh(false), 5000); // Polling ogni 5s
				}
			});
	
			onUnmounted(() => {
				if (refreshInterval.value) clearInterval(refreshInterval.value);
			});
	
			return {
				// Refs
				conversationRef, vueConversationMenu, vueGroupMenu,
				// State
				sv, showInfo, showUserGroup, usersGroups, ievent, itype, errormsg, token, login,
				showList, selectedconvs, searchSelectedConv, selectedusers, convs,
				optionsconversation, optionsgroup,
				// Methods
				listNext, listPrevious, addSelectedUsers, leaveGroup, 
				usersSearchAgain, usersSearch, usersGroupSearch,
				optionClicked, handleClick, info, reloadInfo, classNumberToRead,
				choose, logout, refresh, changeSelectedConv, changeSelectedConvByid, convertTimestamp,
				setCoordsGroups
			};
		}
	}
	</script>
	
	<template>
		<div class="container-fluid main-wrapper">
			<div class="row h-100 g-0">
				
				<nav id="sidebarMenu" class="col-md-4 col-lg-3 d-md-block sidebar collapse">
					<div id="mycont">
						
						<ul class="no-bullets">
							<li v-for="cvs in convs.items" :key="cvs.conid">
								<div class="conversation" 
									 :class="{ active: selectedconvs && selectedconvs.conid === cvs.conid }"
									 @click="choose(cvs.conid)">
									
									<div class="chat">
										<img class="chatimg" v-if="!cvs.photo" src="/standard2.gif"/>
										<img class="chatimg" v-else :src="'data:image/gif;base64,' + cvs.photo"/>
										
										<div class="details">
											<div class="firstline">
												<span class="titlename">{{ cvs.name }}</span>
												<span class="titledate">{{ convertTimestamp(cvs.mtimestamp) }}</span>
											</div>
											<div class="message">
												<b v-if="cvs.sender">{{cvs.sender}}:</b> 
												<span v-if="cvs.mcategory == 'text'">{{ cvs.mtext }}</span>
												<span v-else>📷 Immagine</span>
											</div>
										</div>
									</div>
									
									<div class="numandmenu">
										<div v-if="cvs.toread > 0" class="number yes">{{ cvs.toread }}</div>
										<div class="menu" @click.prevent.stop="handleClick($event, cvs)"></div>
									</div>
								</div>
							</li>
						</ul>
	
					</div>
				</nav>
	
				<main class="col-md-8 ms-sm-auto col-lg-9 p-0 chat-area">
					<conversation
						ref="conversationRef"
						:sv="sv"
						:selectedconvs="selectedconvs"
						:login="login"
						:token="token"
						@list-next="listNext"
						@list-previous="listPrevious"
						@change-selected-conv="changeSelectedConv"
						@change-selected-conv-byid="changeSelectedConvByid"
						@show-user-info="info"
						@refresh="refresh"
						@logout="logout"
					/>
				</main>
	
				<div v-if="showList" class="modal-overlay">
					<div class="modal-box">
						<listuser
							@search-again="usersSearchAgain"
							@send-users-list="addSelectedUsers"    
							@search-cancel="showList = false"
							:users="selectedusers" 
							:searchBar="true"
							:showOk="true"
						/>
					</div>
				</div>
	
				<div v-if="showUserGroup" class="modal-overlay">
					<div class="modal-box">
						<listuser
							@send-users-list="showUserGroup = false"
							@search-cancel="showUserGroup = false"
							:users="usersGroups"
							:searchBar="false"
							:showOk="false"
						/>
					</div>
				</div>
	
				<div v-if="showInfo" class="overlayInfo" :class="{ hiddenInfo: !showInfo }">
					<setinfo 
						:key="ievent"
						:itype="itype"
						:ievent="ievent"
						@close-info="showInfo = false"
						@set-info="reloadInfo"
					/>
				</div>
			</div> 
	
			<vue-simple-context-menu element-id="menu-conversation" :options="optionsconversation" ref="vueConversationMenu" @option-clicked="optionClicked" class="context-menu-dark" />
			<vue-simple-context-menu element-id="menu-group" :options="optionsgroup" ref="vueGroupMenu" @option-clicked="optionClicked" class="context-menu-dark" />
			
		</div>
	</template>
	
	<style>
	/* --- GLOBALI & LAYOUT --- */
	html, body {
	  margin: 0; padding: 0; height: 100vh;
	  background-color: #0f172a; color: #f1f5f9;
	  overflow: hidden;
	}
	
	.main-wrapper { height: 100vh; padding: 0; overflow: hidden; }
	
	/* --- SIDEBAR --- */
	#sidebarMenu {
		background-color: #1e293b;
		border-right: 1px solid #334155;
		padding: 0;
		height: 100vh;
		display: flex;
		flex-direction: column;
	}
	
	#mycont {
		flex-grow: 1;
		overflow-y: auto;
		padding: 10px 0;
	}
	#mycont::-webkit-scrollbar { width: 6px; }
	#mycont::-webkit-scrollbar-thumb { background-color: #334155; border-radius: 3px; }
	
	ul.no-bullets { list-style: none; padding: 0; margin: 0; }
	
	/* --- CONVERSATION CARD --- */
	div.conversation {
		width: 95%; height: 72px;
		display: flex; align-items: center;
		background-color: transparent;
		border-bottom: 1px solid #334155;
		margin: 0 auto;
		padding: 5px 10px;
		border-radius: 8px;
		transition: background 0.2s;
	}
	
	div.conversation:hover, div.conversation.active {
		background-color: #334155;
		border-color: transparent;
		cursor: pointer;
	}
	
	div.chat { flex-grow: 1; display: flex; align-items: center; overflow: hidden; }
	
	img.chatimg {
		height: 48px; width: 48px; min-width: 48px;
		border-radius: 50%; object-fit: cover; margin-right: 12px;
		border: 1px solid #475569;
	}
	
	div.details {
		display: flex; flex-direction: column; justify-content: center;
		overflow: hidden; width: 100%;
	}
	
	.firstline { display: flex; justify-content: space-between; align-items: baseline; margin-bottom: 2px; }
	
	.titlename { 
		font-size: 0.95rem; font-weight: 600; color: #f1f5f9;
		white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
	}
	.titledate { font-size: 0.75rem; color: #94a3b8; white-space: nowrap; }
	
	.message {
		font-size: 0.85rem; color: #cbd5e1;
		white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
	}
	.message b { color: #5eead4; margin-right: 4px; font-weight: 500; }
	
	/* --- MENU & BADGE --- */
	.numandmenu {
		display: flex; flex-direction: column; align-items: center;
		width: 25px; margin-left: 5px;
	}
	
	.menu:after {
		content: '⋮'; font-size: 1.4rem; color: #64748b; margin-top: 5px; cursor: pointer;
	}
	div.conversation:hover .menu:after { color: #f1f5f9; }
	
	.number {
		min-width: 18px; height: 18px; line-height: 18px;
		border-radius: 9px; text-align: center; font-size: 0.7rem;
		padding: 0 4px; margin-bottom: 2px;
	}
	.yes { background-color: #2dd4bf; color: #0f172a; font-weight: bold; }
	
	/* --- AREA CHAT --- */
	.chat-area {
		background-color: #0f172a;
		height: 100vh;
		border-left: 1px solid #334155;
		position: relative;
	}
	
	/* --- MODALI CENTRATE --- */
	.modal-overlay {
		position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
		background: rgba(0,0,0,0.7); z-index: 1050;
		display: flex; justify-content: center; align-items: center;
	}
	.modal-box {
		background: #1e293b; padding: 20px; border-radius: 10px;
		width: 400px; max-width: 95%; max-height: 85vh;
		border: 1px solid #475569; display: flex; flex-direction: column;
	}
	
	/* --- OVERLAY INFO --- */
	.overlayInfo {
		z-index: 2000; transition: 0.3s opacity;
		background: rgba(15, 23, 42, 0.95);
		width: 100%; height: 100%; position: fixed; top: 0; left: 0;
		display: flex; align-items: center; justify-content: center;
	}
	.hiddenInfo { opacity: 0; pointer-events: none; }
	
	/* --- CONTEXT MENU DARK --- */
	.vue-simple-context-menu {
		background-color: #1e293b !important;
		border: 1px solid #334155 !important;
		border-radius: 6px !important;
		box-shadow: 0 5px 15px rgba(0,0,0,0.5) !important;
		padding: 5px 0 !important;
	}
	.vue-simple-context-menu__item { color: #f1f5f9 !important; padding: 8px 15px !important; }
	.vue-simple-context-menu__item:hover { background-color: #334155 !important; color: #2dd4bf !important; }
	.vue-simple-context-menu__divider { background-color: #475569 !important; }
	</style>