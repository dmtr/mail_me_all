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
          <v-card class="pa-md-4 mx-md-auto">
            <v-text-field label="Subscription title"></v-text-field>
            <v-autocomplete
              v-model="select"
              :loading="loading"
              :items="twitterUsers"
              :search-input.sync="search"
              :item-text="getItemText"
              item-value="id"
              label="User name or screen name"
              hide-no-data
              no-filter
              auto-select-first
            >
              <template v-slot:item="data">
                <v-list-item-avatar>
                  <img :src="data.item.profile_image_url" />
                </v-list-item-avatar>
                <v-list-item-content>
                  <v-list-item-title>{{data.item.name}}</v-list-item-title>
                  <v-list-item-title>{{data.item.screen_name}}</v-list-item-title>
                </v-list-item-content>
              </template>
            </v-autocomplete>
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
import _ from "lodash";
import Subscription from "./Subscription";
import axios from "axios";

export default {
  name: "SubscriptionsList",
  components: { Subscription },
  props: { subscriptionsList: Array },

  data: () => ({
    dialog: false,
    loading: false,
    search: null,
    select: null,
    twitterUsers: [],
    query: null
  }),
  watch: {
    search(val) {
      if (val && val !== this.select) {
        this.query = _.trim(val);
        this.debouncedQuery();
      }
    }
  },
  created: function() {
    this.debouncedQuery = _.debounce(this.querySelections, 200);
  },
  methods: {
    getItemText(item) {
      return `Name: ${item.name}, screen name: ${item.screen_name}`;
    },
    querySelections() {
      if (!this.query) {
        return;
      }
      var this_ = this;
      this.loading = true;
      const res = axios
        .get("api/twitter-users?q=" + this_.query)
        .then(function(response) {
          this_.twitterUsers = response.data.users;
          this_.loading = false;
        })
        .catch(function(error) {
          console.log(error);
          this_.loading = false;
        });
    }
  }
};
</script>