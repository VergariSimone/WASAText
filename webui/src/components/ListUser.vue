<script>
    import { ref } from 'vue';
    
    export default {
        name: 'ListUser',
        props: {
            users: { type: Array, default: () => [] },
            searchBar: { type: Boolean, default: true },
            showOk: { type: Boolean, default: false }, // true =selezione multipla x Gruppi, false = singola x nuova chat)
        },
        emits: [
            'send-users-list',
            'search-again',
            'search-cancel',
        ],
        setup(props, { emit }) {
            //state
            const selected = ref([]);
            const searchText = ref("");
    
            // metodi
    
            // conferma selezione (solo x gruppi)
            const ok = () => {
                if (selected.value.length > 0) {
                    emit('send-users-list', selected.value);    
                } else {
                    emit('search-cancel');
                }
            };
    
            const cancel = () => {
                emit('search-cancel');
            };
    
            const searchAgain = () => {
                emit('search-again', searchText.value);
            };
    
            const userChoose = (login) => {
                // selezione singola (nuova chat)
                if (!props.showOk) {
                    emit('send-users-list', [login]);
                    return;
                }
    
                // selezione multipla (gruppi)
                //tToggle selezione (aggiungi/rimuovi array)
                if (selected.value.includes(login)) {
                    const index = selected.value.indexOf(login);
                    if (index > -1) selected.value.splice(index, 1);
                } else {
                    selected.value.push(login);
                }
            };
    
            // per capire se un utente è selezionato (CSS)
            const isSelected = (login) => {
                return selected.value.includes(login);
            };
    
            return {
                selected,
                searchText,
                ok,
                cancel,
                searchAgain,
                userChoose,
                isSelected
            };
        }
    }
    </script>
    
    <template>
        <div id="listuser">
            <div class="list-header">
                <div class="d-flex justify-content-between align-items-center mb-2">
                    <button type="button" @click="cancel" class="btn-custom">Annulla</button>
                    <h5 class="title-modal" v-if="showOk">Seleziona partecipanti</h5>
                    <h5 class="title-modal" v-else>Seleziona utente</h5>
                    <button v-if="showOk" type="button" @click="ok" class="btn-custom btn-ok">Fatto</button>
                    <div v-else style="width: 60px;"></div> </div>
    
                <div v-if="searchBar" class="search-container">
                    <input 
                        type="text" 
                        v-model="searchText" 
                        @keyup.enter="searchAgain"
                        class="listusersearch icon" 
                        placeholder="Cerca..."
                    />
                </div>
            </div>
    
            <div class="listcontext">
                <div v-if="!users || users.length === 0" class="no-results">
                    <div class="user-subcard justify-content-center">
                        <span class="text-muted">Nessun utente trovato</span>
                    </div>
                </div>
    
                <div v-for="u in users" :key="u.login">
                    <div v-if="u && u.login != '0'" 
                         class="user-card" 
                         :class="{ 'user-card-selected': isSelected(u.login) }"
                         @click.stop="userChoose(u.login)">
                        
                        <div class="user-subcard">
                            <img class="user-card-img" v-if="!u.photo" src="/standard.gif"/>
                            
                            <img class="user-card-img" v-else v-bind:src="'data:' + (u.photo.startsWith('/9j/') ? 'image/jpeg' : (u.photo.startsWith('iVBOR') ? 'image/png' : 'image/gif')) + ';base64,' + u.photo"/>
                            
                            <div class="user-details">
                                <div class="user-line">
                                    <span class="user-name">
                                        {{ u.name }}
                                        <span class="user-login" v-if="u.categoryinfo=='user'">@{{ u.login }}</span>
                                    </span>
                                </div>
                            </div>
    
                            <div v-if="isSelected(u.login)" class="check-icon">✓</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </template>
    
    <style scoped>


    /* contenitore principale  */
    #listuser {
       
        background-color: #1e293b;
        border: 1px solid #334155;
        border-radius: 12px;
        padding: 15px;
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
        width: 100%; 
        height: 100%;
        display: flex;
        flex-direction: column;
    }
    
    /* header*/
    .list-header {
        margin-bottom: 10px;
    }
    .title-modal {
        color: #f1f5f9;
        font-size: 1rem;
        margin: 0;
    }
    
    /* bottoni*/
    .btn-custom {
        background-color: transparent;
        border: 1px solid #475569;
        color: #94a3b8;
        border-radius: 20px;
        padding: 5px 15px;
        font-size: 0.85rem;
        cursor: pointer;
        transition: all 0.2s ease;
    }
    
    .btn-custom:hover {
        border-color: #cbd5e1;
        color: #fff;
        background-color: rgba(255, 255, 255, 0.05);
    }
    
    .btn-ok {
        border-color: #0d9488;
        color: #2dd4bf;
        font-weight: 600;
    }
    
    .btn-ok:hover {
        background-color: rgba(45, 212, 191, 0.1);
        color: white;
    }
    
    /* input ricerca */
    .search-container {
        width: 100%;
    }
    .listusersearch {
        background-color: #0f172a;
        border-radius: 8px;
        border: 1px solid #334155;
        font-size: 0.9rem;
        height: 38px;
        color: #f1f5f9;
        width: 100%;
        outline: none;
        transition: border-color 0.2s;
        padding-left: 35px; 
    }
    
    .listusersearch:focus {
        border-color: #2dd4bf;
    }
    
    .icon {
        background: url("/search.svg") no-repeat 10px center;
        background-size: 16px;
        filter: invert(1) brightness(2); 
    }
    
    /* lista scrollabile*/
    .listcontext {
        overflow-y: auto;
        flex-grow: 1; /* occupa lo spazio rimanente */
        margin-top: 5px;
        border: 1px solid #334155;
        border-radius: 8px;
        background-color: #0f172a;
        padding: 5px;
        min-height: 200px;
        max-height: 60vh;
    }
    
    /* scrollbar */
    .listcontext::-webkit-scrollbar { width: 6px; }
    .listcontext::-webkit-scrollbar-thumb { background-color: #334155; border-radius: 4px; }
    
    /* card utente*/
    .user-card {
        width: 98%;
        cursor: pointer;
        height: 56px;
        display: flex;
        background-color: #1e293b;
        border: 1px solid #334155;
        border-radius: 8px;
        margin: 5px auto;
        transition: all 0.1s ease;
        align-items: center;
    }
    
    .user-card:hover {
        border-color: #64748b;
        background-color: #263344;
    }
    
    /* selezione attivata */
    .user-card-selected {
        background-color: rgba(13, 148, 136, 0.2) !important; /* Verde scuro trasparente */
        border-color: #2dd4bf !important;
    }
    
    .user-subcard {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        padding: 0 10px;
        position: relative;
    }
    
    .check-icon {
        position: absolute;
        right: 15px;
        color: #2dd4bf;
        font-weight: bold;
        font-size: 1.2rem;
    }
    
    /* immagini e testo */
    .user-card-img {
        height: 38px;
        width: 38px;
        border-radius: 50%;
        object-fit: cover;
        margin-right: 12px;
        border: 1px solid #475569;
    }
    
    .user-details {
        display: flex;
        flex-direction: column;
        justify-content: center;
        overflow: hidden;
    }
    
    .user-name { 
        font-family: 'Segoe UI', sans-serif;
        font-size: 0.95rem;
        font-weight: 500;
        color: #e2e8f0;
    }
    
    .user-login { 
        font-size: 0.8rem;
        color: #94a3b8;
        margin-left: 5px;
    }
    
    .no-results {
        padding: 20px;
        text-align: center;
        color: #64748b;
    }
    </style>