<template>
  <div>
    <v-text-field :label="getSubscriptionTitle"></v-text-field>
    <TwUserList v-bind:userList="selectedUsers" v-on:removeUser="removeUser" />
    <v-autocomplete
      v-model="selected"
      :loading="loading"
      :items="twitterUsers"
      :search-input.sync="search"
      item-text="name"
      item-value="id"
      label="User name or screen name"
      hide-no-data
      no-filter
      auto-select-first
      @input="inputHandler"
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
  </div>
</template>

<script>
import axios from "axios";
import _ from "lodash";
import TwUserList from "./TwUserList";

export default {
  name: "Subscription",
  components: {
    TwUserList
  },
  props: {
    subscription: Object
  },
  computed: {
    getSubscriptionTitle: function() {
      return this.subscription ? this.subscription.title : "Subscription title";
    }
  },
  data: () => ({
    dialog: false,
    loading: false,
    search: null,
    selected: null,
    twitterUsers: [],
    query: null,
    selectedUsers: []
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
    },
    inputHandler(e) {
      const user = _.chain(this.twitterUsers)
        .filter({ id: e })
        .head()
        .value();

      if (user && -1 === _.findIndex(this.selectedUsers, ["id", user.id])) {
        this.selectedUsers.push(user);
        this.$nextTick(() => {
          this.selected = null;
        });
      }
    },
    removeUser: function(user) {
      this.selectedUsers = _.filter(this.selectedUsers, e => e.id !== user.id);
    }
  }
};
</script>