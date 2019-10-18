<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app>
      <v-list dense>
        <v-list-item>
          <v-list-item-action>
            <v-icon>mdi-home</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-list-item-title>Home</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app color="indigo" dark>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Updater</v-toolbar-title>
    </v-app-bar>
    <v-content>
      <Welcome v-if="notSignedIn" />
    </v-content>
    <v-footer color="indigo" app>
      <span class="white--text">&copy; 2019</span>
    </v-footer>
  </v-app>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import Welcome from "./components/Welcome";

export default {
  name: "App",
  computed: {
    ...mapGetters(["isUserSignedIn"]),
    notSignedIn() {
      return !this.isUserSignedIn;
    }
  },
  components: {
    Welcome
  },
  data: () => ({
    drawer: null
  })
};
</script>