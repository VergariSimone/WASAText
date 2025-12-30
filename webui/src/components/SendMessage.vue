<script>
  import { ref, nextTick } from 'vue';
  import EmojiPicker from 'vue3-emoji-picker';
  import 'vue3-emoji-picker/css';
  import ImageMessage from '../components/ImageMessage.vue';
  
  export default {
      components: {
          'emojipicker': EmojiPicker,
          'imagemessage': ImageMessage,
      },
      props: [
          'conversation',
      ],
      emits: [
          'send-message',
      ],
      setup(props, { emit }) {
          
          // --- State ---
          const message = ref(null);
          const isSMPickerOpen = ref(false);
          const isEnabled = ref(true);
          const showImageMessage = ref(false);
          const image = ref(null);
  
          // --- Methods ---
  
          const toggleEnable = () => {
              var input = document.getElementById("inputRef");
              if (isEnabled.value) {
                  isEnabled.value = false;
                  if(input) input.disabled = true;
              } else {
                  isEnabled.value = true;
                  if(input) input.disabled = false;
              }
          };
  
          const toggleImageMessage = () => {
              if ((!isEnabled.value) && (!showImageMessage.value)) return;
              
              var input = document.getElementById('inputRef');
              if(input) input.value = '';
              
              showImageMessage.value = !showImageMessage.value;
              
              if (showImageMessage.value) {
                  message.value = null;
                  if(input) input.disabled = true;
              } else {
                  if(input) input.disabled = false;
              }
          };
  
          const sendmessage = () => {
              
              if ((!isEnabled.value) || (showImageMessage.value)) return;
  
              emit('send-message', message.value, image.value);
              message.value = "";
              image.value = null;
          };
  
          const imageMessageCreate = (img, text) => {
              if (img || (text != '')) {
                  image.value = img;
                  message.value = text;
                  showImageMessage.value = !showImageMessage.value;
                  
                  var input = document.getElementById("inputRef");
                  isEnabled.value = true;
                  if(input) input.disabled = false;
                  
                  
                  emit('send-message', message.value, image.value);
                  
                  message.value = "";
                  image.value = null;
              }
          };
  
          const toggleEmoji = async () => {
              if ((isEnabled.value) || (showImageMessage.value)) {
                  isSMPickerOpen.value = !isSMPickerOpen.value;
              }
              
              await nextTick(); // Aspetta che il DOM si aggiorni se v-if cambia
  
              if (isSMPickerOpen.value) {
                  var grinBtn = document.getElementById('grinemoji');
                  var element = document.getElementById('overlay');
                  
                  if (grinBtn && element) {
                      var rect = grinBtn.getBoundingClientRect();
                      element.style.display = 'block';
                      
                      var shift = 0;
                      
                      if (showImageMessage.value) {
                          shift = 300; 
                      }
                      
                      element.style.top = (rect.top - 420).toString() + 'px'; 
                      element.style.left = (rect.left + shift).toString() + 'px';
                  }
              } else {
                  var element = document.getElementById('overlay');
                  if(element) element.style.display = 'none';
              }
          };
  
          const selectInput = () => {
              const ipt = document.getElementById('inputRef');
              if(ipt) {
                  const length = ipt.value.length;
                  ipt.focus();
                  ipt.setSelectionRange(length, length);
              }
          };
  
          const onMessageImageSelectEmoji = (emoji) => {
              var xtext = document.getElementById("xtext"); 
              if(!xtext) return;
  
              const currentCursorPosition = xtext.selectionStart;
              const currentValue = xtext.value;
  
              const newValue = currentValue.slice(0, currentCursorPosition) + emoji.i + currentValue.slice(currentCursorPosition);
              xtext.value = newValue;
              xtext.focus();
              
              setTimeout(() => {
                  const newCursorPosition = currentCursorPosition + emoji.i.length;
                  xtext.setSelectionRange(newCursorPosition, newCursorPosition);
              }, 0);
          };
  
          const onMessageSelectEmoji = (emoji) => {
              if (showImageMessage.value) {
                  onMessageImageSelectEmoji(emoji);
                  return;
              } 
              
              var input = document.getElementById("inputRef");
              if(!input) return;
              
              const currentCursorPosition = input.selectionStart;
              const currentValue = input.value;
  
              const newValue = currentValue.slice(0, currentCursorPosition) + emoji.i + currentValue.slice(currentCursorPosition);
  
              input.value = newValue;
              message.value = input.value; // Sync v-model manuale
  
              input.focus();
              setTimeout(() => {
                  const newCursorPosition = currentCursorPosition + emoji.i.length;
                  input.setSelectionRange(newCursorPosition, newCursorPosition);
              }, 0);  
          };
  
          return {
              message,
              isSMPickerOpen,
              isEnabled,
              showImageMessage,
              image,
              toggleEnable,
              toggleImageMessage,
              imageMessageCreate,
              toggleEmoji,
              selectInput,
              sendmessage,
              onMessageSelectEmoji
          };
      }
  }
  </script>
  
  <template>
      <div class="chatbox-input" v-if="conversation">
          <i id="grinemoji" class="grin" @click="toggleEmoji">
              <div id="overlay" @click.stop>
                  <emojipicker 
                      v-if="isSMPickerOpen"
                      :native="true"
                      :hide-search="false"
                      :disable-skin-tones="true"
                      :hide-group-names="true"
                      :disabled-groups='["food_drink", "activities", "travel_places",  "objects",  "symbols",  "flags"]' 
                      @select="onMessageSelectEmoji"
                      class="picker-custom"
                  />
              </div>
          </i>
  
          <i class="paperclip" @click="toggleImageMessage"></i>
          
          <imagemessage 
              v-if="showImageMessage"
              @add-image-cancel="toggleImageMessage"
              @add-image="imageMessageCreate"
          />
  
          <input id="inputRef" type="text" placeholder="Type a message..." 
              v-model="message" 
              @keyup.enter="sendmessage">
          
          <i class="paperplane" @click="sendmessage"></i>
      </div>
  </template>
  
  <style scoped>
  /* --- STILI BLU NOTTE / DARK MODE --- */
  
  /* Contenitore input principale (Footer Chat) */
  .chatbox-input {
      position: relative;
      width: 100%;
      height: 70px;
      background-color: #1e293b; /* Grigio Bluastro */
      border-top: 1px solid #334155;
      display: grid;
      grid-template-columns: 50px 50px 1fr 60px; /* Layout griglia preciso */
      align-items: center;
      justify-content: center;
      padding: 0 10px;
      box-shadow: 0 -4px 6px -1px rgba(0, 0, 0, 0.1);
      z-index: 10;
  }
  
  /* Icone generali */
  .chatbox-input i {
      cursor: pointer;
      font-size: 1.5em;
      color: #94a3b8; /* Colore icone inattive */
      width: 40px;
      height: 40px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 50%;
      transition: all 0.2s ease;
  }
  
  .chatbox-input i:hover {
      background-color: rgba(255, 255, 255, 0.05);
      color: #2dd4bf; /* Verde Acqua al passaggio */
  }
  
  /* Overlay Picker */
  #overlay {
      display: none; /* Gestito da JS */
      position: fixed; 
      z-index: 2000; 
      width: 360px;
      /* Ombra profonda per il popup che galleggia */
      filter: drop-shadow(0 10px 8px rgb(0 0 0 / 0.04)) drop-shadow(0 4px 3px rgb(0 0 0 / 0.1));
  }
  
  /* Personalizzazione picker (se supportata dalla libreria tramite classi o deep) */
  :deep(.v3-emoji-picker) {
      background: #0f172a !important;
      border: 1px solid #334155 !important;
      --ep-color-bg: #0f172a;
      --ep-color-text: #f1f5f9;
      --ep-color-border: #334155;
      --ep-color-hover: #1e293b;
  }
  
  /* Input Text */
  .chatbox-input input {
      position: relative;
      width: 100%;
      padding: 12px 20px;
      border-radius: 24px;
      font-size: 0.95em;
      border: 1px solid #334155;
      outline: none;
      background-color: #0f172a; /* Blu Notte profondo */
      color: #f1f5f9; /* Testo chiaro */
      transition: border-color 0.2s;
  }
  
  .chatbox-input input::placeholder {
      color: #64748b;
  }
  
  .chatbox-input input:focus {
      border-color: #2dd4bf; /* Focus Verde Acqua */
  }
  
  .chatbox-input input:disabled {
      background-color: #1e293b;
      color: #475569;
      cursor: not-allowed;
  }
  
  /* Specifiche Icone SVG tramite pseudo-elementi */
  /* Filtro CSS per colorare le SVG (assumendo siano nere di default) */
  .grin::after, .paperclip::after, .paperplane::after {
      display: inline-block;
      width: 24px;
      height: 24px;
      vertical-align: middle;
      /* Inverte colori per dark mode e aggiusta luminosità */
      filter: invert(63%) sepia(11%) saturate(406%) hue-rotate(174deg) brightness(93%) contrast(89%); 
  }
  
  .grin:hover::after, .paperclip:hover::after, .paperplane:hover::after {
      /* Verde Acqua #2dd4bf filter generator approx */
      filter: invert(73%) sepia(35%) saturate(827%) hue-rotate(124deg) brightness(93%) contrast(88%);
  }
  
  .grin::after {
      content: url('/emoji.svg');
  }
  .paperclip::after {
      content: url('/attach.svg');
  }
  .paperplane::after {
      content: url('/send_message.svg');
      margin-left: 5px; /* Centratura visiva */
  }
  
</style>