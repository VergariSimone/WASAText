<script>
  import { ref } from 'vue';
  import EmojiPicker from 'vue3-emoji-picker';
  import 'vue3-emoji-picker/css';
  
  export default {
      components: {
          'emojipicker': EmojiPicker,
      },
      props: [
          'msgid',
          'conid',
          'topX',
          'topY',
      ],
      emits: [
          'send-comment'
      ],
      setup(props, { emit }) {
          
          // --- State ---
          const comment = ref(null);
          const isSCPickerOpen = ref(false); 
  
          // --- Methods ---
  
          const sendComment = () => {
              emit('send-comment', comment.value, props.msgid, props.conid);
          };
  
          const close = () => {
              comment.value = null;
              // Invia null per segnalare chiusura/annullamento al genitore
              emit('send-comment', comment.value, null, null);
          };
  
          const onSelectEmoji = (emoji) => {
              comment.value = emoji;
              emit('send-comment', comment.value, props.msgid, props.conid);
          };
  
          return {
              comment,
              isSCPickerOpen,
              sendComment,
              close,
              onSelectEmoji
          };
      }
  }
  </script>
  
  <template>
      <div class="container">
          <i>&nbsp;&nbsp;&nbsp;&nbsp;</i>
          <span class="container-span" @click="close()">close</span>
          <i>&nbsp;&nbsp;&nbsp;&nbsp;</i>
          
          <i/> 
  
          <emojipicker 
              :native="true"
              :hide-search="false"
              :disable-skin-tones="true"
              :hide-group-names="true"
              :disabled-groups='["food_drink", "activities", "travel_places", "objects", "symbols", "flags"]'
              @select="onSelectEmoji"
              class="custom-picker"
          />
          
          <i/>
          <i>&nbsp;&nbsp;&nbsp;&nbsp;</i>
          <i>&nbsp;&nbsp;&nbsp;&nbsp;</i>
          <i>&nbsp;&nbsp;&nbsp;&nbsp;</i>
      </div>
  </template>
  
  <style scoped>
  /* --- STILI BLU NOTTE / DARK MODE --- */
  
  .container {
      display: grid;
      grid-template-columns: auto auto auto;
      justify-content: center;
      align-items: center;
      
      position: relative;
      top: 50px;
      
      /* Tema Scuro */
      background-color: #1e293b; /* Grigio Bluastro */
      border: 1px solid #334155;
      
      padding: 8px;
      width: fit-content;
      border-radius: 12px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5); /* Ombra profonda */
  }
  
  .container-span {
      position: relative;
      margin: 5px 0 10px 0;
      
      /* Stile Bottone Close */
      background-color: transparent;
      color: #94a3b8;
      border: 1px solid #475569;
      border-radius: 20px;
      
      padding: 2px 15px;
      text-align: center;
      cursor: pointer;
      font-size: 0.85rem;
      font-family: sans-serif;
      transition: all 0.2s ease;
  }
  
  .container-span:hover {
      background-color: rgba(255, 255, 255, 0.1);
      color: #f1f5f9;
      border-color: #cbd5e1;
  }
  
  /* Stili per le spaziature (i tags) */
  i {
      display: block;
  }
  
  /* --- OVERRIDE STILI EMOJI PICKER (DARK MODE) --- */
  /* Utilizzo :deep() per penetrare nel componente figlio e sovrascrivere le variabili CSS */
  :deep(.v3-emoji-picker) {
      width: 360px !important;
      z-index: 1000;
      box-shadow: none !important;
      border-radius: 8px;
      
      /* Variabili interne della libreria vue3-emoji-picker per il tema scuro */
      --ep-color-bg: #0f172a;       /* Sfondo interno scuro */
      --ep-color-text: #f1f5f9;     /* Testo bianco */
      --ep-color-border: #334155;   /* Bordi grigio scuro */
      --ep-color-hover: #1e293b;    /* Hover items */
      --ep-color-active: #2dd4bf;   /* Active items (Verde Acqua) */
      --ep-search-bg: #1e293b;      /* Sfondo barra ricerca */
      --ep-search-text: #f1f5f9;    /* Testo ricerca */
      border: 1px solid #334155;
  }
  
  </style>