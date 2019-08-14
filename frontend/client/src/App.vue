<template>
  <div id="app">
    <div>
      <SignIn />
      <Subscription />
    </div>
  </div>
</template>

<script>
import Subscription from "./components/Subscription.vue";
import SignIn from "./components/SignIn.vue";

export default {
  name: "app",

  components: {
    Subscription,
    SignIn
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
}
</style>
