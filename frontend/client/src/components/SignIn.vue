<template>
  <div>
    <button @click="authenticate('twitter')">auth Twitter</button>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: "SignIn",
  components: {},
  methods: {
    ...mapActions(["getUser"])
  }
};
</script>

<style></style>
