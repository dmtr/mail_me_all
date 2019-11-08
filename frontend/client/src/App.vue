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
    <v-app-bar app color="blue" dark>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Updater</v-toolbar-title>
    </v-app-bar>
    <v-content>
      <div v-if="isUserLoaded">
        <div v-if="isUserSignedIn">
          <SubscriptionsList v-bind:subscriptions="subscriptions" />
        </div>
        <div v-else>
          <Welcome />
        </div>
      </div>
      <div v-else>Loading...</div>
    </v-content>
    <v-footer color="blue" app>
      <span class="white--text">&copy; 2019</span>
    </v-footer>
  </v-app>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import Welcome from "./components/Welcome";
import SubscriptionsList from "./components/SubscriptionsList";

export default {
  name: "App",
  computed: {
    ...mapGetters(["isUserSignedIn", "isUserLoaded", "subscriptions"])
  },
  components: {
    Welcome,
    SubscriptionsList
  },
  methods: {
    ...mapActions(["getUser", "getSubscriptions"]),
    async loadData() {
      await this.getUser();
      await this.getSubscriptions();
    }
  },
  created: function() {
    this.loadData();
  },
  data: () => ({
    drawer: null
  })
};
</script>