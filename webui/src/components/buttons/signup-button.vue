<template>
    <button class="button__sign-up" @click="handleSignUp">Sign Up</button>
  </template>
  
  <script setup>

/** avvia la registrazione forzando la schermata di signup e reindirizza l'utente su /profile */
  
  const handleSignUp = () => {
    loginWithRedirect({
      appState: {
        target: "/profile",
      },
      authorizationParams: {
        screen_hint: "signup",
      }
    });
  };
  </script>