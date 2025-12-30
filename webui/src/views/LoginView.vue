<template id="login-view">
  <div class="login-wrapper">
    <div class="login-left">
      <img src="/wasatext_logo.png" alt="Logo WASA TEXT" />
    </div>

    <form class="login-right" @submit.prevent>
      <h2>Login</h2>

      <div class="form-group">
        <input
          type="text"
          placeholder="Username"
          v-model="username"
        />
      </div>


      <div class="button-area">
        <button @click="login">
          Login
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import axios from "../services/axios";

export default {
  name: "LoginView",
  template: "#login-view",

  data() {
    return {
      username: "",
    };
  },

  beforeMount() {
    // animazione iniziale (DOM manipulation)
    setTimeout(() => {
      const wrapper = document.querySelector(".login-wrapper");
      if (wrapper) {
        wrapper.classList.add("open");
      }
    }, 1000);
  },

  methods: {
    async login() {
      sessionStorage.removeItem("token");

      try {
        const response = await axios.post("/session", {
          name: this.username,
        });

        if (response.status === 201) {
          sessionStorage.setItem("token", response.data.identifier);
          sessionStorage.setItem("login", this.username);
          this.$router.push("/");
        } else {
          console.error("Login fallito:", response.status);
        }
      } catch (err) {
        console.error("Errore login:", err);
      }
    },
  },
};
</script>

<style scoped>
* { box-sizing: border-box; }

html, body {
  margin: 0;
  padding: 0;
  height: 100%;
  font-family: "Lato", sans-serif;
  background: linear-gradient(243.87deg, #22bf64 30.6%, #371933 130.6%);
}

.login-wrapper {
  display: flex;
  width: 800px;
  height: 440px;
  margin: 50px auto;
  border-radius: 4px;
  overflow: hidden;
  background-color: #38BDF8;
  position: relative;
}

.login-left, .login-right {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.7s cubic-bezier(0.51, 0.04, 0.12, 0.99);
}

.login-left img {
  max-height: 300px;
  object-position: left;
  transition: all 0.7s cubic-bezier(0.51, 0.04, 0.12, 0.99);
}

.login-right {
  flex-direction: column;
  padding: 40px;
  position: absolute;
  top: 0;
  right: 0;
  width: 800px;
  transform: translateX(1400px);
}

.login-wrapper.open .login-left {
  transform: translateX(-800px);
}

.login-wrapper.open .login-left img {
  transform: translateX(210px);
}

.login-wrapper.open .login-right {
  transform: translateX(0);
}

h2 {
  margin-bottom: 20px;
  color: #371933;
}

.form-group {
  margin-bottom: 30px;
}

input[type="text"] {
  width: 100%;
  max-width: 400px;
  height: 40px;
  padding: 0 10px;
  font-size: 16px;
  border: none;
  border-bottom: 1px solid #aaa;
  outline: none;
}

.button-area {
  display: flex;
  justify-content: flex-end;
}

button {
  padding: 10px 20px;
  font-size: 16px;
  background-color: #0F172A;
  color: #fff;
  border: none;
  cursor: pointer;
  border-radius: 4px;
}

button:hover {
  background-color: #0F172A;
}
</style>
