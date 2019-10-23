<template>
  <v-row>
    <v-col cols="12" sm="6" offset-sm="3">
      <v-card>
        <v-toolbar color="light-blue" light extended>
          <v-toolbar-title class="white--text">My subscriptions</v-toolbar-title>
          <template v-slot:extension>
            <v-btn fab color="cyan accent-2" bottom left absolute @click="dialog = !dialog">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </template>
        </v-toolbar>
        <v-list>
          <v-list-item
            v-for="subscription in subscriptionsList"
            :key="subscription.id"
            :subscription="subscription"
          >
            <v-list-item-content>
              <v-list-item-title v-text="subscription.title"></v-list-item-title>
              <!-- <Subscription v-bind:subscription="subscription" /> -->
            </v-list-item-content>
            <v-list-item-action>
              <v-btn icon>
                <v-icon color="grey lighten-1">mdi-playlist-edit</v-icon>
              </v-btn>
            </v-list-item-action>
            <v-list-item-action>
              <v-btn icon>
                <v-icon color="grey lighten-1">mdi-delete</v-icon>
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
        <v-dialog v-model="dialog" max-width="500px">
          <v-card>
            <v-card-text>
              <v-text-field label="Subscription title"></v-text-field>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="dialog = false">Cancel</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import Subscription from "./Subscription";
import axios from "axios";

export default {
  name: "SubscriptionsList",
  components: { Subscription },
  props: { subscriptionsList: Array },

  data: () => ({
    dialog: false
  })
};
</script>