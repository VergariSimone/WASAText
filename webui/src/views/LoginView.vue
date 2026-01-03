<script>
  import { ref, onMounted, getCurrentInstance } from 'vue';
  import { useRouter } from 'vue-router';
  
  export default {
      name: "LoginForm",
      setup() {
          const { proxy } = getCurrentInstance();
          const $axios = proxy.$axios;
          const router = useRouter();
  
          // state
          const username = ref("");
  
          // metodi
  
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
              alert("");
          };
  
          
          onMounted(() => {
              let init = 1000;
              setTimeout(function() {
                  // animazione di apertura
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
                  <div class="h2">Bentornato!</div>
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
  * {
      box-sizing: border-box;
  }
  
  /* backgrpund pagina intera */
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
  
  /* card per login*/
  .login-wrapper {
      width: 800px;
      height: 440px;
      background-color: #1e293b; 
      border-radius: 16px;
      box-shadow: 0 20px 50px rgba(0, 0, 0, 0.5); 
      overflow: hidden;
      position: relative;
      border: 1px solid #334155;
  }
  
  /* lato sx del login*/
  .login-left {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center; 
      transition: all 770ms cubic-bezier(0.51, 0.04, 0.12, 0.99);
      overflow: hidden;
      background-color: #0f172a; 
      position: relative;
      z-index: 2;
  }
  
  /* lato sx login immagine */
  .login-left img {
      display: block;
      transition: all 770ms cubic-bezier(0.51, 0.04, 0.12, 0.99);
  }
  
  /* animazzione */
  .open .login-left {
      transform: translateX(-50%); /* verso sx*/
  }
  
  .open .login-left img {
      transform: translateX(-50px) scale(0.8); /* effettp parallasse */
      opacity: 0.1; 
  }
  
  /* lato dx login form  */
  .login-right {
      padding: 60px 40px;
      position: absolute;
      top: 0;
      right: 0;
      width: 50%; 
      height: 100%;
      transform: translateX(100%); /* nascosto inizialmente a dx */
      transition: all 770ms cubic-bezier(0.51, 0.04, 0.12, 0.99);
      display: flex;
      flex-direction: column;
      justify-content: center;
      background-color: #1e293b;
      z-index: 1;
  }
  
  .open .login-right {
      transform: translateX(0); /* entra nella finestra */
  }
  
  .h2 {
      color: #2dd4bf; 
      font-size: 2rem;
      font-weight: 700;
      margin-bottom: 30px;
      text-align: center;
  }
  
  /* campi di imput  */
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
  
  /* label fluttuante  */
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
  
  /* input ha focus o non è vuoto (u */
  input:focus + label,
  input:not(:placeholder-shown) + label {
      top: 0;
      font-size: 0.8rem;
      color: #2dd4bf;
      background-color: #1e293b; 
  }
  
  /*area bottoni */
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