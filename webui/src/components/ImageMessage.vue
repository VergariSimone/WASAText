<script>
  import { ref } from 'vue';
  
  export default {
      // eventi emessi
      emits: [
          'add-image',
          'add-image-cancel',
      ],
      
      setup(props, { emit }) {
          
          // state
          const image = ref(null);
          const filled = ref(false);
          const text = ref(null);
  
          // metodi

          // scegli file
          const chooseFiles = () => {
              document.getElementById("xmedia").click();
          };
  
          // imposta foto
          const setPhoto = (e) => {
              var files = e.target.files || e.dataTransfer.files;
              if (!files.length) return;
  
              if (files[0]) {
                  image.value = files[0];
                  var reader = new FileReader();
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
              
              var element = document.getElementById("xtext");
              text.value = element.value;
  
              if (image.value || (text.value != '')) {
                  emit('add-image', image.value, text.value);
              } else {
                  emit('add-image-cancel');
              }
          };
  

          // cancella immagine 
          const cancel = () => {
              emit('add-image-cancel');
          };
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
              <button type="button" @click="cancel" class="btn-custom">Annulla</button>
              <i></i>
              <button type="button" @click="chooseFiles()" id="upload" class="btn-custom">Carica</button>
              <i></i>
              <button type="button" @click="ok" class="btn-custom btn-send">Invia</button>
              <i></i>
              <i class="allspan">&nbsp;</i>
          </div>
          <div class="listcontext">
              <div id="xPhoto" class="camera">
                  <div class="camera-inner">
                      <img id="pimg" class="camera-img" :class="{cornice:filled}"/>
                      <span v-if="!filled" class="placeholder-text">Anteprima immagine </span>
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
          <input id="xtext" class="textimage" type="text" placeholder="Scrivi il tuo messaggio..."/>
      </div>
  </template>
  
  <style scoped>
  
  /* div principale */
  #listimage {
      align-items: center;
      background-color: #1e293b; 
      border: 1px solid #334155;
      border-radius: 12px;
      padding: 15px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5); 
      position: fixed;
      width: 340px;
      height: fit-content;
      z-index: 1001; /* sopra tutto */
      top: 200px;
      left: 50%;
      transform: translateX(-50%); /* centrato orizzontalmente */
  }
  
  /* input text*/
  .textimage {
      border-radius: 8px;
      background-color: #0f172a; 
      border: 1px solid #334155;
      color: #f1f5f9; 
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
      border-color: #2dd4bf; 
  }
  
  .textimage::placeholder {
      color: #64748b;
  }
  
  /* griglia pulsanti*/
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
  
  /* pulsanti */
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
      color: #2dd4bf; 
      font-weight: 600;
  }
  
  .btn-send:hover {
      background-color: rgba(45, 212, 191, 0.1);
      border-color: #2dd4bf;
      color: #fff;
      box-shadow: 0 0 10px rgba(45, 212, 191, 0.3);
  }
  
  /* area anteprima imamgine  */
  .listcontext {
      display: flex;
      justify-content: center;
      align-items: center;
      height: 220px;
      margin-bottom: 15px;
      border-radius: 8px;
      background-color: #0f172a;
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
  
  /* placeholder quando non c'è immagine */
  .placeholder-text {
      color: #475569;
      font-size: 0.8rem;
      font-style: italic;
  }
  
  /* cornice quando l'immagine è caricata */
  .cornice {
      border: 2px solid #2dd4bf; /* Verde Acqua */
      box-shadow: 0 0 15px rgba(45, 212, 191, 0.2);
  }
  
  i {
      display: block;
  }
</style>