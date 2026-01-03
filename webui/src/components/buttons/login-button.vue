<template>
    <button class="button__login" @click="handleLogin">Log In</button>
  </template>
  
  <script setup>

    /**
     * gestisce click del pulsante login ,la pagina '/profile'è la destinazione dopo l'accesso riuscito
     */
  
  const handleLogin = () => {
    loginWithRedirect({
      appState: {
        target: "/profile",
      },
    });
  };
  </script>