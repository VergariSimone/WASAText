<template>
    <button class="button__sign-up" @click="handleSignUp">Sign Up</button>
  </template>
  
  <script setup>

  
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