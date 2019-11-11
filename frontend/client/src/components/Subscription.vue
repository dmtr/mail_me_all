<template>
  <v-card flat>
    <v-text-field v-model="title" :rules="titleRules" label="Subscription title"></v-text-field>
    <v-text-field v-model="email" :rules="emailRules" label="E-mail"></v-text-field>
    <v-select v-model="day" :items="days" label="Subscription delivery day"></v-select>
    <TwUserList v-bind:userList="userList" v-on:removeUser="removeUser" />
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
    <v-alert dense border="right" type="warning" v-if="!valid">{{ validationErrors}}</v-alert>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn text color="primary" @click="saveSubscription">Save</v-btn>
      <v-btn text color="primary" @click="cancelSubscriptionEdit">Cancel</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script>
import axios from "axios";
import _ from "lodash";
import { mapActions } from "vuex";
import TwUserList from "./TwUserList";

const days = [
  "Monday",
  "Tuedasy",
  "Wensday",
  "Thursday",
  "Friday",
  "Saturday",
  "Sunday"
];

const re = /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;

const validateEmail = e => {
  return re.test(e.toLowerCase());
};

const validateSubscription = s => {
  var errors = [];
  if (s.title.length === 0) {
    errors.push({ field: "title", msg: "Empty title" });
  }
  if (!validateEmail(s.email)) {
    errors.push({ field: "email", msg: "Invalid email" });
  }
  if (s.userList.length === 0) {
    errors.push({ field: "userList", msg: "Empty user list" });
  }
  if (_.indexOf(days, s.day) === -1) {
    errors.push({ field: "day", msg: "Invalid day" });
  }
  return errors;
};

export default {
  name: "Subscription",
  components: {
    TwUserList
  },
  props: {
    subscription: {
      type: Object,
      default: function() {
        return { id: null, title: "", email: "", day: null, userList: [] };
      }
    }
  },
  computed: {
    valid: {
      get: function() {
        return this.validationErrors.length > 0 ? false : true;
      },
      set: function(errors) {
        this.validationErrors = _.chain(errors)
          .map("msg")
          .join(", ")
          .value();
      }
    }
  },
  data: () => ({
    loading: false,
    search: null,
    selected: null,
    twitterUsers: [],
    query: null,
    id: null,
    title: "",
    day: null,
    days: days,
    userList: [],
    email: "",
    emailRules: [
      value => {
        if (value.length > 0) {
          return validateEmail(value) || "Invalid e-mail.";
        } else {
          return "";
        }
      }
    ],
    titleRules: [
      value => {
        if (value.length > 0) {
          return true;
        } else {
          return "Invalid title";
        }
      }
    ],
    validationErrors: []
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
  beforeMount: function() {
    this.init();
  },
  methods: {
    ...mapActions(["createSubscription", "getSubscriptions"]),

    init() {
      const customizer = (objValue, srcValue) => {
        return _.isArray(srcValue) ? [...srcValue] : srcValue;
      };
      _.assignWith(this, this.subscription, customizer);
    },

    querySelections() {
      if (!this.query) {
        return;
      }
      var this_ = this;
      this.loading = true;
      axios
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

      if (user && -1 === _.findIndex(this.userList, ["id", user.id])) {
        this.userList.push(user);
        this.$nextTick(() => {
          this.selected = null;
        });
      }
    },

    removeUser: function(user) {
      this.userList = _.filter(this.userList, e => e.id !== user.id);
    },

    getSubscriptionDetails: function({ id, title, email, day, userList }) {
      return { id, title, email, day, userList };
    },

    saveSubscription: async function() {
      var s = this.getSubscriptionDetails(this);
      this.valid = validateSubscription(s);

      if (this.valid) {
        var res = await this.createSubscription(s);
        if (!res.error) {
          await this.getSubscriptions();
          this.$emit("cancelSubscriptionEdit");
        } else {
          this.validationErrors = res.error.message;
        }
      }
    },

    cancelSubscriptionEdit: function() {
      this.init();
      this.valid = [];
      this.$emit("cancelSubscriptionEdit");
    }
  }
};
</script>