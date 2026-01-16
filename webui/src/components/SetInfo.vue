<script>
    import { ref, getCurrentInstance, watch } from 'vue';
    
    export default {
        props: ['itype', 'ievent'],
        emits: ['set-info', 'close-info'],
        setup(props, { emit }) {
            const { proxy } = getCurrentInstance();
            const $axios = proxy.$axios;
            const name = ref(null);
            const previewImage = ref(null);
    
            watch(() => props.ievent, () => { previewImage.value = null; });
    
            const chooseFiles = () => { document.getElementById("xmedia").click(); };
    
            const errorName = (status) => {
                if (status == 406 || status == 409 || status == 400) {
                    alert("Attenzione: Nome non valido o già in uso.");
                }
            };
    
            const setName = async () => {
                name.value = document.getElementById("iname").value;
                let config = { headers: { Authorization: "Bearer " + sessionStorage.getItem('token') } };
                try {
                    let url = props.itype == 'group' 
                        ? `/conversations/byid/${props.ievent.conid}/name?name=${name.value}`
                        : "/info/username";
                    
                    await $axios.put(url, name.value, config);
                    emit('set-info', props.itype, props.ievent.conid);
                } catch (error) {
                    if (error.response) errorName(error.response.status);
                }
            };
    
            // --- FUNZIONE UPLOAD FLESSIBILE ---
            const setPhoto = async (e) => {
                var files = e.target.files || e.dataTransfer.files;
                if (!files.length) return;
    
                const file = files[0];
    
                // 1. Validazione Frontend: Accetta GIF, PNG, JPEG
                const validTypes = ['image/gif', 'image/png', 'image/jpeg'];
                if (!validTypes.includes(file.type)) {
                    alert("Formato non supportato. Usa GIF, PNG o JPG.");
                    return;
                }
    
                // 2. Anteprima Immediata
                const reader = new FileReader();
                reader.onload = (ev) => { previewImage.value = ev.target.result; };
                reader.readAsDataURL(file);
    
                // 3. Header Dinamico (Dice al server il tipo reale del file)
                let config = {
                    headers: {
                        Authorization: "Bearer " + sessionStorage.getItem('token'),
                        'Content-Type': file.type 
                    }
                };
    
                try {
                    let url = props.itype == 'group' 
                        ? `/conversations/byid/${props.ievent.conid}/photo`
                        : "/info/userphoto";
                        
                    await $axios.put(url, file, config);
                    emit('set-info', props.itype, props.ievent.conid);
                } catch (e) {
                    console.error(e);
                    if (e.response && (e.response.status === 400 || e.response.status === 406)) {
                        alert(`Errore Backend: Il server ha rifiutato il file ${file.type}. Assicurati di aver aggiornato e ricompilato il Backend.`);
                    } else {
                        alert("Errore caricamento: " + e.message);
                    }
                    previewImage.value = null;
                }
            };
    
            const close = () => { emit('close-info'); };
    
            return { name, previewImage, chooseFiles, setName, setPhoto, close };
        }
    }
    </script>
    
    <template>
        <div class="icontainer">
            <div class="header-spacer"></div>
            <h5>Impostazioni Profilo</h5>
            <div class="header-spacer"></div>
        
            <div class="infoname">Nome:</div>
            <div class="input-wrapper">
                <input maxlength="32" type="text" id="iname" class="input-custom" :value="ievent && ievent.name"/>
            </div> 
            <div class="action">
                <button type="button" class="btn-custom btn-save" @click="setName()">Salva</button>
            </div>
    
            <div class="infoname">Foto:</div>  
            <div id="xPhoto" class="camera Photo">
                <div class="img-wrapper">
                    <img id="pimg" v-if="previewImage" :src="previewImage" />
                    <img id="pimg" v-else-if="ievent && ievent.photo" 
                         :src="'data:' + (ievent.photo.startsWith('/9j/') ? 'image/jpeg' : (ievent.photo.startsWith('iVBOR') ? 'image/png' : 'image/gif')) + ';base64,' + ievent.photo"/>
                    <img id="pimg" v-else src="/standard.gif"/>
                </div>
            </div> 
            
            <div class="file-container action action-col">
                <input type="file" id="xmedia" accept="image/png, image/gif, image/jpeg" hidden @change="(event) => setPhoto(event)"/>
                <div class="btn-custom marginal" @click="chooseFiles()" id="upload">carica foto</div>
            </div>
    
            <div></div>
            <div><div class="btn-custom semimarginal" @click="close()" id="closeinfo">chiudi</div></div>
            <div></div>
        </div>    
    </template>
    
    <style scoped>
    /* Stili originali invariati */
    .icontainer { display: grid; grid-template-columns: 20% 50% 30%; position: relative; top: 50px; background-color: #1e293b; border: 1px solid #334155; border-radius: 12px; padding: 20px; box-shadow: 0 10px 25px rgba(0, 0, 0, 0.5); color: #f1f5f9; width: fit-content; margin: 0 auto; align-items: center; gap: 10px; }
    .icontainer h5 { color: #2dd4bf; font-weight: 700; grid-column: 2; text-align: center; margin-bottom: 20px; }
    .infoname { text-align: right; font-weight: 600; color: #94a3b8; padding-right: 10px; }
    .input-wrapper { width: 100%; }
    .input-custom { background-color: #0f172a; width: 100%; border: 1px solid #334155; color: #f1f5f9; padding: 8px 12px; border-radius: 6px; outline: none; transition: border-color 0.2s; }
    .input-custom:focus { border-color: #2dd4bf; }
    .btn-custom { background-color: transparent; border: 1px solid #475569; color: #94a3b8; border-radius: 20px; padding: 5px 15px; font-size: 0.85rem; cursor: pointer; transition: all 0.2s ease; text-align: center; display: inline-block; }
    .btn-custom:hover { border-color: #cbd5e1; color: #fff; background-color: rgba(255, 255, 255, 0.05); }
    .btn-save { border-color: #0d9488; color: #2dd4bf; }
    .btn-save:hover { background-color: rgba(45, 212, 191, 0.1); border-color: #2dd4bf; color: #fff; }
    .camera { width: 120px; height: 120px; border-radius: 50%; overflow: hidden; margin: 10px auto; border: 3px solid #334155; background-color: #0f172a; display: flex; justify-content: center; align-items: center; }
    .Photo { grid-column: 2; }
    .img-wrapper { width: 100%; height: 100%; display: flex; justify-content: center; align-items: center; }
    .camera img { width: 100%; height: 100%; object-fit: cover; }
    .action { display: flex; align-items: center; }
    .action-col { flex-direction: column; justify-content: center; }
    .semimarginal { margin-top: 20px; width: 100%; } 
    .marginal { margin-top: 0; white-space: nowrap; }
    </style>