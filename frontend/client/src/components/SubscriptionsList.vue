<template>
  <v-list subheader>
    <v-subheader>My Subscriptions</v-subheader>
    <v-list-item
      v-for="subscription in subscriptionsList"
      :key="subscription.id"
      :subscription="subscription"
    >
      <v-list-item-content>
        <Subscription v-bind:subscription="subscription" />
      </v-list-item-content>
    </v-list-item>
  </v-list>
</template>

<script>
import Subscription from "./Subscription";

export default {
  name: "SubscriptionsList",
  components: { Subscription },
  props: { subscriptionsList: Array }
};
</script>