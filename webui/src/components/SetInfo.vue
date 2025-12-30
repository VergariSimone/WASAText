<script>
  import { ref, getCurrentInstance } from 'vue';
  
  export default {
      props: [
          'itype',
          'ievent',
      ],
      emits: [
          'set-info',
          'close-info',
      ],
      setup(props, { emit }) {
          // Accesso all'istanza per usare $axios
          const { proxy } = getCurrentInstance();
          const $axios = proxy.$axios;
  
          // --- State ---
          const name = ref(null);
  
          // --- Methods ---
  
          const chooseFiles = () => {
              document.getElementById("xmedia").click();
          };
  
          const errorName = (status) => {
              if (status == 406) {
                  alert("Username not acceptable or already used");
                  return;
              }
          };
  
          const setName = async () => {
              
              name.value = document.getElementById("iname").value;
              
              let config = {
                  headers: {
                      Authorization: "Bearer " + sessionStorage.getItem('token'),
                  }
              };
  
              try {
                  if (props.itype == 'group') {
                      await $axios.put(`/conversations/byid/${props.ievent.conid}/name?name=${name.value}`, name.value, config);
                      emit('set-info', props.itype, props.ievent.conid);
                  } else {
                      await $axios.put("/info/username", name.value, config);
                      emit('set-info', props.itype, props.ievent.conid);
                  }
              } catch (error) {
                  if (error.response) {
                      errorName(error.response.status);
                  }
              }
          };
  
          const setPhoto = async (e) => {
              var files = e.target.files || e.dataTransfer.files;
              if (!files.length) return;
  
              let config = {
                  headers: {
                      Authorization: "Bearer " + sessionStorage.getItem('token'),
                      'content-type': 'image/gif'
                  }
              };
  
              try {
                  if (props.itype == 'group') {
                      await $axios.put(`/conversations/byid/${props.ievent.conid}/photo`, files[0], config);
                      emit('set-info', props.itype, props.ievent.conid);
                  } else {
                      await $axios.put("/info/userphoto", files[0], config);
                      emit('set-info', props.itype, props.ievent.conid);
                  }
              } catch (e) {
                  console.error(e);
              }
          };
  
          const close = () => {
              emit('close-info');
          };
  
          return {
              name,
              chooseFiles,
              setName,
              setPhoto,
              close
          };
      }
  }
  </script>
  
  <template>
      <div class="icontainer">
          <div class="header-spacer"></div>
          <h5>View/Set Information</h5>
          <div class="header-spacer"></div>
      
          <div class="infoname">Name:</div>
          <div class="input-wrapper">
              <input 
                  maxlength="32"
                  type="text" 
                  id="iname" 
                  class="input-custom"
                  :value="ievent && ievent.name"
              />
          </div> 
          <div class="action">
              <button type="button" class="btn-custom btn-save" @click="setName()">Save</button>
          </div>
  
          <div class="infoname">Photo:</div>  
          <div id="xPhoto" class="camera Photo">
              <div class="img-wrapper">
                  <img id="pimg" v-if="ievent && ievent.photo == ''" src="/standard2.gif"/>
                  <img id="pimg" v-if="ievent && ievent.photo != ''" v-bind:src="'data:image/gif;base64,' + ievent.photo"/>
              </div>
          </div> 
          
          <div class="file-container action action-col">
              <input
                  type="file"
                  id="xmedia"
                  accept="image/gif"
                  hidden
                  @change="(event) => setPhoto(event)"
              />
              <div class="btn-custom marginal" @click="chooseFiles()" id="upload">Upload image</div>
          </div>
  
          <div></div>
          <div>
              <div class="btn-custom semimarginal" @click="close()" id="closeinfo">Close</div>
          </div>
          <div></div>
      </div>    
  </template>
  
  <style scoped>
  /* --- TEMA DARK MODE / BLU NOTTE --- */
  
  .icontainer {
      display: grid;
      grid-template-columns: 20% 50% 30%; /* Layout a griglia aggiustato */
      position: relative;
      top: 50px;
      
      background-color: #1e293b; /* Grigio Bluastro */
      border: 1px solid #334155;
      border-radius: 12px;
      padding: 20px;
      box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5);
      
      color: #f1f5f9;
      width: fit-content;
      margin: 0 auto; /* Centratura */
      align-items: center;
      gap: 10px;
  }
  
  .icontainer h5 {
      color: #2dd4bf; /* Verde Acqua */
      font-weight: 700;
      grid-column: 2;
      text-align: center;
      margin-bottom: 20px;
  }
  
  .infoname {
      text-align: right;
      font-weight: 600;
      color: #94a3b8; /* Grigio chiaro */
      padding-right: 10px;
  }
  
  /* Input Styles */
  .input-wrapper {
      width: 100%;
  }
  
  .input-custom {
      background-color: #0f172a; /* Blu Notte profondo */
      width: 100%;
      border: 1px solid #334155;
      color: #f1f5f9;
      padding: 8px 12px;
      border-radius: 6px;
      outline: none;
      transition: border-color 0.2s;
  }
  
  .input-custom:focus {
      border-color: #2dd4bf;
  }
  
  /* Button Styles */
  .btn-custom {
      background-color: transparent;
      border: 1px solid #475569;
      color: #94a3b8;
      border-radius: 20px;
      padding: 5px 15px;
      font-size: 0.85rem;
      cursor: pointer;
      transition: all 0.2s ease;
      text-align: center;
      display: inline-block;
  }
  
  .btn-custom:hover {
      border-color: #cbd5e1;
      color: #fff;
      background-color: rgba(255, 255, 255, 0.05);
  }
  
  .btn-save {
      border-color: #0d9488;
      color: #2dd4bf;
  }
  
  .btn-save:hover {
      background-color: rgba(45, 212, 191, 0.1);
      border-color: #2dd4bf;
      color: #fff;
  }
  
  /* Photo Section */
  .camera {
      width: 120px;
      height: 120px;
      border-radius: 50%;
      overflow: hidden;
      margin: 10px auto;
      border: 3px solid #334155;
      background-color: #0f172a;
      display: flex;
      justify-content: center;
      align-items: center;
  }
  
  .Photo {
      grid-column: 2; /* Centrato nella griglia */
  }
  
  .img-wrapper {
      width: 100%;
      height: 100%;
      display: flex;
      justify-content: center;
      align-items: center;
  }
  
  .camera img {
      width: 100%;
      height: 100%;
      object-fit: cover;
  }
  
  /* Layout Utilities */
  .action {
      display: flex;
      align-items: center;
  }
  
  .action-col {
      flex-direction: column;
      justify-content: center;
  }
  
  .semimarginal {
      margin-top: 20px;
      width: 100%;
  } 
  
  .marginal {
      margin-top: 0; /* Resettato per flexbox alignment */
      white-space: nowrap;
  }
  
  </style>