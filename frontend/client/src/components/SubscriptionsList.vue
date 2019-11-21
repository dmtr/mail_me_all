<template>
  <v-row>
    <v-col cols="12" sm="6" offset-sm="3">
      <v-card>
        <v-toolbar color="light-blue" light extended>
          <v-toolbar-title class="white--text">My subscriptions</v-toolbar-title>
          <template v-slot:extension>
            <v-btn fab color="cyan accent-2" bottom left absolute @click="newSubscription">
              <v-icon>mdi-plus</v-icon>
            </v-btn>
          </template>
        </v-toolbar>
        <v-list>
          <v-list-item
            v-for="subscription in subscriptions"
            :key="subscription.id"
            :subscription="subscription"
          >
            <v-list-item-content>
              <v-list-item-title v-text="subscription.title"></v-list-item-title>
            </v-list-item-content>
            <v-list-item-action>
              <v-btn @click="editSubscription(subscription)" icon>
                <v-icon color="grey lighten-1">mdi-playlist-edit</v-icon>
              </v-btn>
            </v-list-item-action>
            <v-list-item-action>
              <v-btn @click.stop="openRemoveDialog(subscription)" icon>
                <v-icon color="grey lighten-1">mdi-delete</v-icon>
              </v-btn>
            </v-list-item-action>
          </v-list-item>
        </v-list>
        <v-dialog v-model="dialog" max-width="500px">
          <v-card class="pa-md-4 mx-md-auto">
            <Subscription
              v-bind:subscription="currentSubscription"
              v-on:cancelSubscriptionEdit="cancelSubscriptionEdit"
              v-on:subscriptionSaved="subscriptionSaved"
              v-if="currentSubscription"
            ></Subscription>
            <Subscription
              v-on:cancelSubscriptionEdit="cancelSubscriptionEdit"
              v-on:subscriptionSaved="subscriptionSaved"
              v-else
            ></Subscription>
          </v-card>
        </v-dialog>
        <v-dialog v-model="removeDialog" max-width="500px">
          <v-card class="pa-md-4 mx-md-auto">
            <v-card-actions>
              Remove subscription?
              <v-spacer></v-spacer>
              <v-btn text color="primary" @click="removeSubscription()">Remove</v-btn>
              <v-btn text color="primary" @click="removeDialog=false">Cancel</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>
      </v-card>
    </v-col>
  </v-row>
</template>

<script>
import _ from "lodash";
import { mapActions } from "vuex";
import Subscription from "./Subscription";

export default {
  name: "SubscriptionsList",
  components: { Subscription },
  props: { subscriptions: Array },
  data: () => ({
    dialog: false,
    currentSubscription: null,
    removeDialog: false,
    toRemove: null
  }),
  methods: {
    ...mapActions(["deleteSubscription", "getSubscriptions"]),
    cancelSubscriptionEdit: function() {
      this.currentSubscription = null;
      this.dialog = false;
    },

    subscriptionSaved: async function() {
      await this.getSubscriptions();
      this.currentSubscription = this.subscriptions[0];
      this.dialog = false;
    },

    editSubscription: function(subscription) {
      console.debug(subscription);
      this.dialog = true;
      this.currentSubscription = _.cloneDeep(subscription);
    },

    newSubscription: function() {
      this.currentSubscription = null;
      this.dialog = true;
    },

    openRemoveDialog: function(subscription) {
      this.toRemove = subscription;
      this.removeDialog = true;
    },

    removeSubscription: async function() {
      if (this.toRemove) {
        await this.deleteSubscription(this.toRemove);
        this.toRemove = null;
      }
      this.removeDialog = false;
    }
  }
};
</script>