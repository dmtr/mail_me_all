<template>
  <div v-if="isUserLoaded">
    <div v-if="isUserSignedIn">
      <SubscriptionsList v-bind:subscriptions="subscriptions" />
    </div>
    <div v-else>
      <Welcome />
    </div>
  </div>
  <div v-else>Loading...</div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import Welcome from "./Welcome";
import SubscriptionsList from "./SubscriptionsList";

export default {
  name: "Home",
  components: {
    Welcome,
    SubscriptionsList
  },
  data: () => ({}),
  computed: {
    ...mapGetters(["isUserSignedIn", "isUserLoaded", "subscriptions"])
  }
};
</script>