<script>
  import { ref, onMounted, getCurrentInstance } from 'vue';
  import { useRouter } from 'vue-router';
  
  export default {
      name: "LoginForm",
      // Template non più necessario come stringa esterna, uso il template del SFC
      setup() {
          const { proxy } = getCurrentInstance();
          const $axios = proxy.$axios;
          const router = useRouter();
  
          // --- State ---
          const username = ref("");
  
          // --- Methods ---
  
          const login = async () => {
              sessionStorage.removeItem('token');
              try {
                  let response = await $axios.post("/session", {
                      name: username.value,
                  });
  
                  if (response.status == 201) {
                      sessionStorage.setItem('token', response.data.identifier);
                      sessionStorage.setItem('login', username.value);
                      router.push('/');
                  }
              } catch (e) {
                  console.error(e.toString());
              }
          };
  
          const register = () => {
              alert("Coming soon ...");
          };
  
          // --- Lifecycle ---
          onMounted(() => {
              let init = 1000;
              setTimeout(function() {
                  // Logica originale per l'animazione di apertura
                  const wrapper = document.querySelector(".login-wrapper");
                  if(wrapper) wrapper.classList.toggle("open");
                  init = 300;
              }, init);
          });
  
          return {
              username,
              login,
              register
          };
      }
  };
  </script>
  
  <template>
      <div id="sandbox">
          <div class="login-wrapper">
              <div class="login-left">
                  <img height="300" src="/wasatext_logo.png" alt="Logo">
              </div>
              <form class="login-right" @submit.prevent>
                  <div class="h2">Welcome Back</div>
                  <div class="form-group">
                      <input type="text" id="username" placeholder=" " v-model="username">
                      <label for="username">Username</label>
                  </div>
                  <div class="button-area">
                      <button class="btn-login" @click="login()">Login</button>
                  </div>
              </form>
          </div>
      </div>
  </template>
  
  <style scoped>
  /* --- STILI DARK MODE / BLU NOTTE --- */
  
  * {
      box-sizing: border-box;
  }
  
  /* Sfondo pagina intera */
  #sandbox {
      font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 100vh;
      /* Gradiente scuro elegante */
      background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
      overflow: hidden;
  }
  
  /* --- CARD LOGIN --- */
  .login-wrapper {
      width: 800px;
      height: 440px;
      background-color: #1e293b; /* Card Bluastro Scuro */
      border-radius: 16px;
      box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5); /* Ombra profonda */
      overflow: hidden;
      position: relative;
      border: 1px solid #334155;
  }
  
  /* --- LATO SINISTRO (Logo/Animazione) --- */
  .login-left {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center; /* Centrato */
      transition: all 770ms cubic-bezier(0.51, 0.04, 0.12, 0.99);
      overflow: hidden;
      background-color: #0f172a; /* Più scuro per contrasto */
      position: relative;
      z-index: 2;
  }
  
  .login-left img {
      display: block;
      transition: all 770ms cubic-bezier(0.51, 0.04, 0.12, 0.99);
      /* Filtro per adattare loghi neri al tema scuro, se necessario. Rimuovere se il logo è già chiaro */
      /* filter: invert(1) brightness(2);  */
  }
  
  /* Animazione Open */
  .open .login-left {
      transform: translateX(-50%); /* Sposta a sinistra */
  }
  
  .open .login-left img {
      transform: translateX(-50px) scale(0.8); /* Effetto parallasse */
      opacity: 0.1; /* Logo svanisce leggermente */
  }
  
  /* --- LATO DESTRO (Form) --- */
  .login-right {
      padding: 60px 40px;
      position: absolute;
      top: 0;
      right: 0;
      width: 50%; /* Occupa metà card */
      height: 100%;
      transform: translateX(100%); /* Nascosto inizialmente a destra */
      transition: all 770ms cubic-bezier(0.51, 0.04, 0.12, 0.99);
      display: flex;
      flex-direction: column;
      justify-content: center;
      background-color: #1e293b;
      z-index: 1;
  }
  
  .open .login-right {
      transform: translateX(0); /* Entra in scena */
  }
  
  .h2 {
      color: #2dd4bf; /* Verde Acqua */
      font-size: 2rem;
      font-weight: 700;
      margin-bottom: 30px;
      text-align: center;
  }
  
  /* --- INPUT FIELDS --- */
  .form-group {
      position: relative;
      margin-bottom: 25px;
      width: 100%;
  }
  
  input {
      width: 100%;
      height: 50px;
      background-color: #0f172a;
      border: 1px solid #334155;
      border-radius: 8px;
      padding: 0 15px;
      color: #f1f5f9;
      font-size: 1rem;
      outline: none;
      transition: border-color 0.3s;
  }
  
  input:focus {
      border-color: #2dd4bf;
  }
  
  /* Label Flottante */
  label {
      position: absolute;
      top: 50%;
      left: 15px;
      transform: translateY(-50%);
      color: #64748b;
      font-size: 1rem;
      pointer-events: none;
      transition: all 0.3s ease;
      background-color: transparent; 
      padding: 0 4px;
  }
  
  /* Quando l'input ha focus o non è vuoto (usa placeholder hack) */
  input:focus + label,
  input:not(:placeholder-shown) + label {
      top: 0;
      font-size: 0.8rem;
      color: #2dd4bf;
      background-color: #1e293b; /* Copre il bordo per leggibilità */
  }
  
  /* Button Area */
  .button-area {
      margin-top: 10px;
      display: flex;
      justify-content: center;
  }
  
  .btn-login {
      width: 100%;
      height: 45px;
      border-radius: 25px;
      border: none;
      background: linear-gradient(135deg, #0d9488 0%, #2dd4bf 100%);
      color: #0f172a;
      font-weight: 700;
      font-size: 1rem;
      cursor: pointer;
      box-shadow: 0 4px 15px rgba(45, 212, 191, 0.3);
      transition: transform 0.2s, box-shadow 0.2s;
  }
  
  .btn-login:hover {
      transform: translateY(-2px);
      box-shadow: 0 6px 20px rgba(45, 212, 191, 0.5);
  }
  
  .btn-login:active {
      transform: translateY(0);
  }
  
  </style>