<template>
  <v-app>
    <v-navigation-drawer v-model="drawer" app>
      <v-list dense>
        <v-list-item>
          <v-list-item-action>
            <v-icon>mdi-home</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-btn text to="/">Home</v-btn>
          </v-list-item-content>
        </v-list-item>
        <v-list-item>
          <v-list-item-action>
            <v-icon>mdi-account-settings</v-icon>
          </v-list-item-action>
          <v-list-item-content>
            <v-btn text to="/settings">Settings</v-btn>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-app-bar app color="blue" dark>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Updater</v-toolbar-title>
    </v-app-bar>
    <v-content>
      <router-view></router-view>
    </v-content>
    <v-footer color="blue" app>
      <span class="white--text">&copy; 2019</span>
    </v-footer>
  </v-app>
</template>

<script>
import { mapActions } from "vuex";

export default {
  name: "App",
  data: () => ({
    drawer: null
  }),
  methods: {
    ...mapActions(["getUser", "getSubscriptions"]),
    async loadData() {
      await this.getUser();
      if (this.isUserSignedIn) {
        await this.getSubscriptions();
      }
    }
  },
  created: function() {
    this.loadData();
  }
};
</script>