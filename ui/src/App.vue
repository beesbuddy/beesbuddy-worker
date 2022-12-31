<template>
  <LoginPage></LoginPage>
</template>

<script>
// import HelloWorld from './components/HelloWorld.vue'
import LoginPage from './pages/auth/LoginPage.vue'

export default {
  name: 'App',
  components: {
    // HelloWorld,
    LoginPage,
  }
}
</script>

<style>
#app {
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
