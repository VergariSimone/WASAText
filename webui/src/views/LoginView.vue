<template>
    <div class="login-wrapper" :class="{ open: isOpen }">
      <div class="login-left">
        <img src="/wasatext_logo.svg" alt="Logo" height="300"/>
      </div>
  
      <form class="login-right" @submit.prevent="login">
        <h2>Login</h2>
  
        <div class="form-group">
          <input type="text" placeholder="Username" v-model="username" required />
        </div>
  
        <div class="button-area">
          <button class="btn btn-primary" type="submit">Login</button>
        </div>
      </form>
    </div>
  </template>
  
  <script>
  import { ref, onMounted } from "vue";
  
  export default {
    name: "LoginView",
    setup(_, { emit, router }) {
      const username = ref("");
      const isOpen = ref(false);
  
      onMounted(() => {
        // Animazione iniziale
        setTimeout(() => {
          isOpen.value = true;
        }, 500);
      });
  
      const login = async () => {
        sessionStorage.removeItem("token");
        try {
          const response = await fetch("/session", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ name: username.value }),
          });
          if (response.status === 201) {
            const data = await response.json();
            sessionStorage.setItem("token", data.identifier);
            sessionStorage.setItem("login", username.value);
            window.location.href = "/";
          } else {
            console.error("Login failed", response.status);
          }
        } catch (e) {
          console.error("Errore login:", e);
        }
      };
  
      return { username, isOpen, login };
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
    background-color: #f3f3f3;
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
    background-color: #22bf64;
    color: #fff;
    border: none;
    cursor: pointer;
    border-radius: 4px;
  }
  
  button:hover {
    background-color: #1fa85c;
  }
  </style>
  