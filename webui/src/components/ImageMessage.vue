<script>
  import { ref } from 'vue';
  
  export default {
      // Definizione degli eventi emessi
      emits: [
          'add-image',
          'add-image-cancel',
      ],
      // Funzione setup esplicita (senza <script setup>)
      setup(props, { emit }) {
          
          // --- State Reattivo ---
          const image = ref(null);
          const filled = ref(false);
          const text = ref(null);
  
          // --- Methods ---
          
          const chooseFiles = () => {
              // Accesso DOM diretto come richiesto
              document.getElementById("xmedia").click();
          };
  
          const setPhoto = (e) => {
              var files = e.target.files || e.dataTransfer.files;
              if (!files.length) return;
  
              if (files[0]) {
                  image.value = files[0];
                  var reader = new FileReader();
                  
                  // Manipolazione diretta DOM
                  var element = document.getElementById("pimg");
                  
                  reader.readAsDataURL(files[0]);
                  filled.value = true;
                  
                  reader.onload = function (evt) {
                      element.src = reader.result;
                  }
                  reader.onerror = function (evt) {
                      alert("Errore caricamento immagine");
                  }
              }
          };
  
          const ok = () => {
              // Recupero valore input text via DOM
              var element = document.getElementById("xtext");
              text.value = element.value;
  
              if (image.value || (text.value != '')) {
                  emit('add-image', image.value, text.value);
              } else {
                  emit('add-image-cancel');
              }
          };
  
          const cancel = () => {
              emit('add-image-cancel');
          };
  
          // --- Return ---
          // Tutto ciò che deve essere usato nel template deve essere restituito qui
          return {
              image,
              filled,
              text,
              chooseFiles,
              setPhoto,
              ok,
              cancel
          };
      }
  }
  </script>
  
  <template>
      <div id="listimage">
          <div class="listbuttons">
              <i class="allspan">&nbsp;</i>
              <i></i>
              <button type="button" @click="cancel" class="btn-custom">Cancel</button>
              <i></i>
              <button type="button" @click="chooseFiles()" id="upload" class="btn-custom">Upload</button>
              <i></i>
              <button type="button" @click="ok" class="btn-custom btn-send">Send</button>
              <i></i>
              <i class="allspan">&nbsp;</i>
          </div>
          <div class="listcontext">
              <div id="xPhoto" class="camera">
                  <div class="camera-inner">
                      <img id="pimg" class="camera-img" :class="{cornice:filled}"/>
                      <span v-if="!filled" class="placeholder-text">Preview Image</span>
                  </div>
              </div>
              <input
                  type="file"
                  id="xmedia"
                  accept="image/*"
                  hidden
                  @change="(event) => setPhoto(event)" />
              <br/>
          </div>
          <input id="xtext" class="textimage" type="text" placeholder="Type text here..."/>
      </div>
  </template>
  
  <style scoped>
  /* --- CONFIGURAZIONE COLORI: Blu Notte / Verde Acqua / Bianco --- */
  
  /* --- CONTENITORE PRINCIPALE --- */
  #listimage {
      align-items: center;
      background-color: #1e293b; /* Grigio Bluastro (Slate) */
      border: 1px solid #334155;
      border-radius: 12px;
      padding: 15px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5); /* Ombra profonda */
      position: fixed;
      width: 340px;
      height: fit-content;
      z-index: 1001; /* Sopra tutto */
      top: 200px;
      left: 50%;
      transform: translateX(-50%); /* Centrato orizzontalmente */
  }
  
  /* --- INPUT DI TESTO --- */
  .textimage {
      border-radius: 8px;
      background-color: #0f172a; /* Blu Notte profondo */
      border: 1px solid #334155;
      color: #f1f5f9; /* Testo chiaro */
      margin-left: auto;
      margin-right: auto;
      display: block;
      margin-bottom: 10px;
      font-size: 0.9rem;
      height: 36px;
      width: 90%;
      padding: 0 10px;
      outline: none;
      transition: border-color 0.2s;
  }
  
  .textimage:focus {
      border-color: #2dd4bf; /* Focus Verde Acqua */
  }
  
  .textimage::placeholder {
      color: #64748b;
  }
  
  /* --- GRIGLIA PULSANTI --- */
  .listbuttons {
      display: grid;
      grid-template-columns: 5% 25% 5% 25% 5% 25% 5%;
      margin-bottom: 15px;
      align-items: center;
  }
  
  .allspan {
      grid-column: 1 / span 7;
      height: 5px;
  }
  
  /* --- STILE PULSANTI --- */
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
  
  .btn-send {
      border-color: #0d9488;
      color: #2dd4bf; /* Verde Acqua */
      font-weight: 600;
  }
  
  .btn-send:hover {
      background-color: rgba(45, 212, 191, 0.1);
      border-color: #2dd4bf;
      color: #fff;
      box-shadow: 0 0 10px rgba(45, 212, 191, 0.3);
  }
  
  /* --- AREA PREVIEW IMMAGINE --- */
  .listcontext {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 220px;
      margin-bottom: 15px;
      border-radius: 8px;
      background-color: #0f172a; /* Sfondo scuro interno */
      border: 1px dashed #334155;
      padding: 10px;
      position: relative;
  }
  
  .camera {
      width: 100%;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
  }
  
  .camera-inner {
      width: 100%;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
      overflow: hidden;
  }
  
  .camera-img {
      max-width: 100%;
      max-height: 100%;
      object-fit: contain;
      display: block;
      border-radius: 4px;
  }
  
  /* Placeholder quando non c'è immagine */
  .placeholder-text {
      color: #475569;
      font-size: 0.8rem;
      font-style: italic;
  }
  
  /* Cornice quando l'immagine è caricata */
  .cornice {
      border: 2px solid #2dd4bf; /* Verde Acqua */
      box-shadow: 0 0 15px rgba(45, 212, 191, 0.2);
  }
  
  i {
      display: block;
  }
</style>