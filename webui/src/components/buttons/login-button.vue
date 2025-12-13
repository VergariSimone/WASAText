<template>
    <button class="button__login" @click="handleLogin">Log In</button>
  </template>
  
  <script setup>
  
  const handleLogin = () => {
    loginWithRedirect({
      appState: {
        target: "/profile",
      },
    });
  };
  </script>