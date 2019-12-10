<template>
  <v-list dense color="#E1F5FE">
    <v-list-item>
      <v-list-item-content
        class="font-weight-bold"
      >Read-it-later.app makes it possible to get tweets via email.</v-list-item-content>
    </v-list-item>
    <v-list-item>
      <v-list-item-content>
        <v-list dense>
          <v-subheader>Three simple steps</v-subheader>
          <v-list-item>
            <v-list-item-icon>
              <v-icon>mdi-numeric-1-circle</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Sign In with Twitter</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-icon>
              <v-icon>mdi-numeric-2-circle</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>Choose users you are interested in</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
          <v-list-item>
            <v-list-item-icon>
              <v-icon>mdi-numeric-3-circle</v-icon>
            </v-list-item-icon>
            <v-list-item-content>
              <v-list-item-title>And get user updates via email</v-list-item-title>
            </v-list-item-content>
          </v-list-item>
        </v-list>
      </v-list-item-content>
    </v-list-item>
    <v-list-item>
      <v-list-item-content>
        <SignIn />
      </v-list-item-content>
    </v-list-item>
  </v-list>
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