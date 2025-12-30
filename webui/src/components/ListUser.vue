<script>
  import { ref } from 'vue';
  
  export default {
      props: [
          "list",
          "users",
          "searchBar",
          "showOk",
      ],
      emits: [
          'send-users-list',
          'search-again',
          'search-cancel',
      ],
      setup(props, { emit }) {
          // --- State ---
          const selected = ref([]);
  
          // --- Methods ---
  
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
              // Accesso diretto al DOM per il valore di ricerca
              const val = document.getElementById('list-search-input').value;
              emit('search-again', val);
          };
  
          const userChoose = (login) => {
              if (!props.showOk) return;
  
              // Manipolazione diretta del DOM per la selezione visiva
              var element = document.getElementById('uc' + login);
              
              if (!selected.value.includes(login)) {
                  selected.value.push(login);
                  // Aggiunta classe manuale
                  if(element) element.classList.add("user-card-selected");
              } else {
                  // Rimozione classe manuale
                  if(element) element.classList.remove("user-card-selected");
                  const index = selected.value.indexOf(login);
                  if (index > -1) {
                      selected.value.splice(index, 1);
                  }
              }
          };
  
          return {
              selected,
              ok,
              cancel,
              searchAgain,
              userChoose
          };
      }
  }
  </script>
  
  <template>
      <div id="listuser">
          <div class="listbuttons">
              <i class="allspan">&nbsp;</i>
              <i></i>
              <button type="button" @click="cancel" class="btn-custom">Cancel</button>
              <i></i>
              <button v-if="showOk" type="button" @click="ok" class="btn-custom btn-ok">Ok</button>
              <i v-if="!showOk"></i>
              <i></i>
              <i class="allspan">&nbsp;</i>
              <input v-if="searchBar" id="list-search-input" type="text" class="listusersearch icon" @change="searchAgain"/>
              <i v-if="!searchBar" class="allspan">&nbsp;</i>
              <i class="allspan">&nbsp;</i>
          </div>
  
          <div class="listcontext">
              <div v-for="u in users" :key="u.login">
                  <div v-if="u && u.login != '0'" :id="'uc' + u.login" class="user-card" @click.stop="userChoose(u.login)">
                      <div class="user-subcard">
                          <img class="user-card-img" v-if="u.photo == ''" src="/standard2.gif"/>
                          <img class="user-card-img" v-if="u.photo != ''" v-bind:src="'data:image/gif;base64,' + u.photo"/>
                          
                          <div class="user-details">
                              <div class="user-line">
                                  <span class="user-name">
                                      {{ u.name }}
                                      <span class="user-login" v-if="u.categoryinfo=='user'">({{ u.login }})</span>
                                  </span>
                              </div>
                          </div>
                      </div>
                  </div>
  
                  <div v-if="u && u.login == '0'" :id="u.login" class="user-card no-hover">
                      <div class="user-subcard">
                          <img class="user-card-img" src="/void.gif"/>
                          <div class="user-details">
                              <div class="user-line">
                                  <span class="user-name">Search produced no results!</span>
                              </div>
                          </div>
                      </div>
                  </div>
              </div>
          </div>
      </div>
  </template>
  
  <style scoped>
  /* --- TEMA BLU NOTTE / VERDE ACQUA --- */
  
  /* --- CONTENITORE PRINCIPALE --- */
  #listuser {
      align-items: center;
      background-color: #1e293b; /* Grigio Bluastro */
      border: 1px solid #334155;
      border-radius: 12px;
      display: none; /* Gestito via JS dal padre */
      padding: 15px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
      position: fixed;
      width: 360px;
      height: fit-content;
      z-index: 100;
  }
  
  /* --- GRIGLIA PULSANTI --- */
  .listbuttons {
      display: grid;
      grid-template-columns: 10% 30% 10% 30% 10%; /* Bilanciato */
      margin-bottom: 10px;
      align-items: center;
  }
  
  .allspan {
      grid-column: 1 / span 5;
      height: 5px;
  }
  
  i {
      height: 10px;
      display: block;
  }
  
  /* --- BOTTONI STYLE --- */
  .btn-custom {
      background-color: transparent;
      border: 1px solid #475569;
      color: #94a3b8;
      border-radius: 20px;
      padding: 5px 10px;
      font-size: 0.85rem;
      cursor: pointer;
      transition: all 0.2s ease;
      width: 100%;
  }
  
  .btn-custom:hover {
      border-color: #cbd5e1;
      color: #fff;
      background-color: rgba(255, 255, 255, 0.05);
  }
  
  .btn-ok {
      border-color: #0d9488;
      color: #2dd4bf; /* Verde Acqua */
      font-weight: 600;
  }
  
  .btn-ok:hover {
      background-color: rgba(45, 212, 191, 0.1);
      border-color: #2dd4bf;
      color: #fff;
      box-shadow: 0 0 10px rgba(45, 212, 191, 0.3);
  }
  
  /* --- INPUT DI RICERCA --- */
  .listusersearch {
      background-color: #0f172a; /* Blu Notte profondo */
      border-radius: 8px;
      border: 1px solid #334155;
      font-size: 0.9rem;
      height: 36px;
      grid-column: 2 / span 3;
      color: #f1f5f9;
      width: 100%;
      outline: none;
      transition: border-color 0.2s;
  }
  
  .listusersearch:focus {
      border-color: #2dd4bf;
  }
  
  .icon {
      padding-left: 35px;
      background: url("/search.svg") no-repeat 10px center;
      background-size: 16px;
      /* Inverte il colore dell'icona SVG (assumendo sia nera) per renderla bianca */
      filter: invert(1) brightness(2); 
  }
  
  /* --- AREA LISTA SCROLLABILE --- */
  .listcontext {
      align-items: center;
      overflow-y: auto;
      height: 250px; /* Leggermente più alto */
      margin-bottom: 10px;
      border: 1px solid #334155;
      border-radius: 8px;
      background-color: #0f172a; /* Interno scuro */
      padding: 5px;
  }
  
  /* Custom Scrollbar */
  .listcontext::-webkit-scrollbar {
      width: 8px;
  }
  .listcontext::-webkit-scrollbar-thumb {
      background-color: #334155;
      border-radius: 4px;
  }
  
  /* --- CARD UTENTE --- */
  .user-card {
      width: 95%;
      cursor: pointer;
      height: 54px;
      display: flex;
      background-color: #1e293b; /* Colore card */
      border: 1px solid #334155;
      border-radius: 8px;
      margin: 5px auto;
      transition: all 0.2s ease;
      align-items: center;
  }
  
  .user-card:hover {
      border-color: #475569;
      transform: translateX(2px);
  }
  
  /* Classe aggiunta dinamicamente via JS */
  .user-card-selected {
      background-color: rgba(45, 212, 191, 0.15) !important; /* Verde acqua trasparente */
      border-color: #2dd4bf !important;
  }
  
  .no-hover {
      cursor: default;
      border-color: transparent;
  }
  .no-hover:hover {
      transform: none;
      border-color: transparent;
  }
  
  .user-subcard {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      padding: 0 10px;
  }
  
  /* --- IMMAGINE UTENTE --- */
  img.user-card-img {
      height: 36px;
      width: 36px;
      min-width: 36px;
      border-radius: 50%;
      object-fit: cover;
      margin-right: 12px;
      border: 1px solid #475569;
  }
  
  /* --- DETTAGLI UTENTE --- */
  div.user-details {
      display: flex;
      flex-direction: column;
      justify-content: center;
      width: 100%;
      overflow: hidden;
  }
  
  .user-line {
      display: block;
      width: 100%;
  }
  
  .user-name { 
      font-family: 'Segoe UI', sans-serif;
      font-size: 0.95rem;
      font-weight: 500;
      color: #e2e8f0; /* Bianco sporco */
      white-space: nowrap; 
      overflow: hidden;
      text-overflow: ellipsis;
      display: block;
  }
  
  .user-login { 
      font-size: 0.8rem;
      font-weight: 400;
      color: #94a3b8; /* Grigio chiaro */
      margin-left: 5px;
  }
  </style>
    