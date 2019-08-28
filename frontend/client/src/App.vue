<template>
  <div id="app">
    <div>
      <SignIn />
    </div>
  </div>
</template>

<script>
import SignIn from "./components/SignIn";

export default {
  name: "app",

  components: {
    SignInn
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
