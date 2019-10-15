<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col class="text-center">
        <div>Updater provides possibility to get Facebook users updates via email.</div>
      </v-col>
    </v-row>
    <v-row align="center" justify="center">
      <v-col>
        <SignIn />
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import SignIn from "./SignIn";

export default {
  name: "Welcome",
  components: { SignIn }
};
</script>

<style>
</style>